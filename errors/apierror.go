package errors

import "fmt"
import "os"
import "strings"

import "github.com/gin-gonic/gin"
import "github.com/bsphere/le_go"


type serverError struct {
    Code int
    Developer string
    User string
}

func (a serverError) Error() string { return a.Developer }

func NewAPIError(code int, dev string, user string, c *gin.Context) serverError {
    s := serverError{
        Code: code,
        Developer: dev,
        User: user,
    }
    
    go sendLog("api", s, c)
    return s
}

func NewDBError(code int, dev string, user string, c *gin.Context) serverError {
    
    s := serverError{
        Code: code,
        Developer: dev,
        User: user,
    }
    
    go sendLog("db", s, c)
    return s
}

func sendLog(errorType string, s serverError, c *gin.Context) {
    apiToken := os.Getenv("LE_TOKEN")
    if len(apiToken) == 0 {
        fmt.Println("Warning: LE_TOKEN not set.  Unable to send application logs.")
    } else {
        le, err := le_go.Connect(apiToken)
        defer le.Close()
        if err != nil {
            fmt.Println("Error connecting to logentries. Msg: %s", err)
        } else {
            user, _ := c.Get("user")
            role, _ := c.Get("role")
            
            logFormat := "level=error type=%s code=%d dev_msg=\"%s\" user=%s role=%s params={%s} route=%s method=%s" 
            le.Printf(logFormat,
                errorType,
                s.Code,
                s.Developer,
                user,
                role,
                parseParams(c),
                c.Request.URL,
                c.Request.Method)
          
        }
     
    }
}

func parseParams(c *gin.Context) string {
    var out []string
    for _, param := range c.Params {
        out = append(out, param.Key + ": " + param.Value)
    }
  
    return strings.Join(out, ", "   ) 
    
}