package shopfloors

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, shopfloor Shopfloor) (Shopfloor, error)
	FindAll(ctx context.Context) ([]Shopfloor, error)
	FindByID(ctx context.Context, id uuid.UUID) (Shopfloor, error)
	FindByCustomerID(ctx context.Context, id uuid.UUID) ([]Shopfloor, error)
	CountByCustomerID(ctx context.Context, id uuid.UUID) (int, error)
	Update(ctx context.Context, shopfloor Shopfloor) (Shopfloor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, shopfloor Shopfloor) (Shopfloor, error) {
	query := `INSERT INTO shopfloors (id, customer_id, name, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, shopfloor.ID, shopfloor.CustomerID, shopfloor.Name, shopfloor.CreatedAt, shopfloor.UpdatedAt)
	if err != nil {
		return Shopfloor{}, err
	}
	return shopfloor, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Shopfloor, error) {
	query := `SELECT id, customer_id, name, created_at, updated_at FROM shopfloors`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shopfloors := []Shopfloor{}
	for rows.Next() {
		var shopfloor Shopfloor
		if err := rows.Scan(&shopfloor.ID, &shopfloor.CustomerID, &shopfloor.Name, &shopfloor.CreatedAt, &shopfloor.UpdatedAt); err != nil {
			return nil, err
		}
		shopfloors = append(shopfloors, shopfloor)
	}
	return shopfloors, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Shopfloor, error) {
	query := `SELECT id, customer_id, name, created_at, updated_at FROM shopfloors WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var shopfloor Shopfloor
	if err := row.Scan(&shopfloor.ID, &shopfloor.CustomerID, &shopfloor.Name, &shopfloor.CreatedAt, &shopfloor.UpdatedAt); err != nil {
		return Shopfloor{}, err
	}
	return shopfloor, nil
}

func(r *repository) FindByCustomerID(ctx context.Context, id uuid.UUID)([]Shopfloor, error){
	query := `SELECT id, customer_id, name, created_at, updated_at FROM shopfloors WHERE customer_id = $1`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shopfloors := []Shopfloor{}
	for rows.Next() {
		var shopfloor Shopfloor
		if err := rows.Scan(&shopfloor.ID, &shopfloor.CustomerID, &shopfloor.Name, &shopfloor.CreatedAt, &shopfloor.UpdatedAt); err != nil {
			return nil, err
		}
		shopfloors = append(shopfloors, shopfloor)
	}
	return shopfloors, nil
}

func (r *repository) CountByCustomerID(ctx context.Context, id uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM shopfloors WHERE customer_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, id).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *repository) Update(ctx context.Context, shopfloor Shopfloor) (Shopfloor, error) {
	query := `UPDATE shopfloors SET customer_id = $2, name = $3, updated_at = $4 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, shopfloor.ID, shopfloor.CustomerID, shopfloor.Name, shopfloor.UpdatedAt)
	if err != nil {
		return Shopfloor{}, err
	}
	return shopfloor, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shopfloors WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
