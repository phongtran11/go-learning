package interfaces

import "modular-fx-fiber/internal/shared/models"

type RefreshTokenRepository interface {
	SaveRefreshToken(token *models.RefreshToken) error
	GetRefreshToken(token string) (*models.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID uint64) error
}
