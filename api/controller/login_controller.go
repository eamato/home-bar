package controller

import (
	"github.com/gin-gonic/gin"
	"home-bar/configs"
	"home-bar/domain"
	"home-bar/internal"
	"log"
	"net/http"
)

type LoginController struct {
	Cfg          *configs.Config
	LoginUsecase domain.LoginUsecase
}

func (lc *LoginController) Login(c *gin.Context) {
	var loginRequest domain.LoginRequest

	if err := c.ShouldBind(&loginRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	loginResult, usecaseError := lc.LoginUsecase.Login(c, loginRequest)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (lc *LoginController) ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (lc *LoginController) LoginGoogle(c *gin.Context) {
	if lc.Cfg.OAuthConfig == nil {
		c.JSON(http.StatusInternalServerError, domain.GetErrorResponse(
			domain.NewCustomError("OAuth config error")))
		return
	}

	url := lc.Cfg.OAuthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (lc *LoginController) ProceedGoogleAuth(c *gin.Context) {
	code := c.Query("code")

	loginResult, usecaseError := lc.LoginUsecase.LoginWithGoogle(c, code)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	loginResponse := domain.LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
	}

	c.JSON(http.StatusOK, loginResponse)
}
