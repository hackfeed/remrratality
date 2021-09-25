package models

type User struct {
	Email    *string `json:"email" validate:"email,required" binding:"required"`
	Password *string `json:"password" validate:"required,min=6" binding:"required"`
}
