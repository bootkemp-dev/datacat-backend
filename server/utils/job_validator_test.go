package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestIsUrl(t *testing.T) {
	fmt.Println(isUrl("http://google.com"))
	fmt.Println(isUrl("http://192.158.0.1"))
}

func TestValidateNewJob(t *testing.T) {
	err := ValidateNewJob("some-name", "google.com")
	if err != nil {
		log.Println(err)
		t.Fail()
	}
}
