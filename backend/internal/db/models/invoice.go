package models

import "time"

var AllFields = []string{
	"user_id",
	"file_id",
	"customer_id",
	"period_start",
	"paid_plan",
	"paid_amount",
	"period_end",
}

type Invoice struct {
	UserID      string
	FileID      string
	CustomerID  uint32
	PeriodStart time.Time
	PaidPlan    string
	PaidAmount  float32
	PeriodEnd   time.Time
}
