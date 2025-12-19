package operators

import (
	"api/internal/customers"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request OperatorRequest) (Operator, error)
	FindByID(ctx context.Context, id string) (Operator, error)
	FindAll(ctx context.Context) ([]Operator, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]Operator, error)
	FindByCode(ctx context.Context, code string) (Operator, error)
	Update(ctx context.Context, id string, request OperatorRequest) (Operator, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository	
	customerService customers.Service
}

func NewService(repo Repository, customerService customers.Service) Service {
	return &service{repo: repo, customerService: customerService}
}

func (s *service) Create(ctx context.Context, request OperatorRequest) (Operator, error) {
	var customerID uuid.UUID
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Operator{}, errors.New("invalid or missing is_admin in context")
	}
	
	if isAdmin {
		// Si és admin, usar el customerID del request
		parsedCustomerID, err := uuid.Parse(request.CustomerID)
		if err != nil {
			return Operator{}, err
		}
		customerID = parsedCustomerID
	} else {
		// Si no és admin, usar el customerID del context
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return Operator{}, errors.New("invalid or missing customer_id in context")
		}
		customerID = customerIDFromCtx
	}

	// Check limits
	customer, err := s.customerService.FindByID(ctx, customerID.String())
	if err != nil {
		return Operator{}, err
	}

	count, err := s.repo.CountByCustomerID(ctx, customerID)
	if err != nil {
		return Operator{}, err
	}

	if count >= customer.MaxOperators {
		return Operator{}, errors.New("max operators limit reached for this customer")
	}
	

	operator := Operator{
		ID:          uuid.New(),
		ShopFloorID: request.ShopFloorID,
		CustomerID:  customerID,
		Code:        request.Code,
		Name:        request.Name,
		Surname:     request.Surname,
		VatNumber:   request.VatNumber,
		IsActive:    request.IsActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = s.repo.Create(ctx, operator)
	if err != nil {
		return Operator{}, err
	}
	return operator, nil
}

func (s *service) FindByID(ctx context.Context, id string) (Operator, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Operator{}, err
	}
	return s.repo.FindByID(ctx, parsedID)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]Operator, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByCustomerID(ctx, parsedCustomerID)
}

func (s *service) FindAll(ctx context.Context) ([]Operator, error) {
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

func (s *service) Update(ctx context.Context, id string, request OperatorRequest) (Operator, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Operator{}, err
	}
	operator, err := s.repo.FindByID(ctx, parsedID)
	if err != nil {
		return Operator{}, err
	}
	operator.ShopFloorID = request.ShopFloorID
	operator.Code = request.Code
	operator.Name = request.Name
	operator.Surname = request.Surname
	operator.VatNumber = request.VatNumber
	operator.IsActive = request.IsActive
	operator.UpdatedAt = time.Now()
	_, err = s.repo.Update(ctx, operator)
	if err != nil {
		return Operator{}, err
	}
	return operator, nil
}

func (s *service) FindByCode(ctx context.Context, code string) (Operator, error) {
	return s.repo.FindByCode(ctx, code)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedID)
}
