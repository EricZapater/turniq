package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Repository interface {
	Create(user User) (User, error)
	GetAll() ([]User, error)
	GetByID(id uuid.UUID) (User, error)
	GetByCustomerID(customerID uuid.UUID) ([]User, error)
	FindByUsername(ctx context.Context, username string) (User, error)
	Update(user User) (User, error)
	Delete(id uuid.UUID) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user User) (User, error) {
	query := `INSERT INTO users (id, username, email, password, customer_id, is_admin, is_active, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CustomerID, user.IsAdmin, user.IsActive, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) GetAll() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CustomerID, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) GetByID(id uuid.UUID) (User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CustomerID, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) GetByCustomerID(customerID uuid.UUID) ([]User, error) {
	query := `SELECT * FROM users WHERE customer_id = $1`
	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CustomerID, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CustomerID, &user.IsAdmin, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Update(user User) (User, error) {
	query := `UPDATE users SET username = $2, email = $3, password = $4, customer_id = $5, is_admin = $6, is_active = $7, updated_at = $8 WHERE id = $1 RETURNING id`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.CustomerID, user.IsAdmin, user.IsActive, user.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *repository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
