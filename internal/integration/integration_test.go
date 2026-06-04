package integration

import (
	"testing"
	"time"

	"github.com/tomzxcode/gh-cached/internal/cache"
	"github.com/tomzxcode/gh-cached/internal/github"
	"github.com/tomzxcode/gh-cached/internal/mockserver"
)

func sampleScenario() *mockserver.Scenario {
	return mockserver.NewScenarioBuilder("acme", "myproject").
		WithSeed(42).
		AddIssue("Bug: crash on start", "App crashes on startup", 10*24*time.Hour,
			mockserver.WithIssueState("CLOSED"),
			mockserver.WithIssueAssignee("bob"),
			mockserver.WithIssueLabels("bug", "p0"),
			mockserver.WithIssueComment("alice", "Reproduced.", 9*24*time.Hour),
			mockserver.WithIssueComment("bob", "Fixed.", 8*24*time.Hour),
		).
		AddIssue("Feature: dark mode", "Add dark mode", 5*24*time.Hour,
			mockserver.WithIssueLabels("enhancement"),
			mockserver.WithIssueComment("carol", "Working on it.", 4*24*time.Hour),
		).
		AddIssue("Docs: update README", "README is outdated", 2*24*time.Hour,
			mockserver.WithIssueLabels("documentation"),
		).
		AddIssue("Fix login bug", "Login fails on Safari", 3*24*time.Hour,
			mockserver.WithIssueLabels("bug"),
			mockserver.WithIssueComment("alice", "Can reproduce.", 3*time.Hour),
		).
		AddIssue("Add CI pipeline", "Need automated tests", 1*24*time.Hour,
			mockserver.WithIssueLabels("enhancement", "ci"),
		).
		AddPR("Fix crash on start", "Fixes the crash.", "fix/crash", 8*24*time.Hour,
			mockserver.WithPRState("MERGED"),
			mockserver.WithPRReview("APPROVED"),
			mockserver.WithPRComment("alice", "LGTM.", 8*24*time.Hour),
		).
		AddPR("Add dark mode", "Adds dark theme.", "feat/dark-mode", 3*24*time.Hour,
			mockserver.WithPRLabels("enhancement"),
			mockserver.WithPRComment("bob", "Looking good!", 2*24*time.Hour),
		).
		AddPR("WIP: refactor auth", "Refactoring.", "refactor/auth", 1*24*time.Hour,
			mockserver.WithPRDraft(true),
		).
		Build()
}

func startServer(t *testing.T, scenario *mockserver.Scenario) *mockserver.Server {
	t.Helper()
	srv := mockserver.NewServer(scenario)
	t.Cleanup(srv.Close)
	return srv
}

func TestClient_FetchAllIssues(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	issues, err := client.FetchAllIssues("acme", "myproject", nil)
	if err != nil {
		t.Fatalf("FetchAllIssues: %v", err)
	}
	if len(issues) != 5 {
		t.Errorf("got %d issues, want 5", len(issues))
	}

	for _, issue := range issues {
		if issue.Number <= 0 {
			t.Errorf("issue number = %d, want > 0", issue.Number)
		}
		if issue.CreatedAt.IsZero() {
			t.Errorf("issue #%d has zero CreatedAt", issue.Number)
		}
	}
}

func TestClient_FetchAllPRs(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	prs, err := client.FetchAllPRs("acme", "myproject", nil)
	if err != nil {
		t.Fatalf("FetchAllPRs: %v", err)
	}
	if len(prs) != 3 {
		t.Errorf("got %d PRs, want 3", len(prs))
	}
}

func TestClient_GetIssue(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	issue, err := client.GetIssue("acme", "myproject", 1)
	if err != nil {
		t.Fatalf("GetIssue: %v", err)
	}
	if issue.Number != 1 {
		t.Errorf("Number = %d, want 1", issue.Number)
	}
	if issue.Title != "Bug: crash on start" {
		t.Errorf("Title = %q", issue.Title)
	}
	if issue.State != "CLOSED" {
		t.Errorf("State = %q, want CLOSED", issue.State)
	}
	if len(issue.Comments) != 2 {
		t.Errorf("Comments = %d, want 2", len(issue.Comments))
	}
}

func TestClient_GetPR(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	pr, err := client.GetPR("acme", "myproject", 6)
	if err != nil {
		t.Fatalf("GetPR: %v", err)
	}
	if pr.Number != 6 {
		t.Errorf("Number = %d, want 6", pr.Number)
	}
	if pr.State != "MERGED" {
		t.Errorf("State = %q, want MERGED", pr.State)
	}
	if pr.BaseRefName != "main" {
		t.Errorf("BaseRefName = %q, want main", pr.BaseRefName)
	}
}

func TestClient_ListIssues(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	issues, err := client.ListIssues("acme", "myproject", github.IssueListOptions{
		State: "open",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListIssues: %v", err)
	}
	if len(issues) != 4 {
		t.Errorf("got %d open issues, want 4", len(issues))
	}
}

func TestClient_ListPRs(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	prs, err := client.ListPRs("acme", "myproject", github.PRListOptions{
		State: "open",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("ListPRs: %v", err)
	}
	if len(prs) != 2 {
		t.Errorf("got %d open PRs, want 2", len(prs))
	}
}

func TestClient_GetIssue_NotFound(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	_, err = client.GetIssue("acme", "myproject", 999)
	if err == nil {
		t.Error("expected error for missing issue")
	}
}

func TestClient_DeltaFetch(t *testing.T) {
	cfg := mockserver.SmallConfig()
	scenario := mockserver.Generate(cfg)
	srv := startServer(t, scenario)

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	since := time.Now().Add(-1 * time.Hour)
	issues, err := client.FetchAllIssues("acme", "testrepo", &since)
	if err != nil {
		t.Fatalf("FetchAllIssues with since: %v", err)
	}

	for _, issue := range issues {
		if issue.UpdatedAt.Before(since) {
			t.Errorf("issue #%d updatedAt=%v is before since=%v", issue.Number, issue.UpdatedAt, since)
		}
	}
}

func TestClient_FetchAndCache(t *testing.T) {
	srv := startServer(t, sampleScenario())

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	store := cache.NewStoreWithPath(t.TempDir())

	issues, err := client.FetchAllIssues("acme", "myproject", nil)
	if err != nil {
		t.Fatalf("FetchAllIssues: %v", err)
	}
	for _, issue := range issues {
		if err := store.SaveIssue("mock", "acme", "myproject", issue); err != nil {
			t.Fatalf("SaveIssue #%d: %v", issue.Number, err)
		}
	}

	prs, err := client.FetchAllPRs("acme", "myproject", nil)
	if err != nil {
		t.Fatalf("FetchAllPRs: %v", err)
	}
	for _, pr := range prs {
		if err := store.SavePR("mock", "acme", "myproject", pr); err != nil {
			t.Fatalf("SavePR #%d: %v", pr.Number, err)
		}
	}

	if err := store.SaveCacheInfo("mock", "acme", "myproject", 60); err != nil {
		t.Fatalf("SaveCacheInfo: %v", err)
	}

	loaded, _, err := store.LoadIssue("mock", "acme", "myproject", 1)
	if err != nil {
		t.Fatalf("LoadIssue: %v", err)
	}
	if loaded.Number != 1 {
		t.Errorf("loaded issue number = %d, want 1", loaded.Number)
	}
	if loaded.Title != "Bug: crash on start" {
		t.Errorf("loaded title = %q", loaded.Title)
	}

	allIssues, err := store.LoadAllIssues("mock", "acme", "myproject")
	if err != nil {
		t.Fatalf("LoadAllIssues: %v", err)
	}
	if len(allIssues) != 5 {
		t.Errorf("cached %d issues, want 5", len(allIssues))
	}

	allPRs, err := store.LoadAllPRs("mock", "acme", "myproject")
	if err != nil {
		t.Fatalf("LoadAllPRs: %v", err)
	}
	if len(allPRs) != 3 {
		t.Errorf("cached %d PRs, want 3", len(allPRs))
	}

	fresh, err := store.IsCacheFresh("mock", "acme", "myproject")
	if err != nil {
		t.Fatalf("IsCacheFresh: %v", err)
	}
	if !fresh {
		t.Error("cache should be fresh")
	}
}

func TestClient_GenerateLargeAndFetch(t *testing.T) {
	cfg := mockserver.SimulationConfig{
		NumUsers:          5,
		Repos:             []string{"acme/testrepo"},
		History:           30 * 24 * time.Hour,
		IssuesPerRepo:     20,
		PRsPerRepo:        15,
		CommentsPerIssue:  3,
		CommentsPerPR:     2,
		CloseRate:         0.5,
		MergeRate:         0.6,
		DraftRate:         0.1,
		ReviewRate:        0.8,
		Seed:              1,
		AssigneesPerIssue: 1,
		AssigneesPerPR:    1,
		LabelsPerItem:     2,
		MilestonesPerRepo: 2,
	}
	scenario := mockserver.Generate(cfg)
	srv := startServer(t, scenario)

	client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
	if err != nil {
		t.Fatalf("NewClientWithURL: %v", err)
	}

	issues, err := client.FetchAllIssues("acme", "testrepo", nil)
	if err != nil {
		t.Fatalf("FetchAllIssues: %v", err)
	}
	if len(issues) != 20 {
		t.Errorf("got %d issues, want 20", len(issues))
	}

	prs, err := client.FetchAllPRs("acme", "testrepo", nil)
	if err != nil {
		t.Fatalf("FetchAllPRs: %v", err)
	}
	if len(prs) != 15 {
		t.Errorf("got %d PRs, want 15", len(prs))
	}

	totalComments := 0
	for _, issue := range issues {
		totalComments += len(issue.Comments)
	}
	for _, pr := range prs {
		totalComments += len(pr.Comments)
	}
	if totalComments == 0 {
		t.Error("expected some comments")
	}
}
