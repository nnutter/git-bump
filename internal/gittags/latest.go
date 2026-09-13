package gittags

import (
	"fmt"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"

	"github.com/nnutter/git-bump/internal/bump"
)

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
