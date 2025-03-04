package assets

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func AssetCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var asset models.Asset
		var err error

		if assetID := chi.URLParam(r, "assetID"); assetID != "" {
			asset, err = database.FindAssetByID(assetID)
		} else {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Resource not found", http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "asset", asset)

		if asset.Type != models.FolderType {
			file, err := database.FindLatestFileByAssetID(asset.ID)
			if err != nil {
				http.Error(w, "Resource not found", http.StatusNotFound)
				return
			}

			ctx = context.WithValue(ctx, "file", file)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
