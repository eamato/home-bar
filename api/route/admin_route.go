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

func NewAdminRouter(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	adminRepository := repository.NewAdminRepository(
		database,
		domain.UsersTableName,
		domain.RolesTableName,
		domain.ProfilesTableName,
		domain.UsersHasRolesTableName,
	)
	roleRepository := repository.NewRoleRepository(database, domain.RolesTableName, domain.UsersHasRolesTableName)

	adminController := controller.AdminController{
		AdminUsecase: usecase.NewAdminUsecase(adminRepository, roleRepository),
	}

	router.GET("/admin/users", adminController.GetUsers)
	router.DELETE("/admin/user", adminController.DeleteUser)
	router.PATCH("/admin/user/role", adminController.AssignUserRole)
}
