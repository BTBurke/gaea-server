package routes

import "github.com/gin-gonic/gin"
import "time"

// Sale represents items offered for sale by GAEA with an opening and closing
// time.  The inventory key is a foreign key used to query from the inventory
// table, representing items for sale.
type Sale struct {
	Inventory string    `json:"inventory"` //UUID pointing to inventory items
	OpenDate  time.Time `json:"open_date"`
	CloseDate time.Time `json:"close_date"`
	SaleType  string    `json:"sale_type"` //Set{'alcohol', 'merchandise'}
}

type sales struct {
	Qty   int    `json:"qty"`
	Sales []Sale `json:"sales"`
}

// GetCurrentSale returns only open sales.  Currently only for testing.
func GetCurrentSale(c *gin.Context) {
	sale1 := Sale{
		Inventory: "test-inventory-string",
		OpenDate:  time.Date(2015, time.June, 5, 0, 0, 0, 0, time.UTC),
		CloseDate: time.Date(2015, time.July, 31, 0, 0, 0, 0, time.UTC),
		SaleType:  "alcohol",
	}

	sale2 := Sale{
		Inventory: "test-inventory-string2",
		OpenDate:  time.Date(2015, time.June, 5, 0, 0, 0, 0, time.UTC),
		CloseDate: time.Date(2015, time.December, 31, 0, 0, 0, 0, time.UTC),
		SaleType:  "merchandise",
	}

	sales := sales{
		Qty:   2,
		Sales: []Sale{sale1, sale2},
	}

	c.JSON(200, sales)

}
