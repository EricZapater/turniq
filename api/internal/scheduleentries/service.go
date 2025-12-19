package scheduleentries

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request ScheduleEntryRequest) (ScheduleEntry, error)
	FindByID(ctx context.Context, id string) (ScheduleEntry, error)
	FindAll(ctx context.Context) ([]ScheduleEntry, error)
	GetPlanning(ctx context.Context, shopfloorID string, date string) ([]ScheduleEntry, error)
	GetOperatorPlanning(ctx context.Context, operatorID string, date string) ([]ScheduleEntry, error)
	Search(ctx context.Context, filter ScheduleFilter) ([]ScheduleEntry, error)
	Update(ctx context.Context, id string, request ScheduleEntryRequest) (ScheduleEntry, error)
	Delete(ctx context.Context, id string) error
	Sync(ctx context.Context, shopfloorID string, date string, requests []ScheduleEntryRequest) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request ScheduleEntryRequest) (ScheduleEntry, error) {
	customerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return ScheduleEntry{}, err
	}
	shopfloorID, err := uuid.Parse(request.ShopfloorID)
	if err != nil {
		return ScheduleEntry{}, err
	}
	shiftID, err := uuid.Parse(request.ShiftID)
	if err != nil {
		return ScheduleEntry{}, err
	}
	
	var workcenterID uuid.NullUUID
	if request.WorkcenterID != "" {
		id, err := uuid.Parse(request.WorkcenterID)
		if err != nil {
			return ScheduleEntry{}, err
		}
		workcenterID = uuid.NullUUID{UUID: id, Valid: true}
	}

	var jobID uuid.NullUUID
	if request.JobID != "" {
		id, err := uuid.Parse(request.JobID)
		if err != nil {
			return ScheduleEntry{}, err
		}
		jobID = uuid.NullUUID{UUID: id, Valid: true}
	}

	var operatorID uuid.NullUUID
	if request.OperatorID != "" {
		id, err := uuid.Parse(request.OperatorID)
		if err != nil {
			return ScheduleEntry{}, err
		}
		operatorID = uuid.NullUUID{UUID: id, Valid: true}
	}

	parsedDate, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		parsedDate, err = time.Parse(time.RFC3339, request.Date)
		if err != nil {
			return ScheduleEntry{}, err
		}
	}

	entry := ScheduleEntry{
		ID:           uuid.New(),
		CustomerID:   customerID,
		ShopfloorID:  shopfloorID,
		ShiftID:      shiftID,
		WorkcenterID: workcenterID,

		JobID:        jobID,
		OperatorID:   operatorID,
		Date:         parsedDate,
		Order:        request.Order,
		StartTime:    request.StartTime,
		EndTime:      request.EndTime,
		IsCompleted:  request.IsCompleted,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	_, err = s.repo.Create(ctx, entry)
	if err != nil {
		return ScheduleEntry{}, err
	}
	return entry, nil
}

func (s *service) FindByID(ctx context.Context, id string) (ScheduleEntry, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ScheduleEntry{}, err
	}
	return s.repo.FindByID(ctx, parsedID)
}

func (s *service) FindAll(ctx context.Context) ([]ScheduleEntry, error) {
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

func (s *service) GetPlanning(ctx context.Context, shopfloorID string, date string) ([]ScheduleEntry, error) {
	parsedShopfloorID, err := uuid.Parse(shopfloorID)
	if err != nil {
		return nil, err
	}
	// date validation?
	return s.repo.FindByShopfloorAndDate(ctx, parsedShopfloorID, date)
}

func (s *service) GetOperatorPlanning(ctx context.Context, operatorID string, date string) ([]ScheduleEntry, error) {
	parsedOperatorID, err := uuid.Parse(operatorID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByOperatorAndDate(ctx, parsedOperatorID, date)
}

func (s *service) Search(ctx context.Context, filter ScheduleFilter) ([]ScheduleEntry, error) {
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
		// Force customer ID
		filter.CustomerID = &customerIDFromCtx
		
		// Remove OperatorID filter as per request "els demés sense operari"
		// This strictly enforces that logic.
		filter.OperatorID = nil
	}
	return s.repo.Search(ctx, filter)
}

func (s *service) Update(ctx context.Context, id string, request ScheduleEntryRequest) (ScheduleEntry, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ScheduleEntry{}, err
	}

	entry, err := s.repo.FindByID(ctx, parsedID)
	if err != nil {
		return ScheduleEntry{}, err
	}

	// Update fields
	if request.CustomerID != "" {
		if customerID, err := uuid.Parse(request.CustomerID); err == nil {
			entry.CustomerID = customerID
		}
	}
	if request.ShopfloorID != "" {
		if shopfloorID, err := uuid.Parse(request.ShopfloorID); err == nil {
			entry.ShopfloorID = shopfloorID
		}
	}
	if request.ShiftID != "" {
		if shiftID, err := uuid.Parse(request.ShiftID); err == nil {
			entry.ShiftID = shiftID
		}
	}
	if request.WorkcenterID != "" {
		if workcenterID, err := uuid.Parse(request.WorkcenterID); err == nil {
			entry.WorkcenterID = uuid.NullUUID{UUID: workcenterID, Valid: true}
		}
	}
	if request.JobID != "" {
		if id, err := uuid.Parse(request.JobID); err == nil {
			entry.JobID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}
	if request.OperatorID != "" {
		if id, err := uuid.Parse(request.OperatorID); err == nil {
			entry.OperatorID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}
	
	if request.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", request.Date)
		if err != nil {
			parsedDate, err = time.Parse(time.RFC3339, request.Date)
			if err != nil {
				return ScheduleEntry{}, err
			}
		}
		entry.Date = parsedDate
	}
	
	entry.Order = request.Order
	entry.StartTime = request.StartTime
	entry.EndTime = request.EndTime
	entry.IsCompleted = request.IsCompleted
	entry.UpdatedAt = time.Now().Format(time.RFC3339)

	_, err = s.repo.Update(ctx, entry)
	if err != nil {
		return ScheduleEntry{}, err
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

func (s *service) Sync(ctx context.Context, shopfloorID string, date string, requests []ScheduleEntryRequest) error {
	parsedShopfloorID, err := uuid.Parse(shopfloorID)
	if err != nil {
		return err
	}

	// Parse sync date
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		// Try RFC3339 if partial fails, or just error out
		parsedDate, err = time.Parse(time.RFC3339, date)
		if err != nil {
			return err
		}
	}

	var entries []ScheduleEntry
	for _, req := range requests {
		customerID, err := uuid.Parse(req.CustomerID)
		if err != nil {
			return err
		}
		sfID, err := uuid.Parse(req.ShopfloorID)
		if err != nil {
			return err
		}
		shiftID, err := uuid.Parse(req.ShiftID)
		if err != nil {
			return err
		}
		
		var workcenterID uuid.NullUUID
		if req.WorkcenterID != "" {
			id, err := uuid.Parse(req.WorkcenterID)
			if err != nil {
				return err
			}
			workcenterID = uuid.NullUUID{UUID: id, Valid: true}
		}

		var jobID uuid.NullUUID
		if req.JobID != "" {
			id, err := uuid.Parse(req.JobID)
			if err != nil {
				return err
			}
			jobID = uuid.NullUUID{UUID: id, Valid: true}
		}

		var operatorID uuid.NullUUID
		if req.OperatorID != "" {
			id, err := uuid.Parse(req.OperatorID)
			if err != nil {
				return err
			}
			operatorID = uuid.NullUUID{UUID: id, Valid: true}
		}

		// Use provided ID if valid, else new
		entryID := uuid.New()
		if req.ID != "" {
			if parsed, err := uuid.Parse(req.ID); err == nil {
				entryID = parsed
			}
		}

		entries = append(entries, ScheduleEntry{
			ID:           entryID,
			CustomerID:   customerID,
			ShopfloorID:  sfID,
			ShiftID:      shiftID,
			WorkcenterID: workcenterID,
			JobID:        jobID,
			OperatorID:   operatorID,
			Date:         parsedDate, // Use parsed date
			Order:        req.Order,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
			IsCompleted:  req.IsCompleted,
			CreatedAt:    time.Now().Format(time.RFC3339),
			UpdatedAt:    time.Now().Format(time.RFC3339),
		})
	}

	return s.repo.Sync(ctx, parsedShopfloorID, date, entries)
}
