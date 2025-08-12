package presentation

import (
	"bytes"
	"encoding/json"
	"io"
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

type RegisterRequest struct {
	Password string `json:"password"`
	*domain.UserPayload
}

func (auth *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var register RegisterRequest
	json.NewDecoder(r.Body).Decode(&register)

	// Store user credentials in a persistent database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Unable to hash password: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var authUser domain.AuthUser = domain.AuthUser{
		Username:       register.Username,
		HashedPassword: string(hashedPassword),
	}
	authUser, err = auth.ur.Create(authUser)
	if err != nil {
		http.Error(w, "Unable to store user credentials: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data that will be sent to the other service
	var payloadUser domain.UserPayload = domain.UserPayload{
		AuthUserID: authUser.ID,
		Username:   register.Username,
		Role:       register.Role,
		Email:      register.Email,
	}
	payloadBytes, err := json.Marshal(payloadUser)
	if err != nil {
		http.Error(w, "Unable to marshal payload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := http.Post("http://localhost:8000/users/", "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to make request: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		err := auth.ur.Delete(authUser)
		if err != nil {
			http.Error(w, "Unable to delete user credentials: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(responseBody)
}
