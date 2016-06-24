package main

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

type slackService struct {
	client *slack.Client
}

func (s *slackService) Has2FA(userID string) (bool, error) {
	user, err := s.client.GetUserInfo(userID)
	if err != nil {
		return false, errors.Wrap(err, "unable to get user info")
	}

	if user == nil {
		return false, errors.Errorf("no user with ID %s", userID)
	}

	fmt.Printf("user info for %s: %+v\n", userID, *user)

	return user.Has2FA, nil
}

func newSlackService(config string) (Service, error) {
	return &slackService{
		client: slack.New(config),
	}, nil
}

func init() {
	registerService("slack", "https://get.slack.help/hc/en-us/articles/204509068-Enabling-two-factor-authentication", newSlackService)
}
