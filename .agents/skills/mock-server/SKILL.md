---
name: mock-server
description: >
  Spin up an in-process mock GitHub GraphQL server with configurable scenarios
  for testing gh-cached without hitting the real GitHub API. Trigger when the
  user wants to test gh-cached, write integration tests, or simulate GitHub
  activity over time.
---

# Mock Server Skill

The mock server is an `httptest.Server` backed by the `internal/mockserver`
package. It accepts the same GraphQL queries that `github.Client` sends,
returning deterministic data from an in-memory scenario. No network access, no
tokens, no rate limits.

---

## Prerequisites

The project must be on disk and buildable:

```bash
go build ./...
```

The relevant packages:

| Package | Purpose |
|---|---|
| `internal/mockserver` | Server, scenario builder, simulation generator |
| `internal/github` | `Client`, `NewClientWithURL` |
| `internal/cache` | `Store`, `NewStoreWithPath` |

---

## Architecture

```
Test code
  |
  |  github.NewClientWithURL(srv.URL(), "test-token", "mock")
  v
github.Client  --(GraphQL POST)-->  mockserver.Server
  |                                    |
  |  internal/github types              |  internal/github types
  |                                    |
  v                                    v
Production code                    Scenario (in-memory)
```

The mock server stores `ScenarioIssue` and `ScenarioPR` structs that embed the
same `github.Issue` / `github.PullRequest` types used everywhere else. No
translation layer is needed.

---

## Workflow

### 1. Choose how to build the scenario

There are three approaches, from simple to full-featured.

#### A. Fluent builder (hand-crafted data)

Best for: small, specific test cases with known titles and numbers.

```go
import "github.com/tomzxcode/gh-cached/internal/mockserver"

scenario := mockserver.NewScenarioBuilder("acme", "myproject").
    WithSeed(42).
    AddIssue("Bug: crash on start", "App crashes", 10*24*time.Hour,
        mockserver.WithIssueState("CLOSED"),
        mockserver.WithIssueAssignee("bob"),
        mockserver.WithIssueLabels("bug", "p0"),
        mockserver.WithIssueComment("alice", "Reproduced.", 9*24*time.Hour),
        mockserver.WithIssueComment("bob", "Fixed.", 8*24*time.Hour),
    ).
    AddPR("Fix crash", "Fixes #1.", "fix/crash", 8*24*time.Hour,
        mockserver.WithPRState("MERGED"),
        mockserver.WithPRReview("APPROVED"),
    ).
    Build()
```

The `age` parameter (e.g. `10*24*time.Hour`) sets how far in the past the item
was created. All timestamps are derived from it.

Available issue options:
- `WithIssueState("OPEN"|"CLOSED")`
- `WithIssueAssignee(login)`
- `WithIssueLabels("bug", "p1", ...)`
- `WithIssueMilestone("v2.0")`
- `WithIssueComment(author, body, age)`

Available PR options:
- `WithPRState("OPEN"|"CLOSED"|"MERGED")`
- `WithPRDraft(true)`
- `WithPRReview("APPROVED"|"CHANGES_REQUESTED"|"REVIEW_REQUIRED")`
- `WithPRLabels("enhancement", ...)`
- `WithPRMilestone("v2.0")`
- `WithPRComment(author, body, age)`

#### B. Simulation generator (parameterized bulk data)

Best for: stress testing, pagination, large-scale scenarios.

```go
scenario := mockserver.Generate(mockserver.SimulationConfig{
    NumUsers:          8,
    Repos:             []string{"acme/platform", "acme/frontend"},
    History:           90 * 24 * time.Hour,
    IssuesPerRepo:     40,
    PRsPerRepo:        30,
    CommentsPerIssue:  4,
    CommentsPerPR:     3,
    AssigneesPerIssue: 2,
    AssigneesPerPR:    2,
    LabelsPerItem:     3,
    MilestonesPerRepo: 4,
    CloseRate:         0.6,   // 60% of issues get closed
    MergeRate:         0.7,   // 70% of PRs get merged
    DraftRate:         0.15,  // 15% of PRs are drafts
    ReviewRate:        0.8,   // 80% of PRs get a review decision
    Seed:              42,
    ActivityBursts:    3,     // 3 high-activity windows
})
```

Use `mockserver.DefaultConfig()` for a sensible preset, or
`mockserver.SmallConfig()` for fast tests. Call
`mockserver.SimulationStats(cfg)` to preview what a config would generate.

All events are placed on a timeline and processed chronologically so
timestamps are internally consistent. The `Seed` makes output deterministic.

`ActivityBursts` concentrates events in short windows (simulates hackathons,
release sprints) by weighting the time sampling distribution.

#### C. Time-advance (simulate activity evolving)

Best for: testing delta updates, cache invalidation, chronological workflows.

```go
now := time.Now()
b := mockserver.NewScenarioBuilder("acme", "myproject").WithNow(now)

scenario := b.
    AdvanceTime(time.Hour).
    NewIssue("First issue", "Body").
    NewIssue("Second issue", "Body").
    CommentOnIssue(1, "bob", "A comment").
    NewPR("First PR", "Body", "feat/x").
    CommentOnPR(3, "alice", "Reviewing").
    MergePR(3).
    Build()
```

`AdvanceTime(delta)` returns a `*TimeAdvance` with methods:
- `NewIssue(title, body, opts...)`
- `NewPR(title, body, headBranch, opts...)`
- `CommentOnIssue(number, author, body)`
- `CommentOnPR(number, author, body)`
- `CloseIssue(number)`
- `MergePR(number)`
- `ClosePR(number)`
- `Build() *Scenario`

Each call advances the internal clock by `delta`. Issue and PR numbers are
sequential across the entire builder (issues and PRs share one counter).

### 2. Start the server

```go
srv := mockserver.NewServer(scenario)
defer srv.Close()
```

The server is a standard `httptest.Server`. It requires an `Authorization`
header on every request (any non-empty value works).

### 3. Connect the real client

```go
client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
if err != nil {
    log.Fatal(err)
}
```

`NewClientWithURL` bypasses token resolution and endpoint computation. The
third argument (`"mock"`) is used as the host for cache directory separation.

All standard `Client` methods work against the mock:

```go
issues, _ := client.FetchAllIssues("acme", "myproject", nil)
issue, _   := client.GetIssue("acme", "myproject", 1)
prs, _     := client.FetchAllPRs("acme", "myproject", nil)
pr, _      := client.GetPR("acme", "myproject", 4)
issues, _  := client.ListIssues("acme", "myproject", github.IssueListOptions{State: "open"})
prs, _     := client.ListPRs("acme", "myproject", github.PRListOptions{State: "merged"})
```

### 4. Use an isolated cache (optional)

```go
store := cache.NewStoreWithPath(t.TempDir())
```

`NewStoreWithPath` avoids polluting `~/.cache/gh-cached` during tests.

### 5. Evolve the scenario at runtime (optional)

```go
srv.UpdateScenario(func(s *mockserver.Scenario) {
    s.Issues = append(s.Issues, mockserver.ScenarioIssue{
        Owner: "acme",
        Repo:  "myproject",
        Issue: github.Issue{
            Number:    100,
            Title:     "Newly created issue",
            State:     "OPEN",
            Author:    github.Actor{Login: "alice"},
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
    })
})
```

`UpdateScenario` is concurrency-safe. The next request to the server will see
the new data.

### 6. CLI usage (end-to-end testing)

The `gh-cached` binary supports `--api-url` and `--cache-dir` flags:

```bash
# In one terminal (or via Go test):
#   Start mock server, print its URL

# In another terminal:
gh-cached --api-url http://127.0.0.1:XXXXX --cache-dir /tmp/test-cache --repo acme/myproject cache
gh-cached --api-url http://127.0.0.1:XXXXX --cache-dir /tmp/test-cache --repo acme/myproject issue list --state all
gh-cached --api-url http://127.0.0.1:XXXXX --cache-dir /tmp/test-cache --repo acme/myproject pr list --json
```

---

## Supported queries

The mock server handles all 7 GraphQL query patterns used by `github.Client`:

| Pattern | Trigger | Notes |
|---|---|---|
| Get issue | `issue(number:` | Full comments included |
| Get PR | `pullRequest(number:` | Full comments included |
| List issues | `issues(first:` (not `states: [OPEN, CLOSED]`) | Summary comments (count only) |
| Fetch all issues | `issues(first:` + `states: [OPEN, CLOSED]` | Full comments, `since` filter |
| List PRs | `pullRequests(first:` (not `states: [OPEN, CLOSED, MERGED]`) | Summary comments |
| Fetch all PRs | `pullRequests(first: 100` + `states: [OPEN, CLOSED, MERGED]` | Full comments |
| Search | `search(query:` | Supports repo/author/label/state/draft/base/head filters |

All responses include `pageInfo` with `hasNextPage` and `endCursor` for
pagination support.

---

## Complete test example

```go
package mypackage_test

import (
    "testing"
    "time"

    "github.com/tomzxcode/gh-cached/internal/cache"
    "github.com/tomzxcode/gh-cached/internal/github"
    "github.com/tomzxcode/gh-cached/internal/mockserver"
)

func TestEndToEnd(t *testing.T) {
    // 1. Build scenario.
    scenario := mockserver.Generate(mockserver.SmallConfig())

    // 2. Start mock server.
    srv := mockserver.NewServer(scenario)
    defer srv.Close()

    // 3. Connect client.
    client, err := github.NewClientWithURL(srv.URL(), "test-token", "mock")
    if err != nil {
        t.Fatal(err)
    }

    // 4. Fetch and cache.
    store := cache.NewStoreWithPath(t.TempDir())

    issues, err := client.FetchAllIssues("acme", "testrepo", nil)
    if err != nil {
        t.Fatal(err)
    }
    for _, issue := range issues {
        store.SaveIssue("mock", "acme", "testrepo", issue)
    }

    prs, err := client.FetchAllPRs("acme", "testrepo", nil)
    if err != nil {
        t.Fatal(err)
    }
    for _, pr := range prs {
        store.SavePR("mock", "acme", "testrepo", pr)
    }

    store.SaveCacheInfo("mock", "acme", "testrepo", 60)

    // 5. Verify cache round-trip.
    loaded, _, err := store.LoadIssue("mock", "acme", "testrepo", issues[0].Number)
    if err != nil {
        t.Fatal(err)
    }
    if loaded.Title != issues[0].Title {
        t.Errorf("title mismatch: %q vs %q", loaded.Title, issues[0].Title)
    }
}
```

---

## Tips

- Use `WithSeed(n)` for deterministic output; same seed always produces the
  same data.
- Number counters are shared across issues and PRs in a single builder. If you
  add 3 issues then a PR, the PR gets number 4.
- The `age` parameter in `AddIssue`/`AddPR` sets `createdAt = now - age`. All
  comments and state changes derive their timestamps from the item's creation.
- `UpdateScenario` holds a write lock. Keep the callback short.
- `Scenario.Summary()` returns a human-readable string useful for logging.
