package auth

import (
	"errors"
	"fmt"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/modules/user"
	"modular-fx-fiber/internal/shared/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserNotActive       = errors.New("user is not active")
)

// Service defines the business logic for authentication
type Service interface {
	Login(dto LoginDTO) (*DataResponseDTO, error)
	Register(dto RegisterDTO) (*DataResponseDTO, error)
	RefreshToken(dto RefreshTokenDTO) (*DataResponseDTO, error)
}

// service implements the Service interface
type service struct {
	config      *config.Config
	userService user.Service
	userRepo    user.Repository
	authRepo    Repository
	logger      *logger.ZapLogger
}

// NewService creates a new auth service
func NewService(
	config *config.Config,
	userService user.Service,
	userRepo user.Repository,
	authRepo Repository,
	logger *logger.ZapLogger,
) Service {
	return &service{
		config:      config,
		userService: userService,
		userRepo:    userRepo,
		authRepo:    authRepo,
		logger:      logger,
	}
}

// Login authenticates a user and returns tokens
func (s *service) Login(dto LoginDTO) (*DataResponseDTO, error) {
	// Get user by email
	u, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error fetching user by email: %v", err))
		return nil, err
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if u.Status != user.USER_STATUS_ACTIVE {
		return nil, ErrUserNotActive
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login timestamp
	now := time.Now()
	u.LastLoginAt = &now
	err = s.userRepo.Update(u)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error updating last login time: %v", err))
		return nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokens(u)
	if err != nil {
		return nil, err
	}

	// Create response
	responseData := DataResponseDTO{
		Success: true,
		Data: AuthResponseDTO{
			Token: *tokens,
		},
	}

	return &responseData, nil
}

// Register creates a new user and returns tokens
func (s *service) Register(dto RegisterDTO) (*DataResponseDTO, error) {
	// Convert RegisterDTO to user.CreateUserDTO
	createUserDto := user.CreateUserDTO{
		Email:       dto.Email,
		Password:    dto.Password,
		PhoneNumber: dto.PhoneNumber,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		DateOfBirth: dto.DateOfBirth,
		Gender:      dto.Gender,
	}

	// Create user
	_, err := s.userService.CreateUser(createUserDto)
	if err != nil {
		return nil, err
	}

	// Get complete user
	createdUser, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error fetching created user: %v", err))
		return nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokens(createdUser)
	if err != nil {
		return nil, err
	}

	// Create response
	responseData := DataResponseDTO{
		Success: true,
		Data: AuthResponseDTO{
			Token: *tokens,
		},
	}

	return &responseData, nil
}

// RefreshToken validates a refresh token and issues new tokens
func (s *service) RefreshToken(dto RefreshTokenDTO) (*DataResponseDTO, error) {
	// Get refresh token from database
	savedToken, err := s.authRepo.GetRefreshToken(dto.RefreshToken)
	if err != nil || savedToken == nil {
		return nil, ErrInvalidRefreshToken
	}

	// Get user
	userId := uuid.UUID{}
	err = userId.Scan(savedToken.UserID)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error parsing user ID: %v", err))
		return nil, err
	}

	userObj, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}
	if userObj == nil {
		return nil, ErrInvalidRefreshToken
	}

	// Remove used refresh token
	err = s.authRepo.DeleteRefreshToken(dto.RefreshToken)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error deleting used refresh token: %v", err))
	}

	// Generate new tokens
	tokens, err := s.generateTokens(userObj)
	if err != nil {
		return nil, err
	}

	responseData := DataResponseDTO{
		Success: true,
		Data: AuthResponseDTO{
			Token: *tokens,
		},
	}

	return &responseData, nil
}

// generateTokens generates JWT access and refresh tokens
func (s *service) generateTokens(user *user.User) (*TokenResponseDTO, error) {
	// Get JWT config
	jwtSecret := []byte(s.config.JWT.Secret)
	accessTokenExpiry := time.Duration(s.config.JWT.AccessExpiryMinutes) * time.Minute
	refreshTokenExpiry := time.Duration(s.config.JWT.RefreshExpiryDays) * 24 * time.Hour

	// Create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = user.ID
	accessClaims["email"] = user.Email
	accessClaims["exp"] = time.Now().Add(accessTokenExpiry).Unix()

	// Sign access token
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error signing access token: %v", err))
		return nil, err
	}

	// Create refresh token
	refreshToken := uuid.New().String()

	// Save refresh token to database
	refreshTokenModel := RefreshToken{
		UserID:    uint64(user.ID),
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}

	err = s.authRepo.SaveRefreshToken(&refreshTokenModel)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Error saving refresh token: %v", err))
		return nil, err
	}

	// Create response
	return &TokenResponseDTO{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int(accessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}
