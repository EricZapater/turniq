package customers

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request CustomerRequest) (Customer, error)
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id string) (Customer, error)
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
	tenantID, err := uuid.Parse(request.TenantID)
	if err != nil {
		return Customer{}, err
	}
	customer := Customer{
		ID:            uuid.New(),
		TenantID:      tenantID,
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
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repository.Create(ctx, customer)
}

func (s *service) GetAll(ctx context.Context) ([]Customer, error) {
	return s.repository.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	return s.repository.GetByID(ctx, parsedId)
}

func (s *service) Update(ctx context.Context, id string, request CustomerRequest) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	customer, err := s.repository.GetByID(ctx, parsedId)
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
