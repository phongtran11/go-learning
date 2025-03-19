package auth

import (
	"time"
)

// LoginDTO represents login credentials
// @Description Login credentials
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"secureP@ssw0rd"`
}

// TokenResponseDTO represents token response data
// @Description Token response data
type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"    example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token"   example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    uint   `json:"expires_in"      example:"3600"` // in seconds
	TokenType    string `json:"token_type"      example:"Bearer"`
}

// RegisterDTO represents registration data
// @Description Registration data
type RegisterDTO struct {
	Email       string     `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string     `json:"password" validate:"required,min=8" example:"secureP@ssw0rd"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,vn_phone" example:"0975234412"`
	FirstName   string     `json:"first_name" validate:"required" example:"John"`
	LastName    string     `json:"last_name" validate:"required" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01T00:00:00Z"`
	Gender      *int8      `json:"gender,omitempty" validate:"omitempty,oneof=1 2" example:"1"`
}

// RegisterSuccessDTO represents registration success response data
// @Description Registration success response data
type RegisterSuccessDTO struct {
	Success bool              `json:"success"`
	Data    *TokenResponseDTO `json:"data"`
}

// RefreshTokenDTO represents refresh token request data
// @Description Refresh token request data
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"   validate:"required"    example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// DataResponseDTO represents a generic response with auth data
// @Description Generic response with auth data
type DataResponseDTO struct {
	Success bool              `json:"success"`
	Data    *TokenResponseDTO `json:"data"`
}
