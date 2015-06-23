package routes

import "github.com/gin-gonic/gin"
import "time"
import "github.com/satori/go.uuid"
import "fmt"

// Order represents a single order transaction on behalf of a user.  It is
// associated with a sale and has a status of open, submit, or deliver.
type Order struct {
	SaleID      string    `json:"sale_id"`
	Status      string    `json:"status"` // Set {Saved, Submit, Paid, Deliver, Complete}
	StatusDate  time.Time `json:"status_date"`
	UserName    string    `json:"user_name"`
	UserID      string    `json:"user_id"`
	OrderID     string    `json:"order_id"`
	SaleType    string    `json:"sale_type"`
	ItemQty     int       `json:"item_qty"`
	AmountTotal int       `json:"amount_total"`
}

// An OrderItem represents a single line transaction in an order for a Qty
// of an InventoryItem.
type OrderItem struct {
	Qty         int       `json:"qty"`
	InventoryID string    `json:"inventory_id"`
	OrderID     string    `json:"order_id"`
	OrderItemID uuid.UUID    `json:"order_item_id"`
	UserID      string    `json:"user_id"`
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
		SaleID:      "266c4743-6ad6-47b3-b2e4-3d1c4f24a35e",
		Status:      "saved",
		StatusDate:  time.Date(2015, time.June, 6, 0, 0, 0, 0, time.UTC),
		UserName:    "ambassadorjs",
		UserID:      "06c0eb9f-92c5-485f-9622-c3f225eb6a95",
		OrderID:     "ambassadorjs-001",
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
func GetOrderItems(c *gin.Context) {
	
	orderID := c.Param("orderID")
	fmt.Sprintf("I'm getting the order items for %s", orderID)
	oItems := orderItems{}
	c.JSON(200, oItems)
}

// POST a new order
func AddOrderItem(c *gin.Context) {
	var newItem OrderItem
	
	
	err := c.Bind(&newItem)
	if err != nil {
		fmt.Println(err)
		c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."}) 
	}
	
	
	newItem.UpdatedAt = time.Now()
	newItem.OrderItemID	= uuid.NewV4()
	
	c.JSON(200, newItem)
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