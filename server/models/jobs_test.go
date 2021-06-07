package models

import (
	"log"
	"testing"
	"time"
)

func TestRemoveJobFromPool(t *testing.T) {
	p := NewPool()
	j1 := NewJob(1, 1, "test1job", "http://google.com", 10000000000)
	j2 := NewJob(2, 1, "test2job", "http://google.com", 10000000000)
	j3 := NewJob(3, 1, "test1job", "http://google.com", 10000000000)
	p.Jobs = append(p.Jobs, j1, j2, j3)
	j1.Run()
	time.Sleep(time.Second * 25)
	log.Println("time done")
	j1.Stop()
	log.Println(j1.GetStatus())

	time.Sleep(10 * time.Minute)
}
