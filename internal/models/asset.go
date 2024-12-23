package models

import (
	"time"
)

type Asset struct {
	ID        string
	OwnerID   string
	Name      string
	ParentID  string
	IsPublic  bool
	IsDir     bool
	Path      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type File struct {
	ID       string
	Size     int64
	Version  int64
	Path     string
	AssetID  string
	MimeType string
}
