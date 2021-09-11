package logger

import (
	"log"
	"testing"

	"github.com/bootkemp-dev/datacat-backend/config"
)

func TestLogger(t *testing.T) {
	config, err := config.NewConfig("./../config.yml")
	if err != nil {
		t.Fail()
	}
	log.Println("Config loaded...")
	log.Println("Log files dir:", config.Logger.DirPath)

	logger, err := NewLogger(config)
	if err != nil {
		t.Fail()
	}

	defer logger.Close()

	message := "some test message"

	err = logger.WriteLogToFile(message)
	if err != nil {
		log.Println(err)
	}
}
