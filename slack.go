package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nlopes/slack"
)

func SubscribeUser(api *slack.Client, msg slack.Msg) error {
	msgArray := strings.Split(msg.Text, ":")
	fmt.Println("msgArray is", len(msgArray))
	if len(msgArray) == 2 {
		githubUserID := strings.Trim(msgArray[1], " ")
		addSubscriber(msg.User, githubUserID)
		channelID, timestamp, err := api.PostMessage(msg.User, fmt.Sprintf("Successfully Subscribed %s to github bulletin. To unsubsribe Write \"Unsubscribe: <Your github User ID>\"", githubUserID), slack.PostMessageParameters{})
		if err != nil {
			return err
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	} else {
		channelID, timestamp, err := api.PostMessage(msg.User, "Oops!! Not Successfully Subscribed to bulletin. Write \"Subscribe: <Your github User ID>\"", slack.PostMessageParameters{})
		if err != nil {
			return err
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	}
	return nil
}

func UnsubscribeUser(api *slack.Client, msg slack.Msg) error {
	msgArray := strings.Split(msg.Text, ":")
	fmt.Println("msgArray is", len(msgArray))
	if len(msgArray) == 2 {
		githubUserID := strings.Trim(msgArray[1], " ")
		removeSubscriber(msg.User, githubUserID)
		channelID, timestamp, err := api.PostMessage(msg.User, fmt.Sprintf("Successfully Unsubscribed %s from github bulletin. To subscribe again Write \"Subscribe: <Your github User ID>\"", githubUserID), slack.PostMessageParameters{})
		if err != nil {
			return err
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	} else {
		channelID, timestamp, err := api.PostMessage(msg.User, "Oops!! Not Successfully Unubscribed to bulletin. Write \"Subscribe: <Your github User ID>\"", slack.PostMessageParameters{})
		if err != nil {
			return err
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	}
	return nil
}
func parseMessage(api *slack.Client, msg slack.Msg) {
	if strings.HasPrefix(msg.Text, "Subscribe: ") {
		err := SubscribeUser(api, msg)
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
	} else if strings.HasPrefix(msg.Text, "Unsubscribe: ") {
		err := UnsubscribeUser(api, msg)
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
	} else {
		channelID, timestamp, err := api.PostMessage(msg.User, "You can write here. But I understand only \"Subscribe: <Your github User ID>\" and \"Unsubscribe: <Your github User ID>\" . Enjoy!", slack.PostMessageParameters{})
		if err != nil {
			fmt.Printf("Github-Bulletin : Slack Error %s", err.Error())
			return
		}
		fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	}

}

func configureSlack() {
	api := slack.New(*slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev.Msg.User)
			fmt.Println("Messgae is ", ev.Msg.Text)
			parseMessage(api, ev.Msg)
			// return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
