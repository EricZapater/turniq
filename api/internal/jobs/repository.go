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
	query := "INSERT INTO jobs (id, tenant_id, shop_floor_id, workcenter_id, job_code, product_code, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	_, err := r.db.ExecContext(ctx, query, job.ID, job.TenantID, job.ShopFloorID, job.WorkcenterID, job.JobCode, job.ProductCode, job.Description, job.CreatedAt, job.UpdatedAt)
	if err != nil {
		return Job{},err
	}
	return job,nil
}

func (r *repository) FindAll(ctx context.Context) ([]Job, error) {
	query := "SELECT * FROM jobs"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.TenantID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Job, error) {
	query := "SELECT * FROM jobs WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, id)
	var job Job
	err := row.Scan(&job.ID, &job.TenantID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.CreatedAt, &job.UpdatedAt)
	if err != nil {
		return Job{}, err
	}
	return job, nil
}

func(r *repository) FindByWorkcenterID(ctx context.Context, workcenterId uuid.UUID) ([]Job, error) {
	query := "SELECT * FROM jobs WHERE workcenter_id = $1"
	rows, err := r.db.QueryContext(ctx, query, workcenterId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.TenantID, &job.ShopFloorID, &job.WorkcenterID, &job.JobCode, &job.ProductCode, &job.Description, &job.CreatedAt, &job.UpdatedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *repository) Update(ctx context.Context, job Job) (Job, error) {
	query := "UPDATE jobs SET tenant_id = $2, shop_floor_id = $3, workcenter_id = $4, job_code = $5, product_code = $6, description = $7, updated_at = $8 WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, job.ID, job.TenantID, job.ShopFloorID, job.WorkcenterID, job.JobCode, job.ProductCode, job.Description, job.UpdatedAt)
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

