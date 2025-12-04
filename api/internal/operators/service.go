package operators

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request OperatorRequest) error
	FindByID(ctx context.Context, id string) (Operator, error)
	FindAll(ctx context.Context) ([]Operator, error)
	Update(ctx context.Context, id string, request OperatorRequest) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request OperatorRequest) error {
	operator := Operator{
		ID:          uuid.New(),
		ShopFloorID: request.ShopFloorID,
		CustomerID:  request.CustomerID,
		Code:        request.Code,
		Name:        request.Name,
		Surname:     request.Surname,
		VatNumber:   request.VatNumber,
		IsActive:    request.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repo.Create(ctx, operator)
}

func (s *service) FindByID(ctx context.Context, id string) (Operator, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Operator{}, err
	}
	return s.repo.FindByID(ctx, parsedID)
}

func (s *service) FindAll(ctx context.Context) ([]Operator, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Update(ctx context.Context, id string, request OperatorRequest) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	operator, err := s.repo.FindByID(ctx, parsedID)
	if err != nil {
		return err
	}
	operator.ShopFloorID = request.ShopFloorID
	operator.CustomerID = request.CustomerID
	operator.Code = request.Code
	operator.Name = request.Name
	operator.Surname = request.Surname
	operator.VatNumber = request.VatNumber
	operator.IsActive = request.IsActive
	operator.UpdatedAt = time.Now()
	return s.repo.Update(ctx, operator)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedID)
}
