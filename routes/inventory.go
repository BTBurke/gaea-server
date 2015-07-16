package routes

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/errors"
import "encoding/csv"
import "time"
import "os"
import "strings"
import "strconv"
import "fmt"

// Inventory represents a single inventory item that is associated with
// an offered sale.  Changes are recorded in the changelog.
type Inventory struct {
	InventoryID int       `json:"inventory_id"`
	SaleID      int       `json:"sale_id"`
	UpdatedAt   time.Time `json:"updated_at"`
	SupplierID  string    `json:"supplier_id"`
	Name        string    `json:"name"`
	Description string    `json:"desc"`
	Abv         string    `json:"abv"`
	Size        string    `json:"size"`
	Year        string    `json:"year"`
	NonmemPrice int       `json:"nonmem_price"` // nonmember price in RMB (int)
	MemPrice    int       `json:"mem_price"`    // member price in RMB (int)
	Types       []string  `json:"types"`
	Origin      []string  `json:"origin"`
	Changelog   []string  `json:"changelog"`
}

// loadInventoryFromCSV loads the inventory from a local or uploaded
// CSV file.  Expects CSV in the following format for columns:
// {id, name, desc, abv, size, year, nonmember, member, type, origin}
func loadInventoryFromCSV(fname string, saleId int) ([]Inventory, error) {

	var out []Inventory

	csvfile, err := os.Open(fname)
	if err != nil {
		return out, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = 10

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return out, err
	}

	for idx, rec := range rawCSVdata {

		if idx == 0 {
			// TODO: add check here for format on headers
			continue
		}

		nonmemPrice, err := strconv.Atoi(rec[6])
		if err != nil {
			return out, err
		}

		memPrice, err := strconv.Atoi(rec[7])
		if err != nil {
			return out, err
		}

		var t Inventory
		t.SaleID = saleId
		t.UpdatedAt = time.Now()
		t.InventoryID = idx
		t.SupplierID = rec[0]
		t.Name = rec[1]
		t.Description = rec[2]
		t.Abv = rec[3]
		t.Size = rec[4]
		t.Year = rec[5]
		t.NonmemPrice = nonmemPrice
		t.MemPrice = memPrice
		t.Types = strings.Split(rec[8], ">")
		t.Origin = strings.Split(rec[9], ">")

		out = append(out, t)

	}

	return out, nil

}

func GetInventory(c *gin.Context) {
	//TODO: Normally, this should take the SaleID as the search parameter
	// For testing, load from test fixture file

	inventory, err := loadInventoryFromCSV("./test/inventory.csv", 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	orderQ := c.Query("order")
	saleQ := c.Query("sale")
	
	var queryName string
	switch {
		case len(orderQ) > 0:
			// TODO: need query that uses order number to look up saleID
			queryName = "order-" + orderQ
		case len(saleQ) > 0:
			queryName = "sale-" + saleQ
		default:
			c.AbortWithError(422, errors.APIError{422, "sale or order ID does not exist in query string", "sale or order ID does not exist in query string"})
			return
	}
	
	c.JSON(200, gin.H{"inventory": inventory, "query": queryName, "qty": len(inventory)})
}
