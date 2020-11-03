package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var version = "0.0.1"

var client *githubv4.Client

var configFile = ".github_token"

type PullRequest struct {
	URL string
}

type Repo struct {
	Owner string
	Name  string
}

func getToken() (string, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		configFile := filepath.Join(homeDir, configFile)
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			return "", err
		}

		token = strings.Trim(string(b), "\r\n")
	}
	if token == "" {
		return "", errors.New("github token is empty")
	}

	return token, nil
}

func getPullRequest(owner, name, id string) (*PullRequest, error) {
	var q struct {
		Repository struct {
			Object struct {
				Commit struct {
					AssociatedPullRequests struct {
						Nodes []PullRequest
					} `graphql:"associatedPullRequests(first: 1)"`
				} `graphql:"... on Commit"`
			} `graphql:"object(expression: $id)"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"id":    githubv4.String(id),
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}

	if err := client.Query(context.Background(), &q, variables); err != nil {
		return nil, err
	}

	if len(q.Repository.Object.Commit.AssociatedPullRequests.Nodes) < 1 {
		return nil, errors.New("not found pull request")
	}

	return &q.Repository.Object.Commit.AssociatedPullRequests.Nodes[0], nil
}

func getOwnerRepo() (*Repo, error) {
	cmd := exec.Command("git", "remote", "get-url", "--push", "origin")
	out, err := cmd.CombinedOutput()

	result := strings.TrimRight(string(out), "\r\n")
	if err != nil {
		return nil, err
	}

	return parserRemote(result)
}

func parserRemote(remote string) (*Repo, error) {
	remote = strings.TrimRight(remote, ".git")
	var ownerRepo []string
	if strings.HasPrefix(remote, "ssh") {
		p := strings.Split(remote, "/")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = p[len(p)-2:]
	} else if strings.HasPrefix(remote, "git") {
		p := strings.Split(remote, ":")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = strings.Split(p[1], "/")
	} else if strings.HasPrefix(remote, "http") || strings.HasPrefix(remote, "https") {
		p := strings.Split(remote, "/")
		if len(p) < 1 {
			return nil, fmt.Errorf("cannot get owner/repo from remote: %s", remote)
		}
		ownerRepo = p[len(p)-2:]
	}

	repo := Repo{
		Owner: ownerRepo[0],
		Name:  ownerRepo[1],
	}

	return &repo, nil
}

func run(args []string, flags *flags) error {
	token, err := getToken()
	if err != nil {
		return errors.New(`cannot get github token
please set GitHub token to GITHUB_TOKEN or $HOME/.github_token`)
	}

	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context.Background(), src)
	client = githubv4.NewClient(httpClient)

	var (
		id    string
		owner string
		name  string
	)
	if len(args) == 1 {
		id = args[0]
		repo, err := getOwnerRepo()
		if err != nil {
			return err
		}
		owner = repo.Owner
		name = repo.Name
	} else if len(args) > 1 {
		id = args[1]
		parts := strings.Split(args[0], "/")
		owner = parts[0]
		name = parts[1]
	}

	pr, err := getPullRequest(owner, name, id)
	if err != nil {
		return err
	}

	fmt.Println(pr.URL)

	if flags.open {
		err := browser.OpenURL(pr.URL)
		if err != nil {
			return fmt.Errorf("Cannot open the URL in the browser, %w", err)
		}
	}
	return nil
}

type flags struct {
	open bool
}

func main() {
	flags := &flags{}

	name := "getpr"
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.BoolVar(&flags.open, "open", false, "Open the generated URL automatically")
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`%[1]s - Get GitHub's Pull Request URL.

VERSION: %s

USAGE:
  $ %[1]s [OWNER/REPO] {commit id}

EXAMPLE:
  $ %[1]s getpr 737302e
  $ %[1]s getpr skanehira/getpr 737302e
`, name, version)
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			return
		}
		os.Exit(1)
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		return
	}

	if err := run(args, flags); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
