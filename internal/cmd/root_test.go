package cmd_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/cmd"
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
		command := cmd.NewRootCommand("test")
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

	command := cmd.NewRootCommand("test")
	command.SetArgs([]string{"--patch", "--no-push"})
	output := &bytes.Buffer{}
	command.SetOut(output)
	require.NoError(t, command.Execute())
	require.Equal(t, "v1.2.4\n", output.String())

	_, err = repo.Tag("v1.2.4")
	require.NoError(t, err, "new tag should exist locally")
}

func TestBootstrapWithoutTags(t *testing.T) {
	testenv.Sterilize(t)

	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "patch", args: []string{"--patch", "--no-push"}, want: "v0.0.1\n"},
		{name: "minor", args: []string{"--minor", "--no-push"}, want: "v0.1.0\n"},
		{name: "major", args: []string{"--major", "--no-push"}, want: "v1.0.0\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testrepo.Init(t)
			work, err := repo.Worktree()
			require.NoError(t, err)
			previous, err := os.Getwd()
			require.NoError(t, err)
			require.NoError(t, os.Chdir(work.Filesystem.Root()))
			t.Cleanup(func() {
				require.NoError(t, os.Chdir(previous))
			})

			command := cmd.NewRootCommand("test")
			command.SetArgs(tt.args)
			output := &bytes.Buffer{}
			command.SetOut(output)
			require.NoError(t, command.Execute())
			require.Equal(t, tt.want, output.String())
		})
	}
}

func TestReleaseCreatesDraftRelease(t *testing.T) {
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

	binDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(binDir, "gh"),
		[]byte("#!/bin/sh\necho \"https://github.com/example/repo/releases/tag/$3\"\n"),
		0o700,
	))
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	command := cmd.NewRootCommand("test")
	command.SetArgs([]string{"--patch", "--no-push", "--release"})
	output := &bytes.Buffer{}
	command.SetOut(output)
	require.NoError(t, command.Execute())
	require.Equal(t, "v1.2.4\nhttps://github.com/example/repo/releases/tag/v1.2.4\n", output.String())

	_, err = repo.Tag("v1.2.4")
	require.NoError(t, err, "new tag should exist locally")
}

func TestReleaseFailure(t *testing.T) {
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

	binDir := t.TempDir()
	require.NoError(t, os.WriteFile(
		filepath.Join(binDir, "gh"),
		[]byte("#!/bin/sh\necho \"release failed\" >&2\nexit 1\n"),
		0o700,
	))
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	command := cmd.NewRootCommand("test")
	command.SetArgs([]string{"--patch", "--no-push", "--release"})
	command.SetOut(&bytes.Buffer{})
	require.ErrorContains(t, command.Execute(), `create draft release for "v1.2.4"`)
}
