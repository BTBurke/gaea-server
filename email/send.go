package email

import (
	"fmt"
	"os"

	"github.com/BTBurke/gaea-server/log"
	mg "github.com/mailgun/mailgun-go"
)

func Send(from string, subj string, body string, to ...string) {

	mgApiKey := os.Getenv("MAILGUN_API_KEY")
	if len(mgApiKey) == 0 {
		fmt.Println("Warning: No Mailgun API keys set.  No emails will be sent.")
		return
	}
	mailgun := mg.NewMailgun("mg.guangzhouaea.org", mgApiKey, "")
	email := mailgun.NewMessage(from, subj, body, to...)
	_, _, err := mailgun.Send(email)
	if err != nil {
		log.Error("msg=failed to send email to=%s subj=%s err=%s", to, subj, err)
		return
	}
	return
}
