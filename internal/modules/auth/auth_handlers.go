package auth

import (
	"modular-fx-fiber/internal/shared/dto/auth_dto"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// Handlers defines the HTTP handlers for authentication
type (
	Handlers interface {
		Login(c *fiber.Ctx) error
		Logout(c *fiber.Ctx) error
		Register(c *fiber.Ctx) error
		RefreshToken(c *fiber.Ctx) error
		VerifyEmail(c *fiber.Ctx) error
	}

	handlers struct {
		service   Service
		validator *validator.Validator
		logger    *logger.ZapLogger
	}
)

// NewHandlers creates a new auth handlers instance
func NewHandlers(s Service, l *logger.ZapLogger, v *validator.Validator) Handlers {
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
// @Param credentials body auth_dto.LoginDTO true "Login credentials"
// @Success 200 {object} auth_dto.LoginSuccessResponseDTO
// @Router /auth/login [post]
func (h *handlers) Login(c *fiber.Ctx) error {
	var loginDto auth_dto.LoginDTO

	// Parse request body
	if err := c.BodyParser(&loginDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validate request body
	errs := h.validator.Validate(&loginDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Login user
	tokens, err := h.service.Login(&loginDto)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	// Return response
	return c.JSON(&auth_dto.LoginSuccessResponseDTO{
		Success: true,
		Data:    tokens,
	})
}

// Logout handles user logout
// @Summary User logout
// @Description Invalidate user tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Router /auth/logout [post]
func (h *handlers) Logout(c *fiber.Ctx) error {
	// Get user ID from context
	userId := c.Locals("user_id").(uint64)

	// Logout user
	err := h.service.Logout(&auth_dto.LogoutDTO{
		UserId: userId,
	})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}

// Register handles user registration
// @Summary User registration
// @Description Register a new user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param user body auth_dto.RegisterDTO true "Registration data"
// @Success 201 {object} auth_dto.RegisterSuccessResponseDTO
// @Router /auth/register [post]
func (h *handlers) Register(c *fiber.Ctx) error {
	var registerDto auth_dto.RegisterDTO

	// Parse request body
	if err := c.BodyParser(&registerDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body
	errs := h.validator.Validate(&registerDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Register user
	tokens, err := h.service.Register(&registerDto)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(&auth_dto.RegisterSuccessResponseDTO{
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
// @Param token body auth_dto.RefreshTokenDTO true "Refresh token"
// @Success 201 {object} auth_dto.RefreshTokenSuccessResponseDTO
// @Router /auth/refresh-token [post]
func (h *handlers) RefreshToken(c *fiber.Ctx) error {
	var refreshDto auth_dto.RefreshTokenDTO

	// Parse request body
	if err := c.BodyParser(&refreshDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body
	errs := h.validator.Validate(refreshDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Refresh token
	tokens, err := h.service.RefreshToken(&refreshDto)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(&auth_dto.RefreshTokenSuccessResponseDTO{
		Success: true,
		Data:    tokens,
	})
}

// VerifyEmail handles email verification
// @Summary Verify email
// @Description Verify user email address
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param code body auth_dto.VerifyEmailDTO true "Verification code"
// @Success 200
// @Router /auth/register/verify-email [post]
func (h *handlers) VerifyEmail(c *fiber.Ctx) error {
	var verifyDto auth_dto.VerifyEmailDTO

	// Parse request body
	if err := c.BodyParser(&verifyDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body
	errs := h.validator.Validate(verifyDto)
	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Get user ID from context
	userId := c.Locals("user_id").(uint64)

	// check verification code
	err := h.service.VerifyEmail(&verifyDto, userId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})

}
