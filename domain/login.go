package domain

import "context"

type LoginRequest struct {
	Username string `form:"username" binding:"validateUsername"`
	Email    string `form:"email" binding:"validateEmail"`
	Password string `form:"password" binding:"required,min=3,max=50"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type LoginUsecase interface {
	Login(context.Context, LoginRequest) (*LoginResult, *UsecaseError)
	LoginWithGoogle(context.Context, string) (*LoginResult, *UsecaseError)
}
