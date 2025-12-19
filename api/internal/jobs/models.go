package jobs

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID          uuid.UUID `json:"id"`
	CustomerID    uuid.UUID `json:"customer_id"`
	ShopFloorID uuid.UUID `json:"shop_floor_id"`
	WorkcenterID uuid.UUID `json:"workcenter_id"`
	JobCode     string    `json:"job_code"`
	ProductCode string    `json:"product_code"`
	Description string    `json:"description"`
	EstimatedDuration int `json:"estimated_duration"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type JobRequest struct {
	CustomerID    string `json:"customer_id"`
	ShopFloorID string `json:"shop_floor_id"`
	WorkcenterID string `json:"workcenter_id"`
	JobCode     string `json:"job_code"`
	ProductCode string `json:"product_code"`
	Description string `json:"description"`
	EstimatedDuration int `json:"estimated_duration"`
}