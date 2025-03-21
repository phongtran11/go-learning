package user

import (
	"errors"
	"modular-fx-fiber/internal/shared/interfaces"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/models"
	"modular-fx-fiber/internal/shared/util"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type (
	service struct {
		logger *logger.ZapLogger
		repo   interfaces.UserRepository
	}

	UserService interface {
		CreateUser(dto *CreateUserDTO) (*models.UserResponseDTO, error)
		ListUsers(page int, pageSize int) ([]*models.UserResponseDTO, int64, error)
	}
)

// NewService creates a new user service
func NewService(
	logger *logger.ZapLogger,
	repo interfaces.UserRepository) UserService {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

// CreateUser creates a new user
func (s *service) CreateUser(dto *CreateUserDTO) (*models.UserResponseDTO, error) {
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

	// Generate random verify email code
	verifyEmailCode := util.GenerateRandomCode(6)

	// Create user
	user := &models.User{
		Email:           dto.Email,
		Password:        string(hashedPassword),
		PhoneNumber:     dto.PhoneNumber,
		FirstName:       dto.FirstName,
		LastName:        dto.LastName,
		DateOfBirth:     dto.DateOfBirth,
		Gender:          dto.Gender,
		Status:          models.USER_STATUS_ACTIVE,
		VerifyEmailCode: &verifyEmailCode,
	}

	s.logger.Info("user", zap.Any("user", user))

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	userResponse := user.ToResponseDTO()
	return userResponse, nil
}

func (s *service) ListUsers(page int, pageSize int) ([]*models.UserResponseDTO, int64, error) {
	users, totalCount, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	userResponse := make([]*models.UserResponseDTO, 0)
	for _, user := range users {
		userResponse = append(userResponse, user.ToResponseDTO())
	}

	return userResponse, totalCount, nil
}
