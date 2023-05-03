package metrics

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Ramp_up_score(personal_token string, url string) (float64, string, string) {
	split := strings.Split(url, "/")
	owner := split[len(split)-2]
	repo := split[len(split)-1]

	tmpDir, err := ioutil.TempDir("", "git-clone") //temp dir
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	_, err = git.PlainClone(tmpDir, false, &git.CloneOptions{ //temp clone
		URL: fmt.Sprintf("https://github.com/%s/%s.git", owner, repo),
		Auth: &http.BasicAuth{
			Username: personal_token,
			Password: "",
		},
	})

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	totalCodeLines := 0
	totalCommentLines := 0
	err = filepath.Walk(tmpDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			codeLines, commentLines, err := CommentsAndCode(path)
			if err != nil {
				return err
			}
			totalCodeLines += codeLines
			totalCommentLines += commentLines
		}
		return nil
	})

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if totalCodeLines == 0 || totalCommentLines == 0 {
		return 0, owner, repo
	}

	ratio := float32(totalCommentLines) / float32(totalCodeLines)

	minRatio := 0
	maxRatio := 1
	if ratio <= float32(minRatio) {
		ratio = 0
	} else if ratio >= float32(maxRatio) {
		ratio = 1
	}

	return float64(ratio), owner, repo
}

func CommentsAndCode(path string) (int, int, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, 0, err
	}

	lines := strings.Split(string(file), "\n")
	codeLines := 0
	commentLines := 0
	multiLineComment := false

	for _, line := range lines {
		slice := strings.TrimSpace(line)
		if strings.HasPrefix(slice, "//") {
			commentLines++
		} else if strings.HasPrefix(slice, "/*") && !strings.HasSuffix(slice, "*/") {
			multiLineComment = true
			commentLines++
		} else if multiLineComment {
			commentLines++
			if strings.HasSuffix(slice, "*/") {
				multiLineComment = false
			}
		} else {
			codeLines++
		}
	}

	return codeLines, commentLines, nil
}
