package utils

import (
	"fmt"
	"log"

	"github.com/bootkemp-dev/datacat-backend/models"
	"github.com/goware/emailx"
)

//Validates the request data from Register endpoint
func NewUserValidator(user models.RegisterRequest) error {
	//check length
	if len(user.Username) <= 0 || len(user.Username) > 50 {
		return fmt.Errorf("Username not valid")
	}

	//regex the username
	err := emailx.Validate(user.Email)
	if err != nil {
		if err == emailx.ErrInvalidFormat {
			return fmt.Errorf("Invalid email format")
		}
		return fmt.Errorf("Email validation failed")
	}

	log.Printf("Password 1: %s\n", user.Password1)
	log.Printf("Password 2: %s\n", user.Password2)

	if user.Password1 != user.Password2 {
		return fmt.Errorf("Passwords do not match")
	}

	return nil
}
