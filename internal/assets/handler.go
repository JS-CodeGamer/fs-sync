package assets

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/js-codegamer/fs-sync/config"
	"github.com/js-codegamer/fs-sync/internal/database"
	"github.com/js-codegamer/fs-sync/internal/models"
	"github.com/js-codegamer/fs-sync/internal/utils"
	"github.com/js-codegamer/fs-sync/pkg/logger"
	"github.com/js-codegamer/fs-sync/pkg/validator"
)

func NewAssetHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name     string `json:"name" validate:"required"`
		ParentID string `json:"parent_id" validate:"required,uuid4"`
		Size     int64  `json:"size"`
		IsDir    bool   `json:"is_dir"`
	}
	user := r.Context().Value("user").(models.User)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// validate request
	err := validator.GetValidator().Struct(request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.Name[len(request.Name)-1] == '/' {
		request.IsDir = true
		request.Name = strings.TrimPrefix(request.Name, "/")
	}

	if !request.IsDir && request.Size == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.ParentID == "" {
		request.ParentID = user.RootDirID
	}

	if !request.IsDir && request.Size > config.GetConfig().Storage.MaxFileSizeBytes {
		http.Error(w, "File size too big", http.StatusBadRequest)
		return
	}

	parent, err := database.FindAssetByID(request.ParentID)
	if err != nil || !parent.IsDir {
		fmt.Println(err)
		http.Error(w, "Parent does not exist", http.StatusConflict)
		return
	}

	asset := models.Asset{
		OwnerID:  user.ID,
		Name:     request.Name,
		ParentID: parent.ID,
		IsDir:    request.IsDir,
		Path:     filepath.Join(parent.Path, request.Name),
	}

	if !asset.IsDir {
		file := models.File{
			Size:    request.Size,
			Version: 1,
			AssetID: asset.ID,
		}
		err = database.CreateFileWithAsset(&file, &asset)
		if err != nil {
			logger.Sugar.Errorw("error creating new asset", "error", err)
			http.Error(w, "Error creating file", http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"message":    "success",
			"upload_url": fmt.Sprintf("/asset/%s", asset.ID),
		})
	} else {
		err = database.CreateAsset(&asset, nil)
		if err != nil {
			logger.Sugar.Errorw("error creating new asset", "error", err)
			http.Error(w, "Error creating user", http.StatusBadRequest)
			return
		}

		os.MkdirAll(asset.Path, 0750)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "success",
		})
	}
}

func MetadataUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Size int64 `json:"size"`
	}
	asset := r.Context().Value("asset").(models.Asset)
	file := r.Context().Value("file").(models.File)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// validate request
	err := validator.GetValidator().Struct(request)
	if err != nil || request.Size == 0 {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if asset.IsDir {
		http.Error(w, "Invalid resource update requested", http.StatusBadRequest)
		return
	}

	if err := utils.MoveToVersionStorage(file); err != nil {
		logger.Sugar.Errorw("error moving old version", "error", err)
		http.Error(w, "Error creating new version", http.StatusInternalServerError)
		return
	}

	newFile := models.File{
		Size:    request.Size,
		Version: file.Version + 1,
		AssetID: asset.ID,
	}
	err = database.CreateFile(&newFile, nil)
	if err != nil {
		logger.Sugar.Errorw("error creating file db object", "error", err)
		http.Error(w, "Error creating new version", http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message":    "success",
		"upload_url": fmt.Sprintf("/upload/%s", newFile.ID),
	})
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	asset := r.Context().Value("asset").(models.Asset)
	fileModel := r.Context().Value("file").(models.File)

	if fileModel.Path != "" {
		http.Error(w, "File already uploaded", http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(asset.Path); err == nil {
		http.Error(w, "File already exists", http.StatusInternalServerError)
		return
	}

	fileModel.Path = asset.Path
	if err := database.UpdateFile(&fileModel, nil); err != nil {
		logger.Sugar.Errorw("error updating file db object", "error", err)
		http.Error(w, "File save error", http.StatusInternalServerError)
		return
	}

	fileOnDisk, err := os.Create(asset.Path)
	if err != nil {
		logger.Sugar.Errorw("error creating file", "error", err)
		http.Error(w, "File creation error", http.StatusInternalServerError)
		return
	}
	defer fileOnDisk.Close()

	_, err = io.CopyN(fileOnDisk, r.Body, fileModel.Size)
	if err != nil {
		logger.Sugar.Errorw("error copying upload to file", "error", err)
		http.Error(w, "File save error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	file := r.Context().Value("file").(models.File)

	if file.Path == "" {
		http.Error(w, "File does not exist", http.StatusInternalServerError)
		return
	}

	fileOnDisk, err := os.Open(file.Path)
	if err != nil {
		logger.Sugar.Errorw("error opening file", "error", err)
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer fileOnDisk.Close()

	re := io.Reader(fileOnDisk)

	size := file.Size
	buff := make([]byte, min(1024, size))
	for {
		n, err := re.Read(buff)
		if err != nil {
			logger.Sugar.Errorw("error reading from file", "error", err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		w.Write(buff)
		size -= int64(n)
		if size == 0 {
			break
		}
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	asset := r.Context().Value("asset").(models.Asset)

	if asset.IsDir {
		if err := os.RemoveAll(asset.Path); err != nil && !os.IsNotExist(err) {
			logger.Sugar.Errorw("error removing dir", "error", err)
			http.Error(w, "Error deleting resource", http.StatusInternalServerError)
			return
		}
	} else {
		files, err := database.FindFilesByAssetID(asset.ID)
		if err != nil {
			logger.Sugar.Errorw("error finding files for asset", "error", err)
			http.Error(w, "Error deleting resource", http.StatusInternalServerError)
			return
		}

		for _, file := range files {
			if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
				logger.Sugar.Errorw("error removing file", "error", err)
				http.Error(w, "Error deleting resource", http.StatusInternalServerError)
				return
			}
		}
	}

	// no need to delete file objects due to foreign key cascade
	if err := database.DeleteAsset(asset, nil); err != nil {
		logger.Sugar.Errorw("error deleting db asset object", "error", err)
		http.Error(w, "Error deleting resource", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "success"})
}

func ListingHandler(w http.ResponseWriter, r *http.Request) {
	dir := r.Context().Value("asset").(models.Asset)

	contents, err := database.GetDirContents(dir)
	if err != nil {
		logger.Sugar.Errorw("error getting directory contents", "error", err)
		http.Error(w, "error getting directory listing", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(contents)
}
