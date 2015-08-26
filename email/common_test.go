package email

import "testing"

func TestRenderFromTemplate(t *testing.T) {
	data := map[string]string{
		"first_name": "test",
		"say_this":   "yo mamma",
	}
	tmpl := "Hello {{.first_name}}, {{.say_this}} just called"
	expect := "Hello test, yo mamma just called"
	result, err := RenderFromTemplate(data, tmpl)
	if err != nil {
		t.Fatalf("Got an error: %s", err)
	}
	if result != expect {
		t.Fatalf("Got: %s, Expected: %s", result, expect)
	}

}
