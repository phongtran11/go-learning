package user

import (
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"modular-fx-fiber/internal/shared/dto/user_dto"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/models"
	"modular-fx-fiber/internal/shared/repositories"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type (
	Service interface {
		CreateUser(dto *user_dto.CreateUserDTO) (*models.UserResponseDTO, error)
		ListUsers(page int, pageSize int) ([]*models.UserResponseDTO, int64, error)
		GetMe(userID uint64) (*models.UserResponseDTO, error)
	}

	service struct {
		logger   *logger.ZapLogger
		userRepo repositories.UserRepository
	}
)

// NewService creates a new user service
func NewService(
	logger *logger.ZapLogger,
	userRepo repositories.UserRepository,
) Service {
	return &service{
		logger:   logger,
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (s *service) CreateUser(dto *user_dto.CreateUserDTO) (*models.UserResponseDTO, error) {
	// Check if email already exists
	existingUser, err := s.userRepo.GetByEmail(dto.Email)
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
	user := &models.User{
		Email:       dto.Email,
		Password:    string(hashedPassword),
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		DateOfBirth: dto.DateOfBirth,
		Gender:      dto.Gender,
		Status:      models.USER_STATUS_ACTIVE,
	}

	s.logger.Info("user", zap.Any("user", user))

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	userResponse := user.ToResponseDTO()
	return userResponse, nil
}

func (s *service) ListUsers(page int, pageSize int) ([]*models.UserResponseDTO, int64, error) {
	users, totalCount, err := s.userRepo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	userResponse := make([]*models.UserResponseDTO, 0)
	for _, user := range users {
		userResponse = append(userResponse, user.ToResponseDTO())
	}

	return userResponse, totalCount, nil
}

func (s *service) GetMe(userID uint64) (*models.UserResponseDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToResponseDTO()
	return userResponse, nil
}
