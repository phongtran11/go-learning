package user

import (
	"errors"
	"modular-fx-fiber/internal/shared/database"
	"modular-fx-fiber/internal/shared/interfaces"
	"modular-fx-fiber/internal/shared/models"

	"gorm.io/gorm"
)

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository
func NewRepository(db *database.Database) interfaces.UserRepository {
	return &repository{db: db.GetDB()}
}

// Create inserts a new user into the database
func (r *repository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// List retrieves a paginated list of users
func (r *repository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var totalCount int64

	offset := (page - 1) * pageSize

	// Get total count
	if err := r.db.Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if err := r.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	return &user, nil
}

// Update updates an existing user
func (r *repository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft-deletes a user
func (r *repository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
