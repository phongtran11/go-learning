package interfaces

// RolePermissionRepository defines the operations for managing role-permission assignments
type RolePermissionRepository interface {
	AssignPermissionsToRole(roleID uint64, permissionIDs []uint64) error
	RemovePermissionsFromRole(roleID uint64, permissionIDs []uint64) error
	GetRolePermissions(roleID uint64) ([]uint64, error)
	GetPermissionRoles(permissionID uint64) ([]uint64, error)
}
