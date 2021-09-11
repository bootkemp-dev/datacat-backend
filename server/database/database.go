package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bootkemp-dev/datacat-backend/config"
	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

func NewDatabase(c config.Config) (*Database, error) {
	db, err := connect(c)
	if err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func connect(c config.Config) (*sql.DB, error) {

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
