package shopfloors

import (
	"api/internal/customers"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request ShopfloorRequest) (Shopfloor, error)
	FindAll(ctx context.Context) ([]Shopfloor, error)
	FindByID(ctx context.Context, id string) (Shopfloor, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]Shopfloor, error)
	Update(ctx context.Context, id string, request ShopfloorRequest) (Shopfloor, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
	customerService customers.Service		
}

func NewService(repository Repository, customerService customers.Service) Service {
	return &service{repository: repository, customerService: customerService}
}

func (s *service) Create(ctx context.Context, request ShopfloorRequest) (Shopfloor, error) {
	var customerID uuid.UUID
	
	// Recuperar i validar is_admin
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Shopfloor{}, errors.New("invalid or missing is_admin in context")
	}
	
	if isAdmin {
		// Si és admin, usar el customerID del request
		parsedCustomerID, err := uuid.Parse(request.CustomerID)
		if err != nil {
			return Shopfloor{}, err
		}
		customerID = parsedCustomerID
	} else {
		// Si no és admin, usar el customerID del context
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return Shopfloor{}, errors.New("invalid or missing customer_id in context")
		}
		customerID = customerIDFromCtx
	}
	
	// Check limits
	customer, err := s.customerService.FindByID(ctx, customerID.String())
	if err != nil {
		return Shopfloor{}, err
	}

	count, err := s.repository.CountByCustomerID(ctx, customerID)
	if err != nil {
		return Shopfloor{}, err
	}

	if count >= customer.MaxShopFloors {
		return Shopfloor{}, errors.New("max shop floors limit reached for this tenant")
	}

	shopfloor := Shopfloor{
		ID:         uuid.New(),		
		CustomerID: customerID,
		Name:       request.Name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.repository.Create(ctx, shopfloor)
}

func (s *service) FindAll(ctx context.Context) ([]Shopfloor, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if isAdmin{
		return s.repository.FindAll(ctx)
	}
	customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return nil, errors.New("invalid or missing customer_id in context")
		}
		customerID := customerIDFromCtx
	return s.repository.FindByCustomerID(ctx, customerID)
}

func (s *service) FindByID(ctx context.Context, id string) (Shopfloor, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return Shopfloor{}, err
	}
	return s.repository.FindByID(ctx, parsedId)
}

func (s *service) FindByCustomerID(ctx context.Context, id string) ([]Shopfloor, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return []Shopfloor{}, err
	}
	return s.repository.FindByCustomerID(ctx, parsedId)
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

