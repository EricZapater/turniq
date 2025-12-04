package operators

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, operator Operator) error
	FindByID(ctx context.Context, id uuid.UUID) (Operator, error)
	FindAll(ctx context.Context) ([]Operator, error)
	Update(ctx context.Context, operator Operator) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, operator Operator) error {
	query := `INSERT INTO operators (id, tenant_id, shop_floor_id, customer_id, code, name, surname, vat_number, is_active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query, operator.ID, operator.TenantID, operator.ShopFloorID, operator.CustomerID, operator.Code, operator.Name, operator.Surname, operator.VatNumber, operator.IsActive, operator.CreatedAt, operator.UpdatedAt)
	return err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Operator, error) {
	query := `SELECT id, tenant_id, shop_floor_id, customer_id, code, name, surname, vat_number, is_active, created_at, updated_at 
	FROM operators WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var operator Operator
	err := row.Scan(&operator.ID, &operator.TenantID, &operator.ShopFloorID, &operator.CustomerID, &operator.Code, &operator.Name, &operator.Surname, &operator.VatNumber, &operator.IsActive, &operator.CreatedAt, &operator.UpdatedAt)
	if err != nil {
		return Operator{}, err
	}
	return operator, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Operator, error) {
	query := `SELECT id, tenant_id, shop_floor_id, customer_id, code, name, surname, vat_number, is_active, created_at, updated_at 
	FROM operators`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	operators := []Operator{}
	for rows.Next() {
		var operator Operator
		err := rows.Scan(&operator.ID, &operator.TenantID, &operator.ShopFloorID, &operator.CustomerID, &operator.Code, &operator.Name, &operator.Surname, &operator.VatNumber, &operator.IsActive, &operator.CreatedAt, &operator.UpdatedAt)
		if err != nil {
			return nil, err
		}
		operators = append(operators, operator)
	}
	return operators, nil
}

func (r *repository) Update(ctx context.Context, operator Operator) error {
	query := `UPDATE operators SET tenant_id = $2, shop_floor_id = $3, customer_id = $4, code = $5, name = $6, surname = $7, vat_number = $8, is_active = $9, updated_at = $10 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, operator.ID, operator.TenantID, operator.ShopFloorID, operator.CustomerID, operator.Code, operator.Name, operator.Surname, operator.VatNumber, operator.IsActive, operator.UpdatedAt)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM operators WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
