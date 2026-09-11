// Package gittags reads and writes repository tags through go-git
// instead of shelling out to the git CLI.
package gittags

import (
	"fmt"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/nnutter/git-bump/internal/bump"
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

// Latest returns the highest semantic version tag in repo.
//
// When pattern is not empty only tags matching it (in the sense of
// git tag -l) are considered. Tags that are not strict semantic
// versions are skipped. When pattern is empty and no usable tag
// remains, Latest returns v0.0.0 so callers can bump from scratch;
// with a non-empty pattern the absence of a usable tag is an error.
func Latest(repo *git.Repository, pattern string) (string, error) {
	if _, err := path.Match(pattern, ""); err != nil {
		return "", fmt.Errorf("invalid pattern %q: %w", pattern, err)
	}

	iter, err := repo.Tags()
	if err != nil {
		return "", fmt.Errorf("list tags: %w", err)
	}
	names := []string{}
	iterErr := iter.ForEach(func(ref *plumbing.Reference) error {
		names = append(names, ref.Name().Short())
		return nil
	})
	if iterErr != nil {
		iter.Close()
		return "", fmt.Errorf("list tags: %w", iterErr)
	}
	iter.Close()

	best := ""
	found := false
	for _, name := range names {
		if pattern != "" {
			matched, err := path.Match(pattern, name)
			if err != nil {
				return "", fmt.Errorf("invalid pattern %q: %w", pattern, err)
			}
			if !matched {
				continue
			}
		}
		if _, err := bump.Parse(name); err != nil {
			continue
		}
		if !found || semver.Compare(name, best) > 0 {
			best, found = name, true
		}
	}
	if !found {
		if pattern == "" {
			return bump.Version{}.String(), nil
		}
		return "", fmt.Errorf("no tags match pattern %q", pattern)
	}
	return best, nil
}

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

// Push pushes a single tag to origin.
func Push(repo *git.Repository, tag string) error {
	ref := config.RefSpec("refs/tags/" + tag + ":refs/tags/" + tag)
	if err := repo.Push(&git.PushOptions{RefSpecs: []config.RefSpec{ref}}); err != nil {
		return fmt.Errorf("push tag %q: %w", tag, err)
	}
	return nil
}
