package user

import (
	"time"

	"gorm.io/gorm"
)

// User status enum
const (
	USER_STATUS_ACTIVE   int8 = 1
	USER_STATUS_INACTIVE int8 = 2
)

// User gender enum
const (
	GENDER_MALE   = 1
	GENDER_FEMALE = 2
)

type User struct {
	ID            int64          `json:"id" gorm:"type:bigserial;primaryKey;autoIncrement"`
	Email         string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber   *string        `json:"phone_number" gorm:"type:varchar(20)"`
	Password      string         `json:"password" gorm:"type:varchar(255);not null"`
	FirstName     string         `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName      string         `json:"last_name" gorm:"type:varchar(100);not null"`
	DateOfBirth   *time.Time     `json:"date_of_birth" gorm:"type:date"`
	Gender        *int8          `json:"gender" gorm:"type:smallint"` // References GENDER constants
	AvatarURL     *string        `json:"avatar_url" gorm:"type:varchar(512)"`
	EmailVerified bool           `json:"email_verified" gorm:"type:boolean;default:false"`
	Status        int8           `json:"status" gorm:"type:smallint;default:1"` // Default to USER_STATUS_ACTIVE (1)
	LastLoginAt   *time.Time     `json:"last_login_at" gorm:"type:timestamp with time zone"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`
}

// Getter method to get full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
