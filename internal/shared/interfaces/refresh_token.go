package interfaces

import "modular-fx-fiber/internal/shared/models"

// Repository defines the data access methods for auth
type RefreshTokenRepository interface {
	SaveRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID int64) error
}
