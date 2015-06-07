package routes

import "github.com/gin-gonic/gin"

type User struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Uuid      string `json:"uuid"`
}

func GetCurrentUser(c *gin.Context) {
    // For testing only
    testUser := User{
        UserName: "ambassadorjs",
        FirstName: "Joe",
        LastName: "Ambassador",
        Email: "AmbassadorJS@state.gov",
        Role: "member",
        Uuid: "06c0eb9f-92c5-485f-9622-c3f225eb6a95",
    }
    c.JSON(200, testUser)
}