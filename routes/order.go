package routes

import "github.com/gin-gonic/gin"
import "time"
import "github.com/satori/go.uuid"
import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/error"


// Order represents a single order transaction on behalf of a user.  It is
// associated with a sale and has a status of open, submit, or deliver.
type Order struct {
	SaleId      string    `json:"sale_id"`
	Status      string    `json:"status"` // Set {Saved, Submit, Paid, Deliver, Complete}
	StatusDate  time.Time `json:"status_date"`
	UserName    string    `json:"user_name"`
	OrderId     string    `json:"order_id"`
	SaleType    string    `json:"sale_type"`
	ItemQty     int       `json:"item_qty"`
	AmountTotal int       `json:"amount_total"`
}

// An OrderItem represents a single line transaction in an order for a Qty
// of an InventoryItem.
type OrderItem struct {
	Qty         int       `json:"qty"`
	InventoryId string    `json:"inventory_id"`
	OrderId     string    `json:"order_id"`
	OrderitemId uuid.UUID    `json:"order_item_id,db:"orderitem_id"`
	UserName    string    `json:"user_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type orders struct {
	Qty    int     `json:"qty"`
	Orders []Order `json:"orders"`
}

type orderItems struct {
	Qty        int         `json:"qty"`
	OrderItems []OrderItem `json:"order_items"`
}

// GET for all orders by user
func GetOrders(c *gin.Context) {
	oOrder := Order{
		SaleId:      "266c4743-6ad6-47b3-b2e4-3d1c4f24a35e",
		Status:      "saved",
		StatusDate:  time.Date(2015, time.June, 6, 0, 0, 0, 0, time.UTC),
		UserName:    "ambassadorjs",
		OrderId:     "ambassadorjs-001",
		SaleType:    "alcohol",
		ItemQty:     0,
		AmountTotal: 0,
	}
	ordersR := orders{
		Qty:    1,
		Orders: []Order{oOrder},
	}
	c.JSON(200, ordersR)
}

// GET ordered items for a particular OrderID
func GetOrderItems(db *sqlx.DB) gin.HandlerFunc {
	
	return func (c *gin.Context) {
		orderID := c.Param("orderID")
		fmt.Sprintf("I'm getting the order items for %s", orderID)
		var oItems []OrderItem
		var count int
		
		countErr := db.Get(&count, "SELECT COUNT(*) from gaea.orderitem WHERE order_id=$1", orderID)
		if countErr != nil {
			fmt.Println(countErr)
			c.AbortWithError(503, error.APIError{503, "failed on getting count of order items", "internal server error"})
			return
		}
		if count > 0 {
			err := db.Select(&oItems, "SELECT * FROM gaea.orderitem WHERE order_id=$1", orderID)
			fmt.Println(oItems)
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(503, error.APIError{503, "failed on getting order items", "internal server error"})
				return
			}
		}
		out := orderItems{
				Qty: count,
				OrderItems: oItems,
		}
		c.JSON(200, out)
	}
}

// POST a new order
func AddOrderItem(db *sqlx.DB) gin.HandlerFunc {
	
	return func (c *gin.Context) {
		var newItem OrderItem
		
		err := c.Bind(&newItem)
		if err != nil {
			fmt.Println(err)
			c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."}) 
			return
		}
		
		newItem.UpdatedAt = time.Now()
		newItem.OrderitemId	= uuid.NewV4()
		
		_, dbErr := db.NamedExec(`INSERT INTO gaea.orderitem
			(orderitem_id, order_id, inventory_id, qty, updated_at, user_name)
			VALUES (:orderitem_id, :order_id, :inventory_id, :qty, :updated_at, :user_name`, &newItem)
		if dbErr != nil {
			c.AbortWithError(503, error.APIError{503, "failed on inserting order items", "internal server error"})
			return
		}
		c.JSON(200, newItem)
	}
}

// DELETE an existing order item
func DeleteOrderItem(c *gin.Context) {
	orderItemID := c.Param("itemID")
	
	
	c.JSON(200, gin.H{"order_item_id": orderItemID})
	
}

// PUT update an existing order item
func UpdateOrderItem(c *gin.Context) {
	var updateItem OrderItem
	
	err := c.Bind(&updateItem)
	if err != nil {
		fmt.Println(err)
		c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."}) 
	}
	
	c.JSON(200, updateItem)
}