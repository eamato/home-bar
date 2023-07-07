package controller

import (
	"github.com/gin-gonic/gin"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
	"net/http"
)

type SignupController struct {
	Cfg           *configs.Config
	SignupUsecase domain.SignupUsecase
}

func (sc *SignupController) Signup(c *gin.Context) {
	var signupRequest domain.SignupRequest

	if err := c.ShouldBind(&signupRequest); err != nil {
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	newUser := domain.User{
		Username: signupRequest.Username,
		Email:    signupRequest.Email,
		Password: signupRequest.Password,
	}

	signupResult, usecaseError := sc.SignupUsecase.Signup(c, &newUser)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	signupResponse := domain.SignupResponse{
		AccessToken:  signupResult.AccessToken,
		RefreshToken: signupResult.RefreshToken,
	}

	c.JSON(http.StatusOK, signupResponse)
}
