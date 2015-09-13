package routes

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/BTBurke/gaea-server/auth"
	"github.com/BTBurke/gaea-server/email"
	"github.com/BTBurke/gaea-server/errors"
	"github.com/BTBurke/gaea-server/log"
	"github.com/elithrar/simple-scrypt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type LoginRequest struct {
	Email  string `json:"user"`
	Pwd    string `json:"pwd"`
	Source string `json:"source"`
}

type LoginResponse struct {
	JWT      string `json:"jwt"`
	User     string `json:"user"`
	Content  string `json:"content"`
	Redirect string `json:"redirect"`
}

type ResetRequest struct {
	Email string `json:"user"`
}

type SetRequest struct {
	Pwd   string `json:"pwd"`
	Token string `json:"token"`
}

type AccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func Login(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginReq LoginRequest

		bindErr := c.Bind(&loginReq)
		if bindErr != nil {
			c.AbortWithError(422, errors.NewAPIError(422, "failed to bind login request", "username or password incorrect", c))
			return
		}

		var user1 User
		dbErr := db.Get(&user1, "SELECT * FROM gaea.user WHERE email=$1", loginReq.Email)
		if dbErr != nil {
			c.AbortWithError(401, errors.NewAPIError(401, "failed to get user from DB on login request", "username or password incorrect", c))
			return
		}

		userPwd, err := user1.Password.MarshalText()
		if err != nil {
			c.AbortWithError(401, errors.NewAPIError(401, "failed to marshal text on user password", "username or password incorrect", c))
			return
		}

		compErr := scrypt.CompareHashAndPassword(userPwd, []byte(loginReq.Pwd))
		if compErr != nil {
			c.AbortWithError(401, errors.NewAPIError(401, "login password does not match", "username or password incorrect", c))
			return
		}

		jwtString, err := IssueJWTForUser(user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to issue new JWT", "internal server error", c))
			return
		}

		var loginResp LoginResponse
		loginResp.User = user1.UserName
		loginResp.JWT = jwtString
		c.Writer.Header().Set("Authorization", strings.Join([]string{"Bearer", jwtString}, " "))
		c.JSON(200, loginResp)
	}
}

func RequestResetEmail(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reset ResetRequest
		err := c.Bind(&reset)
		if err != nil {
			c.AbortWithError(422, errors.NewAPIError(503, "unable to bind on password reset", "malformed request", c))
			return
		}

		var user User
		dbErr := db.Get(&user, "SELECT * FROM gaea.user WHERE email=$1", reset.Email)
		if dbErr != nil {
			switch dbErr {
			case sql.ErrNoRows:
				// TODO: probably need to log these phantom resets somewhere
				log.Warn("msg=attempted password reset not a valid user email=%s", reset.Email)
				c.JSON(200, gin.H{"status": "ok"})
				return
			default:
				log.Error("msg=database error reseting password email=%s err=%s", reset.Email, dbErr.Error())
				c.JSON(200, gin.H{"status": "ok"})
				return
			}
		}

		pwdJwt, err := IssuePwdJWTForUser(user)
		if err != nil {
			// TODO: Needs to log to LE to troubleshoot later
			c.JSON(200, gin.H{"status": "ok"})
			return
		}
		body, err := email.PasswordResetEmail(user.FirstName, pwdJwt)
		if err != nil {
			c.JSON(200, gin.H{"status": "ok"})
			return
		}

		go email.Send("GAEA Accounts <help@guangzhouaea.org>", "GAEA Password Reset", body, user.Email)
		c.JSON(200, gin.H{"status": "ok"})
	}
}

func TestAuth(c *gin.Context) {
	if ok := auth.MustUser(c, "yomama"); !ok {
		c.AbortWithStatus(401)
		return
	}
	user, _ := c.Get("user")
	role, _ := c.Get("role")
	iss, _ := c.Get("iss")
	exp, _ := c.Get("exp")
	oldjwt, _ := c.Get("jwt")

	c.JSON(200, gin.H{"user": user, "role": role, "iss": iss, "exp": exp, "oldjwt": oldjwt})
}

func SetPassword(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req SetRequest
		if err := c.Bind(&req); err != nil {
			log.Error("msg=password reset failed dev=could not bind request err=%s", err)
			c.AbortWithStatus(401)
			return
		}

		token, err := ValidateJWT(req.Token)
		if err != nil {
			log.Error("msg=password reset failed dev=token not validated err=%s", err)
			c.AbortWithStatus(401)
			return
		}

		username := token.Claims["user"]
		setToken := token.Claims["jwt"]

		if token.Claims["role"] != "pwd" {
			log.Error("msg=password reset failed user=%s dev=token role not pwd", username)
			c.AbortWithStatus(401)
			return
		}

		hashPwd, err := scrypt.GenerateFromPassword([]byte(req.Pwd), scrypt.DefaultParams)
		if err != nil {
			log.Error("msg=password reset failed user=%s dev=could not hash pwd err=%s", username.(string), err)
			c.AbortWithStatus(401)
			return
		}

		// prevent replay attack by checking that token not already used to reset password
		var checkUserToken string
		if err := db.Get(&checkUserToken, "SELECT update_token FROM gaea.user WHERE user_name=$1", username.(string)); err != nil {
			log.Error("msg=password reset failed user=%s dev=db insert failed err=%s", username.(string), err)
			c.AbortWithStatus(401)
			return
		}
		if checkUserToken == setToken.(string) {
			log.Error("msg=password reset failed user=%s dev=replay attack against token err=%s", username.(string), err)
			c.AbortWithStatus(401)
			return
		}

		var user1 User
		if err := db.Get(&user1, "UPDATE gaea.user SET password=$1, update_token=$2 WHERE user_name=$3 RETURNING *", hashPwd, setToken.(string), username.(string)); err != nil {
			log.Error("msg=password reset failed user=%s dev=db insert failed err=%s", username, err)
			c.AbortWithStatus(401)
			return
		}

		jwtString, err := IssueJWTForUser(user1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to issue new JWT", "internal server error", c))
			return
		}

		var loginResp LoginResponse
		loginResp.User = user1.UserName
		loginResp.JWT = jwtString
		c.Writer.Header().Set("Authorization", strings.Join([]string{"Bearer", jwtString}, " "))
		c.JSON(200, loginResp)

	}
}
