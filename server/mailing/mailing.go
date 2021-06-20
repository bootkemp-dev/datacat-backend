package mailing

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/bootkemp-dev/datacat-backend/config"
)

var smtpHost string
var smtpPort int
var password string
var fromEmail string
var baseURL string

func init() {
	config, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	smtpHost = config.Smtp.Host
	smtpPort = config.Smtp.Port
	password = config.Smtp.Password
	fromEmail = config.Smtp.ResetEmail
	baseURL = config.Server.SiteURL
}

func connectToSMTP() (*smtp.Client, error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return c, nil
}
