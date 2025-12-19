package shifts

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context,request ShiftRequest) (Shift, error)
	FindByID(ctx context.Context, id string) (Shift, error)
	FindByShopfloorID(ctx context.Context,shopfloorID string) ([]Shift, error)
	FindByCustomerID(ctx context.Context,customerID string) ([]Shift, error)
	Update(ctx context.Context,id string, request ShiftRequest) (Shift, error)
	Delete(ctx context.Context,shiftID string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FindByID(ctx context.Context, id string) (Shift, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Shift{}, err
	}
	return s.repo.FindByID(ctx, parsedID)
}

func (s *service) Create(ctx context.Context,request ShiftRequest) (Shift, error) {
	parsedCustomerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Shift{}, err
	}
	var parsedShopfloorID uuid.NullUUID
	if request.ShopfloorID != "" {
		id, err := uuid.Parse(request.ShopfloorID)
		if err != nil {
			return Shift{}, err
		}
		parsedShopfloorID = uuid.NullUUID{UUID: id, Valid: true}
	} else {
		parsedShopfloorID = uuid.NullUUID{Valid: false}
	}

	parsedStartTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return Shift{}, err
	}
	parsedEndTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return Shift{}, err
	}
	
	// Validate Overlap. Scope: If Shopfloor defined, use Shopfloor. Else Use Tenant.
	var existingShifts []Shift
	if parsedShopfloorID.Valid {
		existingShifts, err = s.repo.FindByShopfloorID(ctx, parsedShopfloorID.UUID)
	} else {
		existingShifts, err = s.repo.FindByCustomerID(ctx, parsedCustomerID)
	}

	filteredShifts := []Shift{}
	for _, s := range existingShifts {
		if s.ShopfloorID == parsedShopfloorID {
			filteredShifts = append(filteredShifts, s)
		}
	}

	if err != nil {
		return Shift{}, err
	}
	if checkOverlap(parsedStartTime, parsedEndTime, filteredShifts, uuid.Nil) {
		return Shift{}, errors.New("shift overlaps with an existing shift")
	}

	return s.repo.Create(ctx,Shift{
		ID: uuid.New(),
		CustomerID: parsedCustomerID,
		ShopfloorID: parsedShopfloorID,
		Name: request.Name,
		Color: request.Color,
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
	return s.repo.FindByShopfloorID(ctx, parsedShopfloorID)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]Shift, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByCustomerID(ctx, parsedCustomerID)
}

func (s *service) Update(ctx context.Context,id string, request ShiftRequest) (Shift, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Shift{}, err
	}
	parsedCustomerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Shift{}, err
	}
	
	var parsedShopfloorID uuid.NullUUID
	if request.ShopfloorID != "" {
		sid, err := uuid.Parse(request.ShopfloorID)
		if err != nil {
			return Shift{}, err
		}
		parsedShopfloorID = uuid.NullUUID{UUID: sid, Valid: true}
	} else {
		parsedShopfloorID = uuid.NullUUID{Valid: false}
	}

	parsedStartTime, err := time.Parse("15:04", request.StartTime)
	if err != nil {
		return Shift{}, err
	}
	parsedEndTime, err := time.Parse("15:04", request.EndTime)
	if err != nil {
		return Shift{}, err
	}
	shift, err := s.repo.FindByID(ctx, parsedID)
	if err != nil {
		return Shift{}, err
	}
	
	// Validate Overlap
	var existingShifts []Shift
	if parsedShopfloorID.Valid {
		existingShifts, err = s.repo.FindByShopfloorID(ctx, parsedShopfloorID.UUID)
	} else {
		existingShifts, err = s.repo.FindByCustomerID(ctx, parsedCustomerID)
	}

	filteredShifts := []Shift{}
	for _, s := range existingShifts {
		if s.ShopfloorID == parsedShopfloorID {
			filteredShifts = append(filteredShifts, s)
		}
	}

	if err != nil {
		return Shift{}, err
	}
	if checkOverlap(parsedStartTime, parsedEndTime, filteredShifts, parsedID) {
		return Shift{}, errors.New("shift overlaps with an existing shift")
	}

	shift.CustomerID = parsedCustomerID
	shift.ShopfloorID = parsedShopfloorID
	shift.Name = request.Name
	shift.Color = request.Color
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

func checkOverlap(start, end time.Time, existing []Shift, excludeID uuid.UUID) bool {
	// Normalize to minutes from 00:00
	toMinutes := func(t time.Time) int {
		return t.Hour()*60 + t.Minute()
	}
	s1 := toMinutes(start)
	e1 := toMinutes(end)

	// Define intervals for new shift
	var intervals [][2]int
	if s1 < e1 {
		intervals = append(intervals, [2]int{s1, e1})
	} else {
		// Crosses midnight: [s1, 1440] and [0, e1]
		intervals = append(intervals, [2]int{s1, 24 * 60})
		intervals = append(intervals, [2]int{0, e1})
	}

	for _, shift := range existing {
		if shift.ID == excludeID {
			continue
		}
		if !shift.IsActive {
			continue
		}
		
		s2 := toMinutes(shift.StartTime)
		e2 := toMinutes(shift.EndTime)
		
		var currentIntervals [][2]int
		if s2 < e2 {
			currentIntervals = append(currentIntervals, [2]int{s2, e2})
		} else {
			currentIntervals = append(currentIntervals, [2]int{s2, 24 * 60})
			currentIntervals = append(currentIntervals, [2]int{0, e2})
		}

		// Check intersection
		for _, i1 := range intervals {
			for _, i2 := range currentIntervals {
				// Overlap if max(start1, start2) < min(end1, end2)
				maxStart := i1[0]
				if i2[0] > maxStart {
					maxStart = i2[0]
				}
				minEnd := i1[1]
				if i2[1] < minEnd {
					minEnd = i2[1]
				}

				if maxStart < minEnd {
					return true
				}
			}
		}
	}
	return false
}
