package mailing

import (
	"fmt"
	"log"
	"net/smtp"
)

func SendResetPasswordEmail(username, toEmail, token string) {
	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)
	toList := []string{toEmail}
	msg := fmt.Sprintf("To reset your password follow this url: %s", "https://datacat.pl/accounts/reset_password?username="+username+"&token="+token)
	body := []byte(msg)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpHost, smtpPort), auth, fromEmail, toList, body)
	if err != nil {
		log.Println(err)
	}
}
