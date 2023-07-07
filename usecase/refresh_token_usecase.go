package usecase

import (
	"context"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
)

type refreshTokenUsecase struct {
	cfg            *configs.Config
	userRepository domain.UserRepository
}

func NewRefreshTokenUsecase(cfg *configs.Config, userRepository domain.UserRepository) domain.RefreshTokenUsecase {
	return &refreshTokenUsecase{
		cfg:            cfg,
		userRepository: userRepository,
	}
}

func (ru *refreshTokenUsecase) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*domain.RefreshTokenResult, *domain.UsecaseError) {
	id, err := ru.extractIDFromToken(refreshToken, ru.cfg.TokenConfig.RefreshTokenSecret)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonUnauthorized)
	}

	user, err := ru.getUserByID(ctx, id)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonUnauthorized)
	}

	accessToken, err := ru.createAccessToken(
		user, ru.cfg.TokenConfig.AccessTokenSecret, ru.cfg.TokenConfig.AccessTokenExpiryHour)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	newRefreshToken, err := ru.createRefreshToken(
		user, ru.cfg.TokenConfig.RefreshTokenSecret, ru.cfg.TokenConfig.RefreshTokenExpiryHour)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonServerError)
	}

	user.RefreshToken = newRefreshToken
	_, err = ru.updateUser(ctx, user)
	if err != nil {
		return nil, domain.NewUsecaseError(err, domain.ReasonDBError)
	}

	return &domain.RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (ru *refreshTokenUsecase) getUserByID(ctx context.Context, id int64) (*domain.User, error) {
	return ru.userRepository.GetUserById(ctx, id)
}

func (ru *refreshTokenUsecase) createAccessToken(user *domain.User, secret string, expiry int) (string, error) {
	return internal.CreateAccessToken(user, secret, expiry)
}

func (ru *refreshTokenUsecase) createRefreshToken(user *domain.User, secret string, expiry int) (string, error) {
	return internal.CreateRefreshToken(user, secret, expiry)
}

func (ru *refreshTokenUsecase) extractIDFromToken(token string, secret string) (int64, error) {
	return internal.ExtractIDFromToken(token, secret)
}

func (ru *refreshTokenUsecase) updateUser(ctx context.Context, user *domain.User) (int64, error) {
	return ru.userRepository.UpdateUser(ctx, user)
}
