package mailing

import "testing"

func TestConnectToSMTP(t *testing.T) {
	c, err := connectToSMTP()
	if err != nil {
		t.Fail()
	}

	defer c.Close()
}
