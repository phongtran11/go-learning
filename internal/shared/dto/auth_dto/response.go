package auth_dto

// TokenResponseDTO represents token response data
// @Description Token response data
type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"    example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token"   example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    uint   `json:"expires_in"      example:"3600"` // in seconds
	TokenType    string `json:"token_type"      example:"Bearer"`
}

// LoginSuccessResponseDTO represents a successful login response
// @Description Response structure for successful login requests
type LoginSuccessResponseDTO struct {
	Success bool              `json:"success"`
	Data    *TokenResponseDTO `json:"data"`
}

// RegisterSuccessResponseDTO represents a successful registration response
// @Description Response structure for successful registration requests
type RegisterSuccessResponseDTO struct {
	Success bool              `json:"success"`
	Data    *TokenResponseDTO `json:"data"`
}

// RefreshTokenSuccessResponseDTO represents a successful token refresh response
// @Description Response structure for successful token refresh requests
type RefreshTokenSuccessResponseDTO struct {
	Success bool              `json:"success"`
	Data    *TokenResponseDTO `json:"data"`
}

// VerifySuccessResponseDTO represents a successful email verification response
// @Description Response structure for successful email verification requests
type VerifySuccessResponseDTO struct {
	Success bool `json:"success"`
}

// LogoutSuccessResponseDTO represents a successful logout response
// @Description Response structure for successful logout requests
type LogoutSuccessResponseDTO struct {
	Success bool `json:"success"`
}
