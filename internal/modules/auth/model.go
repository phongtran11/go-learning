package auth

import (
	"time"
)

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"index"`
	Token     string `gorm:"uniqueIndex;size:255"`
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
