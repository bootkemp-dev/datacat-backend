package utils

import (
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
