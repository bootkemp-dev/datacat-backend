package mailing

import (
	"fmt"
	"html/template"
	"log"
)

func SendResetPasswordEmail(username, toEmail, token string) {
	//prepare template
	t, err := template.ParseFiles("./templates/reset_password_email.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Username  string
		ResetLink string
	}{
		Username:  username,
		ResetLink: fmt.Sprintf("%s/change_passoword?username=%s?token=%s", baseURL, username, token),
	}

	t.Parse(data)
}

func send(from, to, subject string) error {
	client, err := connectToSMTP()
	if err != nil {
		log.Println(err)
		return err
	}

	defer client.Close()

	client.Mail(from)
	client.Rcpt(to)
}
