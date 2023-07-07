package usecase

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
)

type SignupStep interface {
	Execute(*domain.User, *domain.SignupResult) *domain.UsecaseError
	SetNext(SignupStep) SignupStep
}

type checkUserExistenceStep struct {
	next           SignupStep
	ctx            context.Context
	userRepository domain.UserRepository
}

func NewCheckUserExistenceStep(ctx context.Context, userRepository domain.UserRepository) SignupStep {
	return &checkUserExistenceStep{
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (cues *checkUserExistenceStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	existingUser, err := cues.userRepository.GetUserByUsernameOrEmail(cues.ctx, user.Username, user.Email)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if existingUser != nil && existingUser.ID != 0 {
		return domain.NewUsecaseError(
			domain.NewCustomError("User already exists with the given username"), domain.ReasonUserDuplicate)
	}

	if cues.next != nil {
		return cues.next.Execute(user, res)
	}

	return nil
}

func (cues *checkUserExistenceStep) SetNext(step SignupStep) SignupStep {
	cues.next = step
	return cues.next
}

type passwordHashingStep struct {
	next SignupStep
}

func NewPasswordHashingStep() SignupStep {
	return &passwordHashingStep{}
}

func (pht *passwordHashingStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	passwordHash := string(encryptedPassword)
	if passwordHash == "" {
		return domain.NewUsecaseError(
			domain.NewCustomError("Password hash in empty"), domain.ReasonServerError)
	}

	user.Password = passwordHash

	if pht.next != nil {
		return pht.next.Execute(user, res)
	}

	return nil
}

func (pht *passwordHashingStep) SetNext(step SignupStep) SignupStep {
	pht.next = step
	return pht.next
}

type savingNewUserStep struct {
	next           SignupStep
	ctx            context.Context
	userRepository domain.UserRepository
}

func NewSavingNewUserStep(ctx context.Context, userRepository domain.UserRepository) SignupStep {
	return &savingNewUserStep{
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (sus *savingNewUserStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	id, err := sus.userRepository.CreateUser(sus.ctx, user)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if id == 0 {
		return domain.NewUsecaseError(
			domain.NewCustomError("Signup CreateUser returned 0"), domain.ReasonDBError)
	}

	user.ID = id

	if sus.next != nil {
		return sus.next.Execute(user, res)
	}

	return nil
}

func (sus *savingNewUserStep) SetNext(step SignupStep) SignupStep {
	sus.next = step
	return sus.next
}

type savingProfileStep struct {
	next              SignupStep
	ctx               context.Context
	profileRepository domain.ProfileRepository
}

func NewSavingProfileStep(ctx context.Context, profileRepository domain.ProfileRepository) SignupStep {
	return &savingProfileStep{
		ctx:               ctx,
		profileRepository: profileRepository,
	}
}

func (sps *savingProfileStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	profile := domain.Profile{
		UserID:   user.ID,
		Nickname: user.Username,
	}

	err := sps.profileRepository.CreateProfile(sps.ctx, &profile)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if sps.next != nil {
		return sps.next.Execute(user, res)
	}

	return nil
}

func (sps *savingProfileStep) SetNext(step SignupStep) SignupStep {
	sps.next = step
	return sps.next
}

type assigningRoleStep struct {
	next           SignupStep
	ctx            context.Context
	roleRepository domain.RoleRepository
}

func NewAssigningRoleStep(ctx context.Context, roleRepository domain.RoleRepository) SignupStep {
	return &assigningRoleStep{
		ctx:            ctx,
		roleRepository: roleRepository,
	}
}

func (ars *assigningRoleStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	err := ars.roleRepository.CreateRole(ars.ctx, user.ID, domain.RoleUser)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if ars.next != nil {
		return ars.next.Execute(user, res)
	}

	return nil
}

func (ars *assigningRoleStep) SetNext(step SignupStep) SignupStep {
	ars.next = step
	return ars.next
}

type creatingTokensAndUpdateUserStep struct {
	next           SignupStep
	ctx            context.Context
	userRepository domain.UserRepository
	cfg            *configs.Config
}

func NewCreatingTokensAndUpdateUserStep(
	ctx context.Context, userRepository domain.UserRepository, cfg *configs.Config) SignupStep {

	return &creatingTokensAndUpdateUserStep{
		ctx:            ctx,
		userRepository: userRepository,
		cfg:            cfg,
	}
}

func (ctus *creatingTokensAndUpdateUserStep) Execute(user *domain.User, res *domain.SignupResult) *domain.UsecaseError {
	accessToken, err := internal.CreateAccessToken(
		user, ctus.cfg.TokenConfig.AccessTokenSecret, ctus.cfg.TokenConfig.AccessTokenExpiryHour)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	refreshToken, err := internal.CreateRefreshToken(
		user, ctus.cfg.TokenConfig.RefreshTokenSecret, ctus.cfg.TokenConfig.RefreshTokenExpiryHour)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	user.RefreshToken = refreshToken

	_, err = ctus.userRepository.UpdateUser(ctus.ctx, user)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	if ctus.next != nil {
		return ctus.next.Execute(user, res)
	}

	return nil
}

func (ctus *creatingTokensAndUpdateUserStep) SetNext(step SignupStep) SignupStep {
	ctus.next = step
	return ctus.next
}
