package customers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request CustomerRequest) (Customer, error)
	FindAll(ctx context.Context) ([]Customer, error)
	FindByID(ctx context.Context, id string) (Customer, error)
	Update(ctx context.Context, id string, request CustomerRequest) (Customer, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, request CustomerRequest) (Customer, error) {
	customer := Customer{
		ID:            uuid.New(),
		Name:          request.Name,
		Email:         request.Email,
		VatNumber:     request.VatNumber,
		Phone:         request.Phone,
		Address:       request.Address,
		City:          request.City,
		State:         request.State,
		ZipCode:       request.ZipCode,
		Country:       request.Country,
		Language:      request.Language,
		ContactName:   request.ContactName,
		Status:        request.Status,
		Plan:          request.Plan,
		BillingCycle:  request.BillingCycle,
		Price:         request.Price,
		TrialEndsAt:   request.TrialEndsAt,
		InternalNotes: request.InternalNotes,
		MaxOperators:  request.MaxOperators,
		MaxWorkcenters: request.MaxWorkcenters,
		MaxShopFloors: request.MaxShopFloors,
		MaxUsers: request.MaxUsers,
		MaxJobs: request.MaxJobs,	
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repository.Create(ctx, customer)
}

func (s *service) FindAll(ctx context.Context) ([]Customer, error) {
	return s.repository.FindAll(ctx)
}

func (s *service) FindByID(ctx context.Context, id string) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	return s.repository.FindByID(ctx, parsedId)
}

func (s *service) Update(ctx context.Context, id string, request CustomerRequest) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	customer, err := s.repository.FindByID(ctx, parsedId)
	if err != nil {
		return Customer{}, err
	}
	customer.Name = request.Name
	customer.Email = request.Email
	customer.VatNumber = request.VatNumber
	customer.Phone = request.Phone
	customer.Address = request.Address
	customer.City = request.City
	customer.State = request.State
	customer.ZipCode = request.ZipCode
	customer.Country = request.Country
	customer.Language = request.Language
	customer.ContactName = request.ContactName
	customer.Status = request.Status
	customer.Plan = request.Plan
	customer.BillingCycle = request.BillingCycle
	customer.Price = request.Price
	customer.TrialEndsAt = request.TrialEndsAt
	customer.InternalNotes = request.InternalNotes
	customer.MaxOperators = request.MaxOperators
	customer.MaxWorkcenters = request.MaxWorkcenters
	customer.MaxShopFloors = request.MaxShopFloors
	customer.MaxUsers = request.MaxUsers
	customer.MaxJobs = request.MaxJobs
	customer.UpdatedAt = time.Now()
	return s.repository.Update(ctx, customer)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(ctx, parsedId)
}
