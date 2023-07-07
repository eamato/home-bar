package usecase

import (
	"context"
	"home-bar/configs"
	"home-bar/database"
	"home-bar/domain"
)

type signupUsecase struct {
	cfg               *configs.Config
	userRepository    domain.UserRepository
	profileRepository domain.ProfileRepository
	roleRepository    domain.RoleRepository
}

func NewSignupUsecase(
	cfg *configs.Config,
	userRepository domain.UserRepository,
	profileRepository domain.ProfileRepository,
	roleRepository domain.RoleRepository) domain.SignupUsecase {

	return &signupUsecase{
		cfg:               cfg,
		userRepository:    userRepository,
		profileRepository: profileRepository,
		roleRepository:    roleRepository,
	}
}

func (su *signupUsecase) Signup(ctx context.Context, user *domain.User) (*domain.SignupResult, *domain.UsecaseError) {
	databaseBackup := database.NewHomeBarDatabaseBackup(su.cfg)
	err := databaseBackup.CreateBackup()
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	signupChain := NewCheckUserExistenceStep(ctx, su.userRepository)
	signupChain.
		SetNext(NewPasswordHashingStep()).
		SetNext(NewSavingNewUserStep(ctx, su.userRepository)).
		SetNext(NewSavingProfileStep(ctx, su.profileRepository)).
		SetNext(NewAssigningRoleStep(ctx, su.roleRepository)).
		SetNext(NewCreatingTokensAndUpdateUserStep(ctx, su.userRepository, su.cfg))

	res := &domain.SignupResult{}
	usecaseError := signupChain.Execute(user, res)
	if usecaseError != nil {
		innerErr := databaseBackup.BackupAndDelete()
		if innerErr != nil {
			return nil, domain.NewUsecaseError(innerErr, domain.ReasonServerError)
		}

		return nil, usecaseError
	}

	err = databaseBackup.DeleteBackup()
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	return res, nil
}
