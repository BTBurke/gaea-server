package email

import (
	"bytes"
	"fmt"
	"text/template"
)

func RenderFromTemplate(data map[string]string, tmpl string) (string, error) {
	var out bytes.Buffer
	tmp, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	renderErr := tmp.Execute(&out, data)
	if renderErr != nil {
		fmt.Println(err)
		return "", err
	}
	return out.String(), nil
}
