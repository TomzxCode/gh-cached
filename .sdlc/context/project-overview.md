# Project Overview

## Purpose

gh-cached is a command-line tool that wraps the GitHub GraphQL API to retrieve issues, pull requests, and comments, caching all results as JSON files on disk.
Its primary goal is to minimise GitHub API calls for users who frequently browse the same repositories, enabling fast offline-capable queries after an initial cache warm-up.
The tool is positioned as a lightweight alternative to `gh issue list` / `gh pr list` for batch or repeated queries on the same data set.

## Key Stakeholders

| Stakeholder | Role | Interest |
|---|---|---|
| Developer / maintainer | Author (Tom Rochette) | Project direction, code quality, feature completeness |
| CLI user | End user | Fast, reliable issue and PR browsing with minimal API usage |
| AI agents / automation | Consumer | Cached data as a low-cost data source for agentic workflows |

## Scope

**In scope:**
- Listing and viewing GitHub issues and pull requests via GraphQL
- On-disk JSON cache with configurable freshness duration
- Delta (incremental) cache updates based on last-cached timestamp
- In-memory filtering of cached data (state, author, assignee, labels, milestone, search)
- Multi-platform binary builds (darwin, linux, windows; amd64 + arm64)
- GitHub Enterprise Server support via `[HOST/]OWNER/REPO` syntax
- Mock GitHub GraphQL server (`mock serve` command) for testing without real API access
- Configurable simulation generator for bulk test scenarios with time-based activity patterns

**Out of scope:**
- Creating, updating, or deleting issues or PRs (read-only)
- Real-time streaming or webhook-based updates
- Authentication management (delegates to `GH_TOKEN`, `GITHUB_TOKEN`, or `gh auth`)
- Caching data beyond issues and PRs (e.g. actions, releases, discussions)

## Key Constraints

- Read-only access to GitHub via GraphQL API; no REST fallback
- Cache is file-based (JSON), stored at `~/.cache/gh-cached/`
- Single binary, no external runtime dependencies beyond a GitHub token
- Go 1.21 minimum version
- GitHub API rate limits apply when cache is cold or stale
- No pagination on comment fetching (capped at 100 comments per item in bulk cache)

