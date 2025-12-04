package customers

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, customer Customer) (Customer, error)
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id uuid.UUID) (Customer, error)
	Update(ctx context.Context, customer Customer) (Customer, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, customer Customer) (Customer, error) {
	query := `INSERT INTO customers (id, tenant_id, name, email, vat_number, phone, address, city, state, zip_code, country, language, contact_name, status, plan, billing_cycle, price, trial_ends_at, internal_notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)`
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.TenantID, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.TrialEndsAt, customer.InternalNotes, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) GetAll(ctx context.Context) ([]Customer, error) {
	query := `SELECT id, tenant_id, name, email, vat_number, phone, address, city, state, zip_code, country, language, contact_name, status, plan, billing_cycle, price, trial_ends_at, internal_notes, created_at, updated_at FROM customers`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.TenantID, &customer.Name, &customer.Email, &customer.VatNumber, &customer.Phone, &customer.Address, &customer.City, &customer.State, &customer.ZipCode, &customer.Country, &customer.Language, &customer.ContactName, &customer.Status, &customer.Plan, &customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, &customer.InternalNotes, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	query := `SELECT id, tenant_id, name, email, vat_number, phone, address, city, state, zip_code, country, language, contact_name, status, plan, billing_cycle, price, trial_ends_at, internal_notes, created_at, updated_at FROM customers WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var customer Customer
	err := row.Scan(&customer.ID, &customer.TenantID, &customer.Name, &customer.Email, &customer.VatNumber, &customer.Phone, &customer.Address, &customer.City, &customer.State, &customer.ZipCode, &customer.Country, &customer.Language, &customer.ContactName, &customer.Status, &customer.Plan, &customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, &customer.InternalNotes, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) Update(ctx context.Context, customer Customer) (Customer, error) {
	query := `UPDATE customers SET tenant_id = $2, name = $3, email = $4, vat_number = $5, phone = $6, address = $7, city = $8, state = $9, zip_code = $10, country = $11, language = $12, contact_name = $13, status = $14, plan = $15, billing_cycle = $16, price = $17, trial_ends_at = $18, internal_notes = $19, updated_at = $20 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.TenantID, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.TrialEndsAt, customer.InternalNotes, customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}