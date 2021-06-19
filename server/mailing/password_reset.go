package mailing

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendResetPasswordEmail(username, toEmail, token string) {

}

func send(from, to, subject string) error {

	// prepare email header
	header := make(map[string]string)
	header["To"] = to
	header["From"] = from
	header["Subject"] = subject

	c, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		log.Println(err)
		return err
	}

	defer c.Close()

	return nil
}
