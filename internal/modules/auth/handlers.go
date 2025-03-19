package auth

import (
	"errors"
	"modular-fx-fiber/internal/modules/user"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// Handlers defines the HTTP handlers for authentication
type Handlers struct {
	service   RefreshTokenService
	validator *validator.Validator
	logger    *logger.ZapLogger
}

// NewHandlers creates a new auth handlers instance
func NewHandlers(l *logger.ZapLogger, v *validator.Validator, s RefreshTokenService) *Handlers {
	return &Handlers{
		service:   s,
		validator: v,
		logger:    l,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate a user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginDTO true "Login credentials"
// @Success 200 {object} DataResponseDTO
// @Router /auth/login [post]
func (h *Handlers) Login(c *fiber.Ctx) error {
	var loginDto LoginDTO

	// Parse request body
	if err := c.BodyParser(&loginDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	err := h.validator.ValidateStruct(loginDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	// Login user
	tokens, err := h.service.Login(loginDto)
	if errors.Is(err, ErrInvalidCredentials) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid email or password",
		})
	} else if errors.Is(err, ErrUserNotActive) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Your account is not active",
		})
	} else if err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	// Return response
	return c.JSON(DataResponseDTO{
		Success: true,
		Data:    tokens,
	})
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterDTO true "Registration data"
// @Success 201 {object} RegisterSuccessDTO
// @Router /auth/register [post]
func (h *Handlers) Register(c *fiber.Ctx) error {
	var registerDto RegisterDTO

	// Parse request body
	if err := c.BodyParser(&registerDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	err := h.validator.ValidateStruct(registerDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	// Register user
	userCreated, err := h.service.Register(registerDto)
	if err != nil {
		if errors.Is(err, user.ErrEmailAlreadyExists) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Email already exists",
			})
		}
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	registerSuccessDTO := &RegisterSuccessDTO{
		Success: true,
		Data:    userCreated,
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(registerSuccessDTO)
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenDTO true "Refresh token"
// @Success 201 {object} DataResponseDTO
// @Router /auth/refresh-token [post]
func (h *Handlers) RefreshToken(c *fiber.Ctx) error {
	var refreshDto RefreshTokenDTO

	// Parse request body
	if err := c.BodyParser(&refreshDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	err := h.validator.ValidateStruct(refreshDto)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	// Refresh token
	response, err := h.service.RefreshToken(refreshDto)
	if errors.Is(err, ErrInvalidRefreshToken) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid or expired refresh token",
		})
	} else if err != nil {
		h.logger.Error(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Internal server error",
		})
	}

	// Return response
	return c.JSON(response)
}
