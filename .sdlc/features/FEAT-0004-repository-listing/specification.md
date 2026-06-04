---
title: "Repository Listing"
status: draft
---

# Specification: Repository Listing

## Overview

The `repo list` command walks the cache directory tree (`~/.cache/gh-cached/`) and aggregates per-repository statistics.
No network calls are made; all data comes from the filesystem.

## Architecture

```
repo list ──► runRepoList()
                   │
                   └─ Store.ListCachedRepos()
                        │
                        ├─ Walk baseDir/<host>/<owner>/<repo>/ directories
                        │
                        ├─ For each repo: LoadCacheInfo(), count issue/PR JSON files
                        │
                        └─ Return []CachedRepo
                                          │
                                          ▼
                                   printRepoList() ──► tabular output
```

## Data Models

### CachedRepo

| Field | Type | Constraints | Description |
|---|---|---|---|
| host | string | not null | GitHub host (e.g. github.com) |
| owner | string | not null | Repository owner |
| repo | string | not null | Repository name |
| info | *CacheInfo | nullable | Cache metadata (nil if no .cache_info.json) |
| issueCount | int | not null | Number of cached issue files |
| prCount | int | not null | Number of cached PR files |

## Technical Decisions

| Decision | Choice | Rationale |
|---|---|---|
| File-count based statistics | Count .json files in issues/ and prs/ dirs | Simple, no index needed |
| Duration formatting | Custom `formatDuration` helper | Human-readable output (e.g. "2h30m", "3d") |

## Risks and Unknowns

(none identified)

## Out of Scope

- Fetching or refreshing repos from this command
- Displaying per-repo configuration or settings
