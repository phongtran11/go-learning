package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service defines the business logic for users
type Service interface {
	CreateUser(dto CreateUserDTO) (*DataResponseDTO, error)
}

// service implements the Service interface
type service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

// CreateUser creates a new user
func (s *service) CreateUser(dto CreateUserDTO) (*DataResponseDTO, error) {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:         dto.Email,
		Password:      string(hashedPassword),
		PhoneNumber:   dto.PhoneNumber,
		FirstName:     dto.FirstName,
		LastName:      dto.LastName,
		DateOfBirth:   dto.DateOfBirth,
		Gender:        dto.Gender,
		AvatarURL:     dto.AvatarURL,
		Status:        USER_STATUS_ACTIVE,
		EmailVerified: false,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	responseDTO := user.ToResponseDTO()
	return &responseDTO, nil
}
