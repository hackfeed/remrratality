package domain

type Invoice struct {
	UserID      string
	FileID      string
	CustomerID  uint32
	PeriodStart string
	PaidPlan    string
	PaidAmount  float32
	PeriodEnd   string
}
