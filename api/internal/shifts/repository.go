package shifts

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context,shift Shift) (Shift, error)
	FindByID(ctx context.Context,shiftID uuid.UUID) (Shift, error)
	FindByShopfloorID(ctx context.Context,shopfloorID uuid.UUID) ([]Shift, error)
	FindByCustomerID(ctx context.Context,customerID uuid.UUID) ([]Shift, error)
	Update(ctx context.Context,shift Shift) (Shift, error)
	Delete(ctx context.Context,shiftID uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context,shift Shift) (Shift, error) {
	query := "INSERT INTO shifts (id, customer_id, shopfloor_id, name, color, start_time, end_time, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING *"
	_, err := r.db.ExecContext(ctx, query, shift.ID, shift.CustomerID, shift.ShopfloorID, shift.Name, shift.Color, shift.StartTime, shift.EndTime, shift.IsActive, shift.CreatedAt, shift.UpdatedAt)
	if err != nil {
		return Shift{}, err
	}
	return shift, nil
}

func (r *repository) FindByID(ctx context.Context,shiftID uuid.UUID) (Shift, error) {
	query := "SELECT id, customer_id, shopfloor_id, name, color, start_time, end_time, is_active, created_at, updated_at FROM shifts WHERE id = $1"
	row := r.db.QueryRowContext(ctx, query, shiftID)
	var shift Shift
	err := row.Scan(&shift.ID, &shift.CustomerID, &shift.ShopfloorID, &shift.Name, &shift.Color, &shift.StartTime, &shift.EndTime, &shift.IsActive, &shift.CreatedAt, &shift.UpdatedAt)
	if err != nil {
		return Shift{}, err
	}
	return shift, nil
}

func (r *repository) FindByShopfloorID(ctx context.Context,shopfloorID uuid.UUID) ([]Shift, error) {
	query := "SELECT id, customer_id, shopfloor_id, name, color, start_time, end_time, is_active, created_at, updated_at FROM shifts WHERE shopfloor_id = $1"
	rows, err := r.db.QueryContext(ctx, query, shopfloorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shifts []Shift
	for rows.Next() {
		var shift Shift
		err := rows.Scan(&shift.ID, &shift.CustomerID, &shift.ShopfloorID, &shift.Name, &shift.Color, &shift.StartTime, &shift.EndTime, &shift.IsActive, &shift.CreatedAt, &shift.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (r *repository) FindByCustomerID(ctx context.Context,customerID uuid.UUID) ([]Shift, error) {
	query := "SELECT id, customer_id, shopfloor_id, name, color, start_time, end_time, is_active, created_at, updated_at FROM shifts WHERE customer_id = $1"
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shifts []Shift
	for rows.Next() {
		var shift Shift
		err := rows.Scan(&shift.ID, &shift.CustomerID, &shift.ShopfloorID, &shift.Name, &shift.Color, &shift.StartTime, &shift.EndTime, &shift.IsActive, &shift.CreatedAt, &shift.UpdatedAt)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (r *repository) Update(ctx context.Context,shift Shift) (Shift, error) {
	query := "UPDATE shifts SET customer_id = $2, shopfloor_id = $3, name = $4, color = $5, start_time = $6, end_time = $7, is_active = $8, updated_at = $9 WHERE id = $10 RETURNING *"
	_, err := r.db.ExecContext(ctx, query, shift.CustomerID, shift.ShopfloorID, shift.Name, shift.Color, shift.StartTime, shift.EndTime, shift.IsActive, shift.UpdatedAt, shift.ID)
	if err != nil {
		return Shift{}, err
	}
	return shift, nil
}

func (r *repository) Delete(ctx context.Context,shiftID uuid.UUID) error {
	query := "DELETE FROM shifts WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, shiftID)
	if err != nil {
		return err
	}
	return nil
}
