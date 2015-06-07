package main

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/routes"

func main() {
    r := gin.Default()
    
    r.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    r.GET("/user", routes.GetCurrentUser)
    

    r.Run(":9000")
}