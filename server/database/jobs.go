package database

import (
	"log"
	"time"

	"github.com/bootkemp-dev/datacat-backend/models"
)

func InsertNewJob(name, url string, frequency int64, userid int) (int, error) {
	stmt, err := db.Prepare(`insert into jobs(id, jobName, jobUrl, frequency, userid, active, created, modified) values(default, $1, $2, $3, $4, $5, $6, $7) returning id`)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var id int
	err = stmt.QueryRow(name, url, frequency, userid, true, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return id, nil
}

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

func GetAllJobsByUserID(userID int) ([]*models.Job, error) {
	rows, err := db.Query(`select * from jobs where userid=$1`, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var jobs []*models.Job

	for rows.Next() {
		var job models.Job
		err := rows.Scan(&job.ID, &job.Name, &job.URL, &job.Frequency, &job.UserID, &job.Active, &job.CreatedAt, &job.ModifiedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func DeleteJob(jobID, userID int) error {
	stmt, err := db.Prepare(`delete from jobs where id=$1 and userid=$2`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(jobID, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UpdateJobActive(active bool, jobID, userID int) error {
	stmt, err := db.Prepare(`update jobs set active = $1 where id = $2 and userid = $3`)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = stmt.Exec(active, jobID, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
