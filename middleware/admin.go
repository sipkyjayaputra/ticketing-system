package middleware

import (
	"net/http"

	"github.com/sipkyjayaputra/ticketing-system/utils"

	"github.com/gin-gonic/gin"
)

func AdminAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get("role")

		if !ok {
			c.JSON(http.StatusForbidden, utils.BuildForbiddenAccessResponse("admin role access only"))
			c.Abort()
			return
		}

		isAdmin := role == "admin" || role == "super_admin"
		if !isAdmin {
			c.JSON(http.StatusForbidden, utils.BuildForbiddenAccessResponse("admin role access only"))
			c.Abort()
			return
		}

		c.Next()
	}
}
