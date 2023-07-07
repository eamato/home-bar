package usecase

import (
	"context"
	"home-bar/domain"
)

type profileUsecase struct {
	profileRepository domain.ProfileRepository
}

func NewProfileUsecase(profileRepository domain.ProfileRepository) domain.ProfileUsecase {
	return &profileUsecase{
		profileRepository: profileRepository,
	}
}

func (pu *profileUsecase) GetProfile(ctx context.Context, userID int64) (*domain.Profile, *domain.UsecaseError) {
	profile, err := pu.profileRepository.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if profile == nil {
		return nil, domain.NewUsecaseError(
			domain.NewCustomError("Profile not found with the given user id"), domain.ReasonNotFound)
	}

	return profile, nil
}
