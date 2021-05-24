package main

import (
	"log"

	"github.com/bootkemp-dev/datacat-backend/config"
	"github.com/bootkemp-dev/datacat-backend/router"
)

func main() {
	c, err := config.NewConfig("./config.yml")
	if err != nil {
		log.Fatalf("loading config file failed: %v\n", err)
	}

	router.Run(c.Server.Port)
}
