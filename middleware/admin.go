package middleware

import (
	"github.com/BTBurke/gaea-server/auth"
	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if ok := auth.MustAdmin(c); !ok {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
