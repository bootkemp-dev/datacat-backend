package utils

import "time"

func InTimeSpan(t time.Time) bool {
	now := time.Now().Local()
	if now.After(t) {
		return false
	} else {
		return true
	}
}
