package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func ramp_up_score(url string) (float64, string, string) {
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]

	req, _ := http.NewRequest("GET", "https://api.github.com/search/issues?q=is:pr+repo:"+owner+"/"+repo, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer res.Body.Close()

	var pullRequests PullRequests
	json.NewDecoder(res.Body).Decode(&pullRequests)

	numPullRequests := pullRequests.TotalCount
	score := float64(numPullRequests)

	maxPullValue := 20.0
	normalizedValue := maxPullValue / score
	if score <= maxPullValue {
		score = 1
	}
	if score > maxPullValue {
		score = normalizedValue
	}

	return score, owner, repo
}
