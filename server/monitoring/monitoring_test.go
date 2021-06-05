package monitoring

import (
	"log"
	"testing"
	"time"
)

func TestAddJobToPool(t *testing.T) {
	p := NewPool()
	j := NewJob(1, 1, "test-job", "http://google.com", 10*time.Second)
	p.Jobs = append(p.Jobs, j)
	log.Println(len(p.Jobs))
}

func TestFindJob(t *testing.T) {
	p := NewPool()
	j1 := NewJob(1, 1, "test-job", "http://google.com", 10*time.Second)
	j2 := NewJob(2, 1, "test-job", "http://google.com", 10*time.Second)

	p.Jobs = append(p.Jobs, j1, j2)

	jobFound, err := p.GetJob(2, 1)
	if err != nil {
		t.Fail()
	}

	log.Println(jobFound.GetStatus())
}
