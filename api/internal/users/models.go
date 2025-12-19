package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`	
	IsAdmin    bool      `json:"is_admin"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UserRequest struct {
	CustomerID string `json:"customer_id"`		
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`	
	IsActive   bool   `json:"is_active"`
}