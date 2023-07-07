package usecase

import (
	"context"
	"home-bar/configs"
	"home-bar/database"
	"home-bar/domain"
)

type loginUsecase struct {
	cfg               *configs.Config
	userRepository    domain.UserRepository
	profileRepository domain.ProfileRepository
	roleRepository    domain.RoleRepository
}

func NewLoginUsecase(
	cfg *configs.Config,
	userRepository domain.UserRepository,
	profileRepository domain.ProfileRepository,
	roleRepository domain.RoleRepository) domain.LoginUsecase {

	return &loginUsecase{
		cfg:               cfg,
		userRepository:    userRepository,
		profileRepository: profileRepository,
		roleRepository:    roleRepository,
	}
}

func (lu *loginUsecase) Login(
	ctx context.Context,
	loginRequest domain.LoginRequest,
) (*domain.LoginResult, *domain.UsecaseError) {
	databaseBackup := database.NewHomeBarDatabaseBackup(lu.cfg)
	err := databaseBackup.CreateBackup()
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	loginChain := NewGetUserStep(ctx, lu.userRepository)
	loginChain.
		SetNext(NewPasswordComparisonStep()).
		SetNext(NewLoginCreatingTokensAndUpdateUserStep(ctx, lu.userRepository, lu.cfg))

	res := &domain.LoginResult{}
	user := &LoginUsersComposition{
		LoginRequest: loginRequest,
	}
	usecaseError := loginChain.Execute(user, res)
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

func (lu *loginUsecase) LoginWithGoogle(ctx context.Context, code string) (*domain.LoginResult, *domain.UsecaseError) {
	databaseBackup := database.NewHomeBarDatabaseBackup(lu.cfg)
	err := databaseBackup.CreateBackup()
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	loginChain := NewGettingGoogleUserStep(ctx, code, lu.cfg)
	loginChain.
		SetNext(NewGetUserStep(ctx, lu.userRepository)).
		SetNext(NewLoginSavingNewUserStep(ctx, lu.userRepository)).
		SetNext(NewLoginSavingProfileStep(ctx, lu.profileRepository)).
		SetNext(NewLoginAssigningRoleStep(ctx, lu.roleRepository)).
		SetNext(NewLoginCreatingTokensAndUpdateUserStep(ctx, lu.userRepository, lu.cfg))

	res := &domain.LoginResult{}
	user := &LoginUsersComposition{
		LoginRequest: domain.LoginRequest{},
	}
	usecaseError := loginChain.Execute(user, res)
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
