# git-bump

Bump the latest semver tag in the current repository.

## Usage

```sh
git-bump --major|--minor|--patch [--pattern <glob>] [--no-push] [--release] [--no-open]
```

Exactly one of `--major`, `--minor`, or `--patch` is required.
`git-bump` finds the highest `vMAJOR.MINOR.PATCH` tag, applies the
bump, creates the new tag, prints it, and pushes it to `origin`.

- `--pattern` limits candidates the way `git tag -l <glob>` does.
  The default considers every tag. Tags that are not strict semver
  are skipped.
- `--no-push` creates the tag locally without pushing it.
- `--release` creates a draft GitHub release for the new tag with
  `gh`, prints its URL after the tag, and opens it in the browser
  with `open` (or `xdg-open` as a fallback) when available.
  `--no-open` skips opening the draft release URL.

## Examples

```sh
git-bump --patch
git-bump --minor --pattern 'v1.*'
git-bump --major --no-push
git-bump --patch --release
```
