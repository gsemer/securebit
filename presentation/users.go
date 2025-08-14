package presentation

import (
	"encoding/json"
	"log"
	"net/http"
	"securebit/domain"

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

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if registerRequest.Username == "" || registerRequest.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

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
	var loginRequest domain.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

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

	w.WriteHeader(http.StatusOK)
}
