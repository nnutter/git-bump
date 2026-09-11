// Package cmd implements the git-bump command line interface.
package cmd

import (
	"context"
	"errors"
	"fmt"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"

	"github.com/nnutter/git-bump/internal/bump"
	"github.com/nnutter/git-bump/internal/gittags"
)

// NewRootCommand builds the git-bump root command.
//
// injectedVersion sets the reported version, normally the value of
// -X main.version at link time. When empty the VCS revision from
// build info is used.
func NewRootCommand(injectedVersion string) *cobra.Command {
	var major bool
	var minor bool
	var patch bool
	var pattern string
	var noPush bool

	cmd := &cobra.Command{
		Use:     "git-bump",
		Short:   "Bump the latest semver tag",
		Version: resolveVersion(injectedVersion, buildSettings()),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var kind bump.Kind
			count := 0
			for _, flag := range []struct {
				set  bool
				kind bump.Kind
			}{
				{major, bump.Major{}},
				{minor, bump.Minor{}},
				{patch, bump.Patch{}},
			} {
				if flag.set {
					kind, count = flag.kind, count+1
				}
			}
			if count != 1 {
				return errors.New("exactly one of --major, --minor, or --patch is required")
			}
			repo, err := gittags.Open(".")
			if err != nil {
				return err
			}
			latest, err := gittags.Latest(repo, pattern)
			if err != nil {
				return err
			}
			next, err := bump.Bump(latest, kind)
			if err != nil {
				return err
			}
			if err := gittags.Create(repo, next); err != nil {
				return err
			}
			if !noPush {
				if err := gittags.Push(repo, next); err != nil {
					return err
				}
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), next)
			return err
		},
	}

	cmd.Flags().BoolVar(&major, "major", false, "Bump the major version")
	cmd.Flags().BoolVar(&minor, "minor", false, "Bump the minor version")
	cmd.Flags().BoolVar(&patch, "patch", false, "Bump the patch version")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Only consider tags matching `git tag -l <pattern>`")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Do not push the new tag")

	return cmd
}

// Execute runs the root command with Fang styling.
func Execute(injectedVersion string) error {
	root := NewRootCommand(injectedVersion)
	return fang.Execute(context.Background(), root, fang.WithVersion(root.Version))
}
