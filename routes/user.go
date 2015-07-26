package routes

import "github.com/gin-gonic/gin"
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/errors"
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
			c.AbortWithError(503, errors.APIError{503, "failed on getting user", "internal server error"})
			return
		}
		c.JSON(200, user1)
	}
}

func GetAllUsers(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Add check for role admin or higher
		
		var users = []User{}
		
		err := db.Select(&users, "SELECT * FROM gaea.user")
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.APIError{503, "failed on getting all users", "failed to get all users"})
			return
		}
		
		c.JSON(200, gin.H{"qty": len(users), "users": users})
		
	}
}

func UpdateUser(db *sqlx.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		
		var user1 User
		err := c.Bind(&user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.APIError{422, "failed on updating user", "failed to update user"})
			return
		}
		
		var retU User
		dbErr := db.Get(&retU, `UPDATE gaea.user SET 
		first_name = $1, last_name = $2, email = $3,
		dip_id = $4, passport = $5, section = $6, updated_at = $7
		WHERE user_name = $8 RETURNING *`,
		user1.FirstName, user1.LastName, user1.Email,
		user1.DipID, user1.Passport, user1.Section, time.Now(),
		user1.UserName)
		
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.APIError{503, "failed on updating user", "failed to update user"})
			return		
		}
		
		c.JSON(200, retU)
	}
}

func DeleteUser(db *sqlx.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var user1 User
		err := c.Bind(&user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.APIError{422, "failed on deleting user", "failed to delete user"})
			return
		}
		
		dbErr := db.MustExec(`DELETE gaea.sale WHERE user_name = $1`, user1.UserName)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.APIError{503, "failed on updating user", "failed to update user"})
			return		
		}
		
		c.JSON(200, gin.H{"user_name": user1.UserName})
	}
}