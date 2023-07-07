package usecase

import (
	"context"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
	"io"
	"net/http"
)

type LoginUsersComposition struct {
	ExistingUser *domain.User
	LoginRequest domain.LoginRequest
	NewUser      *domain.User
}

type LoginStep interface {
	Execute(*LoginUsersComposition, *domain.LoginResult) *domain.UsecaseError
	SetNext(LoginStep) LoginStep
}

type getUserStep struct {
	next           LoginStep
	ctx            context.Context
	userRepository domain.UserRepository
}

func NewGetUserStep(ctx context.Context, userRepository domain.UserRepository) LoginStep {
	return &getUserStep{
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (gus *getUserStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	existingUser, err := gus.userRepository.GetUserByUsernameOrEmail(
		gus.ctx, user.LoginRequest.Username, user.LoginRequest.Email)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	user.ExistingUser = existingUser

	if gus.next != nil {
		return gus.next.Execute(user, res)
	}

	return nil
}

func (gus *getUserStep) SetNext(step LoginStep) LoginStep {
	gus.next = step
	return gus.next
}

type passwordComparisonStep struct {
	next LoginStep
}

func NewPasswordComparisonStep() LoginStep {
	return &passwordComparisonStep{}
}

func (pcs *passwordComparisonStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if user.ExistingUser == nil {
		return domain.NewUsecaseError(
			domain.NewCustomError("User not found with the given username or email"), domain.ReasonNotFound)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.ExistingUser.Password), []byte(user.LoginRequest.Password))
	if err != nil {
		return domain.NewUsecaseError(
			domain.NewCustomError("Invalid credentials"), domain.ReasonUnauthorized)
	}

	user.NewUser = user.ExistingUser

	if pcs.next != nil {
		pcs.next.Execute(user, res)
	}

	return nil
}

func (pcs *passwordComparisonStep) SetNext(step LoginStep) LoginStep {
	pcs.next = step
	return pcs.next
}

type loginCreatingTokensAndUpdateUserStep struct {
	next           LoginStep
	ctx            context.Context
	userRepository domain.UserRepository
	cfg            *configs.Config
}

func NewLoginCreatingTokensAndUpdateUserStep(
	ctx context.Context, userRepository domain.UserRepository, cfg *configs.Config) LoginStep {

	return &loginCreatingTokensAndUpdateUserStep{
		ctx:            ctx,
		userRepository: userRepository,
		cfg:            cfg,
	}
}

func (lctus *loginCreatingTokensAndUpdateUserStep) Execute(
	user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if user.NewUser == nil {
		return domain.NewUsecaseError(domain.NewCustomError("New user is nil"), domain.ReasonServerError)
	}

	accessToken, err := internal.CreateAccessToken(
		user.NewUser, lctus.cfg.TokenConfig.AccessTokenSecret, lctus.cfg.TokenConfig.AccessTokenExpiryHour)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	refreshToken, err := internal.CreateRefreshToken(
		user.NewUser, lctus.cfg.TokenConfig.RefreshTokenSecret, lctus.cfg.TokenConfig.RefreshTokenExpiryHour)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	user.NewUser.RefreshToken = refreshToken

	_, err = lctus.userRepository.UpdateUser(lctus.ctx, user.NewUser)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	res.AccessToken = accessToken
	res.RefreshToken = refreshToken

	if lctus.next != nil {
		return lctus.next.Execute(user, res)
	}

	return nil
}

func (lctus *loginCreatingTokensAndUpdateUserStep) SetNext(step LoginStep) LoginStep {
	lctus.next = step
	return lctus.next
}

type gettingGoogleUserStep struct {
	next LoginStep
	ctx  context.Context
	code string
	cfg  *configs.Config
}

func NewGettingGoogleUserStep(ctx context.Context, code string, cfg *configs.Config) LoginStep {
	return &gettingGoogleUserStep{
		ctx:  ctx,
		code: code,
		cfg:  cfg,
	}
}

func (ggus *gettingGoogleUserStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if ggus.code == "" {
		return domain.NewUsecaseError(domain.NewCustomError("Query code error"), domain.ReasonServerError)
	}

	token, err := ggus.cfg.OAuthConfig.Exchange(ggus.ctx, ggus.code)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	err = json.Unmarshal(contents, &user.LoginRequest)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	if user.LoginRequest.Email == "" {
		return domain.NewUsecaseError(
			domain.NewCustomError("No valid email in google client"), domain.ReasonServerError)
	}

	if ggus.next != nil {
		return ggus.next.Execute(user, res)
	}

	return nil
}

func (ggus *gettingGoogleUserStep) SetNext(step LoginStep) LoginStep {
	ggus.next = step
	return ggus.next
}

type loginSavingNewUserStep struct {
	next           LoginStep
	ctx            context.Context
	userRepository domain.UserRepository
}

func NewLoginSavingNewUserStep(ctx context.Context, userRepository domain.UserRepository) LoginStep {
	return &loginSavingNewUserStep{
		ctx:            ctx,
		userRepository: userRepository,
	}
}

func (lsus *loginSavingNewUserStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if user.ExistingUser != nil {
		if user.NewUser == nil {
			user.NewUser = user.ExistingUser
		}

		if lsus.next != nil {
			return lsus.next.Execute(user, res)
		} else {
			return nil
		}
	}

	user.NewUser = &domain.User{
		Username: user.LoginRequest.Email,
		Email:    user.LoginRequest.Email,
	}

	id, err := lsus.userRepository.CreateUser(lsus.ctx, user.NewUser)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if id == 0 {
		return domain.NewUsecaseError(
			domain.NewCustomError("Login CreateUser returned 0"), domain.ReasonDBError)
	}

	user.NewUser.ID = id

	if lsus.next != nil {
		return lsus.next.Execute(user, res)
	}

	return nil
}

func (lsus *loginSavingNewUserStep) SetNext(step LoginStep) LoginStep {
	lsus.next = step
	return lsus.next
}

type loginSavingProfileStep struct {
	next              LoginStep
	ctx               context.Context
	profileRepository domain.ProfileRepository
}

func NewLoginSavingProfileStep(ctx context.Context, profileRepository domain.ProfileRepository) LoginStep {
	return &loginSavingProfileStep{
		ctx:               ctx,
		profileRepository: profileRepository,
	}
}

func (lsps *loginSavingProfileStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if user.ExistingUser != nil {
		if user.NewUser == nil {
			user.NewUser = user.ExistingUser
		}

		if lsps.next != nil {
			return lsps.next.Execute(user, res)
		} else {
			return nil
		}
	}

	if user.NewUser == nil {
		return domain.NewUsecaseError(domain.NewCustomError("New user is nil"), domain.ReasonServerError)
	}

	profile := domain.Profile{
		UserID:   user.NewUser.ID,
		Nickname: user.NewUser.Username,
	}

	err := lsps.profileRepository.CreateProfile(lsps.ctx, &profile)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if lsps.next != nil {
		return lsps.next.Execute(user, res)
	}

	return nil
}

func (lsps *loginSavingProfileStep) SetNext(step LoginStep) LoginStep {
	lsps.next = step
	return lsps.next
}

type loginAssigningRoleStep struct {
	next           LoginStep
	ctx            context.Context
	roleRepository domain.RoleRepository
}

func NewLoginAssigningRoleStep(ctx context.Context, roleRepository domain.RoleRepository) LoginStep {
	return &loginAssigningRoleStep{
		ctx:            ctx,
		roleRepository: roleRepository,
	}
}

func (lars *loginAssigningRoleStep) Execute(user *LoginUsersComposition, res *domain.LoginResult) *domain.UsecaseError {
	if user.ExistingUser != nil {
		if user.NewUser == nil {
			user.NewUser = user.ExistingUser
		}

		if lars.next != nil {
			return lars.next.Execute(user, res)
		} else {
			return nil
		}
	}

	if user.NewUser == nil {
		return domain.NewUsecaseError(domain.NewCustomError("New user is nil"), domain.ReasonServerError)
	}

	err := lars.roleRepository.CreateRole(lars.ctx, user.NewUser.ID, domain.RoleUser)
	if err != nil {
		return domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	if lars.next != nil {
		return lars.next.Execute(user, res)
	}

	return nil
}

func (lars *loginAssigningRoleStep) SetNext(step LoginStep) LoginStep {
	lars.next = step
	return lars.next
}
