package timeentries

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TimeEntryFilter struct {
	CustomerID *uuid.UUID
	OperatorID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
}

type Repository interface {
	Create(ctx context.Context, entry TimeEntry) (TimeEntry, error)
	FindByID(ctx context.Context, id uuid.UUID) (TimeEntry, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID)([]TimeEntry, error)
	FindByOperatorID(ctx context.Context, operatorID uuid.UUID)([]TimeEntry, error)
	FindCurrent(ctx context.Context, operatorID uuid.UUID)(TimeEntry, error)
	FindAll(ctx context.Context) ([]TimeEntry, error)
	Search(ctx context.Context, filter TimeEntryFilter) ([]TimeEntry, error)
	Update(ctx context.Context, entry TimeEntry) (TimeEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entry TimeEntry) (TimeEntry, error) {
	query := `INSERT INTO time_entries (
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.OperatorID, entry.WorkcenterID,
		entry.CheckIn, entry.CheckOut, entry.CreatedAt, entry.UpdatedAt,
	)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	var entry TimeEntry
	err := row.Scan(
		&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
		&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (r *repository) FindAll(ctx context.Context) ([]TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		err := rows.Scan(
			&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
			&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries WHERE customer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		err := rows.Scan(
			&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
			&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) Search(ctx context.Context, filter TimeEntryFilter) ([]TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries WHERE 1=1`

	var args []interface{}
	argId := 1

	if filter.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argId)
		args = append(args, *filter.CustomerID)
		argId++
	}
	if filter.OperatorID != nil {
		query += fmt.Sprintf(" AND operator_id = $%d", argId)
		args = append(args, *filter.OperatorID)
		argId++
	}
	// check_in range
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND check_in >= $%d", argId)
		args = append(args, *filter.StartDate)
		argId++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND check_in <= $%d", argId)
		args = append(args, *filter.EndDate)
		argId++
	}

	query += " ORDER BY check_in DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		err := rows.Scan(
			&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
			&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) FindByOperatorID(ctx context.Context, operatorID uuid.UUID) ([]TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries WHERE operator_id = $1`

	rows, err := r.db.QueryContext(ctx, query, operatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var entry TimeEntry
		err := rows.Scan(
			&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
			&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) FindCurrent(ctx context.Context, operatorID uuid.UUID) (TimeEntry, error) {
	query := `SELECT 
		id, operator_id, workcenter_id, check_in, check_out, created_at, updated_at
	FROM time_entries WHERE operator_id = $1 AND check_out IS NULL`

	row := r.db.QueryRowContext(ctx, query, operatorID)
	var entry TimeEntry
	err := row.Scan(
		&entry.ID, &entry.OperatorID, &entry.WorkcenterID,
		&entry.CheckIn, &entry.CheckOut, &entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (r *repository) Update(ctx context.Context, entry TimeEntry) (TimeEntry, error) {
	query := `UPDATE time_entries SET 
		operator_id = $2, workcenter_id = $3, check_in = $4, check_out = $5, updated_at = $6
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.OperatorID, entry.WorkcenterID,
		entry.CheckIn, entry.CheckOut, entry.UpdatedAt,
	)
	if err != nil {
		return TimeEntry{}, err
	}
	return entry, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM time_entries WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
