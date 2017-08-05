package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

func giveNotification() {
	ticker := time.NewTicker(time.Second * 120)
	defer ticker.Stop()
	for range ticker.C {
		go issueFramework()
	}
}

func compareLabels(oldLabels, newLabels []github.Label) ([]string, []string) {
	var added, removed []string
	for i := 0; i < len(newLabels); i++ {
		found := 0
		for j := 0; j < len(oldLabels); j++ {
			if newLabels[i].GetName() == oldLabels[j].GetName() {
				found = 1
			}
		}
		if found == 0 {
			added = append(added, newLabels[i].GetName())
		}
	}
	for i := 0; i < len(oldLabels); i++ {
		found := 0
		for j := 0; j < len(newLabels); j++ {
			if newLabels[j].GetName() == oldLabels[i].GetName() {
				found = 1
			}
		}
		if found == 0 {
			removed = append(removed, oldLabels[i].GetName())
		}
	}
	return added, removed
}

func compareAssignees(oldAssignees, newAssignees []*github.User) ([]string, []string) {
	var added, removed []string
	for i := 0; i < len(newAssignees); i++ {
		found := 0
		for j := 0; j < len(oldAssignees); j++ {
			if newAssignees[i].GetName() == oldAssignees[j].GetName() {
				found = 1
			}
		}
		if found == 0 {
			added = append(added, newAssignees[i].GetName())
		}
	}
	for i := 0; i < len(oldAssignees); i++ {
		found := 0
		for j := 0; j < len(newAssignees); j++ {
			if newAssignees[j].GetName() == oldAssignees[i].GetName() {
				found = 1
			}
		}
		if found == 0 {
			removed = append(removed, oldAssignees[i].GetName())
		}
	}
	return added, removed
}

func findDifference(old, new []*github.Issue) {
	var i int
	if len(old) != 0 {
		for i = 0; i < len(new); i++ {
			if len(old) > i {
				if old[i].GetNumber() == new[i].GetNumber() {
					if old[i].GetState() != new[i].GetState() {
						fmt.Printf("The state changed for this issue : %s", new[i].GetHTMLURL())
					}
					if old[i].GetBody() != new[i].GetBody() {
						fmt.Printf("The body changed for this issue : %s", new[i].GetHTMLURL())
					}
					if old[i].GetTitle() != new[i].GetTitle() {
						fmt.Printf("The title changed for this issue : %s", new[i].GetHTMLURL())
					}
					if old[i].GetClosedAt() != new[i].GetClosedAt() {
						fmt.Printf("Issue got closed by %s", new[i].ClosedBy.GetLogin())
					}
					if old[i].Milestone != nil && new[i].Milestone != nil && old[i].Milestone.GetTitle() != new[i].Milestone.GetTitle() {
						fmt.Printf("The milestone changed for this issue : %s", new[i].GetHTMLURL())
					}
					if old[i].Milestone == nil && new[i].Milestone != nil {
						fmt.Printf("%s milestone added for this issue : %s", new[i].Milestone.GetTitle(), new[i].GetHTMLURL())
					}
					addedLabels, removedLabels := compareLabels(old[i].Labels, new[i].Labels)
					if len(addedLabels) != 0 {
						fmt.Printf("The following labels were added for this issue : %s,%s", new[i].GetHTMLURL(), strings.Join(addedLabels, ","))
					}
					if len(removedLabels) != 0 {
						fmt.Printf("The following labels were removed for this issue : %s,%s", new[i].GetHTMLURL(), strings.Join(removedLabels, ","))
					}

					addedAssignees, removedAssignees := compareAssignees(old[i].Assignees, new[i].Assignees)
					if len(addedAssignees) != 0 {
						fmt.Printf("The following assignees were added for this issue : %s,%s", new[i].GetHTMLURL(), strings.Join(addedAssignees, ","))
					}
					if len(removedAssignees) != 0 {
						fmt.Printf("The following assignees were removed for this issue : %s,%s", new[i].GetHTMLURL(), strings.Join(removedAssignees, ","))
					}
				}

			} else {
				fmt.Printf("A new issue is added: %s", new[i].GetHTMLURL())
			}

		}
	}
}

func findIssuesByAssignee(issues []*github.Issue, assignee string) (subscriberIssues []*github.Issue) {
	for _, issue := range issues {
		if issue.Assignee.GetLogin() == assignee {
			subscriberIssues = append(subscriberIssues, issue)
		}
	}
	return
}

func issueFramework() {
	fmt.Println("Hello World")
	filterOptions := CreateFilterOptions()
	issues, err := findIssues(filterOptions)
	if err != nil {
		fmt.Println("Github Bulletin Error: Error in listing by organization ", err)
		return
	}
	for _, value := range subscribers {
		issuesOfSubscriberNew := findIssuesByAssignee(issues, value)
		findDifference(issuesOfSubscriberNew, subsriberIssueMap[value])
		copy(subsriberIssueMap[value], issuesOfSubscriberNew)
	}

}
