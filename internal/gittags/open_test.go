package gittags_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/gittags"
	"github.com/nnutter/git-bump/internal/testrepo"
)

func TestOpenLinkedWorktree(t *testing.T) {
	main, workDir := testrepo.InitWorktree(t, "v1.2.3")

	repo, err := gittags.Open(workDir)
	require.NoError(t, err)

	mainHead, err := main.Head()
	require.NoError(t, err)
	head, err := repo.Head()
	require.NoError(t, err)
	require.Equal(t, mainHead.Hash(), head.Hash(), "worktree HEAD should resolve through commondir")

	latest, err := gittags.Latest(repo, "")
	require.NoError(t, err)
	require.Equal(t, "v1.2.3", latest, "shared tags should be visible from the worktree")

	cfg, err := repo.Config()
	require.NoError(t, err)
	require.Contains(t, cfg.Remotes, "origin")

	require.NoError(t, gittags.Create(repo, "v1.2.4"))
	_, err = main.Tag("v1.2.4")
	require.NoError(t, err, "tag created from the worktree should land in shared refs")
}
