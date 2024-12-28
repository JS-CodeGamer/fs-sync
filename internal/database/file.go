package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func CreateFileWithAsset(file *models.File, asset *models.Asset) error {
	txn, err := db.Begin()
	if err != nil {
		return fmt.Errorf("Error starting transaction: %w", err)
	}

	err = CreateAsset(asset, txn)
	if err != nil {
		txn.Rollback()
		return err
	}

	file.AssetID = asset.ID

	err = CreateFile(file, txn)

	if err != nil {
		txn.Rollback()
		return err
	} else if err = txn.Commit(); err != nil {
		return err
	}

	return nil
}

func CreateFile(file *models.File, txn *sql.Tx) error {
	file.ID = uuid.New().String()

	query := `
	INSERT INTO files
			(id, size, version, path, asset, mime_type)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query,
			file.ID,
			file.Size,
			file.Version,
			file.Path,
			file.AssetID,
			file.MimeType,
		)
	} else {
		_, err = txn.Exec(query,
			file.ID,
			file.Size,
			file.Version,
			file.Path,
			file.AssetID,
			file.MimeType,
		)
	}

	return err
}

func UpdateFile(file *models.File, txn *sql.Tx) error {
	query := `
	UPDATE files
		SET
			path = ?, mime_type = ?
		WHERE
			id = ?
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query,
			file.Path,
			file.MimeType,
			file.ID,
		)
	} else {
		_, err = txn.Exec(query,
			file.Path,
			file.MimeType,
			file.ID,
		)
	}

	return err
}

func FindLatestFileByAssetID(assetID string) (models.File, error) {
	query := `
	SELECT id, size, version, path, asset, mime_type
		FROM files
		WHERE
			asset = ?
		ORDER BY version DESC
		LIMIT 1
	`

	var file models.File
	err := db.QueryRow(query, assetID).Scan(
		&file.ID,
		&file.Size,
		&file.Version,
		&file.Path,
		&file.AssetID,
		&file.MimeType,
	)
	if err != nil {
		return models.File{}, err
	}

	return file, nil
}

func FindAllFilesByAssetID(assetID string) ([]models.File, error) {
	query := `
	SELECT id, size, version, path, asset, mime_type
		FROM files
		WHERE
			asset = ?
		ORDER BY version DESC
	`

	var files []models.File
	rows, err := db.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var file models.File
		err := rows.Scan(
			&file.ID,
			&file.Size,
			&file.Version,
			&file.Path,
			&file.AssetID,
			&file.MimeType,
		)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}

	return files, nil
}
