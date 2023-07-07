package route

import (
	"github.com/gin-gonic/gin"
	"home-bar/api/controller"
	"home-bar/configs"
	"home-bar/database"
	"home-bar/domain"
	"home-bar/repository"
	"home-bar/usecase"
)

func NewRefreshTokenRouter(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(database, domain.UsersTableName)

	refreshTokenController := controller.RefreshTokenController{
		Cfg:                 config,
		RefreshTokenUsecase: usecase.NewRefreshTokenUsecase(config, userRepository),
	}

	router.POST("/refresh", refreshTokenController.RefreshToken)
}
