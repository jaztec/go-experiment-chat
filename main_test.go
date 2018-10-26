package main

import "testing"

func TestPrintMicros(t *testing.T) {
	start := startStampMicros()
	t.Logf("Start stamp: %d", start)
	stamp := printMicros(start)
	if stamp < start {
		t.Fail()
	}
}

func TestInitializeServer(t *testing.T) {
	s, n, err := initializeServer()
	// defer s.Close()
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if s == nil {
		t.Fail()
	}
	if n == nil {
		t.Fail()
	}
}

func TestInitializeClient(t *testing.T) {
	c, n, err := initializeClient()
	defer c.Close()
	if err != nil {
		t.Fatalf("Error %v", err)
	}
	if c == nil {
		t.Fail()
	}
	if n == nil {
		t.Fail()
	}
}
