package users

import (
	"api/internal/customers"
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, request UserRequest) (User, error)
	CreateAdmin(ctx context.Context)error
	FindAll(ctx context.Context) ([]User, error)	
	FindByID(ctx context.Context, id string) (User, error)
	FindByCustomerID(ctx context.Context, customerID string) ([]User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, id string, request UserRequest) (User, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	customerService customers.Service
}

func NewService(repo Repository, customerService customers.Service) Service {
	return &service{repo: repo, customerService: customerService}
}

func (s *service) Create(ctx context.Context, request UserRequest) (User, error) {
	var customerID uuid.UUID
	
	// Recuperar i validar is_admin
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return User{}, errors.New("invalid or missing is_admin in context")
	}
	
	if isAdmin {
		// Si és admin, usar el customerID del request
		parsedCustomerID, err := uuid.Parse(request.CustomerID)
		if err != nil {
			return User{}, err
		}
		customerID = parsedCustomerID
	} else {
		// Si no és admin, usar el customerID del context
		customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return User{}, errors.New("invalid or missing customer_id in context")
		}
		customerID = customerIDFromCtx
	}
	
	// Check limits
	customer, err := s.customerService.FindByID(ctx, customerID.String())
	if err != nil {
		return User{}, err
	}
	slog.Info("Customer found", "customer", customer)
	count, err := s.repo.CountByCustomerID(ctx, customerID)
	if err != nil {
		return User{}, err
	}

	if count >= customer.MaxUsers {
		return User{}, errors.New("max users limit reached for this tenant")
	}
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	
	user := User{
		ID:         uuid.New(),
		Username:   request.Username,
		Email:      request.Email,
		Password:   string(hashedPassword),
		CustomerID: customerID,
		IsActive:   true,
		IsAdmin:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, user)
}
func (s *service) CreateAdmin(ctx context.Context) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("ADMIN_PASSWORD")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//Crear customer 00000000-0000-0000-0000-000000000000
	CustomerRequest := customers.CustomerRequest{
		Name:       "System",
		Language: "Catalan",
		Plan: "System",
		MaxUsers:   1,
		MaxJobs:    1000,
		MaxWorkcenters: 1000,
		MaxShopFloors: 1000,
		MaxOperators: 1000,
	}
	customer, err := s.customerService.Create(ctx, CustomerRequest)
	if err != nil {
		return err
	}


	user := User{
		ID:         uuid.New(),
		Username:   "admin",
		Email:      "admin@turniq.com",
		Password:   string(hashedPassword),
		CustomerID: customer.ID,
		IsActive:   true,
		IsAdmin:    true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = s.repo.Create(ctx, user)
	return err
}

func (s *service) FindAll(ctx context.Context) ([]User, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if isAdmin{
		return s.repo.FindAll(ctx)
	}
	customerIDVal := ctx.Value("customer_id")
		customerIDFromCtx, ok := customerIDVal.(uuid.UUID)
		if !ok {
			return nil, errors.New("invalid or missing customer_id in context")
		}
		customerID := customerIDFromCtx
	return s.repo.FindByCustomerID(ctx, customerID)
}

func (s *service) FindByID(ctx context.Context, id string) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	return s.repo.FindByID(ctx, parsedId)
}

func (s *service) FindByCustomerID(ctx context.Context, customerID string) ([]User, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByCustomerID(ctx, parsedCustomerID)
}

func (s *service) FindByEmail(ctx context.Context, email string) (User, error) {
	return s.repo.FindByEmail(ctx, email)
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
	user, err := s.repo.FindByID(ctx, parsedId)
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
