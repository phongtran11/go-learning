package user_dto

import "time"

// CreateUserDTO represents the data needed to create a new user
// @Description Data for creating a new user
type CreateUserDTO struct {
	Email       string     `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string     `json:"password" validate:"required,min=8" example:"secureP@ssw0rd"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+12125551234"`
	FirstName   string     `json:"first_name" validate:"required" example:"John"`
	LastName    string     `json:"last_name" validate:"required" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender      *uint8     `json:"gender,omitempty" validate:"omitempty,oneof=1 2" example:"1"`
	AvatarURL   *string    `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// UpdateUserDTO represents the data for updating a user
// @Description Data for updating an existing user
type UpdateUserDTO struct {
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,e164" example:"+12125551234"`
	FirstName   *string    `json:"first_name,omitempty" example:"John"`
	LastName    *string    `json:"last_name,omitempty" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01"`
	Gender      *uint8     `json:"gender,omitempty" validate:"omitempty,oneof=1 2" example:"1"`
	AvatarURL   *string    `json:"avatar_url,omitempty" validate:"omitempty,url" example:"https://example.com/avatar.jpg"`
}

// ChangePasswordDTO represents the data for changing a user's password
// @Description Data for changing a user's password
type ChangePasswordDTO struct {
	CurrentPassword string `json:"current_password" validate:"required" example:"oldP@ssw0rd"`
	NewPassword     string `json:"new_password" validate:"required,min=8,nefield=CurrentPassword" example:"newSecureP@ssw0rd"`
}
