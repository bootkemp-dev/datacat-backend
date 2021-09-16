package models

import (
	"fmt"
	"log"
	"testing"
	"time"
)

/*
func TestRemoveJobFromPool(t *testing.T) {
	p := NewPool()
	j1, err := NewJob(1, 1, "test1job", "http://google.com", 10000000000)
	j2, err := NewJob(2, 1, "test2job", "http://google.com", 10000000000)

	var wg sync.WaitGroup
	wg.Add(2)

	p.AddJob(j1)
	p.AddJob(j2)

	fmt.Println("Size of the pool:", p.GetPoolSize())
	j1.Run()
	j2.Run()

	time.Sleep(20 * time.Second)
	fmt.Println("Status of job 1:", j1.GetStatus())
	fmt.Println("Status of job 2:", j2.Status)
	time.Sleep(10 * time.Second)
	fmt.Println("Job 1 Active: ", j1.Active)
	fmt.Println("Stopping job 1 ...")
	j1.Stop()
	fmt.Println("Status of job 1:", j1.GetStatus())
	fmt.Println("Job 1 Active: ", j1.Active)
	fmt.Println("Deleting job 1 from the pool...")
	err = p.RemoveJob(1, 1)
	if err != nil {
		t.Fail()
	}
	fmt.Println("Size of the pool:", p.GetPoolSize())

	wg.Wait()
}
*/
func TestPinger(t *testing.T) {
	p := NewPool()
	job, err := NewJob(1, 1, "test", "google.com", 1)
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	p.AddJob(job)
	job.Run()
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("job done")
				return
			default:
				fmt.Println(job.GetPing())
				fmt.Println("Status: ", job.GetStatus())
			}
		}
	}()

	time.Sleep(time.Second * 10)
	job.Stop()
	fmt.Println("Job stoped")
	done <- true

}
