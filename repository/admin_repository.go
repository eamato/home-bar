package repository

import (
	"context"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map"
	"home-bar/database"
	"home-bar/domain"
	"home-bar/internal"
)

type adminRepository struct {
	database               database.Database
	usersTableName         string
	rolesTableName         string
	profilesTableName      string
	usersHasRolesTableName string
}

func NewAdminRepository(
	database database.Database,
	usersTableName string,
	rolesTableName string,
	profilesTableName string,
	usersHasRolesTableName string,
) domain.AdminRepository {
	return &adminRepository{
		database:               database,
		usersTableName:         usersTableName,
		rolesTableName:         rolesTableName,
		profilesTableName:      profilesTableName,
		usersHasRolesTableName: usersHasRolesTableName,
	}
}

func (ar *adminRepository) GetUsersFullInfo(_ context.Context, take int, skip int) ([]*domain.UsersListResponseUser, error) {
	var usersListResponseUser domain.UsersListResponseUser

	join := orderedmap.New()
	join.Set(ar.usersHasRolesTableName, fmt.Sprintf(
		"%s.user_id = %s.id", ar.usersHasRolesTableName, ar.usersTableName))
	join.Set(ar.rolesTableName, fmt.Sprintf(
		"%s.role_id = %s.id", ar.usersHasRolesTableName, ar.rolesTableName))
	join.Set(ar.profilesTableName, fmt.Sprintf(
		"%s.user_id = %s.id", ar.profilesTableName, ar.usersTableName))

	res, err := ar.database.Collection(ar.usersTableName).FindMany(
		database.WithFields([]string{
			"users.id",
			"users.username",
			"users.email",
			"users.password",
			"users.refresh_token",
			"roles.role",
			"profiles.id AS profile_id",
			"profiles.user_id",
			"profiles.nickname",
		}),
		database.WithJoin(join),
		database.WithPagination(take, skip),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUsersFullInfo error: %s", err.Error()))
	}

	results, err := res.Decode(&usersListResponseUser)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUsersFullInfo Decode error: %s", err.Error()))
	}

	users := internal.MapArray(results, func(value interface{}) *domain.UsersListResponseUser {
		castResult, ok := value.(*domain.UsersListResponseUser)
		if ok {
			return castResult
		}
		return nil
	})

	return internal.UnNilArray(users), err
}

func (ar *adminRepository) DeleteUser(_ context.Context, userID int64) (bool, error) {
	rowsAffected, err := ar.database.
		Collection(ar.usersTableName).
		DeleteOne(database.WithDeleteWhere(fmt.Sprintf("%s.id = %d", ar.usersTableName, userID)))
	if err != nil {
		return false, err
	}

	return rowsAffected.(int64) == 1, nil
}

func (ar *adminRepository) UpdateUserRole(_ context.Context, userID int64, roleID int64) (int64, error) {
	type Role struct {
		ID     int64 `db:"id"`
		UserID int64 `db:"user_id"`
		RoleID int64 `db:"role_id"`
	}
	var role Role
	res, err := ar.database.Collection(ar.usersHasRolesTableName).FindOne(
		database.WithWhere(fmt.Sprintf("%s.user_id = %d", ar.usersHasRolesTableName, userID)))
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpdateUserRole FindOne error: %s", err.Error()))
	}

	err = res.Decode(&role)
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpdateUserRole Decode error: %s", err.Error()))
	}

	if role.RoleID == roleID {
		return role.ID, nil
	}
	id, err := ar.database.
		Collection(ar.usersHasRolesTableName).
		UpsertOne(
			database.WithFieldsValues(map[string]interface{}{
				"id":      role.ID,
				"user_id": role.UserID,
				"role_id": roleID,
			}))
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpdateUserRole error: %s", err.Error()))
	}

	return id.(int64), nil
}
