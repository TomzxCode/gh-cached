---
title: "Repository Listing"
status: draft
---

# Requirements: Repository Listing

## Overview

Provides a `repo list` command that scans the local cache directory and displays all repositories that have been previously cached, along with their issue/PR counts, cache age, and freshness status.

## Stakeholders

| Stakeholder | Interest |
|---|---|
| CLI user | See which repositories are cached and whether they are stale |

## Functional Requirements

Order rows by priority: Must first, then Should, then May.

| ID | Priority | Requirement |
|---|---|---|
| FR-01 | Must | The system shall provide a `repo list` command that displays all locally cached repositories |
| FR-02 | Must | For each repository, the system shall display issue count, PR count, cache age, and freshness status |
| FR-03 | Must | The system shall display the repository as `OWNER/REPO` for github.com hosts and `HOST/OWNER/REPO` for others |

## Non-Functional Requirements

Order rows by priority: Must first, then Should, then May.

| ID | Priority | Category | Requirement |
|---|---|---|---|
| NFR-01 | Must | Performance | Listing shall complete in under 1 second for up to 50 cached repos |

## Constraints

- Only shows repositories with data on disk; no network calls
- Cache age formatting rounds to minutes

## Acceptance Criteria

- [ ] FR-01: `gh-cached repo list` shows all cached repos in tabular format
- [ ] FR-02: Each row shows issue count, PR count, age, and fresh/stale status
- [ ] FR-03: GitHub.com repos omit the host prefix; GHE repos include it

## Open Questions

(none)
