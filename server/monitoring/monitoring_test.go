package monitoring

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	p := NewPool()
	j := NewJob(1, 1, "test-job", "http://google.com", 10*time.Second)
	p.AddJob(j)
	var wg sync.WaitGroup
	wg.Add(1)
	j.Run()

	pj, err := p.GetJob(1, 1)
	if err != nil {
		log.Println(err)
	}

	log.Println(pj)

	fmt.Println(p.jobs)

	wg.Wait()

}
