package presentation

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"securebit/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	ur domain.UserRepository
}

func NewAuthHandler(ur domain.UserRepository) *AuthHandler {
	return &AuthHandler{
		ur: ur,
	}
}

type RegisterRequest struct {
	Password string `json:"password"`
	*domain.UserPayload
}

func (auth *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var register RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic validation
	if register.Username == "" || register.Password == "" || register.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	authUser := domain.AuthUser{
		Username:       register.Username,
		HashedPassword: string(hashedPassword),
	}

	createdUser, err := auth.ur.Create(authUser)
	if err != nil {
		log.Printf("User creation error: %v", err)
		http.Error(w, "Unable to store user credentials", http.StatusInternalServerError)
		return
	}

	payloadUser := domain.UserPayload{
		AuthUserID: createdUser.ID,
		Username:   register.Username,
		Role:       register.Role,
		Email:      register.Email,
	}

	payloadBytes, err := json.Marshal(payloadUser)
	if err != nil {
		log.Printf("Payload marshal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Post("http://localhost:8000/users/", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil || response.StatusCode >= 400 {
		// Rollback user creation on failure
		rollbackErr := auth.ur.Delete(createdUser)
		if rollbackErr != nil {
			log.Printf("Failed to rollback user ID %d: %v", createdUser.ID, rollbackErr)
		}

		if err != nil {
			http.Error(w, "Failed to make request: "+err.Error(), http.StatusInternalServerError)
		} else {
			defer response.Body.Close()
			responseBody, _ := io.ReadAll(response.Body)
			http.Error(w, "User service error: "+string(responseBody), response.StatusCode)
		}
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var register LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if register.Username == "" || register.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	authUser, err := auth.ur.Get(register.Username)
	if err != nil {
		log.Printf("No such a user exists: %v", register.Username)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authUser.HashedPassword), []byte(register.Password)); err != nil {
		log.Printf("Password and hashed password do not match")
		http.Error(w, "Wrong password for user "+authUser.Username, http.StatusBadRequest)
		return
	}

	authUserBytes, err := json.Marshal(authUser)
	if err != nil {
		log.Printf("Auth user marshal error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(authUserBytes)
}
