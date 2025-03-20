package auth

import (
	"errors"
	"modular-fx-fiber/internal/modules/user"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// Handlers defines the HTTP handlers for authentication
type (
	Handlers interface {
		Login(c *fiber.Ctx) error
		Register(c *fiber.Ctx) error
		RefreshToken(c *fiber.Ctx) error
		VerifyEmail(c *fiber.Ctx) error
	}

	handlers struct {
		service   AuthService
		validator *validator.Validator
		logger    *logger.ZapLogger
	}
)

// NewHandlers creates a new auth handlers instance
func NewHandlers(l *logger.ZapLogger, v *validator.Validator, s AuthService) Handlers {
	return &handlers{
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
// @Success 200 {object} LoginSuccessResponseDTO
// @Router /auth/login [post]
func (h *handlers) Login(c *fiber.Ctx) error {
	var loginDto LoginDTO

	// Parse request body
	if err := c.BodyParser(&loginDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	errs := h.validator.Validate(loginDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
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
	return c.JSON(&LoginSuccessResponseDTO{
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
// @Success 201 {object} RegisterSuccessResponseDTO
// @Router /auth/register [post]
func (h *handlers) Register(c *fiber.Ctx) error {
	var registerDto RegisterDTO

	// Parse request body
	if err := c.BodyParser(&registerDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	errs := h.validator.Validate(registerDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Register user
	tokens, err := h.service.Register(registerDto)
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

	// Return response
	return c.Status(fiber.StatusCreated).JSON(&RegisterSuccessResponseDTO{
		Success: true,
		Data:    tokens,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh token
// @Description Refresh access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body RefreshTokenDTO true "Refresh token"
// @Success 201 {object} RefreshTokenSuccessResponseDTO
// @Router /auth/refresh-token [post]
func (h *handlers) RefreshToken(c *fiber.Ctx) error {
	var refreshDto RefreshTokenDTO

	// Parse request body
	if err := c.BodyParser(&refreshDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	errs := h.validator.Validate(refreshDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Refresh token
	tokens, err := h.service.RefreshToken(refreshDto)
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
	return c.Status(fiber.StatusCreated).JSON(&RefreshTokenSuccessResponseDTO{
		Success: true,
		Data:    tokens,
	})
}

func (h *handlers) VerifyEmail(c *fiber.Ctx) error {
	return nil
}
