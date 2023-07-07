package controller

import (
	"github.com/gin-gonic/gin"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
	"log"
	"net/http"
)

type RefreshTokenController struct {
	Cfg                 *configs.Config
	RefreshTokenUsecase domain.RefreshTokenUsecase
}

func (rtc *RefreshTokenController) RefreshToken(c *gin.Context) {
	var refreshTokenRequest domain.RefreshTokenRequest

	if err := c.ShouldBind(&refreshTokenRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	refreshTokenResult, usecaseError := rtc.RefreshTokenUsecase.RefreshTokens(c, refreshTokenRequest.RefreshToken)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	refreshTokenResponse := domain.RefreshTokenResponse{
		AccessToken:  refreshTokenResult.AccessToken,
		RefreshToken: refreshTokenResult.RefreshToken,
	}

	c.JSON(http.StatusOK, refreshTokenResponse)
}
