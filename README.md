# gh-cached

A GitHub CLI that calls the GitHub GraphQL API to retrieve issues, pull requests, and comments, caching all results to disk to minimise API calls.

Cache lives at `~/.cache/gh-cached/<host>/<owner>/<repo>`.

## Installation

```bash
go install github.com/tomzxcode/gh-cached@latest
```

Or download a pre-built binary from the [releases page](https://github.com/TomzxCode/gh-cached/releases/tag/latest).

## Authentication

Set a GitHub personal access token in your environment:

```bash
export GH_TOKEN=ghp_...
```

Alternatively, if you have the [GitHub CLI](https://cli.github.com) installed and authenticated, `gh-cached` will use it as a fallback (`gh auth token`).

## Usage

All commands accept `--repo [HOST/]OWNER/REPO`. When omitted, the repository is detected from the `origin` remote of the current directory.

### Cache

Pre-populate the local cache with all issues and PRs (including comments):

```bash
gh-cached cache
gh-cached cache --repo cli/cli
gh-cached cache --cache-duration 120   # treat cache as fresh for 2 hours
gh-cached cache --cache-duration 0     # always re-fetch
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
```

## How caching works

| Command | Cache behaviour |
|---|---|
| `cache` | Fetches everything (all states, with comments) and writes one JSON file per issue/PR. Skips if cache is younger than `--cache-duration`. |
| `issue list` / `pr list` | Reads all cached files and filters in-memory when cache is fresh; falls back to the GitHub API with server-side filters otherwise. |
| `issue view` / `pr view` | Serves from the individual cached file when it is less than 60 minutes old; fetches from the API and updates the cache otherwise. |

## Flag reference

### Global

| Flag | Description |
|---|---|
| `--repo [HOST/]OWNER/REPO` | Target repository (default: detected from `git remote origin`) |

### `cache`

| Flag | Default | Description |
|---|---|---|
| `--cache-duration int` | `60` | Minutes before the cache is considered stale |

### `issue list`

| Flag | Description |
|---|---|
| `--app string` | Filter by GitHub App author |
| `-a, --assignee string` | Filter by assignee |
| `-A, --author string` | Filter by author |
| `-l, --label strings` | Filter by label (repeat for AND logic) |
| `-L, --limit int` | Maximum results (default 30) |
| `--mention string` | Filter by mention |
| `-m, --milestone string` | Filter by milestone number or title |
| `-S, --search string` | Search query |
| `-s, --state string` | `open` \| `closed` \| `all` (default `open`) |

### `issue view`

| Flag | Description |
|---|---|
| `-c, --comments` | Show comments |

### `pr list`

| Flag | Description |
|---|---|
| `--app string` | Filter by GitHub App author |
| `-a, --assignee string` | Filter by assignee |
| `-A, --author string` | Filter by author |
| `-B, --base string` | Filter by base branch |
| `-d, --draft` | Show only draft PRs |
| `-H, --head string` | Filter by head branch |
| `-l, --label strings` | Filter by label (repeat for AND logic) |
| `-L, --limit int` | Maximum results (default 30) |
| `-S, --search string` | Search query |
| `-s, --state string` | `open` \| `closed` \| `merged` \| `all` (default `open`) |

### `pr view`

| Flag | Description |
|---|---|
| `-c, --comments` | Show comments |
