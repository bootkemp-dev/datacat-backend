package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	_ "github.com/lib/pq"
)

type Logger struct {
	file *os.File
}

func NewLogger(c *config.Config) (*Logger, error) {
	var file *os.File
	if _, err := os.Stat(c.Logger.DirPath); os.IsNotExist(err) {
		file, err = os.Create("auth.log")
		if err != nil {
			return nil, err
		}
	} else {
		file, err = os.Open(c.Logger.DirPath)
		if err != nil {
			return nil, err
		}
	}

	log.Println("Logger connected to the database")

	return &Logger{
		file: file,
	}, nil

}

func (l *Logger) prepareMessage(message string) string {
	timeNow := time.Now()
	timeFormatted := timeNow.Format("2006-01-02 15:04:05")

	return fmt.Sprintf("[%s] %s", timeFormatted, message)
}

func (l *Logger) WriteLogToFile(message string) error {
	preparedMessage := l.prepareMessage(message)
	w := bufio.NewWriter(l.file)
	_, err := w.WriteString(preparedMessage)
	if err != nil {
		log.Println(err)
		return err
	}
	w.Flush()
	return nil
}

func (l *Logger) Close() {
	l.file.Close()
}
