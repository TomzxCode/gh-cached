# gh-cached

A GitHub CLI that calls the GitHub GraphQL API to retrieve issues, pull requests, and comments, caching all results to disk to minimise API calls.

Cache lives at `~/.cache/gh-cached/<host>/<owner>/<repo>`.

## Installation

```bash
go install github.com/tomzxcode/gh-cached@main
```

Or download a pre-built binary from the [releases page](https://github.com/TomzxCode/gh-cached/releases/tag/latest).

## Authentication

Set a GitHub personal access token in your environment:

```bash
export GH_TOKEN=ghp_...
```

`GITHUB_TOKEN` is also supported and checked as a fallback after `GH_TOKEN`.

Alternatively, if you have the [GitHub CLI](https://cli.github.com) installed and authenticated, `gh-cached` will use it as a fallback (`gh auth token`).

## Usage

All commands accept `--repo [HOST/]OWNER/REPO`. When omitted, the repository is detected from the `origin` remote of the current directory.

### Cache

Pre-populate the local cache with all issues and PRs (including comments):

```bash
gh-cached cache
gh-cached cache --repo cli/cli
gh-cached cache --cache-duration 120   # treat cache as fresh for 2 hours
gh-cached cache --cache-duration 0     # always re-fetch (delta)
gh-cached cache --force                  # force full re-fetch
```

Example output (first-time run):

```
Caching issues for octocat/hello-world...
Cached 42 issue(s).
Caching pull requests for octocat/hello-world...
Cached 15 pull request(s).
Cache updated. Valid for 60 minute(s).
```

Example output (subsequent run, delta fetch):

```
Fetching issues updated since 2025-01-15 10:30 for octocat/hello-world...
Cached 3 issue(s).
Fetching PRs updated since 2025-01-15 10:30 for octocat/hello-world...
Cached 1 pull request(s).
Cache updated. Valid for 60 minute(s).
```

List and view commands serve from the cache when it is fresh, and fall back to the GitHub API otherwise.

### Issues

```bash
# List open issues (default)
gh-cached issue list

# Filtering
gh-cached issue list --state all
gh-cached issue list --state closed
gh-cached issue list --author alice
gh-cached issue list --assignee bob
gh-cached issue list --label bug
gh-cached issue list --label bug --label p1   # AND logic
gh-cached issue list --milestone v2.0
gh-cached issue list --mention carol
gh-cached issue list --app dependabot
gh-cached issue list --search "memory leak"
gh-cached issue list --limit 10

# View a single issue
gh-cached issue view 42
gh-cached issue view 42 --comments
gh-cached issue view 42 --refresh   # force fetch from GitHub and update cache
```

Example output (`issue list`):

```
#42  Fix login bug                          bug, priority     5  2025-01-15
#41  Update README                          documentation     2  2025-01-14
#40  Very long title that will be truncat…  enhancement       0  2025-01-13
```

Example output (`issue view 42 --comments`):

```
#42 Fix login bug
OPEN • opened by octocat • 5 comment(s)

Labels:    bug, priority
Assignees: alice
Milestone: v1.0
Created:   2025-01-10 14:30
URL:       https://github.com/octocat/hello-world/issues/42

The login page crashes when...

── Comment 1 by reviewer1 (2025-01-12 09:15) ──

I can reproduce this on Firefox.

── Comment 2 by octocat (2025-01-12 10:30) ──

Thanks, I'll look into it.
```

### Pull Requests

```bash
# List open PRs (default)
gh-cached pr list

# Filtering
gh-cached pr list --state all
gh-cached pr list --state merged
gh-cached pr list --author alice
gh-cached pr list --assignee bob
gh-cached pr list --label enhancement
gh-cached pr list --base main
gh-cached pr list --head feat/dark-mode
gh-cached pr list --draft
gh-cached pr list --app dependabot
gh-cached pr list --search "fix crash"
gh-cached pr list --limit 10

# View a single PR
gh-cached pr view 10
gh-cached pr view 10 --comments
gh-cached pr view 10 --refresh      # force fetch from GitHub and update cache
```

Example output (`pr list`):

```
#10  Add new feature         feature → main    approved          3  2025-01-15
#9   Fix edge case [draft]   fix-bug → dev     review required   1  2025-01-14
```

Example output (`pr view 10 --comments`):

```
#10 Add new feature
OPEN • opened by octocat • 3 comment(s)

Branch:    feature → main
Review:    approved
Labels:    enhancement
Assignees: alice, bob
Created:   2025-01-08 11:00
URL:       https://github.com/octocat/hello-world/pull/10

This PR adds a new feature...

── Comment 1 by reviewer1 (2025-01-09 08:30) ──

Looks good, approved!

── Comment 2 by octocat (2025-01-09 09:00) ──

Thanks for the review.
```

### Repositories

List locally cached repositories:

```bash
gh-cached repo list
```

Example output:

```
REPO                      ISSUES  PRS  CACHED AGE  STATUS
octocat/hello-world       42      15   2h30m       fresh
alice/another-repo        10      3    1d          stale
```

## How caching works

| Command | Cache behaviour |
|---|---|
| `cache` | Fetches everything (all states, with comments) and writes one JSON file per issue/PR. Skips if cache is younger than `--cache-duration`. |
| `issue list` / `pr list` | Reads all cached files and filters in-memory when cache is fresh; falls back to the GitHub API with server-side filters otherwise. |
| `issue view` / `pr view` | Serves from the individual cached file when it is less than 60 minutes old; fetches from the API and updates the cache otherwise. `--refresh` bypasses all cache checks, fetches from the API, and updates the cache. |

## Flag reference

### Global

| Flag | Description |
|---|---|
| `--repo [HOST/]OWNER/REPO` | Target repository (default: detected from `git remote origin`) |

### `cache`

| Flag | Default | Description |
|---|---|---|
| `--cache-duration int` | `60` | Minutes before the cache is considered stale |
| `--force` | `false` | Re-fetch even if the cache is still fresh (full fetch, not delta) |

### `issue list`

| Flag | Description |
|---|---|
| `--app string` | Filter by GitHub App author |
| `-a, --assignee string` | Filter by assignee |
| `-A, --author string` | Filter by author |
| `--json` | Output as JSON |
| `-l, --label strings` | Filter by label (repeat for AND logic) |
| `-L, --limit int` | Maximum results (default 1000) |
| `--mention string` | Filter by mention |
| `-m, --milestone string` | Filter by milestone number or title |
| `--no-truncate` | Don't truncate long titles |
| `-S, --search string` | Search query |
| `-s, --state string` | `open` \| `closed` \| `all` (default `open`) |

### `issue view`

| Flag | Description |
|---|---|
| `-c, --comments` | Show comments |
| `--json` | Output as JSON |
| `--refresh` | Force fetch from GitHub and update cache |

### `pr list`

| Flag | Description |
|---|---|
| `--app string` | Filter by GitHub App author |
| `-a, --assignee string` | Filter by assignee |
| `-A, --author string` | Filter by author |
| `-B, --base string` | Filter by base branch |
| `-d, --draft` | Show only draft PRs |
| `-H, --head string` | Filter by head branch |
| `--json` | Output as JSON |
| `-l, --label strings` | Filter by label (repeat for AND logic) |
| `-L, --limit int` | Maximum results (default 1000) |
| `--no-truncate` | Don't truncate long titles |
| `-S, --search string` | Search query |
| `-s, --state string` | `open` \| `closed` \| `merged` \| `all` (default `open`) |

### `pr view`

| Flag | Description |
|---|---|
| `-c, --comments` | Show comments |
| `--json` | Output as JSON |
| `--refresh` | Force fetch from GitHub and update cache |
