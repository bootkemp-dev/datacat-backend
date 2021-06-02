package monitoring

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Job struct {
	JobID     int
	UserID    float64
	Name      string
	URL       string
	Frequency time.Duration
	status    string
	running   bool
	done      chan bool
}

func (j Job) Run() {
	log.Printf("Starting job | ID: %d | Name: %s\n", j.JobID, j.Name, j.URL)
	go func() {
		for {
			select {
			case <-j.done:
				j.running = false
				log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.JobID)
				return
			default:
				log.Printf("Job with ID: %d checking status of: %s ", j.JobID, j.URL)
				err := j.Status()
				if err != nil {
					j.status = "DOWN"
					log.Println(err)
					//insert into jobLog
				}

				time.Sleep(j.Frequency)
			}
		}
	}()
}

func (j Job) Status() error {
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
}
