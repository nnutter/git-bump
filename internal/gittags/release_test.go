package gittags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/gittags"
	"github.com/nnutter/git-bump/internal/testrepo"
)

func stubGh(t *testing.T, script string) {
	t.Helper()
	binDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(binDir, "gh"), []byte(script), 0o700))
	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func TestCreateDraftRelease(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")
	require.NoError(t, gittags.Create(repo, "v1.2.4"))
	work, err := repo.Worktree()
	require.NoError(t, err)
	root := work.Filesystem.Root()

	pwdFile := filepath.Join(t.TempDir(), "pwd")
	argsFile := filepath.Join(t.TempDir(), "args")
	t.Setenv("FAKE_GH_PWD_FILE", pwdFile)
	t.Setenv("FAKE_GH_ARGS_FILE", argsFile)
	stubGh(t, "#!/bin/sh\necho \"$PWD\" > \"$FAKE_GH_PWD_FILE\"\necho \"$@\" > \"$FAKE_GH_ARGS_FILE\"\necho \"https://github.com/example/repo/releases/tag/v1.2.4\"\n")

	url, err := gittags.CreateDraftRelease(repo, "v1.2.4")
	require.NoError(t, err)
	require.Equal(t, "https://github.com/example/repo/releases/tag/v1.2.4", url)

	pwd, err := os.ReadFile(pwdFile)
	require.NoError(t, err)
	require.Equal(t, root+"\n", string(pwd), "gh should run in the worktree root")

	args, err := os.ReadFile(argsFile)
	require.NoError(t, err)
	require.Equal(t, "release create v1.2.4 --draft --title v1.2.4 --generate-notes\n", string(args))
}

func TestCreateDraftReleaseInvalid(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	_, err := gittags.CreateDraftRelease(repo, "release-foo")
	require.ErrorContains(t, err, "invalid semver tag")
}

func TestCreateDraftReleaseFailure(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	stubGh(t, "#!/bin/sh\necho \"release not found\" >&2\nexit 1\n")

	_, err := gittags.CreateDraftRelease(repo, "v1.2.4")
	require.ErrorContains(t, err, `create draft release for "v1.2.4"`)
}
