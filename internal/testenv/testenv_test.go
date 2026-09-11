package testenv_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/testenv"
)

func TestSterilize(t *testing.T) {
	t.Setenv("GIT_DIR", "/tmp/should-be-hidden")
	t.Setenv("GIT_WORK_TREE", "/tmp/should-be-hidden")

	testenv.Sterilize(t)

	_, ok := os.LookupEnv("GIT_DIR")
	require.False(t, ok, "GIT_DIR should be unset")
	_, ok = os.LookupEnv("GIT_WORK_TREE")
	require.False(t, ok, "GIT_WORK_TREE should be unset")
}
