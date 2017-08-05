package main

import "fmt"

func addSubscriber(slackUserID, githubUserID string) {
	subscribers[slackUserID] = githubUserID
}

func removeSubscriber(slackUserID, githubUserID string) error {
	if subscribers[slackUserID] != githubUserID {
		return fmt.Errorf("You never subscribed with %s", githubUserID)
	}
	delete(subscribers, slackUserID)
	return nil
}
