package jobs

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, job Job) (Job, error)
	FindAll(ctx context.Context) ([]Job, error)
	FindByID(ctx context.Context, id uuid.UUID) (Job, error)
	FindByWorkcenterID(ctx context.Context, workcenterId uuid.UUID) ([]Job, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Job, error)
	FindByShopFloorID(ctx context.Context, shopFloorID uuid.UUID) ([]Job, error)
	CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int, error)
	Update(ctx context.Context, job Job) (Job,error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB	
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, job Job) (Job,error) {
	query := `INSERT INTO jobs (id, customer_id, shop_floor_id, workcenter_id, 
								job_code, product_code, description, 
								estimated_duration, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, job.ID, job.CustomerID, job.ShopFloorID, job.WorkcenterID, 
								job.JobCode, job.ProductCode, job.Description, 
								job.EstimatedDuration, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return Job{},err
	}
	return job,nil
}

func (r *repository) FindAll(ctx context.Context) ([]Job, error) {
	query := "SELECT id, customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration, created_at, updated_at FROM jobs"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.CustomerID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.EstimatedDuration, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Job, error) {
	query := "SELECT id, customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration, created_at, updated_at FROM jobs WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)
	var job Job
	err := row.Scan(&job.ID, &job.CustomerID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.EstimatedDuration, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func(r *repository) FindByWorkcenterID(ctx context.Context, workcenterId uuid.UUID) ([]Job, error) {
	query := "SELECT id, customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration, created_at, updated_at FROM jobs WHERE workcenter_id = $1"
	rows, err := r.db.QueryContext(ctx, query, workcenterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.CustomerID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.EstimatedDuration, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Job, error) {
	query := "SELECT id, customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration, created_at, updated_at FROM jobs WHERE customer_id = $1"
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.CustomerID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.EstimatedDuration, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) FindByShopFloorID(ctx context.Context, shopFloorID uuid.UUID) ([]Job, error) {
	query := "SELECT id, customer_id, shop_floor_id, workcenter_id, job_code, product_code, description, estimated_duration, created_at, updated_at FROM jobs WHERE shop_floor_id = $1"
	rows, err := r.db.QueryContext(ctx, query, shopFloorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.CustomerID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.EstimatedDuration, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int, error) {
	query := "SELECT COUNT(*) FROM jobs WHERE customer_id = $1"
	var count int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, job Job) (Job, error) {
	query := "UPDATE jobs SET tenant_id = $2, shop_floor_id = $3, workcenter_id = $4, job_code = $5, product_code = $6, description = $7, estimated_duration = $8, updated_at = $9 WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, job.ID, job.CustomerID, job.ShopFloorID, job.WorkcenterID, job.JobCode, job.ProductCode, job.Description, job.EstimatedDuration, job.UpdatedAt)
	if err != nil {
		return Job{},err
	}
	return job, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM jobs WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

