package models

import (
	"testing"
)

func TestRemoveJobFromPool(t *testing.T) {
	p := NewPool()
	j1 := NewJob(1, 1, "test1job", "http://google.com", 10000000000)
	j2 := NewJob(2, 1, "test2job", "http://google.com", 10000000000)
	j3 := NewJob(3, 1, "test1job", "http://google.com", 10000000000)
	p.Jobs = append(p.Jobs, j1, j2, j3)
}
