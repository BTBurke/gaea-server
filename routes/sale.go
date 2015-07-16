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
	SaleId    int64     `json:"sale_id" db:"sale_id"`
	Status    string    `json:"status"` //Set('open', 'closed', 'deliver', 'complete')
	Salescopy string    `json:"sales_copy" db:"salescopy"`
}

type sales struct {
	Qty   int    `json:"qty"`
	Sales []Sale `json:"sales"`
}

// GetCurrentSale returns only open sales.  Currently only for testing.
func GetCurrentSale(c *gin.Context) {

	sale1 := Sale{
		OpenDate:  time.Date(2015, time.June, 5, 0, 0, 0, 0, time.UTC),
		CloseDate: time.Date(2015, time.July, 31, 0, 0, 0, 0, time.UTC),
		SaleType:  "alcohol",
		Status:    "open",
		SaleId:    1,
		Salescopy: "This is a test of the sales copy.  It will appear in announcements.",
	}

	sale2 := Sale{
		OpenDate:  time.Date(2015, time.June, 5, 0, 0, 0, 0, time.UTC),
		CloseDate: time.Date(2015, time.December, 31, 0, 0, 0, 0, time.UTC),
		SaleType:  "merchandise",
		Status:    "complete",
		SaleId:    2,
		Salescopy: "no copy",
	}

	sales := sales{
		Qty:   2,
		Sales: []Sale{sale1, sale2},
	}

	c.JSON(200, sales)

}

// UpdateSale will update the details of a sale
func UpdateSale(c *gin.Context) {
	var update Sale

	err := c.Bind(&update)
	if err != nil {
		c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."})
		return
	}

	c.JSON(200, update)
}

// CreateSale creates a new sale in the database
func CreateSale(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSale Sale

		err := c.Bind(&newSale)
		saleID, dbErr := db.NamedExec("INSERT INTO gaea.sale (sale_id, sale_type, open_date, close_date, salescopy) VALUES (DEFAULT, :sale_type, :open_date, :close_date, :salescopy) RETURNING sale_id", &newSale)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.APIError{503, "failed on inserting a new sale", "internal server error"})
		}
		fmt.Println(saleID)
		insertID, _ := saleID.LastInsertId()
		fmt.Println("insertID", insertID)
		newSale.SaleId = insertID
		if err != nil {
			c.JSON(422, gin.H{"error": "Data provided in wrong format, unable to complete request."})
			return
		}

		c.JSON(200, newSale)
	}
}
