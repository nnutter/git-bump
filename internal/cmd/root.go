// Package cmd implements the git-bump command line interface.
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nnutter/git-bump/internal/bump"
	"github.com/nnutter/git-bump/internal/gittags"
)

// options holds the flag values for the root command.
type options struct {
	major   bool
	minor   bool
	patch   bool
	pattern string
	noPush  bool
	release bool
}

// NewRootCommand builds the git-bump root command.
//
// injectedVersion sets the reported version, normally the value of
// -X main.version at link time. When empty the VCS revision from
// build info is used.
func NewRootCommand(injectedVersion string) *cobra.Command {
	var opts options

	runE := func(cmd *cobra.Command, _ []string) error {
		var kind bump.Kind
		count := 0
		for _, flag := range []struct {
			set  bool
			kind bump.Kind
		}{
			{opts.major, bump.Major{}},
			{opts.minor, bump.Minor{}},
			{opts.patch, bump.Patch{}},
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
		latest, err := gittags.Latest(repo, opts.pattern)
		if err != nil {
			return err
		}
		next, err := bump.Tag(latest, kind)
		if err != nil {
			return err
		}
		if err := gittags.Create(repo, next); err != nil {
			return err
		}
		if !opts.noPush {
			if err := gittags.Push(repo, next); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), next); err != nil {
			return err
		}
		if opts.release {
			url, err := gittags.CreateDraftRelease(repo, next)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), url)
			return err
		}
		return nil
	}

	cmd := &cobra.Command{
		Use:     "git-bump",
		Short:   "Bump the latest semver tag",
		Version: resolveVersion(injectedVersion, buildSettings()),
		RunE:    runE,
	}

	cmd.Flags().BoolVar(&opts.major, "major", false, "Bump the major version")
	cmd.Flags().BoolVar(&opts.minor, "minor", false, "Bump the minor version")
	cmd.Flags().BoolVar(&opts.patch, "patch", false, "Bump the patch version")
	cmd.Flags().StringVar(&opts.pattern, "pattern", "", "Only consider tags matching `git tag -l <pattern>`")
	cmd.Flags().BoolVar(&opts.noPush, "no-push", false, "Do not push the new tag")
	cmd.Flags().BoolVar(&opts.release, "release", false, "Create a draft GitHub release for the new tag")

	return cmd
}
