package customers

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	VatNumber     string     `json:"vat_number"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       string     `json:"zip_code"`
	Country       string     `json:"country"`
	Language      string     `json:"language"`
	ContactName   string     `json:"contact_name"`
	Status        string     `json:"status"`
	Plan          string     `json:"plan"`           
	BillingCycle  string     `json:"billing_cycle"`  
	Price         float64    `json:"price"`          
	TrialEndsAt   *time.Time `json:"trial_ends_at"`
	InternalNotes string     `json:"internal_notes"`
	MaxOperators int `json:"max_operators"`
	MaxWorkcenters int `json:"max_workcenters"`
	MaxShopFloors int `json:"max_shop_floors"`
	MaxUsers int `json:"max_users"`
	MaxJobs int `json:"max_jobs"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type CustomerRequest struct {
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	VatNumber     string     `json:"vat_number"`
	Phone         string     `json:"phone"`
	Address       string     `json:"address"`
	City          string     `json:"city"`
	State         string     `json:"state"`
	ZipCode       string     `json:"zip_code"`
	Country       string     `json:"country"`
	Language      string     `json:"language"`
	ContactName   string     `json:"contact_name"`
	Status        string     `json:"status"`
	Plan          string     `json:"plan"`           
	BillingCycle  string     `json:"billing_cycle"`  
	Price         float64    `json:"price"`          
	TrialEndsAt   *time.Time `json:"trial_ends_at"`
	InternalNotes string     `json:"internal_notes"`
	MaxOperators int `json:"max_operators"`
	MaxWorkcenters int `json:"max_workcenters"`
	MaxShopFloors int `json:"max_shop_floors"`
	MaxUsers int `json:"max_users"`
	MaxJobs int `json:"max_jobs"`
}