package interfaces

import (
	"modular-fx-fiber/internal/shared/models"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint64) (*models.User, error)
	List(page int, pageSize int) ([]models.User, int64, error)
	Delete(id uint64) error
}
