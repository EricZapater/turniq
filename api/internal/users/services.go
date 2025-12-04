package users

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(request UserRequest) (User, error)
	GetAll() ([]User, error)
	GetByID(id string) (User, error)
	GetByCustomerID(customerID string) ([]User, error)
	Update(id string, request UserRequest) (User, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(request UserRequest) (User, error) {
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
		Username:   request.Username,
		Email:      request.Email,
		Password:   string(hashedPassword),
		CustomerID: customerID,
		IsActive:   request.IsActive,
		IsAdmin:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.Create(user)
}

func (s *service) GetAll() ([]User, error) {
	return s.repo.GetAll()
}

func (s *service) GetByID(id string) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	return s.repo.GetByID(parsedId)
}

func (s *service) GetByCustomerID(customerID string) ([]User, error) {
	parsedCustomerID, err := uuid.Parse(customerID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetByCustomerID(parsedCustomerID)
}

func (s *service) Update(id string, request UserRequest) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	customerID, err := uuid.Parse(request.CustomerID)
	if err != nil {
		return User{}, err
	}
	user, err := s.repo.GetByID(parsedId)
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
	return s.repo.Update(user)
}

func (s *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(parsedId)
}
