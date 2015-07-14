package routes

import "github.com/gin-gonic/gin"
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/error"
import "github.com/guregu/null/zero"
import "fmt"

import "time"

type User struct {
	UserName    string      `json:"user_name" db:"user_name"`
	FirstName   string      `json:"first_name" db:"first_name"`
	LastName    string      `json:"last_name" db:"last_name"`
	Email       string      `json:"email" db:"email"`
	Role        string      `json:"role" db:"role"`
	Password    zero.String `json:"-" db:"password"`
	DipID       zero.String `json:"dip_id" db:"dip_id"`
	Passport    zero.String `json:"passport" db:"passport"`
	Section     zero.String `json:"section" db:"section"`
	UpdatedAt   time.Time   `json:"-" db:"updated_at"`
	UpdateToken zero.String `json:"-" db:"update_token"`
}

func GetCurrentUser(db *sqlx.DB) gin.HandlerFunc {
	// For testing only

	return func(c *gin.Context) {
		var user1 = User{}

		var userName = "burkebt" // should get this from token

		err := db.Get(&user1, "SELECT * from gaea.user WHERE user_name=$1", userName)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, error.APIError{503, "failed on getting user", "internal server error"})
			return
		}
		c.JSON(200, user1)
	}
}
