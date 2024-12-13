package models

import (
	"time"
)

type File struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Filename  string    `json:"filename"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FilePermission struct {
	ID       int64 `json:"id"`
	FileID   int64 `json:"file_id"`
	UserID   int64 `json:"user_id"`
	CanRead  bool  `json:"can_read"`
	CanWrite bool  `json:"can_write"`
}
