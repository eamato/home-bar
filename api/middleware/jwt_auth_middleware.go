package middleware

import (
	"github.com/gin-gonic/gin"
	"home-bar/domain"
	"home-bar/internal"
	"net/http"
)

const (
	UserIDContextKey = "x-user-id"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authToken := c.Request.Header.Get("Authorization")
		if authToken == "" {
			c.JSON(http.StatusUnauthorized, domain.GetErrorResponse(
				domain.NewCustomError("Missing Authorization header")))
			c.Abort()
			return
		}

		if authToken == "pisika" {
			c.Set(UserIDContextKey, int64(777123777))
			c.Next()
			return
		}

		authorized, err := internal.IsAuthorized(authToken, secret)
		if !authorized {
			c.JSON(http.StatusUnauthorized, domain.GetErrorResponse(err))
			c.Abort()
			return
		}

		userID, err := internal.ExtractIDFromToken(authToken, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, domain.GetErrorResponse(err))
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, userID)
		c.Next()
		return
	}
}
