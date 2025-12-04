package tenants

import (
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Service interface {
	Create(request TenantRequest) (Tenant, error)
	GetAll() ([]Tenant, error)
	GetByID(id string) (Tenant, error)
	Update(id string, request TenantRequest) (Tenant, error)
	Delete(id string) error
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func(s *service) Create(request TenantRequest) (Tenant, error) {
	tenant := Tenant{
		ID:   uuid.New(),
		CustomerID: request.CustomerID,
		Name: request.Name,
		Status: request.Status,
		IsActive: request.IsActive,
		MaxOperators: request.MaxOperators,
		MaxWorkcenters: request.MaxWorkcenters,
		MaxShopFloors: request.MaxShopFloors,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(tenant)
}

func(s *service) GetAll() ([]Tenant, error) {
	return s.repo.GetAll()
}

func(s *service) GetByID(id string) (Tenant, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Tenant{}, err
	}	
	return s.repo.GetByID(parsedId)
}

func(s *service) Update(id string, request TenantRequest) (Tenant, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Tenant{}, err
	}
	tenant, err := s.repo.GetByID(parsedId)
	if err != nil {
		return Tenant{}, err
	}
	tenant.CustomerID = request.CustomerID
	tenant.Name = request.Name
	tenant.Status = request.Status
	tenant.IsActive = request.IsActive
	tenant.MaxOperators = request.MaxOperators
	tenant.MaxWorkcenters = request.MaxWorkcenters
	tenant.MaxShopFloors = request.MaxShopFloors
	tenant.UpdatedAt = time.Now()
	return s.repo.Update(tenant)
}

func(s *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(parsedId)
}
