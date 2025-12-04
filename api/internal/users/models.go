package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	TenantID   uuid.UUID `json:"tenant_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	CustomerID uuid.UUID `json:"customer_id"`
	IsAdmin    bool      `json:"is_admin"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRequest struct {
	TenantID   string `json:"tenant_id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CustomerID string `json:"customer_id"`	
	IsActive   bool   `json:"is_active"`
}