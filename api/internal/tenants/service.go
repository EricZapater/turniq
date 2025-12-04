package tenants

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

type Service interface {
	Create(ctx context.Context, request TenantRequest) (Tenant, error)
	GetAll(ctx context.Context) ([]Tenant, error)
	GetByID(ctx context.Context, id string) (Tenant, error)
	Update(ctx context.Context, id string, request TenantRequest) (Tenant, error)
	Delete(ctx context.Context, id string) error
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func(s *service) Create(ctx context.Context, request TenantRequest) (Tenant, error) {
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
	return s.repo.Create(ctx, tenant)
}

func(s *service) GetAll(ctx context.Context) ([]Tenant, error) {
	return s.repo.GetAll(ctx)
}

func(s *service) GetByID(ctx context.Context, id string) (Tenant, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Tenant{}, err
	}	
	return s.repo.GetByID(ctx, parsedId)
}

func(s *service) Update(ctx context.Context, id string, request TenantRequest) (Tenant, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Tenant{}, err
	}
	tenant, err := s.repo.GetByID(ctx, parsedId)
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
	return s.repo.Update(ctx, tenant)
}

func(s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedId)
}
