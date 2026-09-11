// Package testrepo builds ephemeral git repositories for tests.
package testrepo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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

// InitWorktree creates a repository with one commit and the given tags
// plus a linked-worktree-style checkout of it.
//
// It returns the main repository and the worktree directory. The
// worktree's .git file points at a gitdir holding HEAD and a commondir
// pointer, mirroring `git worktree add` without shelling out to the
// git CLI.
func InitWorktree(t *testing.T, tags ...string) (*git.Repository, string) {
	t.Helper()

	main := Init(t, tags...)

	remoteDir := t.TempDir()
	_, err := git.PlainInit(remoteDir, true)
	require.NoError(t, err)
	_, err = main.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{remoteDir}})
	require.NoError(t, err)

	work, err := main.Worktree()
	require.NoError(t, err)
	mainGitDir := filepath.Join(work.Filesystem.Root(), ".git")

	head, err := main.Head()
	require.NoError(t, err)

	workDir := t.TempDir()
	wtGitDir := filepath.Join(workDir, "gitdir")
	require.NoError(t, os.MkdirAll(wtGitDir, 0o755))
	require.NoError(t, os.WriteFile(
		filepath.Join(wtGitDir, "HEAD"),
		[]byte("ref: "+head.Name().String()+"\n"),
		0o600,
	))
	require.NoError(t, os.WriteFile(filepath.Join(wtGitDir, "commondir"), []byte(mainGitDir+"\n"), 0o600))
	require.NoError(t, os.WriteFile(filepath.Join(workDir, ".git"), []byte("gitdir: "+wtGitDir+"\n"), 0o600))

	return main, workDir
}
