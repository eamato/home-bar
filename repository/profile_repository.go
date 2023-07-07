package repository

import (
	"context"
	"fmt"
	"home-bar/database"
	"home-bar/domain"
)

type profileRepository struct {
	database  database.Database
	tableName string
}

func NewProfileRepository(database database.Database, tableName string) domain.ProfileRepository {
	return &profileRepository{
		database:  database,
		tableName: tableName,
	}
}

func (pr *profileRepository) CreateProfile(_ context.Context, profile *domain.Profile) error {
	_, err := pr.database.Collection(pr.tableName).UpsertOne(
		database.WithFieldsValues(map[string]interface{}{
			"user_id":  profile.UserID,
			"nickname": profile.Nickname,
		}))

	if err != nil {
		return domain.NewCustomError(fmt.Sprintf("CreateProfile error: %s", err.Error()))
	}

	return nil
}

func (pr *profileRepository) GetProfileByUserID(_ context.Context, userID int64) (*domain.Profile, error) {
	var profile domain.Profile
	res, err := pr.database.Collection(pr.tableName).FindOne(
		database.WithWhere(fmt.Sprintf("user_id = %d", userID)),
	)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetProfileByUserID error: %s", err.Error()))
	}

	err = res.Decode(&profile)
	if err != nil {
		return nil, domain.NewCustomError(fmt.Sprintf("GetProfileByUserID Decode error: %s", err.Error()))
	}

	return &profile, err
}
