package models

import "github.com/hackfeed/remrratality/backend/internal/domain"

type ResponseSuccessLoadFiles struct {
	Message string        `json:"message" example:"Files are loaded"`
	Files   []domain.File `json:"files"`
}

type ResponseFailLoadFiles struct {
	Message string `json:"message" example:"Unable to determine logged in user"`
}
