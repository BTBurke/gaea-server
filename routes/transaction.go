package routes

import (
	"time"
	"fmt"
	"database/sql"

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

func GetTransactionsByOrderId(orderID int, db *sqlx.DB) ([]Transaction, error) {
	var trans []Transaction
	if dbErr := db.Select(&trans, "SELECT * FROM gaea.transaction WHERE order_id=$1", orderID); dbErr != nil {
		switch {
			case dbErr == sql.ErrNoRows:
				return []Transaction{}, nil
			default:
				return []Transaction{}, dbErr
		}
	}
	return trans, nil
}

func GetTransactionsBySaleId(saleID int, db *sqlx.DB) ([]Transaction, error) {
	var trans []Transaction
	if dbErr := db.Select(&trans, "SELECT * FROM gaea.transaction WHERE sale_id=$1", saleID); dbErr != nil {
		switch {
			case dbErr == sql.ErrNoRows:
				return []Transaction{}, nil
			default:
				return []Transaction{}, dbErr
		}
	}
	return trans, nil
}

func GetTransactionsByUser(user string, db *sqlx.DB) ([]Transaction, error) {
	var trans []Transaction
	if dbErr := db.Select(&trans, "SELECT * FROM gaea.transaction WHERE user_name=$1", user); dbErr != nil {
		switch {
			case dbErr == sql.ErrNoRows:
				return []Transaction{}, nil
			default:
				return []Transaction{}, dbErr
		}
	}
	return trans, nil
}

func UpdateTransactionAsPaid(trans Transaction, db *sqlx.DB) (Transaction, error) {
	var updTrans Transaction
	if err := db.Get(&updTrans, `UPDATE gaea.transaction SET status='paid', pay_date=$1, notes=$2, track=$3, 
	updated_at=$4, authorized_by1=$5 WHERE transaction_id=$6 RETURNING *`,
	time.Now(),
	trans.Notes,
	trans.Track,
	time.Now(),
	trans.AuthorizedBy1,
	trans.TransactionID); err != nil {
		return trans, err
	}
	return updTrans, nil
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
