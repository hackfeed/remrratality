package domain

import (
	"time"
)

type File struct {
	Name       string    `json:"name"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type User struct {
	UserID       string
	Email        string
	Password     string
	Token        string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Files        []File
}
