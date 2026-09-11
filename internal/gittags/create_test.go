package gittags_test

import (
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/gittags"
	"github.com/nnutter/git-bump/internal/testenv"
	"github.com/nnutter/git-bump/internal/testrepo"
)

func TestCreate(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	require.NoError(t, gittags.Create(repo, "v1.2.4"))

	ref, err := repo.Tag("v1.2.4")
	require.NoError(t, err)
	require.Equal(t, "refs/tags/v1.2.4", ref.Name().String())

	tag, err := repo.TagObject(ref.Hash())
	require.NoError(t, err)
	require.Equal(t, "v1.2.4\n", tag.Message)
}

func TestCreateWithoutAmbientIdentity(t *testing.T) {
	// Simulate CI, where the global and system git configs carry no user
	// identity: only the repository's own identity can supply a tagger.
	testenv.Sterilize(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", home)

	repo := testrepo.Init(t, "v1.2.3")
	require.NoError(t, gittags.Create(repo, "v1.2.4"))
}

func TestCreateInvalid(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	require.Error(t, gittags.Create(repo, "release-foo"))
}

func TestPush(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	remoteDir := t.TempDir()
	_, err := git.PlainInit(remoteDir, true)
	require.NoError(t, err)
	_, err = repo.CreateRemote(&config.RemoteConfig{Name: "origin", URLs: []string{remoteDir}})
	require.NoError(t, err)

	require.NoError(t, gittags.Create(repo, "v1.2.4"))
	require.NoError(t, gittags.Push(repo, "v1.2.4"))

	remote, err := git.PlainOpen(remoteDir)
	require.NoError(t, err)
	_, err = remote.Tag("v1.2.4")
	require.NoError(t, err, "pushed tag should exist on the remote")
}

func TestPushWithoutRemote(t *testing.T) {
	repo := testrepo.Init(t, "v1.2.3")

	require.NoError(t, gittags.Create(repo, "v1.2.4"))
	require.ErrorContains(t, gittags.Push(repo, "v1.2.4"), `push tag "v1.2.4"`)
}
