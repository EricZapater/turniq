package customers

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, customer Customer) (Customer, error)
	FindAll(ctx context.Context) ([]Customer, error)
	FindByID(ctx context.Context, id uuid.UUID) (Customer, error)
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
	query := `INSERT INTO customers (id, name, email, vat_number, phone, address, city, state, zip_code, 
									country, language, contact_name, status, plan, billing_cycle, price, 
									trial_ends_at, internal_notes, max_operators, max_workcenters,
									 max_shop_floors, max_users, max_jobs, created_at, updated_at) 
									 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 
									 $10, $11, $12, $13, $14, $15, $16, 
									 $17, $18, $19, $20, 
									 $21, $22, $23, $24, $25)`
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.TrialEndsAt, customer.InternalNotes, customer.MaxOperators, customer.MaxWorkcenters, customer.MaxShopFloors, customer.MaxUsers, customer.MaxJobs, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Customer, error) {
	query := `SELECT id, name, email, vat_number, 
					phone, address, city, state, 
					zip_code, country, language, 
					contact_name, status, plan, 
					billing_cycle, price, trial_ends_at, 
					internal_notes, max_operators, 
					max_workcenters, max_shop_floors, max_users, 
					max_jobs, created_at, updated_at FROM customers`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.VatNumber, 
						&customer.Phone, &customer.Address, &customer.City, &customer.State,
					 	&customer.ZipCode, &customer.Country, &customer.Language, 
						&customer.ContactName, &customer.Status, &customer.Plan, 
						&customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, 
						&customer.InternalNotes, &customer.MaxOperators, 
						&customer.MaxWorkcenters, &customer.MaxShopFloors, &customer.MaxUsers, 
						&customer.MaxJobs, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (Customer, error) {
	query := `SELECT id, name, email, vat_number, 
					phone, address, city, state, 
					zip_code, country, language, 
					contact_name, status, plan, 
					billing_cycle, price, trial_ends_at, 
					internal_notes, max_operators, 
					max_workcenters, max_shop_floors, max_users, 
					max_jobs, created_at, updated_at FROM customers WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var customer Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.VatNumber, 
						&customer.Phone, &customer.Address, &customer.City, &customer.State,
					 	&customer.ZipCode, &customer.Country, &customer.Language, 
						&customer.ContactName, &customer.Status, &customer.Plan, 
						&customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, 
						&customer.InternalNotes, &customer.MaxOperators, 
						&customer.MaxWorkcenters, &customer.MaxShopFloors, &customer.MaxUsers, 
						&customer.MaxJobs, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) Update(ctx context.Context, customer Customer) (Customer, error) {
	query := `UPDATE customers SET name = $2, email = $3, vat_number = $4, 
						phone = $5, address = $6, city = $7, state = $8, 
						zip_code = $9, country = $10, language = $11, 
						contact_name = $12, status = $13, plan = $14, 
						billing_cycle = $15, price = $16, 
						max_operators = $17, max_workcenters = $18, 
						max_shop_floors = $19, max_users = $20, 
						max_jobs = $21, trial_ends_at = $22, internal_notes = $23, updated_at = $24 
						WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.MaxOperators, customer.MaxWorkcenters, customer.MaxShopFloors, customer.MaxUsers, customer.MaxJobs, customer.TrialEndsAt, customer.InternalNotes, customer.UpdatedAt)
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