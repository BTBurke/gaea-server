package errors

import "github.com/gin-gonic/gin"

type serverError struct {
    Code int
    Developer string
    User string
}

func (a serverError) Error() string { return a.Developer }

func NewAPIError(code int, dev string, user string, c *gin.Context) serverError {
    fmt.Println(c)
    return serverError{
        Code: code,
        Developer: dev,
        User: user,
    }
}

func NewDBError(code int, dev string, user string, c *gin.Context) serverError {
    
    fmt.Println(c)
    return serverError{
        Code: code,
        Developer: dev,
        User: user,
    }
}