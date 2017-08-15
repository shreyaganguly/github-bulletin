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

func mapper() func(string) map[string]int {
	entityMap := make(map[string]int)
	return func(x string) map[string]int {
		entityMap[x]++
		return entityMap
	}
}

func comparer(m map[string]int) func(string) []string {
	var diff []string
	return func(x string) []string {
		if m[x] > 0 {
			delete(m, x)
		} else {
			diff = append(diff, x)
		}
		return diff
	}
}

func addedEntity(m map[string]int) []string {
	var added []string
	for v := range m {
		added = append(added, v)
	}
	return added
}

func compareLabels(oldLabels, newLabels []github.Label) ([]string, []string) {
	var removed []string
	newMap := make(map[string]int)
	m := mapper()
	for _, newLabel := range newLabels {
		newMap = m(newLabel.GetName())
	}
	d := comparer(newMap)
	for _, oldLabel := range oldLabels {
		removed = d(oldLabel.GetName())
	}
	return addedEntity(newMap), removed
}

func compareAssignees(oldAssignees, newAssignees []*github.User) ([]string, []string) {
	var removed []string
	m := mapper()
	newMap := make(map[string]int)
	for _, newAssignee := range newAssignees {
		newMap = m(newAssignee.GetLogin())
	}
	d := comparer(newMap)
	for _, oldAssignee := range oldAssignees {
		removed = d(oldAssignee.GetLogin())
	}

	return addedEntity(newMap), removed
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
						message = fmt.Sprintf("\n%s\n\"%s\" milestone added for this issue : %s", message, new[i].Milestone.GetTitle(), new[i].GetHTMLURL())
					}
					if old[i].Milestone != nil && new[i].Milestone == nil {
						message = fmt.Sprintf("\n%s\n\"%s\" milestone removed for this issue : %s", message, old[i].Milestone.GetTitle(), new[i].GetHTMLURL())
					}
					addedLabels, removedLabels := compareLabels(old[i].Labels, new[i].Labels)
					if len(addedLabels) != 0 {
						message = fmt.Sprintf("\n%s\nThe following labels were added for this issue : %s, \"%s\"", message, new[i].GetHTMLURL(), strings.Join(addedLabels, ","))
					}
					if len(removedLabels) != 0 {
						message = fmt.Sprintf("\n%s\nThe following labels were removed for this issue : %s, \"%s\"", message, new[i].GetHTMLURL(), strings.Join(removedLabels, ","))
					}

					addedAssignees, removedAssignees := compareAssignees(old[i].Assignees, new[i].Assignees)
					if len(addedAssignees) != 0 {
						message = fmt.Sprintf("\n%s\nThe following assignees were added for this issue : %s, \"%s\"", message, new[i].GetHTMLURL(), strings.Join(addedAssignees, ","))
					}

					if len(removedAssignees) != 0 {
						message = fmt.Sprintf("\n%s\nThe following assignees were removed for this issue : %s, \"%s\"", message, new[i].GetHTMLURL(), strings.Join(removedAssignees, ","))
					}
				}

			} else {
				message = fmt.Sprintf("\n%s\nA new issue is added: %s by \"%s\"", message, new[i].GetHTMLURL(), new[i].User.GetLogin())
			}
		}
	} else {
		for _, v := range new {
			if v.GetState() != "closed" {
				message = fmt.Sprintf("\n%s\n%s issue is still open which was assigned to you by \"%s\"", message, v.GetHTMLURL(), v.User.GetLogin())
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
