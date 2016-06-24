package main

import "testing"

func TestNewService(t *testing.T) {
	service, err := NewService("fake", "")
	if err != nil {
		t.Fatalf("NewService shouldn't return error, but did: %v", err)
	}

	if service == nil {
		t.Error("NewService shouldn't return nil service, but did")
	}
}
