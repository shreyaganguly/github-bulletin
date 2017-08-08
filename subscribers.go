package main

import (
	"fmt"

	"github.com/google/go-github/github"
)

type Subscription struct {
	GithubUserID string
	SlackUserID  string
	Issues       []*github.Issue
}

var (
	subscriptionList []*Subscription
)

func (s *Subscription) addSubscriber() {
	for _, v := range subscriptionList {
		if v.SlackUserID == s.SlackUserID && v.GithubUserID != s.GithubUserID {
			v.GithubUserID = s.GithubUserID
		}
	}
	subscriptionList = append(subscriptionList, s)
}

func (s *Subscription) removeSubscriber() error {
	for i, v := range subscriptionList {
		if v.GithubUserID == s.GithubUserID && v.SlackUserID == s.SlackUserID {
			subscriptionList = append(subscriptionList[:i], subscriptionList[i+1:]...)
		} else {
			return fmt.Errorf("You never subscribed with %s", s.GithubUserID)
		}
	}
	return nil
}
