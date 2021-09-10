package mailing

import "testing"

func TestConnectToSMTP(t *testing.T) {
	c, err := connectToSMTP("", 9999)
	if err != nil {
		t.Fail()
	}

	defer c.Close()
}
