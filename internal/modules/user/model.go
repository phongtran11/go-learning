package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email         string     `json:"email" gorm:"uniqueIndex;not null"`
	PhoneNumber   *string    `json:"phone_number" gorm:"uniqueIndex"`
	PasswordHash  string     `json:"-" gorm:"not null"`
	FirstName     string     `json:"first_name" gorm:"not null"`
	LastName      string     `json:"last_name" gorm:"not null"`
	DateOfBirth   *time.Time `json:"date_of_birth"`
	Gender        *string    `json:"gender"`
	AvatarURL     *string    `json:"avatar_url"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified bool       `json:"phone_verified" gorm:"default:false"`
	Status        string     `json:"status" gorm:"default:inactive"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     *time.Time `json:"deleted_at" gorm:"index"`
}

// Getter method to get full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}
