package tenants

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {	
	ID   uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Name string    `json:"name"`
	Status string `json:"status"`
	IsActive bool `json:"is_active"`
	MaxOperators int `json:"max_operators"`
	MaxWorkcenters int `json:"max_workcenters"`
	MaxShopFloors int `json:"max_shop_floors"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TenantRequest struct {
	CustomerID uuid.UUID `json:"customer_id"`
	Name string `json:"name"`
	Status string `json:"status"`
	IsActive bool `json:"is_active"`
	MaxOperators int `json:"max_operators"`
	MaxWorkcenters int `json:"max_workcenters"`
	MaxShopFloors int `json:"max_shop_floors"`
}
