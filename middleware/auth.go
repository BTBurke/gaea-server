package middleware

import (
	"strings"

	"github.com/BTBurke/gaea-server/routes"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtComplete := c.Request.Header.Get("Authorization")
		if len(jwtComplete) == 0 {
			c.AbortWithStatus(401)
			return
		}

		jwtString := strings.Split(jwtComplete, " ")
		if len(jwtString) != 2 || jwtString[0] != "Bearer" {
			c.AbortWithStatus(401)
			return
		}

		token, err := routes.ValidateJWT(jwtString[1])
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", token.Claims["user"])
		c.Set("jwt", jwtString[1])
		c.Set("role", token.Claims["role"])
		c.Set("exp", token.Claims["exp"])
		c.Set("iss", token.Claims["iss"])

		newJwt, err := routes.RenewJWTfromJWT(jwtString[1])
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		c.Writer.Header().Set("Authorization", strings.Join([]string{"Bearer", newJwt}, " "))

		c.Next()

	}
}
