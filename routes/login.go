package routes

import "github.com/jmoiron/sqlx"
import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/errors"
import "github.com/elithrar/simple-scrypt"

type LoginRequest struct {
	Email string `json:"user"`
	Pwd  string `json:"pwd"`
}

type LoginResponse struct {
	JWT      string `json:"jwt"`
	User     string `json:"user"`
	Content  string `json:"content"`
	Redirect string `json:"redirect"`
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
             c.AbortWithError(503, errors.NewAPIError(503, "failed to issue new JWT", "internal server error", c))
			return
        
        }
        
        var loginResp LoginResponse
        loginResp.User = user1.UserName
        loginResp.JWT = jwtString
        c.JSON(200, loginResp)
		
	}
}
