package shopfloors

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, shopfloor Shopfloor) error
	FindAll(ctx context.Context) ([]Shopfloor, error)
	FindByID(ctx context.Context, id int) (Shopfloor, error)
	Update(ctx context.Context, shopfloor Shopfloor) error
	Delete(ctx context.Context, id int) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, shopfloor Shopfloor) error {
	query := `INSERT INTO shopfloors (id, tenant_id, customer_id, name, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, shopfloor.ID, shopfloor.TenantID, shopfloor.CustomerID, shopfloor.Name, shopfloor.CreatedAt, shopfloor.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]Shopfloor, error) {
	query := `SELECT * FROM shopfloors`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shopfloors := []Shopfloor{}
	for rows.Next() {
		var shopfloor Shopfloor
		if err := rows.Scan(&shopfloor.ID, &shopfloor.TenantID, &shopfloor.CustomerID, &shopfloor.Name, &shopfloor.CreatedAt, &shopfloor.UpdatedAt); err != nil {
			return nil, err
		}
		shopfloors = append(shopfloors, shopfloor)
	}
	return shopfloors, nil
}

func (r *repository) FindByID(ctx context.Context, id int) (Shopfloor, error) {
	query := `SELECT * FROM shopfloors WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var shopfloor Shopfloor
	if err := row.Scan(&shopfloor.ID, &shopfloor.TenantID, &shopfloor.CustomerID, &shopfloor.Name, &shopfloor.CreatedAt, &shopfloor.UpdatedAt); err != nil {
		return Shopfloor{}, err
	}
	return shopfloor, nil
}

func (r *repository) Update(ctx context.Context, shopfloor Shopfloor) error {
	query := `UPDATE shopfloors SET tenant_id = $2, customer_id = $3, name = $4, updated_at = $5 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, shopfloor.ID, shopfloor.TenantID, shopfloor.CustomerID, shopfloor.Name, shopfloor.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM shopfloors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
