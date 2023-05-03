package metrics

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func Responseviness_score(personal_token string, owner string, repo string) float64 {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personal_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opt := &github.IssueListByRepoOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repo, opt)
	if err != nil {
		log.Println("Error:", err)
		return 0
	}

	tot_count := 0
	empty_comment_count := 0
	for _, issue := range issues {
		tot_count++
		// check if the issue has no comments
		if *issue.Comments == 0 {
			empty_comment_count++
		}
	}
	//sugar_logger.Debugf("Number of total issues on the first page (max is 100): %d", tot_count)
	//sugar_logger.Debugf("Number of issues with empty comments: %d", empty_comment_count)
	if tot_count == 0 {
		return 0.0
	}
	score, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", float64(empty_comment_count)/float64(tot_count)), 64)
	return 1 - score
}
