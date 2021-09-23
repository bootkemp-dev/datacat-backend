package database

import (
	"log"
	"time"
)

func (db *Database) CheckIfUsernameExists(username string) error {
	err := db.QueryRow(`select username from users where username=$1`, username).Scan(&username)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *Database) CheckIfEmailExists(email string) error {
	err := db.QueryRow(`select email from users where email=$1`, email).Scan(&email)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (db *Database) InsertUser(username, email, password string) error {
	stmt, err := db.Prepare(`insert into users(id, username, email, passwordHash, created, modified, confirmed) values(default, $1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, email, password, time.Now(), time.Now(), false)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetIDAndPasswordHash(username string) (string, int, error) {
	var id int
	var password string

	err := db.QueryRow(`select id, passwordHash from users where username=$1`, username).Scan(&id, &password)
	if err != nil {
		return "", 0, err
	}

	return password, id, nil
}

func (db *Database) GetUserEmail(username string) (string, error) {
	var email string

	err := db.QueryRow(`select email from users where username=$1`, username).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}

func (db *Database) UpdateResetPasswordToken(username string, token string, expiration time.Time) error {
	stmt, err := db.Prepare(`update users SET passwordResetToken=$1, passwordResetTokenExpDate=$2, modified=$3 where username=$4`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(token, expiration, time.Now(), username)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetResetPasswordTokenExpiration(username, token string) (*time.Time, error) {
	var exp time.Time
	err := db.QueryRow(`select passwordResetTokenExpDate from users where username=$1 and passwordResetToken=$2`, username, token).Scan(&exp)
	if err != nil {
		return nil, err
	}

	return &exp, nil

}

func (db *Database) UpdatePasswordHash(username, passwordHash string) error {
	stmt, err := db.Prepare(`update users set passwordHash = $1 where username=$2`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(passwordHash, username)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
