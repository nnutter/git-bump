package cmd

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveVersion(t *testing.T) {
	revision := "5ed0fe7d46723453308c9ef7cf6312399d87ad66"
	settings := func(pairs ...string) []debug.BuildSetting {
		var out []debug.BuildSetting
		for i := 0; i < len(pairs); i += 2 {
			out = append(out, debug.BuildSetting{Key: pairs[i], Value: pairs[i+1]})
		}
		return out
	}

	tests := []struct {
		name     string
		injected string
		settings []debug.BuildSetting
		want     string
	}{
		{name: "injected wins", injected: "v1.2.3", settings: settings("vcs.revision", revision), want: "v1.2.3"},
		{name: "revision shortens", settings: settings("vcs.revision", revision), want: "5ed0fe7d4672"},
		{
			name:     "modified suffix",
			settings: settings("vcs.revision", revision, "vcs.modified", "true"),
			want:     "5ed0fe7d4672-dirty",
		},
		{name: "clean tree", settings: settings("vcs.revision", revision, "vcs.modified", "false"), want: "5ed0fe7d4672"},
		{name: "no revision", settings: settings("vcs.modified", "true"), want: "dev"},
		{name: "no settings", want: "dev"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, resolveVersion(tt.injected, tt.settings))
		})
	}
}
