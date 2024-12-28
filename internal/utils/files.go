package utils

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/js-codegamer/fs-sync/config"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func GetRootDir(user models.User) string {
	base := config.GetConfig().Storage.BasePath
	return filepath.Join(base, user.Username)
}

func CreateRootDir(user models.User) error {
	base := GetRootDir(user)
	return os.MkdirAll(base, 0750)
}

func DestroyRootDir(user models.User) error {
	base := GetRootDir(user)
	return os.RemoveAll(base)
}

func GetVersionPath(file models.File) string {
	return filepath.Join(config.GetConfig().Storage.BasePath, config.GetConfig().Storage.VersionPath, file.AssetID, strconv.FormatInt(file.Version, 10))
}
