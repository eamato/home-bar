package route

import (
	"github.com/gin-gonic/gin"
	"home-bar/api/middleware"
	"home-bar/configs"
	"home-bar/database"
	"home-bar/domain"
	"home-bar/repository"
)

func Setup(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	initPublicRouters(config, database, router)
	initProtectedRouters(config, database, router)
	initPrivateRouters(config, database, router)
}

func initPublicRouters(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	publicRouterGroup := router.Group("")
	NewSignupRouter(config, database, publicRouterGroup)
	NewLoginRouter(config, database, publicRouterGroup)
	NewRefreshTokenRouter(config, database, publicRouterGroup)
}

func initProtectedRouters(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	protectedRouterGroup := router.Group("")
	protectedRouterGroup.Use(middleware.JWTAuthMiddleware(config.TokenConfig.AccessTokenSecret))
	NewProfileRouter(config, database, protectedRouterGroup)
}

func initPrivateRouters(config *configs.Config, database database.Database, router *gin.RouterGroup) {
	roleRepository := repository.NewRoleRepository(database, domain.RolesTableName, domain.UsersHasRolesTableName)

	privateRouterGroup := router.Group("")
	privateRouterGroup.Use(middleware.JWTAuthMiddleware(config.TokenConfig.AccessTokenSecret))
	privateRouterGroup.Use(middleware.AdminAccessMiddleware(roleRepository))
	NewAdminRouter(config, database, privateRouterGroup)
}
