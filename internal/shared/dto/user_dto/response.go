package user_dto

import "modular-fx-fiber/internal/shared/models"

// PaginatedUsersResponse represents a paginated list of users
// @Description Paginated list of users
type PaginatedUsersResponse struct {
	Items      []*models.UserResponseDTO `json:"items"`
	TotalCount int64                     `json:"total_count" example:"42"`
	Page       int                       `json:"page" example:"1"`
	PageSize   int                       `json:"page_size" example:"10"`
	TotalPages int64                     `json:"total_pages" example:"5"`
}

// CreateUserSuccessResponseDTO represents a successful user creation response
// @Description Response structure for successful user creation requests
type CreateUserSuccessResponseDTO struct {
	Success bool                    `json:"success"`
	Data    *models.UserResponseDTO `json:"data"`
}

type ListUsersSuccessResponseDTO struct {
	Success bool                    `json:"success"`
	Data    *PaginatedUsersResponse `json:"data"`
}

// GetMeSuccessResponseDTO represents a successful get me response
// @Description Response structure for successful get me requests
type GetMeSuccessResponseDTO struct {
	Success bool                    `json:"success"`
	Data    *models.UserResponseDTO `json:"data"`
}
