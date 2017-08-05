package main

import "errors"

func addSubscriber(slackUserID, githubUserID string) {
	subscribers[slackUserID] = githubUserID
}

func removeSubscriber(slackUserID string) error {
	_, ok := subscribers[slackUserID]
	if ok {
		delete(subscribers, slackUserID)
		return nil
	}
	return errors.New("You never subscribed")
}
