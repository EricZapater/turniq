package tenants

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Create(ctx context.Context, tenant Tenant) (Tenant, error)
	GetAll(ctx context.Context) ([]Tenant, error)
	GetByID(ctx context.Context, id uuid.UUID) (Tenant, error)
	Update(ctx context.Context, tenant Tenant) (Tenant, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func(r *repository) Create(ctx context.Context, tenant Tenant) (Tenant, error) {
	query := `INSERT INTO tenants (id, customer_id, name, status, is_active, max_operators, max_workcenters, max_shop_floors, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.ExecContext(ctx, query, tenant.ID, tenant.CustomerID, tenant.Name, tenant.Status, tenant.IsActive, tenant.MaxOperators, tenant.MaxWorkcenters, tenant.MaxShopFloors, tenant.CreatedAt, tenant.UpdatedAt)
	if err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) GetAll(ctx context.Context) ([]Tenant, error) {
	query := `SELECT * FROM tenants`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	tenants := []Tenant{}
	for rows.Next() {
		var tenant Tenant
		if err := rows.Scan(&tenant.ID, &tenant.CustomerID, &tenant.Name, &tenant.Status, &tenant.IsActive, &tenant.MaxOperators, &tenant.MaxWorkcenters, &tenant.MaxShopFloors, &tenant.CreatedAt, &tenant.UpdatedAt); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func(r *repository) GetByID(ctx context.Context, id uuid.UUID) (Tenant, error) {
	query := `SELECT * FROM tenants WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	var tenant Tenant
	if err := row.Scan(&tenant.ID, &tenant.CustomerID, &tenant.Name, &tenant.Status, &tenant.IsActive, &tenant.MaxOperators, &tenant.MaxWorkcenters, &tenant.MaxShopFloors, &tenant.CreatedAt, &tenant.UpdatedAt); err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) Update(ctx context.Context, tenant Tenant) (Tenant, error) {
	query := `UPDATE tenants SET customer_id = $2, name = $3, status = $4, is_active = $5, max_operators = $6, max_workcenters = $7, max_shop_floors = $8, updated_at = $9 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, tenant.ID, tenant.CustomerID, tenant.Name, tenant.Status, tenant.IsActive, tenant.MaxOperators, tenant.MaxWorkcenters, tenant.MaxShopFloors, tenant.UpdatedAt)
	if err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
	
