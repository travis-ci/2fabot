package main

import "testing"

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
