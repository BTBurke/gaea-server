package auth

import "github.com/gin-gonic/gin"

func MustUser(c *gin.Context, user string) bool {

	// shortcut for admin override
	if ok := MustAdmin(c); ok {
		return true
	}

	userFromJwt, exists := c.Get("user")
	if !exists {
		return false
	}
	if userFromJwt.(string) != user {
		return false
	}
	return true
}

func MustAdmin(c *gin.Context) bool {
	roleFromJwt, exists := c.Get("role")
	if !exists {
		return false
	}
	if roleFromJwt.(string) == "admin" || roleFromJwt.(string) == "superadmin" {
		return true
	} else {
		return false
	}
}
