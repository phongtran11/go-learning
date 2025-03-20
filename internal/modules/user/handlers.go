package user

import (
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// Handlers defines the HTTP handlers for users
type (
	handlers struct {
		service   UserService
		validator *validator.Validator
		logger    *logger.ZapLogger
	}

	Handlers interface {
		Create(c *fiber.Ctx) error
	}
)

// NewHandlers creates a new user handlers instance
func NewHandlers(l *logger.ZapLogger, v *validator.Validator, s UserService) Handlers {
	return &handlers{
		service:   s,
		validator: v,
		logger:    l,
	}
}

// Create handles user creation
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param   user body CreateUserDTO true "User details"
// @Success 201 {object} CreateUserSuccessResponseDTO
// @Router /users [post]
func (h *handlers) Create(c *fiber.Ctx) error {
	createUserDto := &CreateUserDTO{}

	// Parse request body
	if err := c.BodyParser(&createUserDto); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate request body
	errs := h.validator.Validate(createUserDto)

	if errs != nil {
		err := h.validator.ParseErrorToString(errs)
		return fiber.NewError(fiber.StatusBadRequest, err)
	}

	// Create user
	user, err := h.service.CreateUser(createUserDto)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Response with created user
	return c.Status(fiber.StatusCreated).JSON(&CreateUserSuccessResponseDTO{
		Success: true,
		Data:    user,
	})
}
