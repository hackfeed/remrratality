package models

import "github.com/hackfeed/remrratality/backend/internal/domain"

type ResponseSuccessLoadFiles struct {
	Message string        `json:"message" example:"Files are loaded"`
	Files   []domain.File `json:"files"`
}

type ResponseSuccessSaveFile struct {
	Message  string `json:"message" example:"Files is uploaded"`
	Filename string `json:"filename"`
}

type Response struct {
	Message string `json:"message"`
}
