package customers

import (
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(customer Customer) (Customer, error)
	GetAll() ([]Customer, error)
	GetByID(id uuid.UUID) (Customer, error)
	Update(customer Customer) (Customer, error)
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(customer Customer) (Customer, error) {
	query := `INSERT INTO customers (id, name, email, vat_number, phone, address, city, state, zip_code, country, language, contact_name, status, plan, billing_cycle, price, trial_ends_at, internal_notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`
	_, err := r.db.Exec(query, customer.ID, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.TrialEndsAt, customer.InternalNotes, customer.CreatedAt, customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) GetAll() ([]Customer, error) {
	query := `SELECT * FROM customers`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.VatNumber, &customer.Phone, &customer.Address, &customer.City, &customer.State, &customer.ZipCode, &customer.Country, &customer.Language, &customer.ContactName, &customer.Status, &customer.Plan, &customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, &customer.InternalNotes, &customer.CreatedAt, &customer.UpdatedAt)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *repository) GetByID(id uuid.UUID) (Customer, error) {
	query := `SELECT * FROM customers WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var customer Customer
	err := row.Scan(&customer.ID, &customer.Name, &customer.Email, &customer.VatNumber, &customer.Phone, &customer.Address, &customer.City, &customer.State, &customer.ZipCode, &customer.Country, &customer.Language, &customer.ContactName, &customer.Status, &customer.Plan, &customer.BillingCycle, &customer.Price, &customer.TrialEndsAt, &customer.InternalNotes, &customer.CreatedAt, &customer.UpdatedAt)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) Update(customer Customer) (Customer, error) {
	query := `UPDATE customers SET name = $2, email = $3, vat_number = $4, phone = $5, address = $6, city = $7, state = $8, zip_code = $9, country = $10, language = $11, contact_name = $12, status = $13, plan = $14, billing_cycle = $15, price = $16, trial_ends_at = $17, internal_notes = $18, updated_at = $19 WHERE id = $1`
	_, err := r.db.Exec(query, customer.Name, customer.Email, customer.VatNumber, customer.Phone, customer.Address, customer.City, customer.State, customer.ZipCode, customer.Country, customer.Language, customer.ContactName, customer.Status, customer.Plan, customer.BillingCycle, customer.Price, customer.TrialEndsAt, customer.InternalNotes, customer.UpdatedAt, customer.ID)
	if err != nil {
		return Customer{}, err
	}
	return customer, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}