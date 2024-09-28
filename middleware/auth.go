package middleware

import (
	"github.com/sipkyjayaputra/ticketing-system/helpers"
	"github.com/sipkyjayaputra/ticketing-system/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, ok := c.Request.Header["Authorization"]
		if !ok {
			c.JSON(http.StatusUnauthorized, utils.BuildUnauthorizedResponse("unathorized", "access token not found"))
			c.Abort()
			return
		}

		authorizationFields := strings.Fields(auth[0])

		if len(authorizationFields) != 2 || authorizationFields[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, utils.BuildUnauthorizedResponse("unathorized", "bearer token malformed"))
			c.Abort()
			return
		}

		token := authorizationFields[1]
		claims, err := helpers.DecodeJWTToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, utils.BuildUnauthorizedResponse("unathorized", err.Error()))
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
