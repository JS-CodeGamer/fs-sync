package auth

import (
	"encoding/json"
	"net/http"

	"github.com/js-codegamer/fs-sync/internal/models"
	"github.com/js-codegamer/fs-sync/pkg/jwt"
	"github.com/js-codegamer/fs-sync/pkg/password"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := password.HashPassword(user.Password)
	if err != nil {
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	if err := CreateUser(user); err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := FindUserByUsername(loginRequest.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !password.CheckPasswordHash(loginRequest.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
