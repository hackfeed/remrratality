package models

type User struct {
	Email    *string `json:"email" validate:"email,required" binding:"required" example:"test@test.com"`
	Password *string `json:"password" validate:"required,min=6" binding:"required" example:"password123"`
}
