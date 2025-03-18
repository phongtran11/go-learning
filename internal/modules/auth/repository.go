package auth

import (
	"modular-fx-fiber/internal/shared/database"

	"gorm.io/gorm"
)

// Repository defines the data access methods for auth
type Repository interface {
	SaveRefreshToken(token *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteUserRefreshTokens(userID uint64) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository
func NewRepository(db *database.Database) Repository {
	return &repository{db: db.DB}
}

// SaveRefreshToken saves a refresh token to the database
func (r *repository) SaveRefreshToken(token *RefreshToken) error {
	return r.db.Create(token).Error
}

// GetRefreshToken retrieves a refresh token by its value
func (r *repository) GetRefreshToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	if err := r.db.Where("token = ? AND expires_at > NOW()", token).First(&refreshToken).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

// DeleteRefreshToken deletes a refresh token
func (r *repository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&RefreshToken{}).Error
}

// DeleteUserRefreshTokens deletes all refresh tokens for a user
func (r *repository) DeleteUserRefreshTokens(userID uint64) error {
	return r.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}
