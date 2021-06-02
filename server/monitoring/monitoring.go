package monitoring

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func NewPool() Pool {
	return Pool{jobs: []Job{}}
}

type Pool struct {
	jobs []Job
}

func (p Pool) AddJob(j Job) {
	p.jobs = append(p.jobs, j)
}

func (p Pool) GetJob(jobID int, userID float64) (*Job, error) {
	return nil, fmt.Errorf("Job not found")
}

func NewJob(jobId int, userID float64, name, url string, freq time.Duration) Job {
	j := Job{
		JobID:     jobId,
		UserID:    userID,
		Name:      name,
		URL:       url,
		Frequency: freq,
		status:    "NA",
		running:   false,
		done:      make(chan bool),
	}

	return j
}

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
					log.Println(err)
					//insert into jobLog
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
