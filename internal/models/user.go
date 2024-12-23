package models

import (
	"time"
)

type User struct {
	ID        string
	Username  string
	Email     string
	Password  string
	RootDirID string
	CreatedAt time.Time
}
