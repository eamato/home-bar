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

func SetupWeb(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(database, domain.UsersTableName)
	profileRepository := repository.NewProfileRepository(database, domain.ProfilesTableName)
	roleRepository := repository.NewRoleRepository(database, domain.RolesTableName, domain.UsersHasRolesTableName)

	loginController := controller.LoginController{
		Cfg:          config,
		LoginUsecase: usecase.NewLoginUsecase(config, userRepository, profileRepository, roleRepository),
	}

	router.GET("/login", loginController.ShowLoginPage)
	router.GET("/login/google", loginController.LoginGoogle)
	router.GET("/auth/google/callback", loginController.ProceedGoogleAuth)
}
