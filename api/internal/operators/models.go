package operators

import (
	"time"

	"github.com/google/uuid"
)

type Operator struct {
	ID          uuid.UUID     `json:"id"`	
	ShopFloorID uuid.UUID     `json:"shop_floor_id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	Code        string        `json:"code"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	VatNumber   string    `json:"vat_number"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OperatorRequest struct {
	ShopFloorID uuid.UUID `json:"shop_floor_id"`
	CustomerID  string `json:"customer_id"`
	Code        string `json:"code"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	VatNumber   string    `json:"vat_number"`
	IsActive    bool      `json:"is_active"`
}