package monitoring

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bootkemp-dev/datacat-backend/database"
)

func NewPool() Pool {
	return Pool{Jobs: []Job{}}
}

type Pool struct {
	Jobs []Job
}

func (p Pool) GetJob(jobID int, userID int) (*Job, error) {
	for _, v := range p.Jobs {
		if v.UserID == userID && v.JobID == jobID {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("Job not found")
}

func NewJob(jobId int, userID int, name, url string, freq time.Duration) Job {
	j := Job{
		JobID:     jobId,
		UserID:    userID,
		Name:      name,
		URL:       url,
		Frequency: freq,
		status:    "UP",
		running:   false,
		done:      make(chan bool),
	}

	return j
}

type Job struct {
	JobID     int
	UserID    int
	Name      string
	URL       string
	Frequency time.Duration
	status    string
	running   bool
	done      chan bool
}

func (j Job) Run() {
	j.running = true
	log.Printf("Starting job | ID: %d | Name: %s\n", j.JobID, j.Name)
	go func() {
		for {
			select {
			case <-j.done:
				j.running = false
				log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.JobID, j.Name, j.URL)
				return
			default:
				log.Printf("Job with ID: %d checking status of: %s ", j.JobID, j.URL)
				err := j.URLStatus()
				if err != nil {
					j.status = "DOWN"
					err = database.InsertNewJobLog(j.JobID, true, time.Now())
					if err != nil {
						log.Println(err)
					}
				}

				time.Sleep(j.Frequency)
			}
		}
	}()
}

func (j Job) URLStatus() error {
	resp, err := http.Get(j.URL)
	if err == nil && resp.StatusCode == 200 {
		return nil
	} else {
		return err
	}
}

func (j Job) Stop() error {
	if j.running == false {
		return fmt.Errorf("Job is not running")
	} else {
		j.done <- true
	}

	return nil
}

func (j Job) GetStatus() string {
	return j.status
}
