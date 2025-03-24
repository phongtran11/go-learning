package repositories

import (
	"errors"
	"gorm.io/gorm"
	"modular-fx-fiber/internal/shared/database"
	"modular-fx-fiber/internal/shared/models"
)

type (
	RefreshTokenRepository interface {
		SaveRefreshToken(token *models.RefreshToken) error
		GetRefreshToken(token string) (*models.RefreshToken, error)
		DeleteRefreshToken(token string) error
		DeleteUserRefreshTokens(userID uint64) error
	}

	// refreshTokenRepository implements the Repository interface
	refreshTokenRepo struct {
		db *gorm.DB
	}
)

// NewRefreshTokenRepository creates a new auth refreshTokenRepo
func NewRefreshTokenRepository(db database.Database) RefreshTokenRepository {
	return &refreshTokenRepo{db: db.GetDB()}
}

// SaveRefreshToken saves a refresh token to the database
func (r *refreshTokenRepo) SaveRefreshToken(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

// GetRefreshToken retrieves a refresh token by its value
func (r *refreshTokenRepo) GetRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where("token = ? AND expires_at > NOW()", token).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *refreshTokenRepo) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

// DeleteUserRefreshTokens deletes all refresh tokens for a user
func (r *refreshTokenRepo) DeleteUserRefreshTokens(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}
