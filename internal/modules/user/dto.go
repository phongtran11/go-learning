package user

import (
	"time"

	"github.com/google/uuid"
)

// CreateUserDTO represents the data needed to create a new user
// @Description Data for creating a new user
type CreateUserDTO struct {
	Email       string     `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string     `json:"password" validate:"required,min=8" example:"secureP@ssw0rd"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+12125551234"`
	FirstName   string     `json:"first_name" validate:"required" example:"John"`
	LastName    string     `json:"last_name" validate:"required" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other" example:"male"`
	AvatarURL   *string    `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// UpdateUserDTO represents the data for updating a user
// @Description Data for updating an existing user
type UpdateUserDTO struct {
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+12125551234"`
	FirstName   *string    `json:"first_name,omitempty" example:"John"`
	LastName    *string    `json:"last_name,omitempty" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender      *string    `json:"gender,omitempty" validate:"omitempty,oneof=male female other" example:"male"`
	AvatarURL   *string    `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// ChangePasswordDTO represents the data for changing a user's password
// @Description Data for changing a user's password
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldP@ssw0rd"`
	NewPassword     string `json:"new_password" validate:"required,min=8,nefield=CurrentPassword" example:"newSecureP@ssw0rd"`
}

// UserResponseDTO represents the user data to be returned in API responses
// @Description User information returned in API responses
type UserResponseDTO struct {
	ID            uuid.UUID  `json:"id" example:"a87ff679-4056-4bee-8c3f-37f9e9a48e1e"`
	Email         string     `json:"email" example:"user@example.com"`
	PhoneNumber   *string    `json:"phone_number,omitempty" example:"+12125551234"`
	FirstName     string     `json:"first_name" example:"John"`
	LastName      string     `json:"last_name" example:"Doe"`
	FullName      string     `json:"full_name" example:"John Doe"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender        *string    `json:"gender,omitempty" example:"male"`
	AvatarURL     *string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	EmailVerified bool       `json:"email_verified" example:"true"`
	PhoneVerified bool       `json:"phone_verified" example:"false"`
	Status        string     `json:"status" example:"active"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" example:"2023-01-01T12:00:00Z"`
	CreatedAt     time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updated_at" example:"2023-01-01T12:34:56Z"`
}

// PaginatedUsersResponse represents a paginated list of users
// @Description Paginated list of users
type PaginatedUsersResponse struct {
	Users      []UserResponseDTO `json:"users"`
	TotalCount int64             `json:"total_count" example:"42"`
	Page       int               `json:"page" example:"1"`
	PageSize   int               `json:"page_size" example:"10"`
	TotalPages int               `json:"total_pages" example:"5"`
}

// ToResponseDTO converts a User model to a UserResponseDTO
func (u User) ToResponseDTO() UserResponseDTO {
	return UserResponseDTO{
		ID:            u.ID,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		FullName:      u.FullName(),
		DateOfBirth:   u.DateOfBirth,
		Gender:        u.Gender,
		AvatarURL:     u.AvatarURL,
		EmailVerified: u.EmailVerified,
		PhoneVerified: u.PhoneVerified,
		Status:        u.Status,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
