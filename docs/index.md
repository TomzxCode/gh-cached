# gh-cached

A GitHub CLI that calls the GitHub GraphQL API to retrieve issues, pull requests, and comments, caching all results to disk to minimise API calls.

Cache lives at `~/.cache/gh-cached/<host>/<owner>/<repo>`.

## Quick start

```bash
go install github.com/tomzxcode/gh-cached@main
export GH_TOKEN=ghp_...
gh-cached cache --repo cli/cli
gh-cached issue list
gh-cached pr list
```

See the [User Guide](user-guide/installation.md) for detailed instructions.
