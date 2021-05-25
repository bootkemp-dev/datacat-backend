package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func CheckIfUsernameExists(username string) error {
	err := db.QueryRow(`select username from users where username=$1`, username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
		return err
	}

	return fmt.Errorf("username exists")
}

func CheckIfEmailExists(email string) error {
	err := db.QueryRow(`select email from users where email=$1`, email).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
		return err
	}

	return fmt.Errorf("email exists")
}

func InsertUser(username, email, password string) error {
	stmt, err := db.Prepare(`insert into users(id, username, email, passwordHash, created, modified) values(default, $1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, email, password, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}
