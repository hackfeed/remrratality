package models

type Period struct {
	Filename    string `json:"filename" binding:"required" example:"filename.csv"`
	PeriodStart string `json:"period_start" binding:"required" example:"2019-01-01"`
	PeriodEnd   string `json:"period_end" binding:"required" example:"2021-01-01"`
}
