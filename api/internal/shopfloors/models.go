package shopfloors

import (
	"time"

	"github.com/google/uuid"
)

type Shopfloor struct {
	ID         uuid.UUID       `json:"id"`
	TenantID   uuid.UUID       `json:"tenant_id"`
	CustomerID uuid.UUID       `json:"customer_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ShopfloorRequest struct {
	TenantID   string    `json:"tenant_id"`
	CustomerID string    `json:"customer_id"`
	Name       string `json:"name"`
}
