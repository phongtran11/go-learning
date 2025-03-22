package interfaces

import (
	"modular-fx-fiber/internal/shared/models"
)

type RoleRepository interface {
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uint64) error
	GetByID(id uint64) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	List(page int, pageSize int) ([]models.Role, int64, error)
	AssignPermissions(roleID uint64, permissionIDs []uint64) error
	RemovePermissions(roleID uint64, permissionIDs []uint64) error
}
