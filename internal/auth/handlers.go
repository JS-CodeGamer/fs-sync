package auth

import (
	"encoding/json"
	"net/http"

	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/models"
	"github.com/js-codegamer/fs-sync/pkg/jwt"
	"github.com/js-codegamer/fs-sync/pkg/logger"
	"github.com/js-codegamer/fs-sync/pkg/password"
	"github.com/js-codegamer/fs-sync/pkg/validator"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerRequest struct {
		Username string `json:"username" validator:"required"`
		Password string `json:"password" validator:"required"`
		Email    string `json:"email" validator:"required,email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := validator.GetValidator().Struct(registerRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := password.HashPassword(registerRequest.Password)
	if err != nil {
		logger.Sugar.Errorw("error hashing password", "error", err)
		http.Error(w, "Password hashing failed", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username: registerRequest.Username,
		Password: hashedPassword,
		Email:    registerRequest.Email,
	}

	if _, err := database.CreateUser(user, nil); err != nil {
		logger.Sugar.Errorw("error creating db user", "error", err)
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		logger.Sugar.Errorw("error creating token", "error", err)
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username" validator:"required"`
		Password string `json:"password" validator:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := validator.GetValidator().Struct(loginRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := database.FindUserByUsername(loginRequest.Username)
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
		logger.Sugar.Errorw("error creating token", "error", err)
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	profile := struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		RootDirID string `json:"root_dir"`
	}{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		RootDirID: user.RootDirID,
	}

	json.NewEncoder(w).Encode(profile)
}

func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	var updateRequest struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if updateRequest.NewPassword != "" {
		if !password.CheckPasswordHash(updateRequest.OldPassword, user.Password) {
			http.Error(w, "Invalid current password", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := password.HashPassword(updateRequest.NewPassword)
		if err != nil {
			logger.Sugar.Errorw("error hashing password", "error", err)
			http.Error(w, "Password hashing failed", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
	}

	if updateRequest.Email != "" {
		user.Email = updateRequest.Email
	}

	if err := database.UpdateUser(user, nil); err != nil {
		logger.Sugar.Errorw("error updating db user", "error", err)
		http.Error(w, "Profile update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	if err := database.DeleteUser(user, nil); err != nil {
		logger.Sugar.Errorw("error deleteing db user", "error", err)
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Profile deleted successfully"})
}
