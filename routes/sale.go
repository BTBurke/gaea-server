package routes

import (
	"fmt"
	"database/sql"
	"strconv"
	"time"

	"github.com/BTBurke/gaea-server/errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

// Sale represents items offered for sale by GAEA with an opening and closing
// time.  The inventory key is a foreign key used to query from the inventory
// table, representing items for sale.
type Sale struct {
	OpenDate  time.Time `json:"open_date" db:"open_date"`
	CloseDate time.Time `json:"close_date" db:"close_date"`
	SaleType  string    `json:"sale_type" db:"sale_type"` //Set{'alcohol', 'merchandise'}
	SaleId    int       `json:"sale_id" db:"sale_id"`
	Status    string    `json:"status" db:"status"` //Set('open', 'closed', 'final', 'deliver', 'complete')
	Salescopy string    `json:"sales_copy" db:"salescopy"`
	RequireFinal bool `json:"require_final" db:"require_final"` //if true, sale requires manual shift to status=final before associated transactions are set to paid
}

// updateSaleStatus returns a new array of sales with updated status and
// also updates their current status in the DB
func updateSaleStatus(db *sqlx.DB, sales []Sale) ([]Sale, error) {

	var retSales []Sale
	for _, sale := range sales {
		switch {
		case sale.Status == "complete":
			retSales = append(retSales, sale)
			continue
		case time.Now().Before(sale.CloseDate):
			// handle corner case hack where sale closes then is extended
			// TODO: could update database here to reflect back to open, but hack works fine
			sale.Status = "open"
			retSales = append(retSales, sale)
			continue
		case time.Now().After(sale.CloseDate) && sale.Status == "final":
			fmt.Println("found someone to update")
			var deliveredOrders int
			var openOrders int
			err := db.Get(&deliveredOrders,
				"SELECT COUNT(*) FROM gaea.order WHERE sale_id = $1 AND status = 'deliver'",
				sale.SaleId)
			if err != nil {
				return []Sale{}, err
			}
			err2 := db.Get(&openOrders,
				"SELECT COUNT(*) FROM gaea.order WHERE sale_id = $1 AND status != 'deliver'",
				sale.SaleId)
			if err2 != nil {
				return []Sale{}, err2
			}
			var newStatus string
			switch {
			case deliveredOrders > 0 && openOrders > 0:
				newStatus = "deliver"
			case deliveredOrders > 0 && openOrders == 0:
				newStatus = "complete"
			case deliveredOrders == 0 && openOrders > 0:
				newStatus = "closed"
			case deliveredOrders == 0 && openOrders == 0:
				// unlikely event that sale opens/closes but
				// with no orders
				newStatus = "complete"
			}
			if newStatus != sale.Status {
				err3 := db.Get(&sale.Status,
					"UPDATE gaea.sale SET status = $1 WHERE sale_id = $2 RETURNING status",
					newStatus,
					sale.SaleId)
				if err3 != nil {
					return []Sale{}, err3
				}
			}
			retSales = append(retSales, sale)
		}
	}
	return retSales, nil
}

// GetCurrentSale returns only open sales.  Currently only for testing.
func GetSales(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		sales := []Sale{}
		err := db.Select(&sales, "SELECT * FROM gaea.sale")
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on getting sales", "internal server error",c))
			return
		}
		updatedSales, err := updateSaleStatus(db, sales)
		if err != nil {
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating sales", "internal server error",c))
			return
		}
		c.JSON(200, gin.H{"qty": len(updatedSales), "sales": updatedSales})
	}
}

// UpdateSale will update the details of a sale
func UpdateSale(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update Sale

		err := c.Bind(&update)
		if err != nil {
			c.AbortWithError(422, errors.NewAPIError(422, "format wrong on sale update", "internal server error",c))
			return
		}

		var updatedSale Sale
		err2 := db.Get(&updatedSale,
			"UPDATE gaea.sale SET open_date=$1, close_date=$2, salescopy=$3 WHERE sale_id=$4 RETURNING *",
			update.OpenDate,
			update.CloseDate,
			update.Salescopy,
			update.SaleId)
		if err2 != nil {
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating sale", "internal server error",c))
			return
		}

		c.JSON(200, updatedSale)
	}
}

// CreateSale creates a new sale in the database
func CreateSale(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSale Sale

		err := c.Bind(&newSale)
		if err != nil {
			c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."})
			return
		}
		var retSale Sale
		dbErr := db.Get(&retSale,
			"INSERT INTO gaea.sale (sale_id, sale_type, open_date, close_date, status, salescopy) VALUES (DEFAULT, $1, $2, $3, $4, $5) RETURNING *",
			newSale.SaleType,
			newSale.OpenDate,
			newSale.CloseDate,
			"open",
			newSale.Salescopy)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on inserting a new sale", "internal server error",c))
			return
		}
		c.JSON(200, retSale)
	}
}

// GetAllOrdersForSale returns a data structure consisting of all orders, users, and items
// associated with the sale
func GetAllOrdersForSale(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		saleIDAsString := c.Param("saleID")
		saleID, err := strconv.Atoi(saleIDAsString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to get sale ID", "unknown sale ID",c))
			return
		}
		
		var allOrdersBeforeTotals []Order
		orderErr := db.Select(&allOrdersBeforeTotals, "SELECT * FROM gaea.order WHERE sale_id=$1", saleID)
		if orderErr != nil {
			switch {
				case orderErr == sql.ErrNoRows:
					c.JSON(200, gin.H{"sale_id": saleID, "users": []User{}, "orders": []Order{}, "items": []OrderItem{}})
					return
				default:
					fmt.Println(err)
					c.AbortWithError(503, errors.NewAPIError(503, "failed to get orders for sale ID", "internal server error",c))
					return
			}
		}
		
		var allOrders []Order
		var allOrderItems []OrderItem
		var allUsers []User
		var user1 User
		var items1 []OrderItem
		var qty int
		var total decimal.Decimal
		var dbErr, calcErr error
		var isMember bool
		for _, order := range allOrdersBeforeTotals {
			dbErr = db.Get(&user1, "SELECT * FROM gaea.user WHERE user_name=$1", order.UserName)
			if dbErr != nil {
				fmt.Println(err)
				c.AbortWithError(503, errors.NewAPIError(503, "failed to get user for order", "internal server error",c))
				return
			}
						
			allUsers = append(allUsers, user1)
			
			dbErr = db.Select(&items1, "SELECT * FROM gaea.orderitem WHERE order_id=$1", order.OrderId)
			if dbErr != nil {
				switch {
					case dbErr == sql.ErrNoRows:
						items1 = []OrderItem{}
						continue
					default:
						fmt.Println(dbErr)
						c.AbortWithError(503, errors.NewAPIError(503, "failed to get user for order", "internal server error",c))
						return
				}
			}
			
			allOrderItems = append(allOrderItems, items1...)
			
			switch {
				case user1.Role == "nonmember":
					isMember = false
					continue
				default:
					isMember = true
			}
			qty, total, calcErr = CalcOrderTotals(order.OrderId, isMember, db)
			if calcErr != nil {
				fmt.Println(calcErr)
				c.AbortWithError(503, errors.NewAPIError(503, "failed to get user for order", "internal server error",c))
				return
			}
			
			order.ItemQty = qty
			order.AmountTotal = total
			
			allOrders = append(allOrders, order)
			
		}
		
		c.JSON(200, gin.H{"sale_id": saleID, "orders": allOrders, "items": allOrderItems, "users": allUsers})
		
	}
}