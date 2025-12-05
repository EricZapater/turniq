package shifts

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	CustomerID  uuid.UUID `json:"customer_id"`
	ShopfloorID uuid.UUID `json:"shopfloor_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShiftRequest struct {	
	TenantID    string `json:"tenant_id"`
	CustomerID  string `json:"customer_id"`
	ShopfloorID string `json:"shopfloor_id"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsActive    bool   `json:"is_active"`	
}