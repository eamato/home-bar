package domain

import "context"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"

	RolesTableName         = "roles"
	UsersHasRolesTableName = "users_has_roles"
)

type RoleRepository interface {
	GetRole(context.Context, int64) (*Role, error)
	CreateRole(context.Context, int64, Role) error
}
