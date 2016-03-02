package routes

import "github.com/gin-gonic/gin"

type LogoutRequest struct {
	User string `json:"user"`
}

func Logout(c *gin.Context) {
	// TODO: Delete JWT secret key for user on logout
	c.JSON(200, gin.H{"status": "ok"})
}
