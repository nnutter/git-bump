package gittags

import (
	"fmt"

	"github.com/go-git/go-git/v5"

	"github.com/nnutter/git-bump/internal/bump"
)

// Create creates an annotated tag on HEAD with the tag as its message.
func Create(repo *git.Repository, tag string) error {
	if _, err := bump.Parse(tag); err != nil {
		return err
	}
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("resolve HEAD: %w", err)
	}
	if _, err := repo.CreateTag(tag, head.Hash(), &git.CreateTagOptions{Message: tag}); err != nil {
		return fmt.Errorf("create tag %q: %w", tag, err)
	}
	return nil
}
