package user

import (
	"modular-fx-fiber/internal/shared/dto/user_dto"
	"modular-fx-fiber/internal/shared/logger"
	"modular-fx-fiber/internal/shared/validator"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type (
	// Handlers defines the HTTP handlers for user management
	Handlers interface {
		Create(c *fiber.Ctx) error
		ListUsers(c *fiber.Ctx) error
		GetMe(c *fiber.Ctx) error
	}

	handlers struct {
		service   Service
		validator *validator.Validator
		logger    *logger.ZapLogger
	}
)

// NewHandlers creates a new user handlers instance
func NewHandlers(l *logger.ZapLogger, v *validator.Validator, s Service) Handlers {
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
// @Param   user body user_dto.CreateUserDTO true "User details"
// @Success 201 {object} user_dto.CreateUserSuccessResponseDTO
// @Router /users [post]
func (h *handlers) Create(c *fiber.Ctx) error {
	createUserDto := &user_dto.CreateUserDTO{}

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
	return c.Status(fiber.StatusCreated).JSON(&user_dto.CreateUserSuccessResponseDTO{
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
	listUsers := &user_dto.PaginatedUsersResponse{
		Items:      users,
		TotalCount: total,
		Page:       pageInt,
		PageSize:   pageSizeInt,
		TotalPages: (total + int64(pageSizeInt) - 1) / int64(pageSizeInt),
	}

	return c.JSON(&user_dto.ListUsersSuccessResponseDTO{
		Success: true,
		Data:    listUsers,
	})
}

// GetMe handles getting the current user
// @Summary Get current user
// @Description Get the current user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} GetMeSuccessResponseDTO
// @Router /users/me [get]
func (h *handlers) GetMe(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint64)

	user, err := h.service.GetMe(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(&user_dto.GetMeSuccessResponseDTO{
		Success: true,
		Data:    user,
	})
}
