package main

// Release: v0.1.0

import "os"
import "fmt"

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/routes"
import "github.com/BTBurke/gaea-server/middleware"
import "github.com/BTBurke/gaea-server/errors"

import _ "github.com/lib/pq"
import "github.com/jmoiron/sqlx"

import "log"

func init() {
	reqdEnv := []string{"LE_TOKEN", "MAILGUN_API_KEY", "POSTGRES_USER", "POSTGRES_PASSWORD"}

	var envValue string
	var exit bool
	for _, envKey := range reqdEnv {
		envValue = os.Getenv(envKey)
		if len(envValue) == 0 {
			fmt.Printf("Warning: %s not set\n", envKey)
			exit = true
		}
	}
	if exit {
		os.Exit(1)
	}
}

func main() {
	r := gin.Default()
	r.Use(middleware.CORS())

	// Connect to database
	pgUser := os.Getenv("POSTGRES_USER")
	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	pgHost := os.Getenv("DB_PORT_5432_TCP_ADDR")
	if len(pgHost) == 0 {
		pgHost = "127.0.0.1"
	}
	fmt.Printf("INFO: Using %s as postgres host connection\n", pgHost)
	pgConnectString := fmt.Sprintf("host=%s user=%s password=%s dbname=db_gaea sslmode=disable", pgHost, pgUser, pgPassword)

	db, err := sqlx.Connect("postgres", pgConnectString)
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
	r.POST("/reset", routes.RequestResetEmail(db))

	auth := r.Group("/", middleware.CORS(), middleware.Auth())
	admin := r.Group("/", middleware.CORS(), middleware.Auth(), middleware.Admin())

	r.POST("/set", routes.SetPassword(db))

	auth.GET("/user", routes.GetCurrentUser(db))
	admin.GET("/users", routes.GetAllUsers(db))
	admin.POST("/users", routes.CreateUser(db))

	auth.GET("/inventory", routes.GetInventory(db))
	admin.POST("/inventory/csv", routes.CreateInventoryFromCSVString(db))
	admin.POST("/inventory", routes.CreateItem(db))
	admin.PUT("/inventory/:invID", routes.UpdateItem(db))
	admin.GET("/inventory/:invID/effects", routes.GetEffects(db))

	auth.GET("/sale", routes.GetSales(db))
	admin.POST("/sale", routes.CreateSale(db))
	admin.PUT("/sale/:saleID", routes.UpdateSale(db))
	admin.GET("/sale/:saleID/all", routes.GetAllOrdersForSale(db))

	auth.GET("/order", routes.GetOrders(db))
	auth.POST("/order", routes.CreateOrder(db))
	auth.PUT("/order/:orderID", routes.UpdateOrderStatus(db))
	auth.GET("/order/:orderID/item", routes.GetOrderItems(db))
	auth.POST("/order/:orderID/item", routes.AddOrderItem(db))
	auth.DELETE("/order/:orderID/item/:itemID", routes.DeleteOrderItem)
	auth.PUT("/order/:orderID/item/:itemID", routes.UpdateOrderItem(db))

	auth.POST("/transaction", routes.CreateTransaction(db))

	auth.GET("/announcement", routes.GetAnnouncements(db))
	admin.POST("/announcement", routes.CreateAnnouncement(db))
	admin.PUT("/announcement/:announcementID", routes.UpdateAnnouncement(db))
	admin.DELETE("/announcement/:announcementID", routes.DeleteAnnouncement(db))

	r.GET("/testauth", middleware.Auth(), routes.TestAuth)
	// When developing on c9
	r.Run(":8080")

}
