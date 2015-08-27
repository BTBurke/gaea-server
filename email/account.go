package email

import "fmt"

const PasswordResetTemplate string = `
Hello {{.first_name}},

To reset your Guangzhou American Employees Association website password, copy and paste the following link into your browser:

https://guangzhouaea.org/reset/{{.link}}

If you didn't request a password reset, please contact us at help@guangzhouaea.org.

Sincerely,
Guangzhou AEA Board Members
`

func PasswordResetEmail(firstName string, pwdJwt string) (string, error) {
	data := map[string]string{
		"first_name": firstName,
		"jwt":        pwdJwt,
	}
	body, err := RenderFromTemplate(data, PasswordResetTemplate)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(body)
	return body, nil
}

const NewAccountPasswordTemplate string = `
Hello {{.first_name}},

Thanks for signing up for the Guangzhou American Employees Association website.  To get started,
please create a password by copying and pasting the following link in your browser:

{{.link}}

If you have any issues, please contact us at help@guangzhouaea.org.

Sincerely,
Guangzhou AEA Board Members
`
