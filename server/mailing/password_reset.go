package mailing

import (
	"log"
)

func SendResetPasswordEmail(username, toEmail, token string) {

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

	return nil
}
