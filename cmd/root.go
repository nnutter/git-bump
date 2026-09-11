// Package cmd implements the git-bump command line interface.
package cmd

import (
	"context"
	"errors"

	"charm.land/fang/v2"
	"github.com/spf13/cobra"
)

// NewRootCommand builds the git-bump root command.
func NewRootCommand() *cobra.Command {
	var major bool
	var minor bool
	var patch bool
	var pattern string
	var noPush bool

	cmd := &cobra.Command{
		Use:   "git-bump",
		Short: "Bump the latest semver tag",
		RunE: func(_ *cobra.Command, _ []string) error {
			count := 0
			for _, bump := range []bool{major, minor, patch} {
				if bump {
					count++
				}
			}
			if count != 1 {
				return errors.New("exactly one of --major, --minor, or --patch is required")
			}
			return errors.New("not implemented")
		},
	}

	cmd.Flags().BoolVar(&major, "major", false, "Bump the major version")
	cmd.Flags().BoolVar(&minor, "minor", false, "Bump the minor version")
	cmd.Flags().BoolVar(&patch, "patch", false, "Bump the patch version")
	cmd.Flags().StringVar(&pattern, "pattern", "", "Only consider tags matching `git tag -l <pattern>`")
	cmd.Flags().BoolVar(&noPush, "no-push", false, "Do not push the new tag")
	cmd.MarkFlagsMutuallyExclusive("major", "minor", "patch")

	return cmd
}

// Execute runs the root command with Fang styling.
func Execute() error {
	return fang.Execute(context.Background(), NewRootCommand())
}
