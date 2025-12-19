package shifts

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID          uuid.UUID      `json:"id"`
	CustomerID  uuid.UUID      `json:"customer_id"`
	ShopfloorID uuid.NullUUID  `json:"shopfloor_id"`
	Name        string    `json:"name"`
	Color       string    `json:"color"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShiftRequest struct {
	CustomerID  string `json:"customer_id"`
	ShopfloorID string `json:"shopfloor_id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsActive    bool   `json:"is_active"`
}