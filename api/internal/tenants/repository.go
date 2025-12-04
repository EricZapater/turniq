package tenants

import (
	"database/sql"

	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
}

type Repository interface {
	Create(tenant Tenant) (Tenant, error)
	GetAll() ([]Tenant, error)
	GetByID(id uuid.UUID) (Tenant, error)
	Update(tenant Tenant) (Tenant, error)
	Delete(id uuid.UUID) error
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func(r *repository) Create(tenant Tenant) (Tenant, error) {
	query := `INSERT INTO tenants (id, customer_id, name, status, is_active, max_operators, max_workcenters, max_shop_floors, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query, tenant.ID, tenant.CustomerID, tenant.Name, tenant.Status, tenant.IsActive, tenant.MaxOperators, tenant.MaxWorkcenters, tenant.MaxShopFloors, tenant.CreatedAt, tenant.UpdatedAt)
	if err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) GetAll() ([]Tenant, error) {
	query := `SELECT * FROM tenants`

	rows, err := r.db.Query(query)
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

func(r *repository) GetByID(id uuid.UUID) (Tenant, error) {
	query := `SELECT * FROM tenants WHERE id = $1`

	row := r.db.QueryRow(query, id)
	var tenant Tenant
	if err := row.Scan(&tenant.ID, &tenant.CustomerID, &tenant.Name, &tenant.Status, &tenant.IsActive, &tenant.MaxOperators, &tenant.MaxWorkcenters, &tenant.MaxShopFloors, &tenant.CreatedAt, &tenant.UpdatedAt); err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) Update(tenant Tenant) (Tenant, error) {
	query := `UPDATE tenants SET customer_id = $2, name = $3, status = $4, is_active = $5, max_operators = $6, max_workcenters = $7, max_shop_floors = $8, updated_at = $9 WHERE id = $1`

	_, err := r.db.Exec(query, tenant.ID, tenant.CustomerID, tenant.Name, tenant.Status, tenant.IsActive, tenant.MaxOperators, tenant.MaxWorkcenters, tenant.MaxShopFloors, tenant.UpdatedAt)
	if err != nil {
		return Tenant{}, err
	}

	return tenant, nil
}

func(r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM tenants WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
	
