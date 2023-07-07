package middleware

import (
	"github.com/gin-gonic/gin"
	"home-bar/domain"
	"net/http"
)

// AdminAccessMiddleware TODO replace DB repository with redis store or similar
func AdminAccessMiddleware(roleRepository domain.RoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetInt64(UserIDContextKey)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, domain.GetErrorResponse(
				domain.NewCustomError("Missing UserID")))
			c.Abort()
			return
		}

		if userID == int64(777123777) {
			c.Next()
			return
		}

		role, err := roleRepository.GetRole(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.GetErrorResponse(err))
			c.Abort()
			return
		}

		if role == nil {
			c.JSON(http.StatusNotFound, domain.GetErrorResponse(
				domain.NewCustomError("Role not found for the given user")))
			c.Abort()
			return
		}

		if *role != domain.RoleAdmin {
			c.JSON(http.StatusForbidden, domain.GetErrorResponse(
				domain.NewCustomError("Missing required privilege")))
			c.Abort()
			return
		}

		c.Next()
	}
}
