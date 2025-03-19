package models

import (
	"time"

	"gorm.io/gorm"
)

// User status enum
const (
	USER_STATUS_ACTIVE   uint8 = 1
	USER_STATUS_INACTIVE uint8 = 2
)

// User gender enum
const (
	GENDER_MALE   uint8 = 1
	GENDER_FEMALE uint8 = 2
)

type User struct {
	ID            uint64         `json:"id" gorm:"type:bigserial;primaryKey;autoIncrement"`
	Email         string         `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PhoneNumber   *string        `json:"phone_number" gorm:"type:varchar(20)"`
	Password      string         `json:"password" gorm:"type:varchar(255);not null"`
	FirstName     string         `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName      string         `json:"last_name" gorm:"type:varchar(100);not null"`
	DateOfBirth   *time.Time     `json:"date_of_birth" gorm:"type:date"`
	Gender        *uint8         `json:"gender" gorm:"type:smallint"` // References GENDER constants
	AvatarURL     *string        `json:"avatar_url" gorm:"type:varchar(512)"`
	EmailVerified bool           `json:"email_verified" gorm:"type:boolean;default:false"`
	Status        uint8          `json:"status" gorm:"type:smallint;default:1"` // Default to USER_STATUS_ACTIVE (1)
	LastLoginAt   *time.Time     `json:"last_login_at" gorm:"type:timestamp with time zone"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp with time zone;index"`

	// One-to-Many relationship with RefreshTokens
	RefreshTokens []RefreshToken `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Roles         []Role         `json:"-" gorm:"many2many:user_roles;"`
}

// Getter method to get full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// UserResponseDTO represents the user data to be returned in API responses
// @Description User information returned in API responses
type UserResponseDTO struct {
	ID            uint64     `json:"id" example:"1"`
	Email         string     `json:"email" example:"user@example.com"`
	PhoneNumber   *string    `json:"phone_number,omitempty" example:"+12125551234"`
	FirstName     string     `json:"first_name" example:"John"`
	LastName      string     `json:"last_name" example:"Doe"`
	FullName      string     `json:"full_name" example:"John Doe"`
	DateOfBirth   *time.Time `json:"date_of_birth,omitempty" example:"1990-01-01"`
	Gender        *uint8     `json:"gender,omitempty" example:"1"`
	AvatarURL     *string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
	EmailVerified bool       `json:"email_verified" example:"true"`
	Status        uint8      `json:"status" example:"1"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty" example:"2023-01-01T12:00:00Z"`
	CreatedAt     time.Time  `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updated_at" example:"2023-01-01T12:34:56Z"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" example:"2023-01-10T00:00:00Z"`
}

func (u *User) ToResponseDTO() *UserResponseDTO {
	return &UserResponseDTO{
		ID:            u.ID,
		Email:         u.Email,
		PhoneNumber:   u.PhoneNumber,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		FullName:      u.FullName(),
		DateOfBirth:   u.DateOfBirth,
		Gender:        u.Gender,
		AvatarURL:     u.AvatarURL,
		EmailVerified: u.EmailVerified,
		Status:        u.Status,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		DeletedAt:     &u.DeletedAt.Time,
	}
}

// HasPermission checks if the user has the specified permission
func (u *User) HasPermission(resourceName, action string) bool {
	for _, role := range u.Roles {
		for _, permission := range role.Permissions {
			if permission.ResourceName == resourceName && permission.Action == action {
				return true
			}
		}
	}
	return false
}

// HasRole checks if the user has the specified role
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
