package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type sentryService struct {
	token   string
	orgName string
}

type sentryOrgMember struct {
	User sentryUser `json:"user"`
}

type sentryUser struct {
	Email  string `json:"email"`
	Has2FA bool   `json:"has2fa"`
}

func (s *sentryService) Has2FA(userID string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://app.getsentry.com/api/0/organizations/%s/members/", s.orgName), nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var users []sentryOrgMember
	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		return false, err
	}

	for _, user := range users {
		if user.User.Email == userID {
			return user.User.Has2FA, nil
		}
	}

	return false, errors.Errorf("no sentry user in that org with the email %s", userID)
}

func newSentryService(config string) (Service, error) {
	configParts := strings.Split(config, ",")

	return &sentryService{
		token:   configParts[0],
		orgName: configParts[1],
	}, nil
}

func init() {
	registerService("sentry", "https://app.getsentry.com/account/settings/2fa/", newSentryService)
}
