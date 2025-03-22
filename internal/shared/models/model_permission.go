package models

import (
	"time"

	"gorm.io/gorm"
)

// Permission defines an action that can be performed on a resource
type Permission struct {
	ID           uint64         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string         `json:"name" gorm:"size:100;uniqueIndex;not null"`
	Description  string         `json:"description" gorm:"size:255"`
	ResourceName string         `json:"resource_name" gorm:"size:100;index;not null"`
	Action       string         `json:"action" gorm:"size:50;not null"` // create, read, update, delete, etc.
	CreatedAt    time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`
}
