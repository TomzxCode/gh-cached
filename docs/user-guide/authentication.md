# Authentication

ghx needs a GitHub personal access token to call the GitHub GraphQL API. It resolves a token in the following order:

1. **`GH_TOKEN`** environment variable
2. **`GITHUB_TOKEN`** environment variable
3. **`gh auth token --hostname <host>`** (GitHub CLI, host-specific)
4. **`gh auth token`** (GitHub CLI, any host)

If none of these are available, ghx exits with an error.

## Personal access token

Create a token at [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens). The token needs the `repo` scope for private repositories.

Set it in your shell:

```bash
export GH_TOKEN=ghp_...
```

## GitHub CLI

If you already use [`gh`](https://cli.github.com/) and have run `gh auth login`, ghx will use it automatically as a fallback.

## GitHub Enterprise

For GitHub Enterprise hosts, pass the full repository path:

```bash
ghx --repo ghe.example.com/org/repo cache
```

ghx resolves the API endpoint as `https://ghe.example.com/api/graphql` and uses `gh auth token --hostname ghe.example.com` when looking for a token.
