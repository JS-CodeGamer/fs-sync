package files

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/js-codegamer/fs-sync/internal/auth"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	user, _ := auth.FindUserByUsername(username)

	r.ParseMultipartForm(10 << 20) // 10 MB max file size
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "File upload error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileModel := models.File{
		UserID:    user.ID,
		Filename:  handler.Filename,
		Size:      handler.Size,
		CreatedAt: time.Now(),
	}

	dst, err := os.Create(filepath.Join("uploads", handler.Filename))
	if err != nil {
		http.Error(w, "File creation error", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "File save error", http.StatusInternalServerError)
		return
	}

	CreateFile(fileModel)
	json.NewEncoder(w).Encode(fileModel)
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	user, _ := auth.FindUserByUsername(username)

	files, err := GetUserFiles(user.ID)
	if err != nil {
		http.Error(w, "Cannot retrieve files", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(files)
}
