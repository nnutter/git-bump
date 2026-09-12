package cmd_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/cmd"
)

func TestVersionFlag(t *testing.T) {
	command := cmd.NewRootCommand("v9.9.9")
	command.SetArgs([]string{"--version"})
	output := &bytes.Buffer{}
	command.SetOut(output)
	require.NoError(t, command.Execute())
	require.Contains(t, output.String(), "v9.9.9")
}
