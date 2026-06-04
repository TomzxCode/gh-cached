---
title: "PR Browsing"
status: draft
---

# Specification: PR Browsing

## Overview

PR browsing mirrors the issue browsing architecture with two Cobra subcommands (`pr list`, `pr view`) in `cmd/pr.go`.
It adds PR-specific fields: branch names, draft status, merge/closed timestamps, and review decision.
The filtering and display logic follows the same cache-first, API-fallback pattern.

## Architecture

```
pr list ──► runPRList()
                 │
                 ├─ Cache fresh? ──Yes──► LoadAllPRs() ──► filterPRs() ──► printPRList()
                 │
                 └─ No ──► NewClient() ──► ListPRs() ──► printPRList()

pr view N ──► runPRView()
                 │
                 ├─ Cache fresh? ──Yes──► LoadPR(N) ──► printPRView()
                 │
                 ├─ Individual file fresh (<60m)? ──Yes──► printPRView()
                 │
                 └─ API ──► GetPR(N) ──► SavePR() ──► printPRView()
```

## Data Models

### PullRequest

| Field | Type | Constraints | Description |
|---|---|---|---|
| number | int | not null | PR number |
| title | string | not null | PR title |
| state | string | not null | OPEN, CLOSED, or MERGED |
| isDraft | bool | not null | Whether the PR is a draft |
| author | Actor | not null | PR author |
| assignees | []Actor | | Assigned users |
| labels | []Label | | PR labels |
| milestone | *Milestone | nullable | Associated milestone |
| baseRefName | string | not null | Target branch |
| headRefName | string | not null | Source branch |
| createdAt | time.Time | not null | Creation timestamp |
| updatedAt | time.Time | not null | Last update timestamp |
| mergedAt | *time.Time | nullable | Merge timestamp |
| closedAt | *time.Time | nullable | Closure timestamp |
| url | string | not null | GitHub URL |
| body | string | | PR body text |
| commentCount | int | | Total comment count |
| comments | []Comment | | Inline comments |
| reviewDecision | string | | APPROVED, CHANGES_REQUESTED, REVIEW_REQUIRED |

## API Contracts

### GraphQL: `listPRsQuery`

Direct listing with server-side filters (states, labels, baseRefName, headRefName).
Used when no author/assignee/app/draft/search filters are active.

### GraphQL: `searchPRsQuery`

Used when author, assignee, app, draft, or search filters are provided.
Builds a GitHub search query string.

### GraphQL: `getPRQuery`

Fetches a single PR with up to 100 inline comments.

## Technical Decisions

| Decision | Choice | Rationale |
|---|---|---|
| Merged as distinct state | `MERGED` is a separate PullRequestState in GraphQL | Allows filtering merged PRs separately from closed |
| Review decision formatting | Human-readable strings in output | `APPROVED` becomes "approved" for tabular display |

## Risks and Unknowns

1. `app` filter cannot be verified from cached data

## Out of Scope

- Creating, updating, or merging PRs
- Review comment threads (only top-level issue-style comments)
- CI/check status on PRs
