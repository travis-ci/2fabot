package main

import (
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type githubService struct {
	client *github.Client
	org    string

	membersOnce  sync.Once
	membersError error
	no2fa        []string
	has2fa       []string
}

func (s *githubService) Has2FA(userID string) (bool, error) {
	s.membersOnce.Do(s.refresh)
	if s.membersError != nil {
		return false, s.membersError
	}

	if stringSliceContainsString(s.has2fa, userID) {
		return true, nil
	}

	if stringSliceContainsString(s.no2fa, userID) {
		return false, nil
	}

	return false, errors.Errorf("no user with the username %s", userID)
}

func (s *githubService) refresh() {
	usersWithout2fa, _, err := s.client.Organizations.ListMembers(s.org, &github.ListMembersOptions{Filter: "2fa_disabled"})
	if err != nil {
		s.membersError = err
		return
	}

	s.no2fa = make([]string, 0, len(usersWithout2fa))
	for _, user := range usersWithout2fa {
		s.no2fa = append(s.no2fa, *user.Login)
	}

	usersWith2fa, _, err := s.client.Organizations.ListMembers(s.org, &github.ListMembersOptions{Filter: "all"})
	if err != nil {
		s.membersError = err
		return
	}

	s.has2fa = make([]string, 0, len(usersWith2fa))
	for _, user := range usersWith2fa {
		s.has2fa = append(s.has2fa, *user.Login)
	}
}

func newGitHubService(config string) (Service, error) {
	configParts := strings.Split(config, ",")

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: configParts[0]},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	return &githubService{client: client, org: configParts[1]}, nil
}

func init() {
	registerService("github", "https://help.github.com/articles/securing-your-account-with-two-factor-authentication-2fa/", newGitHubService)
}

func stringSliceContainsString(haystack []string, needle string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}

	return false
}
