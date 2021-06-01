package database

import (
	"log"
	"time"
)

func InsertNewJob(name, url string, frequency uint, userid float64) (int, error) {
	stmt, err := db.Prepare(`insert into jobs(id, jobName, jobUrl, frequency, userid, created, modified) values(default, $1, $2, $3, $4, $5, $6) returning id`)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	var id int
	err = stmt.QueryRow(name, url, frequency, userid, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	return id, nil
}
