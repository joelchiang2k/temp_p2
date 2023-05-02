package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
type Dependencies map[string]string

type PackageJson struct {
	Dependencies Dependencies `json:"dependencies"`
}

func Dependency_score(owner string, repo string) float64 {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/package.json", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err, "line 34")
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err, "line 41")
		os.Exit(1)
	}

	var packageJson PackageJson

	err = json.Unmarshal(body, &packageJson)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err, "line 48")
		os.Exit(1)
	}

	count := 0
	for _, versionStr := range packageJson.Dependencies {
		if MajorMinor(versionStr) {
			count++
		}
	}
	if len(packageJson.Dependencies) == 0 {
		return 1
	}

	preround := float64(count) / float64(len(packageJson.Dependencies))
	rounding := math.Pow(10, 2)
	return math.Round(preround*rounding) / rounding //double formatted return

}

func MajorMinor(versionStr string) bool {
	parts := strings.Split(versionStr, ".")
	if len(parts) < 2 {
		return false
	}
	_, err := strconv.Atoi(parts[0])
	if err == nil {
		return true
	}

	_, err = strconv.Atoi(parts[1])
	if err == nil {
		return true
	}

	return false

}
