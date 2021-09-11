package logger

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bootkemp-dev/datacat-backend/config"
	_ "github.com/lib/pq"
)

type Logger struct {
	file *os.File
	db   *sql.DB
}

func NewLogger(c *config.Config) (*Logger, error) {
	file, err := os.OpenFile("./log_files/jobs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	db, err := connectToDB(*c)
	if err != nil {
		return nil, err
	}

	log.Println("Logger connected to the database")

	return &Logger{
		file: file,
		db:   db,
	}, nil

}

func connectToDB(c config.Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	log.Println(psqlInfo)

	database, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("sql.Open failed: %v\n", err)
		return nil, err
	}

	err = database.Ping()
	if err != nil {
		log.Printf("database.Ping failed: %v\n", err)
		return nil, err
	}

	return database, nil
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
