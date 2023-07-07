package domain

import "context"

type SignupRequest struct {
	Username string `form:"username" binding:"required,min=3,max=50"`
	Email    string `form:"email" binding:"required,email,min=3,max=50"`
	Password string `form:"password" binding:"required,min=3,max=50"`
}

type SignupResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupResult struct {
	AccessToken  string
	RefreshToken string
}

type SignupUsecase interface {
	Signup(context.Context, *User) (*SignupResult, *UsecaseError)
}
