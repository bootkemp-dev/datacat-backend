package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
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
