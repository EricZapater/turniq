package scheduleentries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type ScheduleFilter struct {
	CustomerID   *uuid.UUID
	ShopfloorID  *uuid.UUID
	ShiftID      *uuid.UUID
	WorkcenterID *uuid.UUID
	JobID        *uuid.UUID
	OperatorID   *uuid.UUID
	StartDate    *string // YYYY-MM-DD
	EndDate      *string // YYYY-MM-DD
}

type Repository interface {
	Create(ctx context.Context, entry ScheduleEntry) (ScheduleEntry, error)
	FindByID(ctx context.Context, id uuid.UUID) (ScheduleEntry, error)
	FindAll(ctx context.Context) ([]ScheduleEntry, error)
	FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]ScheduleEntry, error)
	FindByShopfloorAndDate(ctx context.Context, shopfloorID uuid.UUID, date string) ([]ScheduleEntry, error)
	FindByOperatorAndDate(ctx context.Context, operatorID uuid.UUID, date string) ([]ScheduleEntry, error)
	Search(ctx context.Context, filter ScheduleFilter) ([]ScheduleEntry, error)
	Update(ctx context.Context, entry ScheduleEntry) (ScheduleEntry, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Sync(ctx context.Context, shopfloorID uuid.UUID, date string, entries []ScheduleEntry) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, entry ScheduleEntry) (ScheduleEntry, error) {
	query := `INSERT INTO schedule_entries (
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`
	
	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.CustomerID, entry.ShopfloorID, entry.ShiftID, entry.WorkcenterID, entry.JobID, entry.OperatorID,
		entry.Date, entry.Order, entry.StartTime, entry.EndTime, entry.IsCompleted, entry.CreatedAt, entry.UpdatedAt,
	)
	if err != nil {
		return ScheduleEntry{}, err
	}
	return entry, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	var entry ScheduleEntry
	err := row.Scan(
		&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
		&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
	)
	if err != nil {
		return ScheduleEntry{}, err
	}
	return entry, nil
}

func (r *repository) FindAll(ctx context.Context) ([]ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ScheduleEntry
	for rows.Next() {
		var entry ScheduleEntry
		err := rows.Scan(
			&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
			&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries WHERE customer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ScheduleEntry
	for rows.Next() {
		var entry ScheduleEntry
		err := rows.Scan(
			&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
			&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) Search(ctx context.Context, filter ScheduleFilter) ([]ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries 
	WHERE 1=1`

	var args []interface{}
	argId := 1

	if filter.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argId)
		args = append(args, *filter.CustomerID)
		argId++
	}
	if filter.ShopfloorID != nil {
		query += fmt.Sprintf(" AND shopfloor_id = $%d", argId)
		args = append(args, *filter.ShopfloorID)
		argId++
	}
	if filter.ShiftID != nil {
		query += fmt.Sprintf(" AND shift_id = $%d", argId)
		args = append(args, *filter.ShiftID)
		argId++
	}
	if filter.WorkcenterID != nil {
		query += fmt.Sprintf(" AND workcenter_id = $%d", argId)
		args = append(args, *filter.WorkcenterID)
		argId++
	}
	if filter.JobID != nil {
		query += fmt.Sprintf(" AND job_id = $%d", argId)
		args = append(args, *filter.JobID)
		argId++
	}
	if filter.OperatorID != nil {
		query += fmt.Sprintf(" AND operator_id = $%d", argId)
		args = append(args, *filter.OperatorID)
		argId++
	}
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND date::date >= $%d::date", argId)
		args = append(args, *filter.StartDate)
		argId++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND date::date <= $%d::date", argId)
		args = append(args, *filter.EndDate)
		argId++
	}

	query += ` ORDER BY date DESC, "order" ASC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ScheduleEntry
	for rows.Next() {
		var entry ScheduleEntry
		err := rows.Scan(
			&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
			&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) Update(ctx context.Context, entry ScheduleEntry) (ScheduleEntry, error) {
	query := `UPDATE schedule_entries SET 
		customer_id = $2, shopfloor_id = $3, shift_id = $4, workcenter_id = $5, job_id = $6, operator_id = $7,
		date = $8, "order" = $9, start_time = $10, end_time = $11, is_completed = $12, updated_at = $13
	WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		entry.ID, entry.CustomerID, entry.ShopfloorID, entry.ShiftID, entry.WorkcenterID, entry.JobID, entry.OperatorID,
		entry.Date, entry.Order, entry.StartTime, entry.EndTime, entry.IsCompleted, entry.UpdatedAt,
	)
	if err != nil {
		return ScheduleEntry{}, err
	}
	return entry, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM schedule_entries WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *repository) FindByShopfloorAndDate(ctx context.Context, shopfloorID uuid.UUID, date string) ([]ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries 
	WHERE shopfloor_id = $1 AND date::date = $2::date
	ORDER BY "order" ASC`

	rows, err := r.db.QueryContext(ctx, query, shopfloorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ScheduleEntry
	for rows.Next() {
		var entry ScheduleEntry
		err := rows.Scan(
			&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
			&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) FindByOperatorAndDate(ctx context.Context, operatorID uuid.UUID, date string) ([]ScheduleEntry, error) {
	query := `SELECT 
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	FROM schedule_entries 
	WHERE operator_id = $1 AND date::date = $2::date
	ORDER BY "order" ASC`

	rows, err := r.db.QueryContext(ctx, query, operatorID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []ScheduleEntry
	for rows.Next() {
		var entry ScheduleEntry
		err := rows.Scan(
			&entry.ID, &entry.CustomerID, &entry.ShopfloorID, &entry.ShiftID, &entry.WorkcenterID, &entry.JobID, &entry.OperatorID,
			&entry.Date, &entry.Order, &entry.StartTime, &entry.EndTime, &entry.IsCompleted, &entry.CreatedAt, &entry.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *repository) Sync(ctx context.Context, shopfloorID uuid.UUID, date string, entries []ScheduleEntry) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing for this shopfloor and date
	deleteQuery := `DELETE FROM schedule_entries WHERE shopfloor_id = $1 AND date::date = $2::date`
	if _, err := tx.ExecContext(ctx, deleteQuery, shopfloorID, date); err != nil {
		return err
	}

	// Insert new ones
	insertQuery := `INSERT INTO schedule_entries (
		id, customer_id, shopfloor_id, shift_id, workcenter_id, job_id, operator_id,
		date, "order", start_time, end_time, is_completed, created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err := stmt.ExecContext(ctx,
			entry.ID, entry.CustomerID, entry.ShopfloorID, entry.ShiftID, entry.WorkcenterID, entry.JobID, entry.OperatorID,
			entry.Date, entry.Order, entry.StartTime, entry.EndTime, entry.IsCompleted, entry.CreatedAt, entry.UpdatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
