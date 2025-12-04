package jobs

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, request JobRequest) (Job, error)
	FindAll(ctx context.Context) ([]Job, error)
	FindByID(ctx context.Context, id string) (Job, error)
	FindByWorkcenterID(ctx context.Context, workcenterId string) ([]Job, error)
	Update(ctx context.Context, id string, request JobRequest) (Job, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(ctx context.Context, request JobRequest) (Job, error) {
	tenantParsedID, err := uuid.Parse(request.TenantID)
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
	job := Job{
		ID:          uuid.New(),
		TenantID:    tenantParsedID,
		ShopFloorID: shopFloorParsedID,
		WorkcenterID: workcenterParsedID,
		JobCode:     request.JobCode,
		ProductCode: request.ProductCode,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repository.Create(ctx, job)
}

func(s *service) FindAll(ctx context.Context) ([]Job, error) {
	return s.repository.FindAll(ctx)
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

func(s *service) Update(ctx context.Context, id string, request JobRequest) (Job, error) {
	jobParsedID, err := uuid.Parse(id)
	if err != nil {
		return Job{}, err
	}
	tenantParsedID, err := uuid.Parse(request.TenantID)
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
	job.TenantID = tenantParsedID
	job.ShopFloorID = shopFloorParsedID
	job.WorkcenterID = workcenterParsedID
	job.JobCode = request.JobCode
	job.ProductCode = request.ProductCode
	job.Description = request.Description
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

