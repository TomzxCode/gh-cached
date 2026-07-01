# Mock server

The `mock` command runs an in-process mock GitHub GraphQL server for testing ghx (or any GitHub GraphQL client) without hitting the real API.
It generates a synthetic dataset of users, repos, issues, PRs, and comments, then serves it over HTTP until you stop it with Ctrl+C.

## Start the server

```bash
ghx mock serve --repos octocat/hello-world
```

On startup the server prints its URL (for example `http://127.0.0.1:PORT`) and a few example commands.

## Point ghx at the mock server

Use `--api-url` and a throwaway `--cache-dir` so mock data does not mix with your real cache:

```bash
ghx --api-url http://127.0.0.1:PORT --cache-dir /tmp/mock-cache --repo octocat/hello-world cache
ghx --api-url http://127.0.0.1:PORT --cache-dir /tmp/mock-cache --repo octocat/hello-world issue list --state all
```

The token is irrelevant for the mock server; any non-empty `GH_TOKEN` works.

## Presets

The `--preset` flag selects a baseline configuration that individual flags override:

| Preset | Dataset |
|---|---|
| `default` | A realistic mid-size dataset: 3 repos, 8 users, 90 days of history |
| `small` | A minimal dataset for fast runs: 1 repo, 3 users, 7 days |
| `none` | No baseline; every parameter must be supplied explicitly |

```bash
ghx mock serve --preset small --repos octocat/hello-world
```

Any flag you pass overrides the preset value; flags you omit fall back to the preset.

## Preview the dataset without serving

Add `--stats` to print a summary of the generated dataset and exit, without starting the server:

```bash
ghx mock serve --repos octocat/hello-world --stats
```

## Reproducible data

Fix the RNG seed to generate the same dataset across runs:

```bash
ghx mock serve --repos octocat/hello-world --seed 42
```

Adjust the timeline with `--history` (a Go duration) and add `--activity-bursts` to concentrate events into high-activity windows (for example, a release sprint):

```bash
ghx mock serve --repos octocat/hello-world --history 720h --activity-bursts 5
```
