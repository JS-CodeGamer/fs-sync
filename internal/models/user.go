package models

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Won't be serialized
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}