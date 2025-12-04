package customers

import (
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(request CustomerRequest) (Customer, error)
	GetAll() ([]Customer, error)
	GetByID(id string) (Customer, error)
	Update(id string, request CustomerRequest) (Customer, error)
	Delete(id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(request CustomerRequest) (Customer, error) {
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
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repository.Create(customer)
}

func (s *service) GetAll() ([]Customer, error) {
	return s.repository.GetAll()
}

func (s *service) GetByID(id string) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	return s.repository.GetByID(parsedId)
}

func (s *service) Update(id string, request CustomerRequest) (Customer, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Customer{}, err
	}
	customer, err := s.repository.GetByID(parsedId)
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
	return s.repository.Update(customer)
}

func (s *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(parsedId)
}
