package user

import (
	"errors"
	"modular-fx-fiber/internal/shared/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the data access methods for users
type Repository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List(page, pageSize int) ([]User, int64, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *database.Database) Repository {
	return &repository{db: db.GetDB()}
}

// Create inserts a new user into the database
func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(id uuid.UUID) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(email string) (*User, error) {
	var user User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *repository) Update(user *User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&User{}, "id = ?", id).Error
}

// List retrieves a paginated list of users
func (r *repository) List(page, pageSize int) ([]User, int64, error) {
	var users []User
	var totalCount int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.Model(&User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}
