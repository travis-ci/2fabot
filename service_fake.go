package main

import (
	"fmt"

	"github.com/pkg/errors"
)

type fakeService struct {
	config string
}

func (s *fakeService) Has2FA(userID string) (bool, error) {
	switch s.config {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "error":
		return false, errors.New("error fetching 2fa status for fake service")
	default:
		panic(fmt.Sprintf("unsupported fakeService config: %s", s.config))
	}
}

func newFakeService(config string) (Service, error) {
	return &fakeService{config: config}, nil
}

func init() {
	registerService("fake", "http://example.com/2fa_docs", newFakeService)
}
