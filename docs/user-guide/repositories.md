# Repositories

## List cached repositories

```bash
gh-cached repo list
```

Lists all repositories that have been cached locally. This command reads the filesystem only and makes no network calls.

Example output:

```
REPO                      ISSUES  PRS  CACHED AGE  STATUS
octocat/hello-world       42      15   2h30m       fresh
alice/another-repo        10      3    1d          stale
```

For GitHub Enterprise hosts, the full `HOST/OWNER/REPO` format is displayed.

The status column shows `fresh` when the cache is within its configured duration, and `stale` otherwise.

## Clean up

Cache files are never automatically removed. To free disk space, delete directories under the cache location:

```bash
rm -rf ~/.cache/gh-cached/<host>/<owner>/<repo>
```

Use `repo list` to identify which repositories are cached and how much data they hold.
