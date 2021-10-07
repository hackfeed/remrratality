package models

import "github.com/hackfeed/remrratality/backend/internal/domain"

type ResponseSuccessLoadFiles struct {
	Message string        `json:"message" example:"Files are loaded"`
	Files   []domain.File `json:"files"`
}

type ResponseSuccessSaveFileContent struct {
	Message  string `json:"message" example:"File is uploaded"`
	Filename string `json:"filename" example:"filename.csv"`
}

type ResponseSuccessAuth struct {
	Message   string `json:"message"`
	IDToken   string `json:"id_token"`
	LocalID   string `json:"local_id"`
	ExpiresAt int64  `json:"expires_at"`
}

type ResponseSuccessAnalytics struct {
	Message string          `json:"message" example:"Analytics is loaded"`
	Months  []string        `json:"months"`
	MRR     domain.TotalMRR `json:"mrr"`
}

type Response struct {
	Message string `json:"message"`
}
