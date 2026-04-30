package github

import (
	"strings"
	"testing"
)

func TestIssueStates(t *testing.T) {
	cases := []struct{ in, want string }{
		{"open", "OPEN"},
		{"closed", "CLOSED"},
		{"all", "OPEN,CLOSED"},
		{"", "OPEN,CLOSED"},
	}
	for _, c := range cases {
		got := strings.Join(issueStates(c.in), ",")
		if got != c.want {
			t.Errorf("issueStates(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPRStates(t *testing.T) {
	cases := []struct{ in, want string }{
		{"open", "OPEN"},
		{"closed", "CLOSED"},
		{"merged", "MERGED"},
		{"all", "OPEN,CLOSED,MERGED"},
		{"", "OPEN,CLOSED,MERGED"},
	}
	for _, c := range cases {
		got := strings.Join(prStates(c.in), ",")
		if got != c.want {
			t.Errorf("prStates(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestBuildIssueSearchQuery(t *testing.T) {
	cases := []struct {
		opts IssueListOptions
		want []string // substrings that must appear
	}{
		{
			IssueListOptions{State: "open", Author: "alice"},
			[]string{"repo:owner/repo", "is:issue", "is:open", "author:alice"},
		},
		{
			IssueListOptions{State: "closed", Assignee: "bob", Labels: []string{"bug", "p1"}},
			[]string{"is:closed", "assignee:bob", `label:"bug"`, `label:"p1"`},
		},
		{
			IssueListOptions{Milestone: "v2.0", Mention: "carol"},
			[]string{`milestone:"v2.0"`, "mentions:carol"},
		},
		{
			IssueListOptions{App: "my-bot", Search: "memory leak"},
			[]string{"author:app/my-bot", "memory leak"},
		},
		{
			IssueListOptions{State: "all"},
			[]string{"is:issue"},
			// "all" must NOT inject a state qualifier
		},
	}
	for _, c := range cases {
		q := buildIssueSearchQuery("owner", "repo", c.opts)
		for _, sub := range c.want {
			if !strings.Contains(q, sub) {
				t.Errorf("buildIssueSearchQuery: %q not found in %q", sub, q)
			}
		}
		if c.opts.State == "all" && (strings.Contains(q, "is:open") || strings.Contains(q, "is:closed")) {
			t.Errorf("buildIssueSearchQuery: state=all should not add is:open/is:closed, got %q", q)
		}
	}
}

func TestBuildPRSearchQuery(t *testing.T) {
	cases := []struct {
		opts PRListOptions
		want []string
	}{
		{
			PRListOptions{State: "open", Author: "alice", Base: "main"},
			[]string{"is:pr", "is:open", "author:alice", "base:main"},
		},
		{
			PRListOptions{State: "merged", Draft: true},
			[]string{"is:merged", "draft:true"},
		},
		{
			PRListOptions{Head: "feature-x", App: "ci-bot"},
			[]string{"head:feature-x", "author:app/ci-bot"},
		},
		{
			PRListOptions{State: "all"},
			[]string{"is:pr"},
		},
	}
	for _, c := range cases {
		q := buildPRSearchQuery("owner", "repo", c.opts)
		for _, sub := range c.want {
			if !strings.Contains(q, sub) {
				t.Errorf("buildPRSearchQuery: %q not found in %q", sub, q)
			}
		}
	}
}
