package routes

import (
	"net/http"

	"github.com/js-codegamer/fs-sync/config"
	"github.com/js-codegamer/fs-sync/internal/auth"
	"github.com/js-codegamer/fs-sync/internal/files"
	"github.com/js-codegamer/fs-sync/pkg/logger"
)

func SetupRoutes(cfg *config.Config) http.HandlerFunc {
	mux := http.NewServeMux()

	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)

	mux.HandleFunc("/upload", auth.AuthMiddleware(files.UploadHandler))
	mux.HandleFunc("/files", auth.AuthMiddleware(files.ListFilesHandler))
	// mux.HandleFunc("/delete", auth.AuthMiddleware(files.DeleteFileHandler))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Sugar.Infof("Request: %s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}
