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

func NewSignupRouter(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	userRepository := repository.NewUserRepository(database, domain.UsersTableName)
	profileRepository := repository.NewProfileRepository(database, domain.ProfilesTableName)
	roleRepository := repository.NewRoleRepository(database, domain.RolesTableName, domain.UsersHasRolesTableName)

	signupController := controller.SignupController{
		Cfg:           config,
		SignupUsecase: usecase.NewSignupUsecase(config, userRepository, profileRepository, roleRepository),
	}

	router.POST("/signup", signupController.Signup)
}
