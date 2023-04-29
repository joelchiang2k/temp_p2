package metrics

import (
	"context"
	"log"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Query struct {
	Repository struct {
		LicenseInfo struct {
			Key string
		} `graphql:"licenseInfo"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func license_score(personal_token string, owner string, repo string) float64 {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personal_token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := githubv4.NewClient(tc)

	var q Query
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}
	err := client.Query(ctx, &q, variables)
	if err != nil {
		log.Println("Error:", err)
		return 0
	}

	//sugar_logger.Debugf("License Name:", q.Repository.LicenseInfo.Key)
	compatible_list := [5]string{"cc0-1.0", "unlicense", "bsd-3-clause", "mit", "lgpl-2.1"}
	for _, element := range compatible_list {
		if q.Repository.LicenseInfo.Key == element {
			return 1.0
		}
	}
	return 0.0
}
