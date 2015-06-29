package routes

import "github.com/gin-gonic/gin"
import "time"
import "github.com/satori/go.uuid"

// Sale represents items offered for sale by GAEA with an opening and closing
// time.  The inventory key is a foreign key used to query from the inventory
// table, representing items for sale.
type Sale struct {
	OpenDate  time.Time `json:"open_date"`
	CloseDate time.Time `json:"close_date"`
	SaleType  string    `json:"sale_type"` //Set{'alcohol', 'merchandise'}
	SaleId    uuid.UUID `json:"sale_id"`
	Status    string    `json:"status"` //Set('open', 'closed', 'deliver', 'complete')
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
		Status: "open",
		SaleId: uuid.NewV4(),
	}

	sale2 := Sale{
		OpenDate:  time.Date(2015, time.June, 5, 0, 0, 0, 0, time.UTC),
		CloseDate: time.Date(2015, time.December, 31, 0, 0, 0, 0, time.UTC),
		SaleType:  "merchandise",
		Status: "open",
		SaleId: uuid.NewV4(),
	}

	sales := sales{
		Qty:   2,
		Sales: []Sale{sale1, sale2},
	}

	c.JSON(200, sales)

}
