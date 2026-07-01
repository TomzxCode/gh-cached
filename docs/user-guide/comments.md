# Comments

ghx can add, reply to, edit, and delete comments on pull requests and issues.
PR comments go beyond the standard `gh` CLI: inline comments on specific lines, file-level comments, replies to review threads, and comments staged in a pending review or a local stash.

## Top-level PR comment

```bash
ghx pr comment 10 --body "Looks good"
ghx pr comment 10 --body-file comment.txt
ghx pr comment 10 --body-file -        # read body from stdin
```

Without `--file`, the comment is added as a top-level PR comment.

## Inline comment on a line or range

```bash
ghx pr comment 10 --file src/main.go --line 42 --body "Nit"
ghx pr comment 10 --file src/main.go --line 42-45 --body "Consider extracting"
```

`--line` requires `--file`. A range is written as `START-END`.

## File-level comment

Omit `--line` to leave a comment on an entire file:

```bash
ghx pr comment 10 --file src/main.go --body "Overall looks clean"
```

## Diff side

Inline comments default to the right-hand side of the diff. Use `--side` to comment on the left-hand side:

```bash
ghx pr comment 10 --file src/main.go --line 42 --side LEFT --body "Before"
```

`--side` requires `--file`.

## Reply to a review thread

```bash
ghx pr comment 10 --reply-thread <thread-id> --body "Agreed"
```

Find thread IDs with `ghx pr threads 10` (see [Reviews](reviews.md)). `--reply-thread` is mutually exclusive with `--file` and `--line`.

## Stage a comment in a pending review

Add `--pending` to place a comment in a pending review instead of posting it immediately.
Pending reviews are submitted later with `ghx pr review submit` (see [Reviews](reviews.md)):

```bash
ghx pr comment 10 --file src/main.go --line 42 --body "Nit" --pending
```

`--pending` requires `--file` or `--reply-thread`; it does not apply to top-level comments.

## Save a comment to a local stash

Add `--stash` to save a comment to a local stash entry instead of calling the API.
Stashes are restored later with `ghx pr review stash pop` (see [Reviews](reviews.md)):

```bash
ghx pr comment 10 --file src/main.go --line 42 --body "Nit" --stash
ghx pr comment 10 --file src/main.go --line 42 --body "Nit" --stash=1   # target stash@{1}
```

`--stash` requires `--file` and is mutually exclusive with `--pending`.

## Issue comments

```bash
ghx issue comment 42 --body "Fixed in #50"
ghx issue comment 42 --body-file -
```

Issue comments are always top-level.

## Edit a comment

```bash
ghx pr comment edit <comment-id> --body "Updated text"
ghx issue comment edit <comment-id> --body "Updated text"
```

`pr comment edit` automatically detects whether the comment is an inline review comment or a top-level comment.

## Delete a comment

```bash
ghx pr comment delete <comment-id>
ghx issue comment delete <comment-id>
```

## Find comment IDs

Use `--ids` to display the node IDs needed for editing and deleting:

```bash
ghx pr threads 10 --ids        # IDs for review-thread comments
ghx issue view 42 --ids        # IDs for issue comments
```
