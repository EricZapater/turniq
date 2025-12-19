package workcenters

import (
	"time"

	"github.com/google/uuid"
)

type Workcenter struct {
	ID        uuid.UUID `json:"id"`	
	CustomerID uuid.UUID `json:"customer_id"`
	ShopFloorID uuid.NullUUID `json:"shop_floor_id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WorkcenterRequest struct {
	Name string `json:"name"`	
	CustomerID string `json:"customer_id"`
	ShopFloorID string `json:"shop_floor_id"`
	IsActive bool `json:"is_active"`
}