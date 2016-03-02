package routes

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Item struct {
	Item      OrderItem `json:"item"`
	Inventory Inventory `json:"inventory"`
}

type UserOrder struct {
	User  User   `json:"user"`
	Order Order  `json:"order"`
	Items []Item `json:"items"`
}

type SaleOrders struct {
	Sale   Sale        `json:"sale"`
	Orders []UserOrder `json:"orders"`
}

func GetAllOrders(db *sqlx.DB, saleID int) (SaleOrders, error) {
	orders, err := GetAllOrdersAsList(db, saleID)
	if err != nil {
		return SaleOrders{}, err
	}
	var sale1 Sale
	if err := db.Get(&sale1, "SELECT * FROM gaea.sale WHERE sale_id=$1", saleID); err != nil {
		return SaleOrders{}, err
	}
	return SaleOrders{sale1, orders}, nil
}

func GetAllOrdersAsList(db *sqlx.DB, saleID int) ([]UserOrder, error) {
	var orders []Order
	if err := db.Select(&orders, "SELECT * FROM gaea.order where sale_id=$1", saleID); err != nil {
		switch {
		case err == sql.ErrNoRows:
			return []UserOrder{}, nil
		default:
			return []UserOrder{}, err
		}
	}

	var user1 User
	var userOrders []UserOrder
	var inv Inventory
	var memberStatus bool
	var amtErr error
	for _, order := range orders {
		var allItems = []Item{}
		var oItems = []OrderItem{}

		if err := db.Get(&user1, "SELECT * FROM gaea.user WHERE user_name=$1", order.UserName); err != nil {
			return []UserOrder{}, err
		}
		if err := db.Select(&oItems, "SELECT * FROM gaea.orderitem WHERE order_id=$1", order.OrderId); err != nil {
			switch {
			case err == sql.ErrNoRows:
				oItems = []OrderItem{}
			default:
				return []UserOrder{}, err
			}
		}
		for _, item := range oItems {
			if err := db.Get(&inv, "SELECT * FROM gaea.inventory WHERE inventory_id=$1", item.InventoryId); err != nil {
				return []UserOrder{}, err
			}
			allItems = append(allItems, Item{item, inv})
		}

		// update Order in place with totals
		switch {
		case user1.Role == "nonmember":
			memberStatus = false
		default:
			memberStatus = true
		}
		order.ItemQty, order.AmountTotal, amtErr = CalcOrderTotals(order.OrderId, memberStatus, db)
		if amtErr != nil {
			return []UserOrder{}, amtErr
		}

		userOrders = append(userOrders, UserOrder{user1, order, allItems})
	}
	return userOrders, nil
}

func orderItemsAsCSVBytes(order UserOrder) ([]byte, error) {
	var csvBuffer = new(bytes.Buffer)

	var records [][]string
	var price decimal.Decimal
	records = append(records, []string{"supplier_id", "name", "qty", "case_size", "price", "currency"})
	for _, item := range order.Items {
		switch {
		case order.User.Role == "nonmember":
			price = item.Inventory.NonmemPrice
		default:
			price = item.Inventory.MemPrice
		}
		var rec = []string{item.Inventory.SupplierID, item.Inventory.Name, strconv.Itoa(item.Item.Qty), strconv.Itoa(item.Inventory.CaseSize), price.String(), item.Inventory.Currency}
		records = append(records, rec)
	}
	w := csv.NewWriter(csvBuffer)
	w.WriteAll(records)
	if err := w.Error(); err != nil {
		return []byte{}, err
	}
	return csvBuffer.Bytes(), nil
}

func AllOrdersAsCSVZip(sale SaleOrders) (string, error) {
	fName := strings.Join([]string{"sale-", strconv.Itoa(sale.Sale.SaleId), ".zip"}, "")
	f, err := os.Create(path.Join("/tmp", fName))
	if err != nil {
		return "", err
	}
	z := zip.NewWriter(f)
	for _, order := range sale.Orders {
		zipFile, err := z.Create(strings.Join([]string{order.User.UserName, ".csv"}, ""))
		csvBytes, err := orderItemsAsCSVBytes(order)
		if err != nil {
			return "", err
		}
		if _, err = zipFile.Write(csvBytes); err != nil {
			return "", err
		}
	}
	// add in all orders as JSON file
	zipFile, err := z.Create("orders.json")
	ordersAsJson, err := json.Marshal(sale)
	if err != nil {
		return "", err
	}
	if _, err := zipFile.Write(ordersAsJson); err != nil {
		return "", err
	}

	if err := z.Close(); err != nil {
		return "", err
	}
	return fName, nil
}
