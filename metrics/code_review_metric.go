package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func code_review_metric(personal_token string, owner string, repo string) float64 {

	type Response struct {
		Data struct {
			Repository struct {
				PullRequests struct {
					Nodes []struct {
						Reviews struct {
							TotalCount int `json:"totalCount"`
						} `json:"reviews"`
					} `json:"nodes"`
				} `json:"pullRequests"`
			} `json:"repository"`
		} `json:"data"`
	}

	url := "https://api.github.com/graphql"
	mySecret := personal_token

	payload := strings.NewReader(fmt.Sprintf("{\"query\":\"{\\nrepository(owner: \\\"%s\\\", name: \\\"%s\\\") {\\npullRequests(last: 100) {\\nnodes {\\nreviews(first: 1) {\\ntotalCount\\n}\\n}\\n}\\n}\\n}\"}", owner, repo))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Error sending API Request:", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+mySecret)

	res, _ := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println("res error", err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var resp Response

	err = json.Unmarshal([]byte(body), &resp)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
	}

	var numReviewed int = 0
	var score float64

	for _, node := range resp.Data.Repository.PullRequests.Nodes {
		if node.Reviews.TotalCount != 0 {
			numReviewed++
		}
	}
	score = (float64(numReviewed) / float64(len(resp.Data.Repository.PullRequests.Nodes)))
	// fmt.Println(score)
	return score
}
