package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseRemote(t *testing.T) {
	want := &Repo{
		Owner: "skanehira",
		Name:  "github-blame",
	}
	tests := []struct {
		name   string
		remote string
		want   Repo
	}{
		{
			name:   "ssh-1",
			remote: "ssh://git@github.com/skanehira/github-blame",
		},
		{
			name:   "ssh-2",
			remote: "ssh://git@github.com/skanehira/github-blame.git",
		},
		{
			name:   "ssh-3",
			remote: "ssh://github.com/skanehira/github-blame.git",
		},
		{
			name:   "ssh-4",
			remote: "ssh://github.com/skanehira/github-blame",
		},
		{
			name:   "git-1",
			remote: "git@github.com:skanehira/github-blame.git",
		},
		{
			name:   "git-2",
			remote: "git@github.com:skanehira/github-blame",
		},
		{
			name:   "http-1",
			remote: "http://github.com/skanehira/github-blame",
		},
		{
			name:   "http-2",
			remote: "http://github.com/skanehira/github-blame.git",
		},
		{
			name:   "https-1",
			remote: "https://github.com/skanehira/github-blame",
		},
		{
			name:   "https-2",
			remote: "https://github.com/skanehira/github-blame.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parserRemote(tt.remote)
			if err != nil {
				t.Fatalf("got error: %s", err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("want=%q, got=%q", want, got)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	os.Setenv("GITHUB_TOKEN", "a")

	t.Run("get token from GITHUB_TOKEN", func(t *testing.T) {
		got, err := getToken()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		want := "a"
		if got != want {
			t.Fatalf("got=%s, want=%s", got, want)
		}

		t.Cleanup(func() {
			os.Setenv("GITHUB_TOKEN", "")
		})
	})

	t.Run("get token from $HOME/.github_token", func(t *testing.T) {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		configFile := filepath.Join(homeDir, ".github_token")
		if err := ioutil.WriteFile(configFile, []byte("a\n"), 0777); err != nil {
			t.Fatalf("got error: %s", err)
		}

		token, err := getToken()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}

		want := "a"
		if token != want {
			t.Fatalf("got=%s, want=%s", token, want)
		}

		t.Cleanup(func() {
			os.Remove(configFile)
		})
	})
}
