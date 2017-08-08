package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func giveNotification() {
	timeInterval, _ := time.ParseDuration(fmt.Sprintf("%ds", *timePeriod))
	ticker := time.NewTicker(timeInterval)
	defer ticker.Stop()
	for range ticker.C {
		go issueFramework()
	}
}

func compareLabels(oldLabels, newLabels []github.Label) ([]string, []string) {
	var added, removed []string
	newMap := make(map[string]int)
	for _, newLabel := range newLabels {
		newMap[newLabel.GetName()]++
	}
	for _, oldLabel := range oldLabels {
		if newMap[oldLabel.GetName()] > 0 {
			delete(newMap, oldLabel.GetName())
		} else {
			removed = append(removed, oldLabel.GetName())
		}
	}
	for addedLabel := range newMap {
		added = append(added, addedLabel)
	}
	return added, removed
}

func compareAssignees(oldAssignees, newAssignees []*github.User) ([]string, []string) {
	var added, removed []string
	newMap := make(map[string]int)
	for _, newAssignee := range newAssignees {
		newMap[newAssignee.GetName()]++
	}
	for _, oldAssignee := range oldAssignees {
		if newMap[oldAssignee.GetName()] > 0 {
			delete(newMap, oldAssignee.GetName())
		} else {
			removed = append(removed, oldAssignee.GetName())
		}
	}
	for addedAssignee := range newMap {
		added = append(added, addedAssignee)
	}
	return added, removed
}

func reverse(issues []*github.Issue) []*github.Issue {
	for i, j := 0, len(issues)-1; i < j; i, j = i+1, j-1 {
		issues[i], issues[j] = issues[j], issues[i]
	}
	return issues
}

func findDifference(old, new []*github.Issue) string {
	var i int
	var message string
	if len(old) != 0 {
		for i = 0; i < len(new); i++ {
			if len(old) > i {
				if old[i].GetNumber() == new[i].GetNumber() {
					if old[i].GetState() != new[i].GetState() {
						message = fmt.Sprintf("\n%s\nThe state changed for this issue : %s, changed from \"%s\" to \"%s\"", message, new[i].GetHTMLURL(), old[i].GetState(), new[i].GetState())
					}
					if old[i].GetBody() != new[i].GetBody() {
						message = fmt.Sprintf("\n%s\nThe body changed for this issue : %s, changed from \"%s\" to \"%s\"", message, new[i].GetHTMLURL(), old[i].GetBody(), new[i].GetBody())
					}
					if old[i].GetTitle() != new[i].GetTitle() {
						message = fmt.Sprintf("\n%s\nThe title changed for this issue : %s, changed from \"%s\" to \"%s\"", message, new[i].GetHTMLURL(), old[i].GetTitle(), new[i].GetTitle())
					}
					if old[i].GetClosedAt() != new[i].GetClosedAt() {
						message = fmt.Sprintf("\n%s\nIssue got closed by %s", message, new[i].ClosedBy.GetLogin())
					}
					if old[i].Milestone != nil && new[i].Milestone != nil && old[i].Milestone.GetTitle() != new[i].Milestone.GetTitle() {
						message = fmt.Sprintf("\n%s\nThe milestone changed for this issue : %s, changed from \"%s\" to \"%s\"", message, new[i].GetHTMLURL(), old[i].Milestone.GetTitle(), new[i].Milestone.GetTitle())
					}
					if old[i].Milestone == nil && new[i].Milestone != nil {
						message = fmt.Sprintf("\n%s\n%s milestone added for this issue : %s", message, new[i].Milestone.GetTitle(), new[i].GetHTMLURL())
					}
					if old[i].Milestone != nil && new[i].Milestone == nil {
						message = fmt.Sprintf("\n%s\n%s milestone removed for this issue : %s", message, old[i].Milestone.GetTitle(), new[i].GetHTMLURL())
					}
					addedLabels, removedLabels := compareLabels(old[i].Labels, new[i].Labels)
					if len(addedLabels) != 0 {
						message = fmt.Sprintf("\n%s\nThe following labels were added for this issue : %s, %s", message, new[i].GetHTMLURL(), strings.Join(addedLabels, ","))
					}
					if len(removedLabels) != 0 {
						message = fmt.Sprintf("\n%s\nThe following labels were removed for this issue : %s, %s", message, new[i].GetHTMLURL(), strings.Join(removedLabels, ","))
					}

					addedAssignees, removedAssignees := compareAssignees(old[i].Assignees, new[i].Assignees)
					if len(addedAssignees) != 0 {
						message = fmt.Sprintf("\n%s\nThe following assignees were added for this issue : %s, %s", message, new[i].GetHTMLURL(), strings.Join(addedAssignees, ","))
					}
					if len(removedAssignees) != 0 {
						message = fmt.Sprintf("\n%s\nThe following assignees were removed for this issue : %s, %s", message, new[i].GetHTMLURL(), strings.Join(removedAssignees, ","))
					}
				}

			} else {
				message = fmt.Sprintf("\n%s\nA new issue is added: %s by %s", message, new[i].GetHTMLURL(), new[i].User.GetLogin())
			}
		}
	}
	return message
}

func findIssuesByAssignee(issues []*github.Issue, assignee string) (subscriberIssues []*github.Issue) {
	for _, issue := range issues {
		if issue.Assignee != nil && issue.Assignee.GetLogin() == assignee {
			subscriberIssues = append(subscriberIssues, issue)
		}
	}
	return
}

func issueFramework() {
	filterOptions := CreateFilterOptions()
	issues, err := findIssues(filterOptions)
	if err != nil {
		fmt.Println("Github Bulletin Error: Error in listing by organization ", err)
		os.Exit(0)
	}
	for _, subscription := range subscriptionList {
		issuesOfSubscriberNew := findIssuesByAssignee(issues, subscription.GithubUserID)
		issuesOfSubscriberNew = reverse(issuesOfSubscriberNew)
		message := findDifference(subscription.Issues, issuesOfSubscriberNew)
		if message != "" {
			err := postMessage(subscription.SlackUserID, message)
			if err != nil {
				fmt.Println("Github Bulletin Error: Slack Error in posting message ", err)
				os.Exit(0)
			}
		}
		subscription.Issues = make([]*github.Issue, len(issuesOfSubscriberNew))
		copy(subscription.Issues, issuesOfSubscriberNew)
	}

}
