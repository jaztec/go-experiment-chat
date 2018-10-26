package network

import (
	"testing"
)

func TestClientCycle(t *testing.T) {
	c, err := NewClient("tcp", "localhost:12356")
	if err != nil {
		t.Fatalf("Error closing: %v\n", err)
	}
	if c == nil {
		t.Fail()
	}
	c.Close()
}
