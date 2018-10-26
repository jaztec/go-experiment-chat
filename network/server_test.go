package network

import (
	"testing"
)

func TestServerCycle(t *testing.T) {
	s := NewServer(ServerConfig{Port: 12356, MessageBufferSize: 1024})
	if s == nil {
		t.Fail()
	}
	err := s.Close()
	if err != nil {
		t.Fatalf("Error closing: %v\n", err)
	}
}
