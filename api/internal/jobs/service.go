package jobs

import (
	"api/internal/customers"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request JobRequest) (Job, error)
	FindAll(ctx context.Context) ([]Job, error)
	FindByID(ctx context.Context, id string) (Job, error)
	FindByWorkcenterID(ctx context.Context, workcenterId string) ([]Job, error)
	FindByShopFloorID(ctx context.Context, shopFloorID string) ([]Job, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]Job, error)
	Update(ctx context.Context, id string, request JobRequest) (Job, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository	
	customerService customers.Service
}

func NewService(repository Repository, customerService customers.Service) Service {
	return &service{repository: repository, customerService: customerService}
}

func (s *service) Create(ctx context.Context, request JobRequest) (Job, error) {
	var customerID uuid.UUID
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return Job{}, errors.New("invalid or missing is_admin in context")
	}
	
	if isAdmin {
		// Si és admin, usar el customerID del request
		parsedCustomerID, err := uuid.Parse(request.CustomerID)
		if err != nil {
			return Job{}, err
		}
		customerID = parsedCustomerID
	} else {
		// Si no és admin, usar el customerID del context
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return Job{}, errors.New("invalid or missing customer_id in context")
		}
		customerID = customerIDFromCtx
	}

	// Check limits
	customer, err := s.customerService.FindByID(ctx, customerID.String())
	if err != nil {
		return Job{}, err
	}

	count, err := s.repository.CountByCustomerID(ctx, customerID)
	if err != nil {
		return Job{}, err
	}

	if count >= customer.MaxJobs {
		return Job{}, errors.New("max jobs limit reached for this customer")
	}

	shopFloorParsedID, err := uuid.Parse(request.ShopFloorID)
	if err != nil {
		return Job{}, err
	}
	workcenterParsedID, err := uuid.Parse(request.WorkcenterID)
	if err != nil {
		return Job{}, err
	}
	job := Job{
		ID:          uuid.New(),
		CustomerID:    customerID,
		ShopFloorID: shopFloorParsedID,
		WorkcenterID: workcenterParsedID,
		JobCode:     request.JobCode,
		ProductCode: request.ProductCode,
		Description: request.Description,
		EstimatedDuration: request.EstimatedDuration,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repository.Create(ctx, job)
}

func(s *service) FindAll(ctx context.Context) ([]Job, error) {
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

func(s *service) FindByID(ctx context.Context, id string) (Job, error) {
	jobParsedID, err := uuid.Parse(id)
	if err != nil {
		return Job{}, err
	}
	return s.repository.FindByID(ctx, jobParsedID)
}

func(s *service) FindByWorkcenterID(ctx context.Context, workcenterId string) ([]Job, error) {	
	workcenterParsedID, err := uuid.Parse(workcenterId)
	if err != nil {
		return []Job{}, err
	}
	return s.repository.FindByWorkcenterID(ctx, workcenterParsedID)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]Job, error) {
	parsedID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repository.FindByCustomerID(ctx, parsedID)
}

func(s *service) FindByShopFloorID(ctx context.Context, shopFloorID string) ([]Job, error) {
	shopFloorParsedID, err := uuid.Parse(shopFloorID)
	if err != nil {
		return nil, err
	}
	return s.repository.FindByShopFloorID(ctx, shopFloorParsedID)
}

func(s *service) Update(ctx context.Context, id string, request JobRequest) (Job, error) {
	jobParsedID, err := uuid.Parse(id)
	if err != nil {
		return Job{}, err
	}
	customerParsedID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return Job{}, err
	}
	shopFloorParsedID, err := uuid.Parse(request.ShopFloorID)
	if err != nil {
		return Job{}, err
	}
	workcenterParsedID, err := uuid.Parse(request.WorkcenterID)
	if err != nil {
		return Job{}, err
	}
	job, err := s.repository.FindByID(ctx, jobParsedID)
	if err != nil {
		return Job{}, err
	}	
	job.CustomerID = customerParsedID
	job.ShopFloorID = shopFloorParsedID
	job.WorkcenterID = workcenterParsedID
	job.JobCode = request.JobCode
	job.ProductCode = request.ProductCode
	job.Description = request.Description
	job.EstimatedDuration = request.EstimatedDuration
	job.UpdatedAt = time.Now()
	return s.repository.Update(ctx, job)
}

func(s *service) Delete(ctx context.Context, id string) error {
	jobParsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(ctx, jobParsedID)
}

