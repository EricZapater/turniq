package workcenters

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, workcenter Workcenter)(Workcenter, error)
	FindAll(ctx context.Context) ([]Workcenter, error)
	FindByID(ctx context.Context, id uuid.UUID) (Workcenter, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Workcenter, error)
	CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int, error)
	Update(ctx context.Context, workcenter Workcenter) (Workcenter, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB	
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, workcenter Workcenter)(Workcenter, error) {
	query := `INSERT INTO workcenters (id, customer_id, shop_floor_id, name, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	_, err := r.db.ExecContext(ctx, query, workcenter.ID, workcenter.CustomerID, workcenter.ShopFloorID, workcenter.Name, workcenter.IsActive, workcenter.CreatedAt, workcenter.UpdatedAt)
	if err != nil {
		return Workcenter{}, err
	}
	return workcenter, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Workcenter, error) {
	query := `SELECT id, customer_id, shop_floor_id, name, is_active, created_at, updated_at FROM workcenters`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workcenters []Workcenter
	for rows.Next() {
		var workcenter Workcenter
		if err := rows.Scan(&workcenter.ID, &workcenter.CustomerID, &workcenter.ShopFloorID, &workcenter.Name, &workcenter.IsActive, &workcenter.CreatedAt, &workcenter.UpdatedAt); err != nil {
			return nil, err
		}
		workcenters = append(workcenters, workcenter)
	}
	return workcenters, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Workcenter, error) {
	query := `SELECT id, customer_id, shop_floor_id, name, is_active, created_at, updated_at FROM workcenters WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var workcenter Workcenter
	if err := row.Scan(&workcenter.ID, &workcenter.CustomerID, &workcenter.ShopFloorID, &workcenter.Name, &workcenter.IsActive, &workcenter.CreatedAt, &workcenter.UpdatedAt); err != nil {
		return Workcenter{}, err
	}
	return workcenter, nil
}

func (r *repository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]Workcenter, error) {
	query := `SELECT id, customer_id, shop_floor_id, name, is_active, created_at, updated_at FROM workcenters WHERE customer_id = $1`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workcenters []Workcenter
	for rows.Next() {
		var workcenter Workcenter
		if err := rows.Scan(&workcenter.ID, &workcenter.CustomerID, &workcenter.ShopFloorID, &workcenter.Name, &workcenter.IsActive, &workcenter.CreatedAt, &workcenter.UpdatedAt); err != nil {
			return nil, err
		}
		workcenters = append(workcenters, workcenter)
	}
	return workcenters, nil
}

func (r *repository) CountByCustomerID(ctx context.Context, customerID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM workcenters WHERE customer_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, workcenter Workcenter) (Workcenter, error) {
	query := `UPDATE workcenters SET customer_id = $2, shop_floor_id = $3, name = $4, is_active = $5, created_at = $6, updated_at = $7 WHERE id = $1 RETURNING id`
	_, err := r.db.ExecContext(ctx, query, workcenter.ID, workcenter.CustomerID, workcenter.ShopFloorID, workcenter.Name, workcenter.IsActive, workcenter.CreatedAt, workcenter.UpdatedAt)
	if err != nil {
		return Workcenter{}, err
	}
	return workcenter, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workcenters WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
