package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/pkg/jwt"
	"github.com/js-codegamer/fs-sync/pkg/logger"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, "Bearer ")
		if len(bearerToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		token := bearerToken[1]
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			logger.Sugar.Error(err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user, err := database.FindUserByUsername(claims.Username)
		if err != nil {
			logger.Sugar.Error(err.Error())
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
