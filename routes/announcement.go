package routes

import (
	"database/sql"
	"fmt"

	"github.com/BTBurke/gaea-server/errors"
	"github.com/gin-gonic/gin"
)
import "time"
import "github.com/jmoiron/sqlx"

type Announcement struct {
	Title          string    `json:"title" db:"title"`
	Markdown       string    `json:"markdown" db:"markdown"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	ShowAt         time.Time `json:"show_at" db:"show_at"`
	ShowUntil      time.Time `json:"show_until" db:"show_until"`
	AnnouncementID string    `json:"announcement_id" db:"announcement_id"`
}

func GetAnnouncements(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcements []Announcement
		dbErr := db.Select(&announcements, "SELECT * FROM gaea.announcement")
		if dbErr != nil {
			switch {
			case dbErr == sql.ErrNoRows:
				c.JSON(200, gin.H{"qty": 0, "announcements": []Announcement{}})
				return
			default:
				fmt.Println(dbErr)
				c.AbortWithError(503, errors.NewAPIError(503, "failed to get announcements from database", "internal server error", c))
				return
			}
		}
		c.JSON(200, gin.H{"qty": len(announcements), "announcements": announcements})
	}
}

func CreateAnnouncement(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ann Announcement
		if err := c.Bind(&ann); err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to bind announcement on create", "internal server error", c))
			return
		}

		var retAnn Announcement
		dbErr := db.Get(&retAnn, `INSERT INTO gaea.announcement (title, markdown,
      updated_at, show_at, show_until, announcement_id) VALUES $1, $2, $3, $4, $5, DEFAULT
      RETURNING *`, ann.Title, ann.Markdown, time.Now(), ann.ShowAt, ann.ShowUntil)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on inserting announcement", "internal server error", c))
			return
		}

		c.JSON(200, retAnn)
	}
}

func UpdateAnnouncement(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ann Announcement
		if err := c.Bind(&ann); err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to bind announcement on create", "internal server error", c))
			return
		}

		var retAnn Announcement
		dbErr := db.Get(&retAnn, `UPDATE gaea.announcement SET title=$1, markdown=$2,
      updated_at=$3, show_at=$4, show_until=$5 WHERE announcement_id=$6
      RETURNING *`, ann.Title, ann.Markdown, time.Now(), ann.ShowAt, ann.ShowUntil, ann.AnnouncementID)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on updating announcement", "internal server error", c))
			return
		}

		c.JSON(200, retAnn)
	}
}

func DeleteAnnouncement(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ann Announcement
		if err := c.Bind(&ann); err != nil {
			fmt.Println(err)
			c.AbortWithError(422, errors.NewAPIError(422, "failed to bind announcement on create", "internal server error", c))
			return
		}

		dbErr := db.MustExec(`DELETE gaea.announcement WHERE announcement_id = $1`, ann.AnnouncementID)
		if dbErr != nil {
			fmt.Println(dbErr)
			c.AbortWithError(503, errors.NewAPIError(503, "failed on deleting announcement", "internal server error", c))
			return
		}
		c.JSON(200, gin.H{"announcement_id": ann.AnnouncementID})
	}
}
