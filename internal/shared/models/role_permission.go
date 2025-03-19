package models

import (
	"time"

	"gorm.io/gorm"
)

// RolePermission defines the many-to-many relationship between roles and permissions
type RolePermission struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleID       uint64         `json:"role_id" gorm:"index;not null"`
	PermissionID uint64         `json:"permission_id" gorm:"index;not null"`
	CreatedAt    time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`

	Role       Role       `json:"role" gorm:"foreignKey:RoleID;references:ID;constraint:OnDelete:CASCADE"`
	Permission Permission `json:"permission" gorm:"foreignKey:PermissionID;references:ID;constraint:OnDelete:CASCADE"`
}
