package mockserver

import "github.com/tomzxcode/ghx/internal/github"

type ScenarioIssue struct {
	Owner string
	Repo  string
	github.Issue
}

type ScenarioPR struct {
	Owner string
	Repo  string
	github.PullRequest
}

type Scenario struct {
	Issues []ScenarioIssue
	PRs    []ScenarioPR
}

type ScenarioOption func(*Scenario)

func WithScenarioIssue(owner, repo string, issue github.Issue) ScenarioOption {
	return func(s *Scenario) {
		s.Issues = append(s.Issues, ScenarioIssue{Owner: owner, Repo: repo, Issue: issue})
	}
}

func WithScenarioPR(owner, repo string, pr github.PullRequest) ScenarioOption {
	return func(s *Scenario) {
		s.PRs = append(s.PRs, ScenarioPR{Owner: owner, Repo: repo, PullRequest: pr})
	}
}

func NewScenario(opts ...ScenarioOption) *Scenario {
	s := &Scenario{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
