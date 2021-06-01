package utils

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
)

func ValidateNewJob(name string, url string) error {
	if len(name) <= 0 || len(name) > 50 {
		return fmt.Errorf("Name is not valid")
	}

	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		url = fmt.Sprintf("http://%s", url)
	}

	if !isUrl(url) {
		return fmt.Errorf("Url is not valid")
	}

	return nil
}

func isUrl(str string) bool {
	url, err := url.ParseRequestURI(str)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	address := net.ParseIP(url.Host)

	if address == nil {
		return strings.Contains(url.Host, ".")
	}

	return true
}
