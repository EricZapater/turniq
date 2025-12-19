package timeentries

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request TimeEntryRequest) (TimeEntry, error)
	FindByID(ctx context.Context, id string) (TimeEntry, error)
	FindAll(ctx context.Context) ([]TimeEntry, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]TimeEntry, error)
	FindByOperatorID(ctx context.Context, operatorID string) ([]TimeEntry, error)
	FindCurrent(ctx context.Context, operatorID string) (TimeEntry, error)
	Search(ctx context.Context, filter TimeEntryFilter) ([]TimeEntry, error)
	Update(ctx context.Context, id string, request TimeEntryRequest) (TimeEntry, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request TimeEntryRequest) (TimeEntry, error) {
	operatorID, err := uuid.Parse(request.OperatorID)
	if err != nil {
		return TimeEntry{}, err
	}
	
	var workcenterID *uuid.UUID
	if request.WorkcenterID != nil && *request.WorkcenterID != "" {
		id, err := uuid.Parse(*request.WorkcenterID)
		if err != nil {
			return TimeEntry{}, err
		}
		workcenterID = &id
	}

	entry := TimeEntry{
		ID:           uuid.New(),
		OperatorID:   operatorID,
		WorkcenterID: workcenterID,
		CheckIn:      request.CheckIn,
		CheckOut:     request.CheckOut,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err = s.repo.Create(ctx, entry)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (s *service) FindByID(ctx context.Context, id string) (TimeEntry, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return TimeEntry{}, err
	}
	return s.repo.FindByID(ctx, parsedID)
}

func (s *service) FindAll(ctx context.Context) ([]TimeEntry, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if isAdmin{
		return s.repo.FindAll(ctx)
	}
	customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return nil, errors.New("invalid or missing customer_id in context")
		}
		customerID := customerIDFromCtx
	return s.repo.FindByCustomerID(ctx, customerID)
}

func(s *service) FindByCustomerID(ctx context.Context, customerID string) ([]TimeEntry, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return []TimeEntry{}, err
	}
	return s.repo.FindByCustomerID(ctx, parsedCustomerID)
}

func(s *service) FindByOperatorID(ctx context.Context, operatorID string) ([]TimeEntry, error) {
	parsedOperatorID, err := uuid.Parse(operatorID)
	if err != nil {
		return []TimeEntry{}, err
	}
	return s.repo.FindByOperatorID(ctx, parsedOperatorID)
}

func(s *service) FindCurrent(ctx context.Context, operatorID string) (TimeEntry, error) {
	parsedOperatorID, err := uuid.Parse(operatorID)
	if err != nil {
		return TimeEntry{}, err
	}
	return s.repo.FindCurrent(ctx, parsedOperatorID)
}

func (s *service) Search(ctx context.Context, filter TimeEntryFilter) ([]TimeEntry, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}

	if !isAdmin {
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return nil, errors.New("invalid or missing customer_id in context")
		}
		filter.CustomerID = &customerIDFromCtx
	}
	return s.repo.Search(ctx, filter)
}

func (s *service) Update(ctx context.Context, id string, request TimeEntryRequest) (TimeEntry, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return TimeEntry{}, err
	}

	entry, err := s.repo.FindByID(ctx, parsedID)
	if err != nil {
		return TimeEntry{}, err
	}

	// Update fields
	if request.OperatorID != "" {
		if operatorID, err := uuid.Parse(request.OperatorID); err == nil {
			entry.OperatorID = operatorID
		}
	}
	
	if request.WorkcenterID != nil {
		if *request.WorkcenterID != "" {
			if workcenterID, err := uuid.Parse(*request.WorkcenterID); err == nil {
				entry.WorkcenterID = &workcenterID
			}
		} else {
			entry.WorkcenterID = nil
		}
	}

	entry.CheckIn = request.CheckIn
	entry.CheckOut = request.CheckOut
	entry.UpdatedAt = time.Now()

	_, err = s.repo.Update(ctx, entry)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedID)
}
