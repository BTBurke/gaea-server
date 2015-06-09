package routes

import "github.com/gin-gonic/gin"
import "time"

// Order represents a single order transaction on behalf of a user.  It is
// associated with a sale and has a status of open, submit, or deliver.
type Order struct {
    SaleID    string    `json:"sale_id"`
    Status    string    `json:"status"` // Set {Open, Submit, Deliver}
    StatusDate time.Time `json:"status_date"`
    UserName  string    `json:"user_name"`
    UserID    string    `json:"user_id"`
    OrderID   string    `json:"order_id"`
    SaleType  string    `json:"sale_type"`
}

// An OrderItem represents a single line transaction in an order for a Qty
// of an InventoryItem.  
type OrderItem struct {
    Qty             int    `json:"qty"`
    InventoryID     string `json:"inventory_id"`
    OrderID         string `json:"order_id"`
    OrderItemID     string `json:"order_item_id"`
    UserID          string `json:"user_id"`
    UpdatedAt       time.Time `json:"updated_at"`
}

type orders struct {
    Qty int `json:"qty"`
    Orders []Order `json:"orders"`
}

type orderItems struct {
    Qty        int         `json:"qty"`
    Order      Order       `json:"order"`
    OrderItems []OrderItem `json:"order_items"`
}



// GET for all orders by user
func GetOrders(c *gin.Context) {
    oOrder := Order{
        SaleID: "266c4743-6ad6-47b3-b2e4-3d1c4f24a35e",
        Status: "open",
        StatusDate: time.Date(2015, time.June, 6, 0, 0, 0, 0, time.UTC),
        UserName: "ambassadorjs",
        UserID: "06c0eb9f-92c5-485f-9622-c3f225eb6a95",
        OrderID: "90ae0ec6-945b-4d27-8cab-a1fcaef25b65",
        SaleType: "alcohol",
    }
    ordersR := orders{
        Qty: 1,
        Orders: []Order{oOrder},
    }
    c.JSON(200, ordersR)
}