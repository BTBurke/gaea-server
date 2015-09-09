package middleware

import (
	"strings"

	"github.com/BTBurke/gaea-server/log"
	"github.com/BTBurke/gaea-server/routes"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtComplete := c.Request.Header.Get("Authorization")
		if len(jwtComplete) == 0 {
			log.Error("failed to receive JWT dev=no JWT in header")
			c.AbortWithStatus(401)
			return
		}

		jwtString := strings.Split(jwtComplete, " ")
		if len(jwtString) != 2 || jwtString[0] != "Bearer" {
			log.Error("failed to receive JWT dev=JWT not right format")
			c.AbortWithStatus(401)
			return
		}

		token, err := routes.ValidateJWT(jwtString[1])
		if err != nil {
			log.Error("failed to validate JWT dev=JWT failed validation err=%s", err)
			c.AbortWithStatus(401)
			return
		}

		c.Set("user", token.Claims["user"].(string))
		c.Set("jwt", jwtString[1])
		c.Set("role", token.Claims["role"].(string))
		c.Set("exp", token.Claims["exp"].(string))
		c.Set("iss", token.Claims["iss"].(string))

		newJwt, err := routes.RenewJWTfromJWT(jwtString[1])
		if err != nil {
			log.Error("failed to validate JWT dev=JWT failed validation err=%s", err)
			c.AbortWithStatus(401)
			return
		}

		c.Writer.Header().Set("Authorization", strings.Join([]string{"Bearer", newJwt}, " "))

		c.Next()

	}
}
