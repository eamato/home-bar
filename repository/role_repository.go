package repository

import (
	"context"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map"
	"home-bar/database"
	"home-bar/domain"
)

type roleRepository struct {
	database               database.Database
	rolesTableName         string
	usersHasRolesTableName string
}

func NewRoleRepository(
	database database.Database,
	rolesTableName string,
	usersHasRolesTableName string) domain.RoleRepository {

	return &roleRepository{
		database:               database,
		rolesTableName:         rolesTableName,
		usersHasRolesTableName: usersHasRolesTableName,
	}
}

func (rr *roleRepository) GetRole(_ context.Context, userID int64) (*domain.Role, error) {
	type Role struct {
		Role domain.Role `db:"role"`
	}
	var role = Role{
		Role: domain.RoleUser,
	}

	join := orderedmap.New()
	join.Set(rr.usersHasRolesTableName, fmt.Sprintf(
		"%s.id = %s.role_id", rr.rolesTableName, rr.usersHasRolesTableName))

	res, err := rr.database.Collection(rr.rolesTableName).FindOne(
		database.WithFields([]string{
			"role",
		}),
		database.WithJoin(join),
		database.WithWhere(fmt.Sprintf("%s.user_id = %d", rr.usersHasRolesTableName, userID)))
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetRole error: %s", err.Error()))
	}

	err = res.Decode(&role)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetRole Decode error: %s", err.Error()))
	}

	return &role.Role, err
}

func (rr *roleRepository) CreateRole(_ context.Context, userID int64, role domain.Role) error {
	type RoleID struct {
		RoleID int64 `db:"id"`
	}
	var roleID RoleID
	res, err := rr.database.Collection(rr.rolesTableName).FindOne(
		database.WithFields([]string{
			"id",
		}),
		database.WithWhere(fmt.Sprintf("%s.role = '%s'", rr.rolesTableName, role)))
	if err != nil {
		return domain.NewCustomError(fmt.Sprintf("CreateRole error: %s", err.Error()))
	}

	err = res.Decode(&roleID)
	if err != nil {
		return domain.NewCustomError(fmt.Sprintf("CreateRole Decode error: %s", err.Error()))
	}

	if roleID.RoleID == 0 {
		return domain.NewCustomError("CreateRole role ID is 0")
	}

	_, err = rr.database.Collection(rr.usersHasRolesTableName).UpsertOne(
		database.WithFieldsValues(map[string]interface{}{
			"user_id": userID,
			"role_id": roleID.RoleID,
		}))

	if err != nil {
		return domain.NewCustomError(fmt.Sprintf("CreateRole error: %s", err.Error()))
	}

	return nil
}
