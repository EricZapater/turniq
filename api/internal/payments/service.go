package payments

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request PaymentRequest) (Payment, error)
	FindAll(ctx context.Context) ([]Payment, error)
	FindById(ctx context.Context, id string) (Payment, error)
	FindByCustomerId(ctx context.Context, customerId string) ([]Payment, error)
	Search(ctx context.Context, filter PaymentFilter) ([]Payment, error)
	Update(ctx context.Context, id string, request PaymentRequest) (Payment, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request PaymentRequest) (Payment, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Payment{}, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return Payment{}, errors.New("user is not admin")
	}
	customerParsedId, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Payment{}, err
	}	
	payment := Payment{
		ID:            uuid.New(),
		CustomerID:    customerParsedId,
		Amount:        request.Amount,
		Currency:      request.Currency,
		PaymentMethod: request.PaymentMethod,
		Status:        request.Status,
		DueDate:       request.DueDate,
		PaidAt:        request.PaidAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return s.repo.Create(ctx, payment)
}

func (s *service) FindAll(ctx context.Context) ([]Payment, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return nil, errors.New("user is not admin")
	}
	return s.repo.FindAll(ctx)
}

func (s *service) FindById(ctx context.Context, id string) (Payment, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Payment{}, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return Payment{}, errors.New("user is not admin")
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Payment{}, err
	}
	return s.repo.FindById(ctx, parsedID)
}

func (s *service) FindByCustomerId(ctx context.Context, customerId string) ([]Payment, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return nil, errors.New("user is not admin")
	}
	customerParsedId, err := uuid.Parse(customerId)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByCustomerId(ctx, customerParsedId)
}

func (s *service) Search(ctx context.Context, filter PaymentFilter) ([]Payment, error) {
	// Only Admin can search payments (Billing Report)
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return nil, errors.New("user is not admin")
	}
	return s.repo.Search(ctx, filter)
}

func (s *service) Update(ctx context.Context, id string, request PaymentRequest) (Payment, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Payment{}, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return Payment{}, errors.New("user is not admin")
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return Payment{}, err
	}
	customerParsedId, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Payment{}, err
	}	
	payment, err := s.repo.FindById(ctx, parsedID)
	if err != nil {
		return Payment{}, err
	}

	payment.CustomerID = customerParsedId
	payment.Amount = request.Amount
	payment.Currency = request.Currency
	payment.PaymentMethod = request.PaymentMethod
	payment.Status = request.Status
	payment.DueDate = request.DueDate
	payment.PaidAt = request.PaidAt
	payment.UpdatedAt = time.Now()
	return s.repo.Update(ctx, payment)
}

func (s *service) Delete(ctx context.Context, id string) error {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin {
		return errors.New("user is not admin")
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedID)
}
