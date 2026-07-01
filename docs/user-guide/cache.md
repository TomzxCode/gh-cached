# Cache

ghx caches GitHub data as individual JSON files on your local disk. After an initial fetch, subsequent commands serve results from cache without hitting the API.

## Cache location

```
~/.cache/ghx/cache/<host>/<owner>/<repo>/
├── .cache_info.json
├── issues/
│   ├── 1.json
│   ├── 2.json
│   └── ...
└── prs/
    ├── 1.json
    ├── 2.json
    └── ...
```

Override the base directory with `--cache-dir`:

```bash
ghx --cache-dir /tmp/gh-cache cache --repo cli/cli
```

## Populate the cache

```bash
ghx cache
ghx cache --repo cli/cli
ghx cache --cache-duration 120   # treat cache as fresh for 2 hours
ghx cache --cache-duration 0     # always re-fetch (delta)
ghx cache --force                # force full re-fetch
```

First run fetches everything. Example output:

```
Caching issues for octocat/hello-world...
Cached 42 issue(s).
Caching pull requests for octocat/hello-world...
Cached 15 pull request(s).
Cache updated. Valid for 60 minute(s).
```

Subsequent runs use a delta fetch, only retrieving items updated since the last cache write:

```
Fetching issues updated since 2025-01-15 10:30 for octocat/hello-world...
Cached 3 issue(s).
Fetching PRs updated since 2025-01-15 10:30 for octocat/hello-world...
Cached 1 pull request(s).
Cache updated. Valid for 60 minute(s).
```

## How freshness works

The `.cache_info.json` file tracks when the cache was last written and the configured duration:

```json
{
  "cachedAt": "2025-01-15T10:30:00Z",
  "duration": 60
}
```

A cache is considered **fresh** when `time.Since(cachedAt) < duration × 1 minute`.

When the cache is fresh, `list` and `view` commands serve entirely from disk with no API calls. When stale, they fall back to the GitHub API.

## Cache behavior per command

| Command | Behavior |
|---|---|
| `cache` | Fetches all issues and PRs (all states, with comments). Skips if cache is younger than `--cache-duration`. Supports delta fetch. |
| `issue list` / `pr list` | Reads cached files and filters in memory when cache is fresh. Falls back to the GitHub API when stale. Does not write to cache. |
| `issue view` / `pr view` | Serves the individual cached file if the full cache is fresh, or if the file is less than 60 minutes old. Otherwise fetches from the API and saves to cache. `--refresh` bypasses all checks. |

## Notes

- Cache files are never automatically cleaned up. Delete directories under `~/.cache/ghx/cache/` to free space.
- The `--mention` and `--app` filters cannot be evaluated from cached data and are silently skipped when serving from cache.
- Bulk cache operations fetch up to 100 comments per item.
