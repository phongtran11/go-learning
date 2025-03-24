package auth_dto

import "time"

// LoginDTO represents login credentials
// @Description Login credentials
type LoginDTO struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"secureP@ssw0rd"`
}

// RegisterDTO represents registration data
// @Description Registration data
type RegisterDTO struct {
	Email       string     `json:"email" validate:"required,email" example:"user@example.com"`
	Password    string     `json:"password" validate:"required,password" example:"secureP@ssw0rd"`
	PhoneNumber *string    `json:"phone_number,omitempty" validate:"omitempty,vn_phone" example:"0912345678"`
	FirstName   string     `json:"first_name" validate:"required" example:"John"`
	LastName    string     `json:"last_name" validate:"required" example:"Doe"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" validate:"omitempty,datetime=1990-01-01T00:00:00Z" example:"1990-01-01T00:00:00Z"`
	Gender      *uint8     `json:"gender,omitempty" validate:"omitempty,oneof=1 2" example:"1"`
}

type VerifyEmailDTO struct {
	Code *string `json:"code" validate:"required,min=6,max=6,alphanum" example:"123456"`
}

// RefreshTokenDTO represents refresh token request data
// @Description Refresh token request data
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token"   validate:"required"    example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type LogoutDTO struct {
	UserId uint64
}
