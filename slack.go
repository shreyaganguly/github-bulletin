package main

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

var api *slack.Client

func postMessage(user, msg string) error {
	channelID, timestamp, err := api.PostMessage(user, msg, slack.PostMessageParameters{})
	if err != nil {
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}

func subscribeUser(api *slack.Client, msg slack.Msg) error {
	msgArray := strings.Split(msg.Text, ":")
	if len(msgArray) == 2 {
		githubUserID := strings.Trim(msgArray[1], " ")
		s := &Subscription{githubUserID, msg.User, []*github.Issue{}}
		s.addSubscriber()
		err := postMessage(msg.User, fmt.Sprintf("Successfully Subscribed %s to github bulletin. To unsubsribe Write \"Unsubscribe: %s\"", githubUserID, githubUserID))
		if err != nil {
			return err
		}
	} else {
		err := postMessage(msg.User, "Oops!! Not Successfully Subscribed to bulletin. Write \"Subscribe: <Your github User ID>\"")
		if err != nil {
			return err
		}
	}
	return nil
}

func unsubscribeUser(api *slack.Client, msg slack.Msg) error {
	msgArray := strings.Split(msg.Text, ":")
	if len(msgArray) == 2 {
		githubUserID := strings.Trim(msgArray[1], " ")
		s := &Subscription{githubUserID, msg.User, []*github.Issue{}}
		errSubscriber := s.removeSubscriber()
		if errSubscriber != nil {
			err := postMessage(msg.User, errSubscriber.Error())
			if err != nil {
				return err
			}
			return nil
		}
		err := postMessage(msg.User, fmt.Sprintf("Successfully Unsubscribed %s from github bulletin. To subscribe again Write \"Subscribe: %s\"", githubUserID, githubUserID))
		if err != nil {
			return err
		}
	} else {
		err := postMessage(msg.User, "Oops!! Not Successfully Unubscribed to bulletin. Write \"Subscribe: <Your github User ID>\"")
		if err != nil {
			return err
		}
	}
	return nil
}

func parseMessage(api *slack.Client, msg slack.Msg) {
	if strings.HasPrefix(msg.Text, "Subscribe: ") {
		err := subscribeUser(api, msg)
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
	} else if strings.HasPrefix(msg.Text, "Unsubscribe: ") {
		err := unsubscribeUser(api, msg)
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
	} else {
		err := postMessage(msg.User, "You can write here. But I understand only \"Subscribe: <Your github User ID>\" and \"Unsubscribe: <Your github User ID>\" . Enjoy!")
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
	}

}

func configureSlack() {
	api = slack.New(*slackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev.Msg.User)
			fmt.Println("Messgae is ", ev.Msg.Text)
			parseMessage(api, ev.Msg)
		default:
		}
	}
}
