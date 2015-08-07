package routes

import "github.com/gin-gonic/gin"
import "time"
import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/errors"



// Order represents a single order transaction on behalf of a user.  It is
// associated with a sale and has a status of open, submit, or deliver.
type Order struct {
	SaleId      int    `json:"sale_id" db:"sale_id"`
	Status      string    `json:"status" db:"status"` // Set {Saved, Submit, Paid, Deliver}
	StatusDate  time.Time `json:"status_date" db:"status_date"`
	UserName    string    `json:"user_name" db:"user_name"`
	OrderId     int    `json:"order_id" db:"order_id"`
	SaleType    string    `json:"sale_type" db:"sale_type"`
	ItemQty     int       `json:"item_qty"`
	AmountTotal int       `json:"amount_total"`
}

// An OrderItem represents a single line transaction in an order for a Qty
// of an InventoryItem.
type OrderItem struct {
	Qty         int       `json:"qty" db:"qty"`
	InventoryId int    `json:"inventory_id" db:"inventory_id"`
	OrderId     int    `json:"order_id" db:"order_id"`
	OrderitemId int    `json:"orderitem_id" db:"orderitem_id"`
	UserName    string    `json:"user_name" db:"user_name"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// type orders struct {
// 	Qty    int     `json:"qty"`
// 	Orders []Order `json:"orders"`
// }

type orderItems struct {
	Qty        int         `json:"qty"`
	OrderItems []OrderItem `json:"order_items"`
}

// GET for all orders by user
func GetOrders (db *sqlx.DB) gin.HandlerFunc {

	return func (c *gin.Context) {
	
		//for testing until JWT implemented
		uName := "burkebt"
		
		var ords []Order
		var qtyOrd int
		
		err1 := db.Get(&qtyOrd, `SELECT COUNT(*) FROM gaea.order WHERE user_name=$1`,
			uName)
		if err1 != nil {
			fmt.Println(err1)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to get orders", "internal server error",c))
			return
		}
		if qtyOrd > 0 {
			err2 := db.Select(&ords, `SELECT * FROM gaea.order WHERE user_name=$1`,
				uName)
			if err2 != nil {
				fmt.Println(err2)
				c.AbortWithError(503, errors.NewAPIError(503, "failed to get orders", "internal server error",c))
				return
			}
		}
		
		c.JSON(200, gin.H{"qty": qtyOrd, "orders": ords})
	}
}
// Post to create a new order
func CreateOrder(db *sqlx.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var ord Order
		
		err := c.Bind(&ord)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to create new order", "internal server error",c))
			return
		}
		
		ord.StatusDate = time.Now()
		ord.Status = "saved"
		
		var returnID int
		dbErr := db.Get(&returnID, `INSERT INTO gaea.order
		(order_id, sale_id, status, status_date, user_name, sale_type) VALUES 
		(DEFAULT, $1, $2, $3, $4, $5) RETURNING order_id`,
			ord.SaleId,
			ord.Status,
			ord.StatusDate,
			ord.UserName,
			ord.SaleType)
		
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to bind new order", "internal server error",c))
			return
		}
		ord.OrderId = returnID
		
		c.JSON(200, ord)
	}
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
			c.AbortWithError(503, errors.NewAPIError(503, "failed on getting count of order items", "internal server error",c))
			return
		}
		if count > 0 {
			err := db.Select(&oItems, "SELECT * FROM gaea.orderitem WHERE order_id=$1", orderID)
			fmt.Println(oItems)
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(503, errors.NewAPIError(503, "failed on getting order items", "internal server error",c))
				return
			}
		}

		c.JSON(200, gin.H{"qty": count, "order_items": oItems, "query": "order-"+orderID})
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
		
		var returnID int
		dbErr := db.Get(&returnID, `INSERT INTO gaea.orderitem
			(orderitem_id, order_id, inventory_id, qty, updated_at, user_name)
			VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING orderitem_id`, newItem.OrderId, newItem.InventoryId, newItem.Qty, newItem.UpdatedAt, newItem.UserName)
		if dbErr != nil {
			fmt.Println("error on db entry")
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on inserting order items", "internal server error",c))
			return
		}
		newItem.OrderitemId = returnID
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