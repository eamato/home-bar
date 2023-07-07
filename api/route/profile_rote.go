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

func NewProfileRouter(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	profileRepository := repository.NewProfileRepository(database, domain.ProfilesTableName)

	profileController := controller.ProfileController{
		ProfileUsecase: usecase.NewProfileUsecase(profileRepository),
	}

	router.GET("/profile", profileController.GetProfile)
}
