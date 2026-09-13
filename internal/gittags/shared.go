package gittags

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

// workDir returns the worktree root for running external commands.
func workDir(repo *git.Repository) (string, error) {
	work, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("resolve worktree: %w", err)
	}
	return work.Filesystem.Root(), nil
}
