package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Correctness_score(personal_token string, owner string, repo string) float64 {
	req, _ := http.NewRequest("GET", "https://api.github.com/search/issues?q=is:pr+repo:"+owner+"/"+repo, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer res.Body.Close()

	var pullRequests SpecCount
	json.NewDecoder(res.Body).Decode(&pullRequests)

	numPullRequests := pullRequests.TotalCount
	score := float64(numPullRequests)
	maxValue := 100.0 //max num of PR
	minValue := 10.0  //min num of PR
	maxScore := 1.0
	minScore := 0.0

	if score <= minValue {
		score = minScore
	} else if score >= maxValue {
		score = maxScore
	} else {
		normalizedValue := (score - minValue) / (maxValue - minValue)
		score = minScore + normalizedValue*(maxScore-minScore)
	}
	return score
}
