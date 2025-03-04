package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/js-codegamer/fs-sync/internal/models"
)

func CreateAsset(asset *models.Asset, txn *sql.Tx) error {
	asset.ID = uuid.New().String()

	query := `
	INSERT INTO assets
			(id, owner, name, parent, type, path)
		VALUES
			(?, ?, ?, ?, ?, ?)
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query,
			asset.ID,
			asset.OwnerID,
			asset.Name,
			asset.ParentID,
			asset.Type,
			asset.Path,
		)
	} else {
		_, err = txn.Exec(query,
			asset.ID,
			asset.OwnerID,
			asset.Name,
			asset.ParentID,
			asset.Type,
			asset.Path,
		)
	}

	if err != nil {
		return fmt.Errorf("Error creating asset: %w", err)
	}

	return err
}

func UpdateAsset(asset *models.Asset, txn *sql.Tx) error {
	query := `
	UPDATE assets
		SET
			name = ?, path = ?
		WHERE
			id = ?
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query,
			asset.Name,
			asset.Path,
			asset.ID,
		)
	} else {
		_, err = txn.Exec(query,
			asset.Name,
			asset.Path,
			asset.ID,
		)
	}

	return err
}

func FindAssetByID(assetID string) (models.Asset, error) {
	query := `
	SELECT id, owner, name, parent, is_public, type, path, created_at, updated_at
		FROM assets
		WHERE id = ?
	`

	var asset models.Asset
	err := db.QueryRow(query, assetID).Scan(
		&asset.ID,
		&asset.OwnerID,
		&asset.Name,
		&asset.ParentID,
		&asset.IsPublic,
		&asset.Type,
		&asset.Path,
		&asset.CreatedAt,
		&asset.UpdatedAt,
	)
	if err != nil {
		return models.Asset{}, err
	}

	return asset, nil
}

func DeleteAsset(asset models.Asset, txn *sql.Tx) error {
	query := `
	DELETE FROM assets
	WHERE
		id = $1 AND owner = $2
	-- cascade will delete all children
	`

	var err error
	if txn == nil {
		_, err = db.Exec(query, asset.ID, asset.OwnerID)
	} else {
		_, err = txn.Exec(query, asset.ID, asset.OwnerID)
	}

	return err
}

func GetDirContents(asset models.Asset) ([]models.Asset, error) {
	query := `
	SELECT id, owner, name, parent, is_public, type, path, created_at, updated_at
		FROM assets
		WHERE parent = ?`

	rows, err := db.Query(query, asset.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID,
			&asset.OwnerID,
			&asset.Name,
			&asset.ParentID,
			&asset.IsPublic,
			&asset.Type,
			&asset.Path,
			&asset.CreatedAt,
			&asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		contents = append(contents, asset)
	}
	return contents, nil
}
