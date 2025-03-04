package auth

import (
	"context"
	"net/http"

	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/utils"
	"github.com/js-codegamer/fs-sync/pkg/jwt"
	"github.com/js-codegamer/fs-sync/pkg/logger"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := utils.GetAuthToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ValidateToken(token, jwt.AccessToken)
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
