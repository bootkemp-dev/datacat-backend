package logger

import (
	"log"
	"os"
)

type Logger struct {
	name string
	file *os.File
}

func NewLogger(name, filepath string) (*Logger, error) {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Logger{name: name, file: file}, nil
}
