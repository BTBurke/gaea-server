package routes

import (
	"time"
	"fmt"

	"github.com/BTBurke/gaea-server/errors"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/zero"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionID int             `json:"transaction_id" db:"transaction_id"`
	SaleID        int             `json"sale_id" db:"sale_id"`
	OrderID       int             `json:"order_id" db:"order_id"`
	UserName      zero.String     `json:"user_name" db:"user_name"`
	From          string          `json:"from" db:"from"`
	To            string          `json:"from" db:"from"`
	Description   string          `json:"description" db:"description"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	Type          string          `json:"type" db:"type"`
	Status        string          `json:"status" db:"status"`
	Track         zero.String     `json:"track" db:"track"`
	Notes         zero.String     `json:"notes" db:"notes"`
	PayDate       time.Time       `json:"pay_date" db:"pay_date"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
	AuthorizedBy1 zero.String     `json:"authorized_by1" db:"authorized_by1"`
	AuthorizedBy2 zero.String     `json:"authorized_by2" db:"authorized_by2"`
}

func CreateTransaction(db *sqlx.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var trans Transaction
        
        err := c.Bind(&trans)
        if err != nil {
            fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to bind transaction", "internal server error",c))
			return
        }
        
        var retTrans Transaction
        dbErr := db.Get(&retTrans, `INSERT INTO gaea.transaction (transaction_id, sale_id,
            order_id, user_name, from, to, description, amount, type, status,
            track, notes, pay_date, updated_at, authorized_by1, authorized_by2) VALUES
            (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
            RETURNING *`, trans.SaleID, trans.OrderID, trans.UserName, trans.From, trans.To,
            trans.Description, trans.Amount, trans.Type, trans.Status, trans.Track,
            trans.Notes, trans.PayDate, time.Now(), trans.AuthorizedBy1, trans.AuthorizedBy2)
        if dbErr != nil {
            fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on transaction insert", "internal server error",c))
			return
        }
    }
}
