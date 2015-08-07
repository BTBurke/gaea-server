package routes

import (
	"fmt"

	"github.com/BTBurke/gaea-server/errors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)
import "time"

// Sale represents items offered for sale by GAEA with an opening and closing
// time.  The inventory key is a foreign key used to query from the inventory
// table, representing items for sale.
type Sale struct {
	OpenDate  time.Time `json:"open_date" db:"open_date"`
	CloseDate time.Time `json:"close_date" db:"close_date"`
	SaleType  string    `json:"sale_type" db:"sale_type"` //Set{'alcohol', 'merchandise'}
	SaleId    int       `json:"sale_id" db:"sale_id"`
	Status    string    `json:"status" db:"status"` //Set('open', 'closed', 'deliver', 'complete')
	Salescopy string    `json:"sales_copy" db:"salescopy"`
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
		case time.Now().After(sale.CloseDate):
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
