package database

import (
	"log"
	"time"
)

func InsertNewJob(name, url string, frequency int64, userid float64) (int, error) {
	stmt, err := db.Prepare(`insert into jobs(id, jobName, jobUrl, frequency, userid, created, modified) values(default, $1, $2, $3, $4, $5, $6) returning id`)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var id int
	err = stmt.QueryRow(name, url, frequency, userid, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

/*
func GetAllJobs() {

}
*/

func InsertNewJobLog(jobID int, down bool, timeChecked time.Time) error {
	stmt, err := db.Prepare(`insert into jobLog(id, jobID, down, timeChecked) values(default, $1, $2, $3)`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(jobID, down, timeChecked)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
