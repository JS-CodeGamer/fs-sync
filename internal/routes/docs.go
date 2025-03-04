package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
)

// GenerateAPIDoc returns a handler that serves the API documentation
func GenerateAPIDoc(router *chi.Mux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doc := docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
			ProjectPath: "github.com/js-codegamer/fs-sync",
			Intro:       "Welcome to the API documentation",
		})
		w.Header().Set("Content-Type", "text/markdown")
		w.Write([]byte(doc))
	}
}

// GenerateAPIJSON returns a handler that serves the API documentation in JSON format
func GenerateAPIJSON(router *chi.Mux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		doc := docgen.JSONRoutesDoc(router)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(doc)
	}
}

// AddDocRoutes adds documentation routes to the router
func AddDocRoutes(r *chi.Mux) {
	// Serve documentation in both Markdown and JSON formats
	r.Get("/docs/markdown", GenerateAPIDoc(r))
	r.Get("/docs/json", GenerateAPIJSON(r))
}
