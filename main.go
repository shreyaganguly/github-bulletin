package main

import (
	"flag"

	"github.com/google/go-github/github"
)

var (
	accessToken  = flag.String("token", "", "Access Token of the github account used for fetching issues of users")
	organization = flag.String("org", "", "Organization for which issues are to be searched")
)
var subscribers map[string]string
var subsriberIssueMap map[string][]*github.Issue

func main() {
	flag.Parse()
	subscribers = make(map[string]string)
	subscribers["test-slack-name"] = "test-github-id"
	subsriberIssueMap = make(map[string][]*github.Issue)
	giveNotification()
}
