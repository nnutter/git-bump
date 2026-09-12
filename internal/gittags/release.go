package gittags

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/nnutter/git-bump/internal/bump"
)

// CreateDraftRelease creates a draft GitHub release for tag with the gh
// CLI so authentication (GH_TOKEN, gh auth login, credential helpers)
// behaves like a manual release. It returns the draft release URL
// printed by gh.
func CreateDraftRelease(repo *git.Repository, tag string) (string, error) {
	if _, err := bump.Parse(tag); err != nil {
		return "", err
	}
	dir, err := workDir(repo)
	if err != nil {
		return "", fmt.Errorf("create draft release for %q: %w", tag, err)
	}
	create := exec.Command("gh", "release", "create", tag, "--draft", "--title", tag, "--generate-notes")
	create.Dir = dir
	out, err := create.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("create draft release for %q: %w: %s", tag, err, out)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	url := strings.TrimSpace(lines[len(lines)-1])
	if url == "" {
		return "", fmt.Errorf("create draft release for %q: no URL in output: %s", tag, out)
	}
	return url, nil
}
