package blame

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ktr0731/go-fuzzyfinder"
)

var (
	mergePattern  = regexp.MustCompile(`Merge pull request #(\d+)`)
	squashPattern = regexp.MustCompile(`\(#(\d+)\)`)
)

type Blame struct {
	ID   string
	Text string
}

func parseBlame(text string) []Blame {
	blames := []Blame{}
	lines := strings.Split(text, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		cols := strings.SplitN(lines[i], " ", 2)
		// id can be ^xxxx
		id := strings.TrimLeft(cols[0], "^")
		blames = append(blames, Blame{
			ID:   id,
			Text: cols[1],
		})
	}
	return blames
}

func GetCommitID(file string) (string, error) {
	cmd := exec.Command("git", "blame", file)
	out, err := cmd.CombinedOutput()
	result := strings.TrimRight(string(out), "\r\n")
	if err != nil {
		return "", err
	}

	blames := parseBlame(result)
	idx, err := fuzzyfinder.Find(
		blames,
		func(i int) string {
			return fmt.Sprintf("%s %s", blames[i].ID, blames[i].Text)
		},
	)

	return blames[idx].ID, nil
}

func GetPRNumber(id string) (string, error) {
	cmd := exec.Command("git", "log", "--oneline", "-n", "1", id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	result := strings.TrimRight(string(out), "\r\n")

	s := mergePattern.FindStringSubmatch(result)
	if len(s) > 0 {
		return s[1], nil
	}

	s = squashPattern.FindStringSubmatch(result)
	if len(s) > 0 {
		return s[1], nil
	}
	return "", fmt.Errorf("cannot get pull request number from commit message %q", result)
}

func GetPRURL(file string) (string, error) {
	id, err := GetCommitID(file)
	if err != nil {
		return "", err
	}
	number, err := GetPRNumber(id)
	if err != nil {
		return "", err
	}
	url, err := GetRepoURL()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/pull/%s", url, number), nil
}

func GetRepoURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "--push", "origin")
	out, err := cmd.CombinedOutput()

	result := strings.TrimRight(string(out), "\r\n")
	if err != nil {
		return "", err
	}

	return RemoteToURL(result), nil
}

func RemoteToURL(remote string) string {
	// https://stackoverflow.com/questions/31801271/what-are-the-supported-git-url-formats
	url := "https://github.com/"
	if strings.HasPrefix(remote, "ssh") {
		parts := strings.Split(remote, "/")
		ownerRepo := strings.Join(parts[len(parts)-2:], "/")
		url += ownerRepo
	} else if strings.HasPrefix(remote, "git") {
		parts := strings.Split(remote, ":")
		if len(parts) < 1 {
			log.Fatalf("invalid remote: %s", remote)
		}
		url += parts[1]
	} else {
		url = remote
	}

	return strings.TrimRight(url, ".git")
}
