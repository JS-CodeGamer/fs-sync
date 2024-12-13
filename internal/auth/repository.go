package auth

import (
	"database/sql"
	"errors"

	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func CreateUser(user models.User) error {
	db := database.GetConnection()
	query := `INSERT INTO users (username, password, email) VALUES (?, ?, ?)`

	_, err := db.Exec(query, user.Username, user.Password, user.Email)
	if err != nil {
		return err
	}
	return nil
}

func FindUserByUsername(username string) (models.User, error) {
	db := database.GetConnection()
	var user models.User

	query := `SELECT id, username, password, email FROM users WHERE username = ?`
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Email)

	if err == sql.ErrNoRows {
		return user, errors.New("user not found")
	}
	if err != nil {
		return user, err
	}

	return user, nil
}
