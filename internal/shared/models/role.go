package models

import (
	"time"

	"gorm.io/gorm"
)

// Role defines a role that can be assigned to users
type Role struct {
	ID          uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string         `json:"name" gorm:"size:100;uniqueIndex;not null"`
	Description string         `json:"description" gorm:"size:255"`
	Permissions []Permission   `json:"permissions" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`
}
