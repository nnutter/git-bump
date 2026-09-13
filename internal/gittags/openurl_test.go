package gittags_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/gittags"
)

func stubOpener(t *testing.T, binDir, name, script string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(binDir, name), []byte(script), 0o700))
}

func TestOpenURLPrefersOpen(t *testing.T) {
	binDir := t.TempDir()
	openCalls := filepath.Join(t.TempDir(), "open-calls")
	xdgCalls := filepath.Join(t.TempDir(), "xdg-calls")
	t.Setenv("FAKE_OPEN_CALLS", openCalls)
	t.Setenv("FAKE_XDG_CALLS", xdgCalls)
	stubOpener(t, binDir, "open", "#!/bin/sh\necho \"$@\" > \"$FAKE_OPEN_CALLS\"\n")
	stubOpener(t, binDir, "xdg-open", "#!/bin/sh\necho \"$@\" > \"$FAKE_XDG_CALLS\"\n")
	t.Setenv("PATH", binDir)

	require.NoError(t, gittags.OpenURL("https://example.com/release"))

	calls, err := os.ReadFile(openCalls)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/release\n", string(calls))
	require.NoFileExists(t, xdgCalls, "xdg-open should not run when open exists")
}

func TestOpenURLFallsBackToXdgOpen(t *testing.T) {
	binDir := t.TempDir()
	xdgCalls := filepath.Join(t.TempDir(), "xdg-calls")
	t.Setenv("FAKE_XDG_CALLS", xdgCalls)
	stubOpener(t, binDir, "xdg-open", "#!/bin/sh\necho \"$@\" > \"$FAKE_XDG_CALLS\"\n")
	t.Setenv("PATH", binDir)

	require.NoError(t, gittags.OpenURL("https://example.com/release"))

	calls, err := os.ReadFile(xdgCalls)
	require.NoError(t, err)
	require.Equal(t, "https://example.com/release\n", string(calls))
}

func TestOpenURLWithoutOpener(t *testing.T) {
	t.Setenv("PATH", t.TempDir())

	require.NoError(t, gittags.OpenURL("https://example.com/release"))
}

func TestOpenURLFailure(t *testing.T) {
	binDir := t.TempDir()
	stubOpener(t, binDir, "open", "#!/bin/sh\necho \"cannot open\" >&2\nexit 1\n")
	t.Setenv("PATH", binDir)

	err := gittags.OpenURL("https://example.com/release")
	require.ErrorContains(t, err, `open "https://example.com/release"`)
}
