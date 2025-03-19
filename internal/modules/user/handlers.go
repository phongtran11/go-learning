package user

import (
	"errors"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"

	"github.com/gofiber/fiber/v2"
)

// Handlers defines the HTTP handlers for users
type Handlers struct {
	service   UserService
	validator *validator.Validator
	logger    *logger.ZapLogger
}

// NewHandlers creates a new user handlers instance
func NewHandlers(l *logger.ZapLogger, v *validator.Validator, s UserService) *Handlers {
	return &Handlers{
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
// @Success 201 {object} DataResponseDTO
// @Router /users [post]
func (h *Handlers) Create(c *fiber.Ctx) error {
	createUserDto := &CreateUserDTO{}

	// Parse request body
	if err := c.BodyParser(&createUserDto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	// Validate request body
	err := h.validator.ValidateStruct(*createUserDto)
	if err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err,
		})
	}

	// Create user
	user, err := h.service.CreateUser(createUserDto)
	if errors.Is(err, ErrEmailAlreadyExists) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Server Internal Error",
		})
	}

	// Response with created user
	return c.Status(fiber.StatusCreated).JSON(DataResponseDTO{
		Success: true,
		Data:    user,
	})
}
