package main

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/routes"
import "github.com/BTBurke/gaea-server/middleware"
import "github.com/BTBurke/gaea-server/errors"

import _ "github.com/lib/pq"
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
	
	

	r.GET("/error", func(c *gin.Context) {
		c.Set("user", "usertest")
		c.Set("role", "admin")
		c.AbortWithError(422, errors.NewAPIError(422, "test error development msg", "test error user message", c))
		return
	})

	r.POST("/login", routes.Login(db))
	r.GET("/user", routes.GetCurrentUser(db))

	r.GET("/announcement", routes.GetAnnouncements)

	r.GET("/inventory", routes.GetInventory(db))
	r.POST("/inventory/csv", routes.CreateInventoryFromCSVString(db))
	r.POST("/inventory", routes.CreateItem(db))
	r.PUT("/inventory/:invID", routes.UpdateItem(db))
	r.GET("/inventory/:invID/effects", routes.GetEffects(db))

	r.GET("/sale", routes.GetSales(db))
	r.POST("/sale", routes.CreateSale(db))
	r.PUT("/sale/:saleID", routes.UpdateSale(db))

	r.GET("/order", routes.GetOrders(db))
	r.POST("/order", routes.CreateOrder(db))
	r.PUT("/order/:orderID", routes.UpdateOrderStatus(db))
	r.GET("/order/:orderID/item", routes.GetOrderItems(db))
	r.POST("/order/:orderID/item", routes.AddOrderItem(db))
	r.DELETE("/order/:orderID/item/:itemID", routes.DeleteOrderItem)
	r.PUT("/order/:orderID/item/:itemID", routes.UpdateOrderItem(db))

	// When developing on c9
	r.Run(":8080")


}
