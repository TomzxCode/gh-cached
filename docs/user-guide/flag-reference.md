# Flag Reference

## Global flags

These flags apply to all commands.

| Flag | Description |
|---|---|
| `--repo [HOST/]OWNER/REPO` | Target repository. When omitted, detected from `git remote origin` in the current directory. |
| `--cache-dir string` | Override the cache directory (default `~/.cache/gh-cached/`). |
| `--api-url string` | Override the GitHub GraphQL API endpoint (for testing). |
| `--version`, `-v` | Print version. |
| `--help`, `-h` | Show help. |

## cache

| Flag | Default | Description |
|---|---|---|
| `--cache-duration int` | `60` | Minutes before the cache is considered stale |
| `--force` | `false` | Re-fetch even if the cache is still fresh (full fetch, not delta) |

## issue list

| Flag | Short | Default | Description |
|---|---|---|---|
| `--state string` | `-s` | `open` | Filter by state: `open`, `closed`, or `all` |
| `--author string` | `-A` | | Filter by author |
| `--assignee string` | `-a` | | Filter by assignee |
| `--label strings` | `-l` | | Filter by label (repeat for AND logic) |
| `--milestone string` | `-m` | | Filter by milestone number or title |
| `--mention string` | | | Filter by mention |
| `--app string` | | | Filter by GitHub App author |
| `--search string` | `-S` | | Search query (case-insensitive substring match on title and body) |
| `--limit int` | `-L` | `1000` | Maximum number of results |
| `--json` | | `false` | Output as JSON |
| `--no-truncate` | | `false` | Don't truncate long titles |

## issue view

| Flag | Short | Default | Description |
|---|---|---|---|
| `--comments` | `-c` | `false` | Show comments |
| `--json` | | `false` | Output as JSON |
| `--refresh` | | `false` | Force fetch from GitHub and update cache |

## pr list

| Flag | Short | Default | Description |
|---|---|---|---|
| `--state string` | `-s` | `open` | Filter by state: `open`, `closed`, `merged`, or `all` |
| `--author string` | `-A` | | Filter by author |
| `--assignee string` | `-a` | | Filter by assignee |
| `--label strings` | `-l` | | Filter by label (repeat for AND logic) |
| `--base string` | `-B` | | Filter by base branch name |
| `--head string` | `-H` | | Filter by head branch name |
| `--draft` | `-d` | `false` | Show only draft PRs |
| `--app string` | | | Filter by GitHub App author |
| `--search string` | `-S` | | Search query (case-insensitive substring match on title and body) |
| `--limit int` | `-L` | `1000` | Maximum number of results |
| `--json` | | `false` | Output as JSON |
| `--no-truncate` | | `false` | Don't truncate long titles |

## pr view

| Flag | Short | Default | Description |
|---|---|---|---|
| `--comments` | `-c` | `false` | Show comments |
| `--json` | | `false` | Output as JSON |
| `--refresh` | | `false` | Force fetch from GitHub and update cache |

## repo list

No command-specific flags. Uses only the global `--cache-dir`.
