package workcenters

import (
	"time"

	"github.com/google/uuid"
)

type Workcenter struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkcenterRequest struct {
	Name string `json:"name"`
	TenantID string `json:"tenant_id"`
	CustomerID string `json:"customer_id"`
	IsActive bool `json:"is_active"`
}