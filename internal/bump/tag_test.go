package bump_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/bump"
	"github.com/nnutter/git-bump/internal/testenv"
)

func TestTag(t *testing.T) {
	testenv.Sterilize(t)

	tests := []struct {
		name string
		tag  string
		kind bump.Kind
		want string
	}{
		{name: "patch", tag: "v1.2.3", kind: bump.Patch{}, want: "v1.2.4"},
		{name: "minor", tag: "v1.2.3", kind: bump.Minor{}, want: "v1.3.0"},
		{name: "major", tag: "v1.2.3", kind: bump.Major{}, want: "v2.0.0"},
		{name: "patch zero", tag: "v0.0.0", kind: bump.Patch{}, want: "v0.0.1"},
		{name: "minor zero", tag: "v0.0.0", kind: bump.Minor{}, want: "v0.1.0"},
		{name: "major zero", tag: "v0.0.0", kind: bump.Major{}, want: "v1.0.0"},
		{name: "multi digit", tag: "v10.20.30", kind: bump.Patch{}, want: "v10.20.31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bump.Tag(tt.tag, tt.kind)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTagInvalid(t *testing.T) {
	testenv.Sterilize(t)

	tests := []struct {
		name string
		tag  string
		kind bump.Kind
	}{
		{name: "missing v prefix", tag: "1.2.3", kind: bump.Patch{}},
		{name: "missing patch", tag: "v1.2", kind: bump.Patch{}},
		{name: "missing minor", tag: "v1", kind: bump.Patch{}},
		{name: "empty", tag: "", kind: bump.Patch{}},
		{name: "not a version", tag: "release-foo", kind: bump.Patch{}},
		{name: "pre-release", tag: "v1.2.3-rc.1", kind: bump.Patch{}},
		{name: "build metadata", tag: "v1.2.3+build", kind: bump.Patch{}},
		{name: "leading text", tag: "release-v1.2.3", kind: bump.Patch{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bump.Tag(tt.tag, tt.kind)
			require.Error(t, err)
		})
	}
}

func TestTagNilKindPanics(t *testing.T) {
	testenv.Sterilize(t)

	require.Panics(t, func() {
		_, _ = bump.Tag("v1.2.3", nil)
	})
}

func TestKindString(t *testing.T) {
	testenv.Sterilize(t)

	require.Equal(t, "major", bump.Major{}.String())
	require.Equal(t, "minor", bump.Minor{}.String())
	require.Equal(t, "patch", bump.Patch{}.String())
}
