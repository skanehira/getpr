package main

import (
	"reflect"
	"testing"
)

func TestRemoteToURL(t *testing.T) {
	want := Repo{
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
			if reflect.DeepEqual(got, want) {
				t.Errorf("want=%q, got=%q", want, got)
			}
		})
	}
}

func TestParseBlame(t *testing.T) {
	tests := []struct {
		line string
		want []Blame
	}{
		{
			line: `^737302e (skanehira 2020-11-02 23:28:33 +0900  1) # github-blame
^737302e (skanehira 2020-11-02 23:28:33 +0900  2) Get GitHub's pull request URL.`,
			want: []Blame{
				{
					ID:   "737302e",
					Text: " (skanehira 2020-11-02 23:28:33 +0900  1) # github-blame",
				},
				{
					ID:   "737302e",
					Text: " (skanehira 2020-11-02 23:28:33 +0900  2) Get GitHub's pull request URL.",
				},
			},
		},
		{
			line: `7ef8fc2f (skanehira         2019-04-05 23:49:29 +0900 43) func main() {
7ef8fc2f (skanehira         2019-04-05 23:49:29 +0900 45)       os.Exit(run())`,
			want: []Blame{
				{
					ID:   "7ef8fc2f",
					Text: " (skanehira         2019-04-05 23:49:29 +0900 43) func main() {",
				},
				{
					ID:   "7ef8fc2f",
					Text: " (skanehira         2019-04-05 23:49:29 +0900 45)       os.Exit(run())",
				},
			},
		},
	}

	for _, tt := range tests {
		got := parseBlame(tt.line)
		if reflect.DeepEqual(tt.want, got) {
			t.Errorf("want=%q, got=%q", tt.want, got)
		}
	}
}
