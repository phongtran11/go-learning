package user

import (
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Service defines the business logic for users
type Service interface {
	CreateUser(dto CreateUserDTO) (*UserResponseDTO, error)
	GetUserByID(id uuid.UUID) (*UserResponseDTO, error)
	UpdateUser(id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error)
	DeleteUser(id uuid.UUID) error
	ListUsers(page, pageSize int) (*PaginatedUsersResponse, error)
	ChangePassword(id uuid.UUID, dto ChangePasswordDTO) error
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
func (s *service) CreateUser(dto CreateUserDTO) (*UserResponseDTO, error) {
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
		PasswordHash:  string(hashedPassword),
		PhoneNumber:   dto.PhoneNumber,
		FirstName:     dto.FirstName,
		LastName:      dto.LastName,
		DateOfBirth:   dto.DateOfBirth,
		Gender:        dto.Gender,
		AvatarURL:     dto.AvatarURL,
		Status:        "active",
		EmailVerified: false,
		PhoneVerified: false,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	responseDTO := user.ToResponseDTO()
	return &responseDTO, nil
}

// GetUserByID retrieves a user by ID
func (s *service) GetUserByID(id uuid.UUID) (*UserResponseDTO, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	responseDTO := user.ToResponseDTO()
	return &responseDTO, nil
}

// UpdateUser updates an existing user
func (s *service) UpdateUser(id uuid.UUID, dto UpdateUserDTO) (*UserResponseDTO, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update fields if provided
	if dto.PhoneNumber != nil {
		user.PhoneNumber = dto.PhoneNumber
	}
	if dto.FirstName != nil {
		user.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		user.LastName = *dto.LastName
	}
	if dto.DateOfBirth != nil {
		user.DateOfBirth = dto.DateOfBirth
	}
	if dto.Gender != nil {
		user.Gender = dto.Gender
	}
	if dto.AvatarURL != nil {
		user.AvatarURL = dto.AvatarURL
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	responseDTO := user.ToResponseDTO()
	return &responseDTO, nil
}

// DeleteUser deletes a user
func (s *service) DeleteUser(id uuid.UUID) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	return s.repo.Delete(id)
}

// ListUsers retrieves a paginated list of users
func (s *service) ListUsers(page, pageSize int) (*PaginatedUsersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	users, totalCount, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	// Convert to response DTOs
	userDTOs := make([]UserResponseDTO, len(users))
	for i, user := range users {
		userDTOs[i] = user.ToResponseDTO()
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))

	return &PaginatedUsersResponse{
		Users:      userDTOs,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ChangePassword changes a user's password
func (s *service) ChangePassword(id uuid.UUID, dto ChangePasswordDTO) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(dto.CurrentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(user)
}
