package email

import "testing"

func TestSend(t *testing.T) {
	err := Send("GAEA Test <test@guangzhouaea.org>", "Test Email", "This is a test email", "btburke@fastmail.com")
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
