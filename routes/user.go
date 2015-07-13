package routes

import "github.com/gin-gonic/gin"
import "time"

type User struct {
	UserName  string `json:"user_name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	password    string
	DipID       string `json:"dip_id"`
	Passport    string `json:"passport"`
	Section     string `json:"section"`
	updatedAt   time.Time
	updateToken string
}

func GetCurrentUser(c *gin.Context) {
	// For testing only
	testUser := User{
		UserName:  "ambassadorjs",
		FirstName: "Joe",
		LastName:  "Ambassador",
		Email:     "AmbassadorJS@state.gov",
		Role:      "admin",
	}
	c.JSON(200, testUser)
}
