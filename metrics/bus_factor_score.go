package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SpecCount struct {
	TotalCount int `json:"total_count"`
}

func Bus_factor_score(personal_token string, owner string, repo string) float32 {
	req, _ := http.NewRequest("GET", "https://api.github.com/search/issues?q=is:pr+repo:"+owner+"/"+repo, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer res.Body.Close()

	var Contributors SpecCount
	json.NewDecoder(res.Body).Decode(&Contributors)

	numContributors := Contributors.TotalCount
	score := float32(numContributors)
	maxBusValue := 100.0 //max num of PR
	minBusValue := 10.0  //min num of PR
	maxBusScore := 1.0
	minBusScore := 0.0

	if float64(score) <= minBusValue {
		score = float32(minBusScore)
	} else if float64(score) >= maxBusValue {
		score = float32(maxBusScore)
	} else {
		normalizedValue := (float64(score) - minBusValue) / (maxBusValue - minBusValue)
		score = float32(minBusScore) + float32(normalizedValue)*(float32(maxBusScore)-float32(minBusScore))
	}
	return score
}
