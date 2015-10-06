package routes

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type UserOrder struct {
	User  User        `json:"user"`
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}

type SaleOrders struct {
	Sale   Sale        `json:"sale"`
	Orders []UserOrder `json:"orders"`
}

func GetAllOrders(db *sqlx.DB, saleID int) ([]UserOrder, error) {
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
	var oItems []OrderItem
	var userOrders []UserOrder
	for _, order := range orders {
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
			userOrders = append(userOrders, UserOrder{user1, order, oItems})
		}
	}
	return userOrders, nil
}
