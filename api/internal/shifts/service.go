package shifts

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context,request ShiftRequest) (Shift, error)
	FindByShopfloorID(ctx context.Context,shopfloorID string) ([]Shift, error)
	Update(ctx context.Context,id string, request ShiftRequest) (Shift, error)
	Delete(ctx context.Context,shiftID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context,request ShiftRequest) (Shift, error) {
	parsedTenantId, err := uuid.Parse(request.TenantID)
	if err != nil {
		return Shift{}, err
	}
	parsedCustomerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Shift{}, err
	}
	parsedShopfloorID, err := uuid.Parse(request.ShopfloorID)
	if err != nil {
		return Shift{}, err
	}
	parsedStartTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return Shift{}, err
	}
	parsedEndTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return Shift{}, err
	}
	return s.repo.Create(ctx,Shift{
		ID: uuid.New(),
		TenantID: parsedTenantId,
		CustomerID: parsedCustomerID,
		ShopfloorID: parsedShopfloorID,
		StartTime: parsedStartTime,
		EndTime: parsedEndTime,
		IsActive: request.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *service) FindByShopfloorID(ctx context.Context, shopfloorID string) ([]Shift, error) {
	parsedShopfloorID, err := uuid.Parse(shopfloorID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByShopfloorID(ctx,parsedShopfloorID)
}

func (s *service) Update(ctx context.Context,id string, request ShiftRequest) (Shift, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Shift{}, err
	}
	parsedTenantId, err := uuid.Parse(request.TenantID)
	if err != nil {
		return Shift{}, err
	}
	parsedCustomerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Shift{}, err
	}
	parsedShopfloorID, err := uuid.Parse(request.ShopfloorID)
	if err != nil {
		return Shift{}, err
	}
	parsedStartTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return Shift{}, err
	}
	parsedEndTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return Shift{}, err
	}
	shift, err := s.repo.FindByID(ctx,parsedID)
	if err != nil {
		return Shift{}, err
	}
	shift.TenantID = parsedTenantId
	shift.CustomerID = parsedCustomerID
	shift.ShopfloorID = parsedShopfloorID
	shift.StartTime = parsedStartTime
	shift.EndTime = parsedEndTime
	shift.IsActive = request.IsActive
	shift.UpdatedAt = time.Now()
	return s.repo.Update(ctx,shift)
}

func (s *service) Delete(ctx context.Context,shiftID string) error {
	parsedShiftID, err := uuid.Parse(shiftID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx,parsedShiftID)
}
