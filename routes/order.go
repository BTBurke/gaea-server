package routes

import "github.com/gin-gonic/gin"
import "time"
import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/BTBurke/gaea-server/errors"
import "github.com/shopspring/decimal"
import "database/sql"



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
	AmountTotal decimal.Decimal       `json:"amount_total"`
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



type orderItems struct {
	Qty        int         `json:"qty"`
	OrderItems []OrderItem `json:"order_items"`
}

// GET for all orders by user
func GetOrders (db *sqlx.DB) gin.HandlerFunc {

	return func (c *gin.Context) {
	
		var user1 User
		userName, exists := c.Get("user")
		if !exists {
			fmt.Println(err1)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to get user", "internal server error",c))
			return
		}
		
		dbErr := db.Get(&user1, "SELECT * FROM gaea.user WHERE user_name=$1", userName)
		if dbErr != nil {
			fmt.Println(err1)
			c.AbortWithError(503, errors.NewAPIError(503, "failed to get user", "internal server error",c))
			return
		}
		
		var memberStatus bool
		switch {
		case user1.Role == "nonmember":
			memberStatus = false
		default:
			memberStatus = true
		}
		
		var ords []Order
		var retOrds []Order
		var qtyOrd int
		
		err1 := db.Get(&qtyOrd, `SELECT COUNT(*) FROM gaea.order WHERE user_name=$1`,
			userName)
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
			
			var amtErr error
			
			for _, order := range ords {
				order.ItemQty, order.AmountTotal, amtErr = CalcOrderTotals(order.OrderId, memberStatus, db)
				if amtErr != nil {
					fmt.Printf("%s", amtErr)
				}
				retOrds = append(retOrds, order)
			}
		}
		
		c.JSON(200, gin.H{"qty": qtyOrd, "orders": retOrds})
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

func UpdateOrderStatus(db *sqlx.DB) gin.HandlerFunc {
	return func (c *gin.Context) {
		var updOrder Order
		err := c.Bind(&updOrder)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, fmt.Sprintf("failed on binding to updated order : %s", err), "internal server error",c))
			return
		}
		
		saleOpen, err := CheckOrderOpen(updOrder.OrderId, db)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, fmt.Sprintf("CheckOpenOrder returned an error : %s", err), "Internal Server Error",c))
			return
		}
		if (!saleOpen) {
			c.AbortWithError(409, errors.NewAPIError(409, "sale not open", "Cannot fufill request because the associated sale is not open.",c))
			return
		}
		
		var retOrder Order
		dbErr := db.Get(&retOrder, "UPDATE gaea.order SET status=$1, status_date=$2 WHERE order_id=$3 RETURNING *",
			updOrder.Status, time.Now(), updOrder.OrderId)
		if dbErr != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, fmt.Sprintf("failed on updating order : %s", err), "internal server error",c))
			return
		}
		
		c.JSON(200, retOrder)
		
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
		
		// check to make sure sale is still open
		saleOpen, err := CheckOrderOpen(newItem.OrderId, db)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, fmt.Sprintf("CheckOpenOrder returned an error : %s", err), "Internal Server Error",c))
			return
		}
		if (!saleOpen) {
			c.AbortWithError(409, errors.NewAPIError(409, "sale not open", "Cannot fufill request because the associated sale is not open.",c))
			return
		}
		
		
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
func UpdateOrderItem(db *sqlx.DB) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		var updateItem OrderItem
	
		err := c.Bind(&updateItem)
		if err != nil {
			fmt.Println(err)
			c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."})
			return
		}
		
		// Check order is still open before updating
		saleOpen, err := CheckOrderOpen(updateItem.OrderId, db)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, fmt.Sprintf("CheckOpenOrder returned an error : %s", err), "Internal Server Error",c))
			return
		}
		if (!saleOpen) {
			c.AbortWithError(409, errors.NewAPIError(409, "sale not open", "Cannot fufill request because the associated sale is not open.",c))
			return
		}
		
		var returnItem OrderItem
		dbErr := db.Get(&returnItem, `UPDATE gaea.orderitem SET qty=$1, updated_at=$2 WHERE
			orderitem_id=$3 RETURNING *`, updateItem.Qty, time.Now(), updateItem.OrderitemId)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating order item", "internal server error",c))
			return
		}
		
		
		c.JSON(200, returnItem)
	}
}

// CheckOrderOpen will check the status of the associated sale, returning true if the time
// is between the open and close dates, and false otherwise
func CheckOrderOpen(orderID int, db *sqlx.DB) (bool, error) {
	
		var saleId int
		dbErr1 := db.Get(&saleId, "SELECT sale_id FROM gaea.order WHERE order_id=$1", orderID);
		if dbErr1 != nil {
			return false, dbErr1
		}
	
		var assocSale Sale
		dbErr := db.Get(&assocSale, "SELECT * FROM gaea.sale WHERE sale_id = $1", saleId)
		if dbErr != nil {
			return false, dbErr
		}
		
		switch {
			case time.Now().Before(assocSale.OpenDate):
				return false, nil
			case time.Now().After(assocSale.CloseDate):
				return false, nil
			default:
				return true, nil
		}
		
		return true, nil
}

// CalcOrderTotals will calculate the total amount and quantity of items for a given order.
// Returns a fixed precision decimal.Decimal number.
func CalcOrderTotals(orderID int, member bool, db *sqlx.DB) (int, decimal.Decimal, error) {
	
	var oItems []OrderItem
	err := db.Select(&oItems, "SELECT * FROM gaea.orderitem WHERE order_id = $1", orderID)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return 0, decimal.NewFromFloat(0), nil
		default:
			return 0, decimal.NewFromFloat(0), err
		}
	}
	
	var dbErr error
	var total decimal.Decimal
	var inv Inventory
	var totalQty int
	var decQty decimal.Decimal
	for _, item := range oItems {
		dbErr = db.Get(&inv, "SELECT * FROM gaea.inventory WHERE inventory_id=$1", item.InventoryId)
		if dbErr != nil {
			return 0, decimal.NewFromFloat(0), dbErr
		}
		switch {
			case member:
				decQty = decimal.New(int64(item.Qty), 0)
				total = total.Add(inv.MemPrice.Mul(decQty))
			default:
				decQty = decimal.New(int64(item.Qty), 0)
				total = total.Add(inv.NonmemPrice.Mul(decQty))
		}
		totalQty = totalQty + item.Qty
	}
	
	return totalQty, total, nil
}