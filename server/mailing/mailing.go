package mailing

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/bootkemp-dev/datacat-backend/config"
)

type Mailing struct {
	client    *smtp.Client
	fromEmail string
	password  string
	smtpHost  string
	smtpPort  int
}

func NewMailing(c config.Config) (*Mailing, error) {
	client, err := connectToSMTP(c.Smtp.Host, c.Smtp.Port)
	if err != nil {
		return nil, err
	}

	return &Mailing{
		client:    client,
		fromEmail: c.Smtp.ResetEmail,
		password:  c.Smtp.Password,
		smtpHost:  c.Smtp.Host,
		smtpPort:  c.Smtp.Port,
	}, nil
}

func connectToSMTP(smtpHost string, smtpPort int) (*smtp.Client, error) {
	c, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpHost, smtpPort))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return c, nil
}
