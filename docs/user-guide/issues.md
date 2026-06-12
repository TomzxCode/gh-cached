# Issues

## List issues

```bash
gh-cached issue list
```

When the cache is fresh, results are served from disk. Otherwise the GitHub API is used.

### Filter by state

```bash
gh-cached issue list                    # open (default)
gh-cached issue list --state closed
gh-cached issue list --state all
```

### Filter by author, assignee, or mention

```bash
gh-cached issue list --author alice
gh-cached issue list --assignee bob
gh-cached issue list --mention carol
```

### Filter by label

Labels use AND logic. Repeat the flag to require all specified labels:

```bash
gh-cached issue list --label bug
gh-cached issue list --label bug --label p1
```

### Filter by milestone

Accepts a milestone title (case-insensitive) or number:

```bash
gh-cached issue list --milestone v2.0
gh-cached issue list --milestone 5
```

### Search

Case-insensitive substring match against title and body:

```bash
gh-cached issue list --search "memory leak"
```

### Limit results

```bash
gh-cached issue list --limit 10
```

When the total exceeds the limit, a summary is printed to stderr:

```
Showing 10 of 42 issues
```

### JSON output

```bash
gh-cached issue list --json
```

Outputs a pretty-printed JSON array of issue objects.

### Disable title truncation

By default, titles are truncated to 60 characters. Use `--no-truncate` to show full titles:

```bash
gh-cached issue list --no-truncate
```

## View a single issue

```bash
gh-cached issue view 42
```

### Show comments

```bash
gh-cached issue view 42 --comments
```

### Force refresh

Bypass the cache and fetch the latest data from GitHub, updating the cache:

```bash
gh-cached issue view 42 --refresh
```

### JSON output

```bash
gh-cached issue view 42 --json
```

## Example output

### issue list

```
#42  Fix login bug                          bug, priority     5  2025-01-15
#41  Update README                          documentation     2  2025-01-14
#40  Very long title that will be truncat…  enhancement       0  2025-01-13
```

### issue view

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
```
