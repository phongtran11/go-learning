package interfaces

type UserRoleRepository interface {
	AssignRolesToUser(userID uint64, roleIDs []uint64) error
	RemoveRolesFromUser(userID uint64, roleIDs []uint64) error
	GetUserRoles(userID uint64) ([]uint64, error)
	GetRoleUsers(roleID uint64) ([]uint64, error)
}
