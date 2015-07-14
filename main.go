package main

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/routes"
import "github.com/BTBurke/gaea-server/middleware"

import _ "github.com/lib/pq"
//import "database/sql"
import "github.com/jmoiron/sqlx"

import "log"

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())
	
	// Connect to database
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=db_gaea sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/401", func(c *gin.Context) {
		c.String(401, "Unauthorized")
	})

	r.GET("/user", routes.GetCurrentUser(db))
	
	
	r.GET("/announcement", routes.GetAnnouncements)
	r.GET("/inventory", routes.GetInventory)
	
	r.GET("/sale", routes.GetCurrentSale)
	r.PUT("/sale/:saleID", routes.UpdateSale)
	
	
	r.GET("/order", routes.GetOrders)
	r.POST("/order", routes.CreateOrder(db))
	r.GET("/order/:orderID/item", routes.GetOrderItems(db))
	r.POST("/order/:orderID/item", routes.AddOrderItem(db))
	r.DELETE("/order/:orderID/item/:itemID", routes.DeleteOrderItem)
	r.PUT("/order/:orderID/item/:itemID", routes.UpdateOrderItem)

	// When developing on c9
	r.Run(":8080")

	// Local development
	//r.Run(":9000")

}
