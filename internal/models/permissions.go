package models

type FilePermission struct {
	ID       string
	FileID   string
	UserID   string
	CanRead  bool
	CanWrite bool
}
