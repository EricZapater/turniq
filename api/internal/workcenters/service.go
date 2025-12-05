package workcenters

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request WorkcenterRequest) (Workcenter, error)
	FindAll(ctx context.Context) ([]Workcenter, error)
	FindByID(ctx context.Context, id string) (*Workcenter, error)
	FindByTenantID(ctx context.Context, tenantID string) ([]Workcenter, error)
	Update(ctx context.Context, id string, request WorkcenterRequest) (Workcenter, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request WorkcenterRequest) (Workcenter, error) {
	tenantParsedId, err := uuid.Parse(request.TenantID)
	if err != nil {
		return Workcenter{}, err
	}
	customerParsedId, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Workcenter{}, err
	}
	workcenter := Workcenter{
		ID:        uuid.New(),
		TenantID:  tenantParsedId,
		CustomerID: customerParsedId,
		Name:      request.Name,
		IsActive:  request.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(ctx, workcenter)
}

func (s *service) FindAll(ctx context.Context) ([]Workcenter, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) FindByID(ctx context.Context, id string) (*Workcenter, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, parsedId)
}

func (s *service) FindByTenantID(ctx context.Context, tenantID string) ([]Workcenter, error) {
	parsedTenantId, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, err
	}	
	return s.repo.FindByTenantID(ctx, parsedTenantId)
}

func (s *service) Update(ctx context.Context, id string, request WorkcenterRequest) (Workcenter, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Workcenter{}, err
	}
	workcenter, err := s.repo.FindByID(ctx, parsedId)
	if err != nil {
		return Workcenter{}, err
	}
	workcenter.Name = request.Name
	workcenter.IsActive = request.IsActive
	workcenter.UpdatedAt = time.Now()
	return s.repo.Update(ctx, workcenter)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedId)
}
