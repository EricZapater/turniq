package jobs

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID `json:"id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	ShopFloorID uuid.UUID `json:"shop_floor_id"`
	WorkcenterID uuid.UUID `json:"workcenter_id"`
	JobCode     string    `json:"job_code"`
	ProductCode string    `json:"product_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JobRequest struct {
	TenantID    string `json:"tenant_id"`
	ShopFloorID string `json:"shop_floor_id"`
	WorkcenterID string `json:"workcenter_id"`
	JobCode     string `json:"job_code"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
}