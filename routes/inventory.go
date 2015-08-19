package routes

import "github.com/gin-gonic/gin"
import "github.com/BTBurke/gaea-server/errors"
import "encoding/csv"
import "time"
import "strings"
import "strconv"
import "fmt"
import "github.com/jmoiron/sqlx"
import "github.com/shopspring/decimal"
import "github.com/guregu/null/zero"

// Inventory represents a single inventory item that is associated with
// an offered sale.  Changes are recorded in the changelog.
type Inventory struct {
	InventoryID int             `json:"inventory_id" db:"inventory_id"`
	SaleID      int             `json:"sale_id" db:"sale_id"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	SupplierID  string          `json:"supplier_id" db:"supplier_id"`
	Name        string          `json:"name" db:"name"`
	Description zero.String     `json:"desc" db:"description"`
	Abv         zero.String     `json:"abv" db:"abv"`
	Size        zero.String     `json:"size" db:"size"`
	Year        zero.String     `json:"year" db:"year"`
	NonmemPrice decimal.Decimal `json:"nonmem_price" db:"nonmem_price"` // nonmember price in USD (7,2)precision
	MemPrice    decimal.Decimal `json:"mem_price" db:"mem_price"`       // member price in USD (7,2)precision
	Types       zero.String     `json:"types" db:"types"`               //String-representation of list, > delimiter
	Origin      zero.String     `json:"origin" db:"origin"`             //String-representation of list, > delimiter
	InStock		bool			`json:"in_stock" db:"in_stock"`
	Changelog   zero.String     `json:"changelog" db:"changelog"`       //String-representation of list, > delimiter
}

type csvInventory struct {
	CSV    string `json:"csv"`
	SaleId int    `json:"sale_id"`
	Header bool   `json:"header"`
}

type updateInventory struct {
	Old Inventory `json:"old"`
	New Inventory `json:"new"`
	ChangeFields []string `json:"change_fields"`
}

// loadInventoryFromCSV loads the inventory from a local or uploaded
// CSV file.  Expects CSV in the following format for columns:
// {id, name, desc, abv, size, year, nonmember, member, type, origin}
func inventoryFromCSV(csvString string, saleId int, hasHeader bool) ([]Inventory, error) {

	var out []Inventory

	reader := csv.NewReader(strings.NewReader(csvString))

	reader.FieldsPerRecord = 10

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return out, err
	}

	for idx, rec := range rawCSVdata {

		if hasHeader && idx == 0 {
			// skip header row
			continue
		}

		nonmemPrice, err := decimal.NewFromString(rec[6])
		if err != nil {
			return out, err
		}

		memPrice, err := decimal.NewFromString(rec[7])
		if err != nil {
			return out, err
		}

		var t Inventory
		t.SaleID = saleId
		t.UpdatedAt = time.Now()
		t.InventoryID = idx
		t.SupplierID = rec[0]
		t.Name = rec[1]
		t.Description = zero.StringFrom(rec[2])
		t.Abv = zero.StringFrom(rec[3])
		t.Size = zero.StringFrom(rec[4])
		t.Year = zero.StringFrom(rec[5])
		t.NonmemPrice = nonmemPrice
		t.MemPrice = memPrice
		t.Types = zero.StringFrom(rec[8])
		t.Origin = zero.StringFrom(rec[9])

		out = append(out, t)

	}

	return out, nil

}

func CreateInventoryFromCSVString(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var inv csvInventory

		err := c.Bind(&inv)

		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on binding inventory", "internal server error", c))
			return
		}

		inventory, err := inventoryFromCSV(inv.CSV, inv.SaleId, inv.Header)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to parse inventory", "internal server error", c))
			return
		}

		var dbErr error
		var invId int
		for _, inv1 := range inventory {
			dbErr = db.Get(&invId,
				`INSERT INTO gaea.inventory
				(inventory_id, sale_id, updated_at, supplier_id, name, description,
				abv, size, year, nonmem_price, mem_price, types, origin, in_stock)
				VALUES (DEFAULT, $1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10, $11, $12, $13) RETURNING inventory_id`,
				inv1.SaleID, inv1.UpdatedAt,
				inv1.SupplierID, inv1.Name, inv1.Description,
				inv1.Abv, inv1.Size, inv1.Year, inv1.NonmemPrice, inv1.MemPrice,
				inv1.Types, inv1.Origin, true)
			if dbErr != nil {
				fmt.Println(inv1)
				fmt.Println(dbErr)
				c.AbortWithError(422, errors.NewAPIError(422, "failed to insert inventory in db", "internal server error",c))
				return
			}
			inv1.InventoryID = invId
		}
		query := "sale-" + strconv.Itoa(inv.SaleId)
		c.JSON(200, gin.H{"qty": len(inventory), "inventory": inventory, "query": query})
	}
}

func GetInventory(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {

		orderQ := c.Query("order")
		saleQ := c.Query("sale")

		var queryName string
		var saleId int
		switch {
		case len(orderQ) > 0:
			queryName = "order-" + orderQ
			err := db.Get(&saleId, "SELECT sale_id FROM gaea.order WHERE order_id=$1", orderQ)
			if err != nil {
				fmt.Println(err)
				c.AbortWithError(422, errors.NewAPIError(422, "order ID does not exist", "order ID does not exist", c))
				return
			}
		case len(saleQ) > 0:
			queryName = "sale-" + saleQ
			saleId, _ = strconv.Atoi(saleQ)
		default:
			c.AbortWithError(422, errors.NewAPIError(422, "sale or order ID does not exist in query string", "sale or order ID does not exist in query string", c))
			return
		}

		var inv []Inventory
		dbErr := db.Select(&inv, "SELECT * FROM gaea.inventory WHERE sale_id=$1", saleId)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(422, errors.NewAPIError(422, "sale ID does not exist", "sale ID does not exist", c))
			return
		}
		c.JSON(200, gin.H{"inventory": inv, "query": queryName, "qty": len(inv)})
	}
}

func CreateItem(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var inv Inventory
		
		err := c.Bind(&inv)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on binding inventory", "internal server error", c))
			return
		}
		
		var invId int
		dbErr := db.Get(&invId,
				`INSERT INTO gaea.inventory
				(inventory_id, sale_id, updated_at, supplier_id, name, description,
				abv, size, year, nonmem_price, mem_price, types, origin, in_stock)
				VALUES (DEFAULT, $1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10, $11, $12, $13) RETURNING inventory_id`,
				inv.SaleID, time.Now(),
				inv.SupplierID, inv.Name, inv.Description,
				inv.Abv, inv.Size, inv.Year, inv.NonmemPrice, inv.MemPrice,
				inv.Types, inv.Origin, inv.InStock)
		if dbErr != nil {
			fmt.Println(inv)
			fmt.Println(dbErr)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to insert inventory in db", "internal server error",c))
			return
		}
		inv.InventoryID = invId
		
		saleId := strconv.Itoa(inv.SaleID)
		c.JSON(200, gin.H{"inventory": inv, "query": "sale-"+saleId})
	}
}

func UpdateItem(db *sqlx.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var chg updateInventory
		
		err := c.Bind(&chg)
		if err != nil {
			fmt.Println(err)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on binding inventory", "internal server error", c))
			return
		}
		
		var inv = chg.New
		var invResult Inventory
		dbErr := db.Get(&invResult,
				`UPDATE gaea.inventory SET sale_id=$1, updated_at=$2, supplier_id=$3, name=$4, description=$5,
				abv=$6, size=$7, year=$8, nonmem_price=$9, mem_price=$10, types=$11, origin=$12, in_stock=$13, changelog=$14
				WHERE inventory_id=$15 RETURNING *`,
				inv.SaleID, time.Now(),
				inv.SupplierID, inv.Name, inv.Description,
				inv.Abv, inv.Size, inv.Year, inv.NonmemPrice, inv.MemPrice,
				inv.Types, inv.Origin, inv.InStock, inv.Changelog, inv.InventoryID)
		if dbErr != nil {
			fmt.Println(inv)
			fmt.Println(dbErr)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to update inventory in db", "internal server error",c))
			return
		}
		saleId := strconv.Itoa(inv.SaleID)
		c.JSON(200, gin.H{"inventory": inv, "query": "sale-"+saleId})
	}
}