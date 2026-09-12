package gittags

import (
	"fmt"
	"os/exec"

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

// Push pushes a single tag to origin with the git CLI so SSH agent,
// ssh config, and credential helpers are honored. go-git pushes
// anonymously, which GitHub rejects with "No anonymous write access"
// on SSH remotes.
func Push(repo *git.Repository, tag string) error {
	dir, err := workDir(repo)
	if err != nil {
		return fmt.Errorf("push tag %q: %w", tag, err)
	}
	ref := "refs/tags/" + tag + ":refs/tags/" + tag
	push := exec.Command("git", "push", "origin", ref)
	push.Dir = dir
	if out, err := push.CombinedOutput(); err != nil {
		return fmt.Errorf("push tag %q: %w: %s", tag, err, out)
	}
	return nil
}
