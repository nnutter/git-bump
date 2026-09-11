// Package testrepo builds ephemeral git repositories for tests.
package testrepo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/testenv"
)

// Init creates a repository with one commit and the given tags.
func Init(t *testing.T, tags ...string) *git.Repository {
	t.Helper()
	testenv.Sterilize(t)

	repo, err := git.PlainInit(t.TempDir(), false)
	require.NoError(t, err)

	// Annotated tags need a tagger. Pin a local identity so tests do not
	// depend on the ambient git config, which CI does not provide.
	cfg, err := repo.Config()
	require.NoError(t, err)
	cfg.User.Name = "test"
	cfg.User.Email = "test@example.com"
	require.NoError(t, repo.SetConfig(cfg))

	work, err := repo.Worktree()
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(work.Filesystem.Root(), "file.txt"), []byte("data"), 0o600))
	_, err = work.Add("file.txt")
	require.NoError(t, err)
	hash, err := work.Commit("initial", &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com"},
	})
	require.NoError(t, err)
	for _, tag := range tags {
		_, err = repo.CreateTag(tag, hash, nil)
		require.NoError(t, err)
	}
	return repo
}
