package utils

import (
	"testing"
	"time"
)

func TestInTimeSpan(t *testing.T) {
	if !InTimeSpan(time.Now().Local().Add(12 * time.Hour)) {
		t.Fail()
	}
}
