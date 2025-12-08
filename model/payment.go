package model

type Payment struct {
	PaymentID   int     `json:"payment_id"`
	BookingID   int     `json:"booking_id"`
	Customer    string  `json:"customer"`
	Driver      string  `json:"driver"`
	Amount      float64 `json:"amount"`
	Method      string  `json:"method"`
	Status      string  `json:"status"`
	PaymentDate string  `json:"payment_date"`
}

type PaymentStats struct {
	TotalPayment      int64 `json:"total_payment"`
	PendingPayment    int64 `json:"pending_payment"`
	TotalTransactions int64 `json:"total_transactions"`
}
