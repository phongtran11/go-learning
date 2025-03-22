package interfaces

import (
	"modular-fx-fiber/internal/shared/models"
)

type PermissionRepository interface {
	Create(permission *models.Permission) error
	Update(permission *models.Permission) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Permission, error)
	GetByName(name string) (*models.Permission, error)
	GetByResourceAndAction(resourceName, action string) (*models.Permission, error)
	List(page, pageSize int) ([]models.Permission, int64, error)
	ListByResourceName(resourceName string, page, pageSize int) ([]models.Permission, int64, error)
	ListByAction(action string, page, pageSize int) ([]models.Permission, int64, error)
}
