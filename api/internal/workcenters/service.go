package workcenters

import (
	"api/internal/customers"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request WorkcenterRequest) (Workcenter, error)
	FindAll(ctx context.Context) ([]Workcenter, error)
	FindByID(ctx context.Context, id string) (Workcenter, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]Workcenter, error)
	Update(ctx context.Context, id string, request WorkcenterRequest) (Workcenter, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	customerService customers.Service
}

func NewService(repo Repository, customerService customers.Service) Service {
	return &service{repo: repo, customerService: customerService}
}

func (s *service) Create(ctx context.Context, request WorkcenterRequest) (Workcenter, error) {
	var customerID uuid.UUID
	
	// Recuperar i validar is_admin
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Workcenter{}, errors.New("invalid or missing is_admin in context")
	}
	
	if isAdmin {
		// Si és admin, usar el customerID del request
		parsedCustomerID, err := uuid.Parse(request.CustomerID)
		if err != nil {
			return Workcenter{}, err
		}
		customerID = parsedCustomerID
	} else {
		// Si no és admin, usar el customerID del context
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return Workcenter{}, errors.New("invalid or missing customer_id in context")
		}
		customerID = customerIDFromCtx
	}
	
	// Check limits
	customer, err := s.customerService.FindByID(ctx, customerID.String())
	if err != nil {
		return Workcenter{}, err
	}

	count, err := s.repo.CountByCustomerID(ctx, customerID)
	if err != nil {
		return Workcenter{}, err
	}

	if count >= customer.MaxWorkcenters {
		return Workcenter{}, errors.New("max workcenters limit reached for this tenant")
	}

	var shopFloorID uuid.NullUUID
	if request.ShopFloorID != "" {
		if id, err := uuid.Parse(request.ShopFloorID); err == nil {
			shopFloorID = uuid.NullUUID{UUID: id, Valid: true}
		}
	}

	workcenter := Workcenter{
		ID:        uuid.New(),		
		CustomerID: customerID,
		ShopFloorID: shopFloorID,
		Name:      request.Name,
		IsActive:  request.IsActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(ctx, workcenter)
}

func (s *service) FindAll(ctx context.Context) ([]Workcenter, error) {
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

func (s *service) FindByID(ctx context.Context, id string) (Workcenter, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Workcenter{}, err
	}
	return s.repo.FindByID(ctx, parsedId)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]Workcenter, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}	
	return s.repo.FindByCustomerID(ctx, parsedCustomerID)
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
	
	if request.ShopFloorID != "" {
		if id, err := uuid.Parse(request.ShopFloorID); err == nil {
			workcenter.ShopFloorID = uuid.NullUUID{UUID: id, Valid: true}
		}
	} else {
		// If empty, set to NULL (allowing unassignment)
		workcenter.ShopFloorID = uuid.NullUUID{Valid: false}
	}

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
