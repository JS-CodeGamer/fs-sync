package models

import (
	"time"
)

type AssetType string

const (
	FolderType AssetType = "folder"
	FileType             = "file"
)

type Asset struct {
	ID        string
	OwnerID   string
	Name      string
	ParentID  string
	IsPublic  bool
	Type      AssetType
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PublicAsset struct {
	ID        string        `json:"id"`
	OwnerID   string        `json:"owner"`
	Name      string        `json:"name"`
	ParentID  string        `json:"parent_id"`
	IsPublic  bool          `json:"is_public"`
	Type      AssetType     `json:"type"`
	Children  []PublicAsset `json:"children"`
	Files     []File        `json:"files"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func ToPublicAsset(asset *Asset) PublicAsset {
	return PublicAsset{
		ID:        asset.ID,
		OwnerID:   asset.OwnerID,
		Name:      asset.Name,
		ParentID:  asset.ParentID,
		IsPublic:  asset.IsPublic,
		Type:      asset.Type,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
		Children:  nil,
		Files:     nil,
	}
}

type File struct {
	ID       string
	Size     int64
	Version  int64
	Path     string
	AssetID  string
	MimeType string
}
