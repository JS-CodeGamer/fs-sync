package utils

import (
	"fmt"
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

func MoveToVersionStorage(file models.File) error {
	if file.Path == "" {
		return fmt.Errorf("File does not exist")
	}

	versionBase := filepath.Join(config.GetConfig().Storage.BasePath, config.GetConfig().Storage.VersionPath)

	assetDir := filepath.Join(versionBase, file.AssetID)

	if _, err := os.Stat(assetDir); err != nil && os.IsNotExist(err) {
		os.MkdirAll(assetDir, 0750)
	}

	err := os.Rename(file.Path, filepath.Join(assetDir, strconv.FormatInt(file.Version, 10)))

	return err
}
