package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/cmd"
	"github.com/nnutter/git-bump/internal/testenv"
	"github.com/nnutter/git-bump/internal/testrepo"
)

func TestRequiresExactlyOneBumpFlag(t *testing.T) {
	testenv.Sterilize(t)

	args := [][]string{
		{},
		{"--major", "--minor"},
		{"--major", "--patch"},
		{"--minor", "--patch"},
		{"--major", "--minor", "--patch"},
	}
	for _, arg := range args {
		command := cmd.NewRootCommand()
		command.SetArgs(arg)
		err := command.Execute()
		require.ErrorContains(t, err, "exactly one of --major, --minor, or --patch is required")
	}
}

func TestPatchWithoutPush(t *testing.T) {
	testenv.Sterilize(t)

	repo := testrepo.Init(t, "v1.2.3")
	work, err := repo.Worktree()
	require.NoError(t, err)
	previous, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(work.Filesystem.Root()))
	t.Cleanup(func() {
		require.NoError(t, os.Chdir(previous))
	})

	command := cmd.NewRootCommand()
	command.SetArgs([]string{"--patch", "--no-push"})
	output := &bytes.Buffer{}
	command.SetOut(output)
	require.NoError(t, command.Execute())
	require.Equal(t, "v1.2.4\n", output.String())

	_, err = repo.Tag("v1.2.4")
	require.NoError(t, err, "new tag should exist locally")
}
