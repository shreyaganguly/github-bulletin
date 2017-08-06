package main

import (
	"flag"

	"github.com/google/go-github/github"
)

var (
	githubToken  = flag.String("github-token", "", "Access Token of the github account used for fetching issues of users")
	slackToken   = flag.String("slack-token", "", "Slack Token of the bot that you will configure to send the github bulletin")
	organization = flag.String("org", "", "Organization for which issues are to be searched")
)
var subscribers map[string]string
var subsriberIssueMap map[string][]*github.Issue

func main() {
	flag.Parse()
	subscribers = make(map[string]string)
	subsriberIssueMap = make(map[string][]*github.Issue)
	go func() {
		giveNotification()
	}()

	configureSlack()
}
