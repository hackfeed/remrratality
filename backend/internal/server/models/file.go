package models

type File struct {
	Name string `json:"name" binding:"required"`
}
