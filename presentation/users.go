package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"securebit/config"
	"securebit/domain"
	"securebit/utils"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	ur          domain.UserRepository
	redisClient *redis.Client
}

func NewAuthHandler(ur domain.UserRepository, redisClient *redis.Client) *AuthHandler {
	return &AuthHandler{
		ur:          ur,
		redisClient: redisClient,
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Request validation
	var registerRequest domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	if registerRequest.Username == "" || registerRequest.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Hash the password & create user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash the password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user := domain.User{
		Username:       registerRequest.Username,
		HashedPassword: string(hashedPassword),
	}
	userDB, err := ah.ur.Create(user)
	if err != nil {
		log.Printf("User creation error: %v", err)
		http.Error(w, "Unable to store user credentials", http.StatusInternalServerError)
		return
	}

	userBytes, err := json.Marshal(userDB)
	if err != nil {
		log.Printf("Failed to marshal the user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(userBytes)
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Request validation
	var loginRequest domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}
	if loginRequest.Username == "" || loginRequest.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// User credentials validation
	user, err := ah.ur.Get(loginRequest.Username)
	if err != nil {
		log.Printf("User %v not found in database", loginRequest.Username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginRequest.Password)); err != nil {
		log.Printf("Password and hashed password do not match")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create access token (short-lived)
	signedAccessToken, err := utils.SignedToken(user, time.Now().Add(5*time.Minute), config.GetEnv("JWT_SECRET_KEY", ""))
	if errors.Is(err, domain.ErrTokenSigningFailed) {
		log.Print("Failed to sign access token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Create refresh token (longer-lived)
	signedRefreshToken, err := utils.SignedToken(user, time.Now().Add(24*time.Hour), config.GetEnv("JWT_REFRESH_SECRET_KEY", ""))
	if errors.Is(err, domain.ErrTokenSigningFailed) {
		log.Print("Failed to sign refresh token")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set refresh token to Cookie
	// Client can't use it directly but browser send it automatically for refresh
	// Javascript can't access it, so it is protected against XSS attacks
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    signedRefreshToken,
		HttpOnly: true,
		Secure:   false,             // Set to true in production with HTTPS
		Path:     "/api/v1/refresh", // Limit cookie to refresh token endpoint
		MaxAge:   24 * 3600,         // 24 hours in seconds
		SameSite: http.SameSiteStrictMode,
	})
	// Store refresh token in Redis
	ah.redisClient.Set(context.Background(), user.ID.String(), signedRefreshToken, 24*time.Hour)

	w.Header().Set("Authorization", "Bearer "+signedAccessToken)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user logged in"))
}
