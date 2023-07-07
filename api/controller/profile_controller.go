package controller

import (
	"github.com/gin-gonic/gin"
	"home-bar/api/middleware"
	"home-bar/domain"
	"home-bar/internal"
	"net/http"
)

type ProfileController struct {
	ProfileUsecase domain.ProfileUsecase
}

func (pc *ProfileController) GetProfile(c *gin.Context) {
	userID := c.GetInt64(middleware.UserIDContextKey)

	profile, usecaseError := pc.ProfileUsecase.GetProfile(c, userID)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	if profile == nil {
		c.JSON(http.StatusNotFound, domain.GetErrorResponse(
			domain.NewCustomError("Profile not found with the given user id")))
		return
	}

	c.JSON(http.StatusOK, profile)
}
