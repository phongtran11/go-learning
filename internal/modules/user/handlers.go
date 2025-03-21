package user

import (
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"
	"strconv"

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
		ListUsers(c *fiber.Ctx) error
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

// ListUsers handles listing users
// @Summary List users
// @Description List users with pagination
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} ListUsersSuccessResponseDTO
// @Router /users [get]
func (h *handlers) ListUsers(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageSize := c.Query("page_size", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page")
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid page size")
	}

	// Limit page size to 100
	if pageSizeInt > 100 {
		pageSizeInt = 100
	}

	users, total, err := h.service.ListUsers(pageInt, pageSizeInt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	listUsers := &PaginatedUsersResponse{
		Items:      users,
		TotalCount: total,
		Page:       pageInt,
		PageSize:   pageSizeInt,
		TotalPages: (total + int64(pageSizeInt) - 1) / int64(pageSizeInt),
	}

	return c.JSON(ListUsersSuccessResponseDTO{
		Success: true,
		Data:    listUsers,
	})
}
