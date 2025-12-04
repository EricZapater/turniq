package shopfloors

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request ShopfloorRequest) (Shopfloor,error)
	FindAll(ctx context.Context) ([]Shopfloor, error)
	FindByID(ctx context.Context, id string) (Shopfloor, error)
	FindByTenantID(ctx context.Context, id string) ([]Shopfloor, error)
	Update(ctx context.Context, id string, request ShopfloorRequest) (Shopfloor,error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, request ShopfloorRequest) (Shopfloor,error){
	parsedTenantId, err := uuid.Parse(request.TenantID)
	if err != nil {
		return Shopfloor{}, err
	}
	parsedCustomerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Shopfloor{}, err
	}

	shopfloor := Shopfloor{
		ID: uuid.New(),
		TenantID: parsedTenantId,
		CustomerID: parsedCustomerID,
		Name: request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repository.Create(ctx, shopfloor)
}

func (s *service) FindAll(ctx context.Context) ([]Shopfloor, error) {
	return s.repository.FindAll(ctx)
}

func (s *service) FindByID(ctx context.Context, id string) (Shopfloor, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Shopfloor{}, err
	}
	return s.repository.FindByID(ctx, parsedId)
}

func (s *service) FindByTenantID(ctx context.Context, id string) ([]Shopfloor, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return []Shopfloor{}, err
	}
	return s.repository.FindByTenantID(ctx, parsedId)
}

func (s *service) Update(ctx context.Context, id string, request ShopfloorRequest) (Shopfloor, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Shopfloor{}, err
	}
	shopfloor, err := s.repository.FindByID(ctx, parsedId)
	if err != nil {
		return Shopfloor{}, err
	}
	shopfloor.Name = request.Name
	return s.repository.Update(ctx, shopfloor)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(ctx, parsedId)
}

