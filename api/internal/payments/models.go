package payments

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `json:"id"`	
	CustomerID    uuid.UUID `json:"customer_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	DueDate       time.Time `json:"due_date"`
	PaidAt        time.Time `json:"paid_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PaymentRequest struct {
	CustomerID    string    `json:"customer_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	DueDate       time.Time `json:"due_date"`
	PaidAt        time.Time `json:"paid_at"`
}