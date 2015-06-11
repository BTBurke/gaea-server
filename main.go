package main

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/routes"
import "github.com/BTBurke/gaea-server/middleware"

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/401", func(c *gin.Context) {
		c.String(401, "Unauthorized")
	})

	r.GET("/user", routes.GetCurrentUser)
	r.GET("/sale", routes.GetCurrentSale)
	r.GET("/order", routes.GetOrders)
	r.GET("/announcement", routes.GetAnnouncements)
	r.GET("/inventory", routes.GetInventory)

	// When developing on c9
	// r.Run(":8080")

	// Local development
	r.Run(":9000")

}
