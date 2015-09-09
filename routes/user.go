package routes

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)
import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/errors"
import "github.com/BTBurke/gaea-server/email"
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
	LastLogin   time.Time   `json:"-" db:"last_login"`
	MemberExp   time.Time   `json:"member_exp" db:"member_exp"`
	MemberType  zero.String `json:"member_type" db:"member_type"`
	StripeToken zero.String `json:"-" db:"stripe_token"`
}

func GetCurrentUser(db *sqlx.DB) gin.HandlerFunc {
	// For testing only

	return func(c *gin.Context) {
		var user1 = User{}

		userName, exists := c.Get("user")
		if !exists {
			c.AbortWithError(503, errors.NewAPIError(503, "failed on getting user from token", "internal server error", c))
			return
		}

		err := db.Get(&user1, "SELECT * from gaea.user WHERE user_name=$1", userName)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on getting user", "internal server error", c))
			return
		}
		c.JSON(200, user1)
	}
}

func GetAllUsers(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var users = []User{}

		err := db.Select(&users, "SELECT * FROM gaea.user")
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on getting all users", "failed to get all users", c))
			return
		}

		c.JSON(200, gin.H{"qty": len(users), "users": users})

	}
}

func CreateUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user1 = User{}
		if err := c.Bind(&user1); err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed on creating user", "failed to update user", c))
			return
		}
		username, err := createUserName(user1, db, 0)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed on creating user", "failed to update user", c))
			return
		}

		var retUser User
		if err := db.Get(&retUser, `INSERT INTO gaea.user
				(user_name, first_name, last_name,
					email, role, password, dip_id, passport,
					section, updated_at, update_token, last_login,
					member_exp, member_type, stripe_token) VALUES
					($1, $2, $3,
					$4, $5, $6, $7, $8,
					$9, $10, $11, $12,
					$13, $14, $15) RETURNING *`,
			username, user1.FirstName, user1.LastName,
			user1.Email, user1.Role, "", user1.DipID, user1.Passport,
			user1.Section, time.Now(), "", user1.LastLogin,
			user1.MemberExp, user1.MemberType, ""); err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on inserting user", "failed to update user", c))
			return
		}

		pwdJwt, err := IssuePwdJWTForUser(retUser)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, "failed on created pwd jwt for new user", "failed to create user", c))
			return
		}

		body, err := email.NewAccountPasswordEmail(retUser.FirstName, pwdJwt)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, "failed to create welcome email", "failed to create user", c))
			return
		}

		go email.Send("GAEA Accounts <accounts@guangzhouaea.org>", "Welcome to the GAEA website", body, retUser.Email)

		c.JSON(200, retUser)

	}
}

// Recursive function to find a username starting with LastName+First Letter then adding
// numbers until you get a unique name (e.g. burkeb, burkeb1, burkeb2)
func createUserName(user1 User, db *sqlx.DB, inc int) (string, error) {
	var rUser User

	firstLetter := strings.Split(user1.FirstName, "")[0]
	switch {
	case inc == 0:
		tryName := strings.ToLower(strings.Join([]string{user1.LastName, firstLetter}, ""))
		err := db.Get(&rUser, "SELECT * FROM gaea.user WHERE user_name=$1", tryName)
		if err == sql.ErrNoRows {
			return tryName, nil
		}
		if err != nil {
			return "", err
		}
		return createUserName(user1, db, inc+1)
	default:
		tryName := strings.ToLower(strings.Join([]string{user1.LastName, firstLetter, strconv.Itoa(inc)}, ""))
		err := db.Get(&rUser, "SELECT * FROM gaea.user WHERE user_name=$1", tryName)
		if err == sql.ErrNoRows {
			return tryName, nil
		}
		if err != nil {
			return "", err
		}
		return createUserName(user1, db, inc+1)
	}
}

func UpdateUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user1 User
		err := c.Bind(&user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed on updating user", "failed to update user", c))
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
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating user", "failed to update user", c))
			return
		}

		c.JSON(200, retU)
	}
}

func DeleteUser(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user1 User
		err := c.Bind(&user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed on deleting user", "failed to delete user", c))
			return
		}

		dbErr := db.MustExec(`DELETE gaea.user WHERE user_name = $1`, user1.UserName)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating user", "failed to update user", c))
			return
		}

		c.JSON(200, gin.H{"user_name": user1.UserName})
	}
}
