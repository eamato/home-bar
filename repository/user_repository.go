package repository

import (
	"context"
	"fmt"
	"home-bar/database"
	"home-bar/domain"
)

type userRepository struct {
	database  database.Database
	tableName string
}

func NewUserRepository(database database.Database, tableName string) domain.UserRepository {
	return &userRepository{
		database:  database,
		tableName: tableName,
	}
}

func (ur *userRepository) CreateUser(_ context.Context, user *domain.User) (int64, error) {
	id, err := ur.database.Collection(ur.tableName).UpsertOne(
		database.WithFieldsValues(map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		}))

	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("CreateUser error: %s", err.Error()))
	}

	return id.(int64), nil
}

func (ur *userRepository) GetUserByName(_ context.Context, username string) (*domain.User, error) {
	var user domain.User
	res, err := ur.database.Collection(ur.tableName).FindOne(
		database.WithWhere(fmt.Sprintf("username = '%s'", username)),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByName error: %s", err.Error()))
	}

	err = res.Decode(&user)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByName Decode error: %s", err.Error()))
	}

	if user.ID <= 0 {
		return nil, err
	}

	return &user, err
}

func (ur *userRepository) UpdateUser(_ context.Context, user *domain.User) (int64, error) {
	id, err := ur.database.Collection(ur.tableName).UpsertOne(
		database.WithFieldsValues(map[string]interface{}{
			"id":            user.ID,
			"username":      user.Username,
			"email":         user.Email,
			"password":      user.Password,
			"refresh_token": user.RefreshToken,
		}))
	if err != nil {
		return 0, domain.NewCustomError(fmt.Sprintf("UpdateUser error: %s", err.Error()))
	}

	return id.(int64), nil
}

func (ur *userRepository) GetUserById(_ context.Context, id int64) (*domain.User, error) {
	var user domain.User
	res, err := ur.database.Collection(ur.tableName).FindOne(
		database.WithWhere(fmt.Sprintf("id = %d", id)),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserById error: %s", err.Error()))
	}

	err = res.Decode(&user)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserById Decode error: %s", err.Error()))
	}

	if user.ID <= 0 {
		return nil, err
	}

	return &user, err
}

func (ur *userRepository) GetUserByEmail(_ context.Context, email string) (*domain.User, error) {
	var user domain.User
	res, err := ur.database.Collection(ur.tableName).FindOne(
		database.WithWhere(fmt.Sprintf("email = '%s'", email)),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByEmail error: %s", err.Error()))
	}

	err = res.Decode(&user)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByEmail Decode error: %s", err.Error()))
	}

	if user.ID <= 0 {
		return nil, err
	}

	return &user, err
}

func (ur *userRepository) GetUserByUsernameOrEmail(
	_ context.Context, username string, email string) (*domain.User, error) {

	var user domain.User
	res, err := ur.database.Collection(ur.tableName).FindOne(
		database.WithWhere(fmt.Sprintf("username = '%s' OR email = '%s'", username, email)),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByUsernameOrEmail error: %s", err.Error()))
	}

	err = res.Decode(&user)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetUserByUsernameOrEmail Decode error: %s", err.Error()))
	}

	if user.ID <= 0 {
		return nil, err
	}

	return &user, err
}
