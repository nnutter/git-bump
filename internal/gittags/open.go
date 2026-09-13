// Package gittags reads and writes repository tags through go-git.
// Push shells out to the git CLI so authentication (SSH agent,
// credential helpers such as gh) behaves like a manual push.
package gittags

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// Open opens the repository at path.
//
// Linked worktrees (.git file plus commondir pointer) are resolved
// through the shared directory so HEAD, tags, and remotes behave as
// they do under the git CLI.
func Open(path string) (*git.Repository, error) {
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{EnableDotGitCommonDir: true})
	if err != nil {
		return nil, fmt.Errorf("open git repository: %w", err)
	}
	return repo, nil
}
