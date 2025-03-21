package auth

import (
	"errors"
	"modular-fx-fiber/internal/core/config"
	"modular-fx-fiber/internal/modules/mailer"
	"modular-fx-fiber/internal/modules/user"
	"modular-fx-fiber/internal/shared/interfaces"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrUserNotActive       = errors.New("user is not active")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidVerifyCode   = errors.New("invalid verification code")
)

type (
	AuthService interface {
		Login(dto LoginDTO) (*TokenResponseDTO, error)
		Register(dto RegisterDTO) (*TokenResponseDTO, error)
		RefreshToken(dto RefreshTokenDTO) (*TokenResponseDTO, error)
		VerifyEmail(token VerifyEmailDTO, userId uint64) error
	}

	authService struct {
		config *config.Config
		logger *logger.ZapLogger

		userService user.UserService
		gmailMailer mailer.GmailMailer

		userRepo interfaces.UserRepository
		authRepo interfaces.RefreshTokenRepository
	}
)

// NewService creates a new auth service
func NewService(
	config *config.Config,
	logger *logger.ZapLogger,
	userService user.UserService,
	gmailMailer mailer.GmailMailer,
	userRepo interfaces.UserRepository,
	authRepo interfaces.RefreshTokenRepository,
) AuthService {
	return &authService{
		config:      config,
		logger:      logger,
		userService: userService,
		gmailMailer: gmailMailer,
		userRepo:    userRepo,
		authRepo:    authRepo,
	}
}

// Login authenticates a user and returns tokens
func (s *authService) Login(dto LoginDTO) (*TokenResponseDTO, error) {
	// Get user by email
	u, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		s.logger.Error("Failed to fetch user by email", zap.String("email", dto.Email), zap.Error(err))
		return nil, err
	}
	if u == nil {
		s.logger.Info("Login attempt with non-existent email", zap.String("email", dto.Email))
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if u.Status != models.USER_STATUS_ACTIVE {
		s.logger.Info("Login attempt with inactive account",
			zap.String("email", dto.Email),
			zap.Uint8("status", u.Status))
		return nil, ErrUserNotActive
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(dto.Password))
	if err != nil {
		s.logger.Info("Failed password verification", zap.String("email", dto.Email))
		return nil, ErrInvalidCredentials
	}

	// Update last login timestamp
	now := time.Now()
	u.LastLoginAt = &now
	err = s.userRepo.Update(u)
	if err != nil {
		s.logger.Error("Failed to update last login time",
			zap.String("email", u.Email),
			zap.Uint64("user_id", u.ID),
			zap.Error(err))
		return nil, err
	}

	// Generate tokens
	tokens, err := s.generateTokens(u)
	if err != nil {
		s.logger.Error("Failed to generate tokens",
			zap.String("email", u.Email),
			zap.Uint64("user_id", u.ID),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("User logged in successfully",
		zap.String("email", u.Email),
		zap.Uint64("user_id", u.ID))
	return tokens, nil
}

// Register creates a new user and returns tokens
func (s *authService) Register(dto RegisterDTO) (*TokenResponseDTO, error) {
	// Convert RegisterDTO to user.CreateUserDTO
	createUserDto := &user.CreateUserDTO{
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
		s.logger.Error("Failed to create user during registration",
			zap.String("email", dto.Email),
			zap.Error(err))
		return nil, err
	}

	// Get complete user
	createdUser, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		s.logger.Error("Failed to fetch created user",
			zap.String("email", dto.Email),
			zap.Error(err))
		return nil, err
	}
	if createdUser == nil {
		s.logger.Error("Created user not found after registration",
			zap.String("email", dto.Email))
		return nil, ErrUserNotFound
	}

	// Generate tokens
	tokens, err := s.generateTokens(createdUser)
	if err != nil {
		s.logger.Error("Failed to generate tokens for new user",
			zap.String("email", createdUser.Email),
			zap.Uint64("user_id", createdUser.ID),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("User registered successfully",
		zap.String("email", createdUser.Email),
		zap.Uint64("user_id", createdUser.ID))
	return tokens, nil
}

// RefreshToken validates a refresh token and issues new tokens
func (s *authService) RefreshToken(dto RefreshTokenDTO) (*TokenResponseDTO, error) {
	// Get refresh token from database
	savedToken, err := s.authRepo.GetRefreshToken(dto.RefreshToken)
	if err != nil {
		s.logger.Warn("Failed to retrieve refresh token",
			zap.String("token", dto.RefreshToken),
			zap.Error(err))
		return nil, ErrInvalidRefreshToken
	}
	if savedToken == nil {
		s.logger.Warn("Refresh token not found", zap.String("token", dto.RefreshToken))
		return nil, ErrInvalidRefreshToken
	}

	// Check if token is expired
	if time.Now().After(savedToken.ExpiresAt) {
		s.logger.Warn("Expired refresh token used",
			zap.String("token", dto.RefreshToken),
			zap.Time("expires_at", savedToken.ExpiresAt))
		// Clean up expired token
		if err := s.authRepo.DeleteRefreshToken(dto.RefreshToken); err != nil {
			s.logger.Error("Failed to delete expired token",
				zap.String("token", dto.RefreshToken),
				zap.Error(err))
		}
		return nil, ErrInvalidRefreshToken
	}

	// Get user
	user, err := s.userRepo.GetByID(savedToken.UserID)
	if err != nil {
		s.logger.Error("Failed to fetch user for refresh token",
			zap.Uint64("user_id", savedToken.UserID),
			zap.Error(err))
		return nil, err
	}
	if user == nil {
		s.logger.Warn("Refresh token used for non-existent user",
			zap.Uint64("user_id", savedToken.UserID))
		return nil, ErrInvalidRefreshToken
	}

	// Check if user is still active
	if user.Status != models.USER_STATUS_ACTIVE {
		s.logger.Warn("Refresh token used for inactive user",
			zap.Uint64("user_id", user.ID),
			zap.Uint8("status", user.Status))
		return nil, ErrUserNotActive
	}

	// Remove used refresh token
	if err = s.authRepo.DeleteRefreshToken(dto.RefreshToken); err != nil {
		s.logger.Error("Failed to delete used refresh token",
			zap.String("token", dto.RefreshToken),
			zap.Error(err))
		// Continue despite error - don't block token refresh
	}

	// Generate new tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("Failed to generate new tokens",
			zap.Uint64("user_id", user.ID),
			zap.Error(err))
		return nil, err
	}

	s.logger.Info("Token refreshed successfully", zap.Uint64("user_id", user.ID))
	return tokens, nil
}

// generateTokens generates JWT access and refresh tokens
func (s *authService) generateTokens(user *models.User) (*TokenResponseDTO, error) {
	// Get JWT config
	jwtSecret := []byte(s.config.JWT.Secret)
	accessTokenExpiry := time.Duration(s.config.JWT.AccessExpiryMinutes) * time.Minute
	refreshTokenExpiry := time.Duration(s.config.JWT.RefreshExpiryDays) * 24 * time.Hour

	// Log token expiry values for debugging
	s.logger.Debug("Token expiry settings",
		zap.Int("access_expiry_minutes", s.config.JWT.AccessExpiryMinutes),
		zap.Int("refresh_expiry_days", s.config.JWT.RefreshExpiryDays),
		zap.Duration("access_duration", accessTokenExpiry),
		zap.Duration("refresh_duration", refreshTokenExpiry))

	// Create access token
	accessToken := jwt.New(jwt.SigningMethodHS256)
	accessClaims := accessToken.Claims.(jwt.MapClaims)
	accessClaims["user_id"] = user.ID
	accessClaims["email"] = user.Email
	accessClaims["exp"] = time.Now().Add(accessTokenExpiry).Unix()

	// Sign access token
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		s.logger.Error("Failed to sign access token", zap.Error(err))
		return nil, err
	}

	// Create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = user.ID
	refreshClaims["email"] = user.Email
	refreshClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	// Sign refresh token
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		s.logger.Error("Failed to sign refresh token", zap.Error(err))
		return nil, err
	}

	// Save refresh token to database
	refreshTokenModel := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	}

	err = s.authRepo.SaveRefreshToken(&refreshTokenModel)
	if err != nil {
		s.logger.Error("Failed to save refresh token to database", zap.Error(err))
		return nil, err
	}

	// Create response
	return &TokenResponseDTO{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    uint(accessTokenExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *authService) VerifyEmail(ved VerifyEmailDTO, userId uint64) error {
	// Get user by ID
	user, err := s.userRepo.GetByID(uint64(userId))

	// Check if there was an error
	if err != nil {
		s.logger.Error("Failed to fetch user by ID", zap.Uint64("user_id", userId), zap.Error(err))
		return err
	}

	// Check if user exists
	if user == nil {
		s.logger.Warn("User not found", zap.Uint64("user_id", userId))
		return ErrUserNotFound
	}

	// Check if user is already verified
	if user.EmailVerified {
		s.logger.Warn("User already verified", zap.Uint64("user_id", userId))
		return nil
	}

	// Check if verification code is correct
	if user.VerifyEmailCode != ved.Code {
		s.logger.Warn("Invalid verification code", zap.Uint64("user_id", userId), zap.Any("code", ved.Code))
		return ErrInvalidVerifyCode
	}

	// Update user
	user.EmailVerified = true
	user.VerifyEmailCode = nil
	err = s.userRepo.Update(user)

	// Check if there was an error
	if err != nil {
		s.logger.Error("Failed to update user", zap.Uint64("user_id", userId), zap.Error(err))
		return err
	}

	// Log success
	s.logger.Info("User email verified", zap.Uint64("user_id", userId))
	return nil
}
