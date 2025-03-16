package user

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handlers defines the HTTP handlers for users
type Handlers struct {
	service Service
}

// NewHandlers creates a new user handlers instance
func NewHandlers(service Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

// Create handles user creation
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserDTO true "User information"
// @Success 201 {object} UserResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func (h *Handlers) Create(c *fiber.Ctx) error {
	var dto CreateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate input
	if dto.Email == "" || dto.Password == "" || dto.FirstName == "" || dto.LastName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	user, err := h.service.CreateUser(dto)
	if err != nil {
		if err == ErrEmailAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Email already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetByID handles fetching a user by ID
// @Summary Get a user by ID
// @Description Retrieve a user by their unique ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} UserResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *Handlers) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
	}

	return c.JSON(user)
}

// Update handles updating a user
// @Summary Update a user
// @Description Update a user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateUserDTO true "Updated user information"
// @Success 200 {object} UserResponseDTO
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *Handlers) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var dto UpdateUserDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := h.service.UpdateUser(id, dto)
	if err != nil {
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(user)
}

// Delete handles deleting a user
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *Handlers) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	if err := h.service.DeleteUser(id); err != nil {
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User successfully deleted",
	})
}

// List handles fetching a paginated list of users
// @Summary List users
// @Description Retrieve a paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param page query integer false "Page number (default: 1)"
// @Param pageSize query integer false "Page size (default: 10, max: 100)"
// @Success 200 {object} PaginatedUsersResponse
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *Handlers) List(c *fiber.Ctx) error {
	// Parse pagination parameters
	pageStr := c.Query("page", "1")
	pageSizeStr := c.Query("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, err := h.service.ListUsers(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve users",
		})
	}

	return c.JSON(users)
}

// ChangePassword handles changing a user's password
// @Summary Change user password
// @Description Change a user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param passwordChange body ChangePasswordDTO true "Password change information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id}/change-password [post]
func (h *Handlers) ChangePassword(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var dto ChangePasswordDTO
	if err := c.BodyParser(&dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if dto.CurrentPassword == "" || dto.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields",
		})
	}

	if len(dto.NewPassword) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "New password must be at least 8 characters",
		})
	}

	if err := h.service.ChangePassword(id, dto); err != nil {
		if err == ErrUserNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User not found",
			})
		}
		if err == ErrInvalidCredentials {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Current password is incorrect",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change password",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password successfully changed",
	})
}
