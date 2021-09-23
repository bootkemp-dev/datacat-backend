package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword("lol")
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	log.Println(hashedPassword)
}

func TestGenerateRandomToken(t *testing.T) {
	token, err := GenerateRandomToken(30)
	if err != nil {
		t.Fail()
	}

	fmt.Println(token)

	if len(token) != 30 {
		t.Fail()
	}
}
