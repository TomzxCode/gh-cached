---
title: "Issue Browsing"
status: draft
---

# Specification: Issue Browsing

## Overview

Issue browsing is implemented through two Cobra subcommands (`issue list`, `issue view`) in `cmd/issue.go`.
Filtering logic is performed in-memory on cached data or via GitHub GraphQL search when the cache is stale.
Display is handled by dedicated print functions that output tabular or JSON format.

## Architecture

```
issue list ──► runIssueList()
                    │
                    ├─ Cache fresh? ──Yes──► LoadAllIssues() ──► filterIssues() ──► printIssueList()
                    │
                    └─ No ──► NewClient() ──► ListIssues() ──► printIssueList()

issue view N ──► runIssueView()
                    │
                    ├─ Cache fresh? ──Yes──► LoadIssue(N) ──► printIssueView()
                    │
                    ├─ Individual file fresh (<60m)? ──Yes──► printIssueView()
                    │
                    └─ API ──► GetIssue(N) ──► SaveIssue() ──► printIssueView()
```

## Data Models

### Issue

| Field | Type | Constraints | Description |
|---|---|---|---|
| number | int | not null | GitHub issue number |
| title | string | not null | Issue title |
| state | string | not null | OPEN or CLOSED |
| author | Actor | not null | Issue author |
| assignees | []Actor | | Assigned users |
| labels | []Label | | Issue labels |
| milestone | *Milestone | nullable | Associated milestone |
| createdAt | time.Time | not null | Creation timestamp |
| updatedAt | time.Time | not null | Last update timestamp |
| closedAt | *time.Time | nullable | Closure timestamp |
| url | string | not null | GitHub URL |
| body | string | | Issue body text |
| commentCount | int | | Total comment count from API |
| comments | []Comment | | Inline comments (populated on view/fetch) |

## API Contracts

### GraphQL: `listIssuesQuery`

Fetches issues with server-side filtering (states, assignee, author, labels, milestone, mention).
Paginated with cursor-based pagination (100 items per page).

### GraphQL: `searchIssuesQuery`

Used when `--app` or `--search` flags are provided.
Builds a GitHub search query string from filter options.

### GraphQL: `getIssueQuery`

Fetches a single issue with up to 100 inline comments.

## Technical Decisions

| Decision | Choice | Rationale |
|---|---|---|
| In-memory filtering | Filter cached data in Go code | Avoids API calls for simple filter combinations |
| Case-insensitive matching | `strings.EqualFold` | GitHub filter behaviour is case-insensitive |
| Label AND logic | All specified labels must be present | Matches GitHub API behaviour |

## Risks and Unknowns

1. The `mention` and `app` filters cannot be evaluated from cached data alone; they are silently skipped when serving from cache

## Out of Scope

- Creating or updating issues
- Streaming or real-time issue updates
- Filtering by reactions, projects, or custom fields
