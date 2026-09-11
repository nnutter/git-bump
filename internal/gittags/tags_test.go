package gittags_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nnutter/git-bump/internal/gittags"
	"github.com/nnutter/git-bump/internal/testrepo"
)

func TestLatest(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		pattern string
		want    string
	}{
		{name: "single", tags: []string{"v1.2.3"}, want: "v1.2.3"},
		{name: "highest wins", tags: []string{"v1.9.0", "v1.10.0", "v2.0.0"}, want: "v2.0.0"},
		{name: "numeric not lexical", tags: []string{"v1.10.0", "v1.9.0"}, want: "v1.10.0"},
		{name: "skips non-semver", tags: []string{"release-foo", "v1.2.3"}, want: "v1.2.3"},
		{name: "pattern filters", tags: []string{"v1.2.3", "v2.0.0"}, pattern: "v1.*", want: "v1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testrepo.Init(t, tt.tags...)
			got, err := gittags.Latest(repo, tt.pattern)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestLatestErrors(t *testing.T) {
	t.Run("no tags", func(t *testing.T) {
		repo := testrepo.Init(t)
		_, err := gittags.Latest(repo, "")
		require.ErrorContains(t, err, "no tags found")
	})

	t.Run("pattern matches nothing", func(t *testing.T) {
		repo := testrepo.Init(t, "v1.2.3")
		_, err := gittags.Latest(repo, "v2.*")
		require.ErrorContains(t, err, `no tags match pattern "v2.*"`)
	})

	t.Run("only non-semver tags", func(t *testing.T) {
		repo := testrepo.Init(t, "release-foo")
		_, err := gittags.Latest(repo, "")
		require.ErrorContains(t, err, "no tags found")
	})

	t.Run("pattern matches only non-semver", func(t *testing.T) {
		repo := testrepo.Init(t, "v1.2.3", "release-1")
		_, err := gittags.Latest(repo, "release-*")
		require.ErrorContains(t, err, `no tags match pattern "release-*"`)
	})

	t.Run("invalid pattern", func(t *testing.T) {
		repo := testrepo.Init(t, "v1.2.3")
		_, err := gittags.Latest(repo, "[")
		require.ErrorContains(t, err, `invalid pattern`)
	})
}
