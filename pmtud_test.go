package pmtud

import (
	"testing"
)

var testAddresses = [...]string{"google.com", "facebook.com", "amazon.com", "127.0.0.1"}

func TestPmtud(t *testing.T) {
	for _, addr := range testAddresses {
		pmtu, err := Pmtud(addr)
		if err != nil {
			t.Fatal(err)
		}
		if pmtu < 28 {
			t.Fatal("pmtu too small")
		}
	}
}
