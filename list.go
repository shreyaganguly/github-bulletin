package main

import (
	"context"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

func CreateFilterOptions() *github.IssueListOptions {
	listOptions := &github.IssueListOptions{Filter: "all", ListOptions: github.ListOptions{PerPage: 100}}
	listOptions.State = "all"
	return listOptions
}

func findIssues(filterOptions *github.IssueListOptions) (issueList []*github.Issue, err error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	i := 1
	for {
		filterOptions.ListOptions.Page = i
		filterOptions.Direction = "asc"
		issues, _, err := client.Issues.ListByOrg(ctx, *organization, filterOptions)
		if err != nil {
			return nil, err
		}
		if len(issues) == 0 {
			break
		}
		i++
		issueList = append(issueList, issues...)
	}
	return
}
