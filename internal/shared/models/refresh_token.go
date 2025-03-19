package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    int64          `json:"user_id" gorm:"index;not null;OnDelete:CASCADE"`
	Token     string         `json:"token" gorm:"uniqueIndex;size:255;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"type:timestamp with time zone;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`

	// Many-to-One relationship with User
	User *User `json:"-" gorm:"foreignKey:UserID;references:ID"`
}
