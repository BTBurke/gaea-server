package routes

import "github.com/gin-gonic/gin"
import "time"

type Announcement struct {
    Title     string `json:"title"`
    Text      string `json:"text"`
    UpdatedAt time.Time `json:"updated_at"`
    DeleteAt  time.Time `json:"delete_at"`
    AnnouncementID string `json:"announcement_id"`
}

type announcements struct {
    Qty           int            `json:"qty"`
    Announcements []Announcement `json:"announcements"`
}

func GetAnnouncements(c *gin.Context) {
    ann1 := Announcement{
        Title: "This is a test announcement",
        Text: "It can be used to announce upcoming events or items of interest to members.",
        UpdatedAt: time.Date(2015, time.June, 04, 0, 0, 0, 0, time.UTC),
        DeleteAt: time.Date(2099, time.January, 01, 0, 0, 0, 0, time.UTC),
        AnnouncementID: "7f9f333c-3a2d-46a6-8ab9-5822feeed58f",
    }
    ann2 := Announcement{
        Title: "This is another test announcement",
        Text: "It can be used to announce upcoming events or items of interest to members.",
        UpdatedAt: time.Date(2015, time.June, 04, 0, 0, 0, 0, time.UTC),
        DeleteAt: time.Date(2099, time.January, 01, 0, 0, 0, 0, time.UTC),
        AnnouncementID: "7f9f333c-3a2d-46a6-8ab9-5822feeed58e",
    }
    anns := announcements{
        Qty: 2,
        Announcements: []Announcement{ann1, ann2},
    }
    c.JSON(200, anns)
}