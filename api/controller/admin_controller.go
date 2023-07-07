package controller

import (
	"github.com/gin-gonic/gin"
	"home-bar/domain"
	"home-bar/internal"
	"log"
	"net/http"
)

type AdminController struct {
	AdminUsecase domain.AdminUsecase
}

func (ac *AdminController) GetUsers(c *gin.Context) {
	var usersListRequest domain.UsersListRequest
	if err := c.ShouldBindQuery(&usersListRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	getUsersResult, usecaseError := ac.AdminUsecase.GetUsers(c, usersListRequest)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	c.JSON(http.StatusOK, getUsersResult)
}

func (ac *AdminController) DeleteUser(c *gin.Context) {
	var userDeleteRequest domain.UserDeleteRequest

	if err := c.ShouldBind(&userDeleteRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	userDeletionResult, usecaseError := ac.AdminUsecase.DeleteUser(c, userDeleteRequest)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	c.JSON(http.StatusOK, userDeletionResult)
}

func (ac *AdminController) AssignUserRole(c *gin.Context) {
	var roleAssignRequest domain.RoleAssignRequest

	if err := c.ShouldBind(&roleAssignRequest); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, domain.GetErrorResponse(err))
		return
	}

	userRoleAssignmentResult, usecaseError := ac.AdminUsecase.AssignRole(c, roleAssignRequest)
	if usecaseError != nil {
		internal.PrintError("", usecaseError.Error)

		code, errorResponse := usecaseError.ParseUsecaseErrorToRest()
		c.JSON(code, errorResponse)
		return
	}

	c.JSON(http.StatusOK, userRoleAssignmentResult)
}
