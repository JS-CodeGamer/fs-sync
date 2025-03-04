package routes

import (
	"net/http"

	"github.com/js-codegamer/fs-sync/internal/assets"
	"github.com/js-codegamer/fs-sync/internal/auth"
	"github.com/js-codegamer/fs-sync/pkg/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestLogger(logger.NewChiLoggerWithZap(logger.Logger)))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	}))

	// Add documentation routes
	AddDocRoutes(r)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", auth.RegisterHandler)
		r.Post("/login", auth.LoginHandler)
		r.Get("/refresh-token", auth.RefreshTokenHandler)

		// me routes
		r.Route("/me", func(r chi.Router) {
			r.Use(auth.AuthMiddleware)
			r.Get("/", auth.GetProfileHandler)
			r.Post("/", auth.UpdateProfileHandler)
			r.Delete("/", auth.DeleteUserHandler)
		})

		// asset routes
		r.Route("/asset", func(r chi.Router) {
			r.With(auth.AuthMiddleware).Post("/", assets.CreateHandler)

			r.Route("/{assetID}", func(r chi.Router) {
				r.Use(assets.AssetCtx)
				r.Get("/", assets.ReadHandler)
				r.Put("/", assets.UpdateHandler)
				r.Patch("/", assets.UploadHandler)
				r.Delete("/", assets.DeleteHandler)
			})
		})
	})

	return r
}
