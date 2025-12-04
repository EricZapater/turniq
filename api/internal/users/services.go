package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, request UserRequest) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByCustomerID(ctx context.Context, customerID string) ([]User, error)
	Update(ctx context.Context, id string, request UserRequest) (User, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request UserRequest) (User, error) {
	tenantID, err := uuid.Parse(request.TenantID)
	if err != nil {
		return User{}, err
	}
	customerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return User{}, err
	}
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	
	user := User{
		ID:         uuid.New(),
		TenantID:   tenantID,
		Username:   request.Username,
		Email:      request.Email,
		Password:   string(hashedPassword),
		CustomerID: customerID,
		IsActive:   request.IsActive,
		IsAdmin:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, user)
}

func (s *service) GetAll(ctx context.Context) ([]User, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetByID(ctx context.Context, id string) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	return s.repo.GetByID(ctx, parsedId)
}

func (s *service) GetByCustomerID(ctx context.Context, customerID string) ([]User, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByCustomerID(ctx, parsedCustomerID)
}

func (s *service) Update(ctx context.Context, id string, request UserRequest) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	customerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return User{}, err
	}
	user, err := s.repo.GetByID(ctx, parsedId)
	if err != nil {
		return User{}, err
	}
	
	user.Username = request.Username
	user.Email = request.Email
	
	// Only hash and update password if a new one is provided
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		user.Password = string(hashedPassword)
	}
	
	user.CustomerID = customerID
	user.IsActive = request.IsActive
	user.IsAdmin = false
	user.UpdatedAt = time.Now()
	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedId)
}
