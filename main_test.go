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

var fakeServiceTests = []struct {
	config   string
	outBool  bool
	outError bool
}{
	{"true", true, false},
	{"false", false, false},
	{"error", false, true},
}

func TestFakeServiceHas2FA(t *testing.T) {
	for _, tt := range fakeServiceTests {
		service, err := NewService("fake", tt.config)
		if err != nil {
			t.Fatalf("couldn't create fake service: %v", err)
		}

		ok, err := service.Has2FA("")

		if ok != tt.outBool {
			t.Errorf("expected Has2FA() to return %v, but returned %v", tt.outBool, ok)
		}
		if tt.outError && err == nil {
			t.Errorf("expected Has2FA() to return an error, but returned nil")
		}
		if !tt.outError && err != nil {
			t.Errorf("expected Has2FA() not to return an error, but returned %v", err)
		}
	}
}
