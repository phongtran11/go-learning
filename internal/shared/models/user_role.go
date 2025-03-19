package models

import (
	"time"

	"gorm.io/gorm"
)

// UserRole defines the many-to-many relationship between users and roles
type UserRole struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint64         `json:"user_id" gorm:"index;not null"`
	RoleID    uint64         `json:"role_id" gorm:"index;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`

	User User `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Role Role `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
}
