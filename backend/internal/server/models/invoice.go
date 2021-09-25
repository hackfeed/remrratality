package models

type Invoice struct {
	CustomerId  uint32  `csv:"customer_id"`
	PeriodStart string  `csv:"period_start"`
	PaidPlan    string  `csv:"paid_plan"`
	PaidAmount  float32 `csv:"paid_amount"`
	PeriodEnd   string  `csv:"period_end"`
}
