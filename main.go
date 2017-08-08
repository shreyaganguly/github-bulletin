package main

import (
	"flag"
	"fmt"
)

var (
	githubToken  = flag.String("github-token", "", "Access Token of the github account used for fetching issues of users")
	slackToken   = flag.String("slack-token", "", "Slack Token of the bot that you will configure to send the github bulletin")
	organization = flag.String("org", "", "Organization for which issues are to be fetched")
	timePeriod   = flag.Int64("t", 180, "Time interval in seconds after which issues will be fetched")
)

func main() {
	flag.Parse()
	fmt.Println("Starting the bulletin")
	go giveNotification()
	configureSlack()
}
