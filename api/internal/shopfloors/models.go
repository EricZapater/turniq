package shopfloors

import (
	"time"

	"github.com/google/uuid"
)

type Shopfloor struct {
	ID         uuid.UUID       `json:"id"`		
	CustomerID uuid.UUID       `json:"customer_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ShopfloorRequest struct {
	CustomerID string    `json:"customer_id"`
	Name       string `json:"name"`
}
