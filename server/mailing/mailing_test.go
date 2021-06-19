package mailing

import "testing"

func TestConnectToSMTP(t *testing.T) {
	err := connectToSMTP()
	if err != nil {
		t.Fail()
	}
}
