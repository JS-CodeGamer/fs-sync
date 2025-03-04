package utils

import (
	"errors"
	"net/http"
	"strings"
)

func GetAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Missing authorization header")
	}

	bearerToken := strings.Split(authHeader, "Bearer ")
	if len(bearerToken) != 2 {
		return "", errors.New("Invalid header format")
	}

	token := bearerToken[1]
	return token, nil
}
