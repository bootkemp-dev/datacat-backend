package models

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/bootkemp-dev/datacat-backend/logger"
	"github.com/go-ping/ping"
)

type NewJobRequest struct {
	JobName   string `json:"name"`
	JobURL    string `json:"url"`
	Frequency int64  `json:"frequency"`
}

func NewPool() Pool {
	return Pool{jobs: []*Job{}}
}

type Pool struct {
	jobs []*Job
}

func (p *Pool) AddJob(job *Job) {
	p.jobs = append(p.jobs, job)
}

func (p *Pool) RemoveJob(jobID int, userID int) error {
	for i := range p.jobs {
		if p.jobs[i].ID == jobID && p.jobs[i].UserID == userID {
			if p.jobs[i].Active == true {
				go p.jobs[i].Stop()
			}
			p.jobs = append(p.jobs[:i], p.jobs[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Job not found in the pool")
}

func (p Pool) GetJob(jobID int, userID int) (*Job, error) {
	for _, v := range p.jobs {
		if v.UserID == userID && v.ID == jobID {
			return v, nil
		}
	}
	return nil, fmt.Errorf("Job not found")
}

func (p Pool) GetJobsByUserID(userID int) []*Job {
	var jobs []*Job
	for i := range p.jobs {
		if p.jobs[i].UserID == userID {
			jobs = append(jobs, p.jobs[i])
		}
	}

	return jobs
}

func (p Pool) GetPoolSize() int {
	return len(p.jobs)
}

type Job struct {
	ID         int           `json:"id"`
	Name       string        `json:"jobName"`
	URL        string        `json:"jobUrl"`
	Frequency  int64         `json:"frequency"`
	UserID     int           `json:"userID"`
	Active     bool          `json:"active"`
	CreatedAt  time.Time     `json:"createdAt"`
	ModifiedAt time.Time     `json:"modifiedAt"`
	Status     string        `json:"status"`
	Done       chan bool     `json:"-"`
	logger     logger.Logger `json:"-"`
	pinger     *ping.Pinger  `json:"-"`
	ping       chan time.Duration
}

func NewJob(jobId int, userID int, name, url string, freq int64, createdAt, modifiedAt time.Time, active bool) (*Job, error) {

	if strings.Contains(url, "http://") {
		url = strings.ReplaceAll(url, "http://", "")
	}
	if strings.Contains(url, "https://") {
		url = strings.ReplaceAll(url, "https://", "")
	}

	j := Job{
		ID:         jobId,
		Name:       name,
		URL:        url,
		Frequency:  freq,
		UserID:     userID,
		Active:     active,
		CreatedAt:  createdAt,
		ModifiedAt: time.Now(),
		Status:     "NA",
		Done:       make(chan bool),
		ping:       make(chan time.Duration),
	}

	p, err := ping.NewPinger(j.URL)
	if err != nil {
		return nil, err
	}

	j.pinger = p
	p.OnRecv = func(pkt *ping.Packet) { j.ping <- pkt.Rtt }

	return &j, nil
}

func (j *Job) Run() {
	log.Printf("Starting job | ID: %d | Name: %s\n", j.ID, j.Name)
	j.SetActive(true)
	go j.run()
}

func (j *Job) run() {
	go j.pinger.Run()
	for {
		select {
		case <-j.Done:
			j.pinger.Stop()
			log.Printf("Ending job | ID: %d | Name: %s | URL: %s\n", j.ID, j.Name, j.URL)
			j.SetStatus("NA")
			j.SetActive(false)
			j.SetModifiedNow()
			return
		default:
			_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:80", j.URL), 1*time.Second)
			if err != nil {
				j.SetStatus("DOWN")
			} else {
				j.SetStatus("UP")
			}
		}
	}

}

func (j *Job) SetModifiedNow() {
	j.ModifiedAt = time.Now()
}

func (j *Job) Stop() {
	j.Done <- true
}

func (j *Job) SetStatus(s string) {
	j.Status = s
}

func (j *Job) SetActive(a bool) {
	j.Active = a
}

func (j *Job) GetStatus() string {
	return j.Status
}

func (j *Job) GetActive() bool {
	return j.Active
}

func (j *Job) GetPing() *time.Duration {
	for x := range j.ping {
		return &x
	}
	return nil
}
