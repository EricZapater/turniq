package payments

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PaymentFilter struct {
	CustomerID *uuid.UUID
	StartDate  *time.Time
	EndDate    *time.Time
}

type Repository interface {
	Create(ctx context.Context, payment Payment) (Payment, error)
	FindAll(ctx context.Context) ([]Payment, error)
	FindById(ctx context.Context, id uuid.UUID) (Payment, error)
	FindByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error)
	Search(ctx context.Context, filter PaymentFilter) ([]Payment, error)
	Update(ctx context.Context, payment Payment) (Payment, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, payment Payment) (Payment, error) {
	query := `INSERT INTO payments (id, customer_id, amount, currency, payment_method, status, due_date, paid_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, payment.ID, payment.CustomerID, payment.Amount, payment.Currency, payment.PaymentMethod, payment.Status, payment.DueDate, payment.PaidAt)
	if err != nil {
		return Payment{}, err
	}
	return Payment{}, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Payment, error) {
	query := `SELECT id, customer_id, amount, currency, payment_method, status, due_date, paid_at FROM payments`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []Payment
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.ID, &payment.CustomerID, &payment.Amount, &payment.Currency, &payment.PaymentMethod, &payment.Status, &payment.DueDate, &payment.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *repository) FindById(ctx context.Context, id uuid.UUID) (Payment, error) {
	query := `SELECT id, customer_id, amount, currency, payment_method, status, due_date, paid_at FROM payments WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	payment := Payment{}
	if err := row.Scan(&payment.ID, &payment.CustomerID, &payment.Amount, &payment.Currency, &payment.PaymentMethod, &payment.Status, &payment.DueDate, &payment.PaidAt); err != nil {
		return Payment{}, err
	}
	return payment, nil
}

func (r *repository) FindByCustomerId(ctx context.Context, customerId uuid.UUID) ([]Payment, error) {
	query := `SELECT id, customer_id, amount, currency, payment_method, status, due_date, paid_at FROM payments WHERE customer_id = $1`
	rows, err := r.db.QueryContext(ctx, query, customerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var payments []Payment
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.ID, &payment.CustomerID, &payment.Amount, &payment.Currency, &payment.PaymentMethod, &payment.Status, &payment.DueDate, &payment.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *repository) Search(ctx context.Context, filter PaymentFilter) ([]Payment, error) {
	query := `SELECT id, customer_id, amount, currency, payment_method, status, due_date, paid_at FROM payments WHERE 1=1`
	var args []interface{}
	argId := 1

	if filter.CustomerID != nil {
		query += fmt.Sprintf(" AND customer_id = $%d", argId)
		args = append(args, *filter.CustomerID)
		argId++
	}
	if filter.StartDate != nil {
		query += fmt.Sprintf(" AND paid_at >= $%d", argId)
		args = append(args, *filter.StartDate)
		argId++
	}
	if filter.EndDate != nil {
		query += fmt.Sprintf(" AND paid_at <= $%d", argId)
		args = append(args, *filter.EndDate)
		argId++
	}

	query += " ORDER BY paid_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		payment := Payment{}
		if err := rows.Scan(&payment.ID, &payment.CustomerID, &payment.Amount, &payment.Currency, &payment.PaymentMethod, &payment.Status, &payment.DueDate, &payment.PaidAt); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, nil
}

func (r *repository) Update(ctx context.Context, payment Payment) (Payment, error) {
	query := `UPDATE payments SET customer_id = $2, amount = $3, currency = $4, payment_method = $5, status = $6, due_date = $7, paid_at = $8 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, payment.ID, payment.CustomerID, payment.Amount, payment.Currency, payment.PaymentMethod, payment.Status, payment.DueDate, payment.PaidAt)
	if err != nil {
		return Payment{}, err
	}
	return Payment{}, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payments WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
