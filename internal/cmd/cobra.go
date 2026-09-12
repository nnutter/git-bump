package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand builds the git-bump root command.
//
// injectedVersion sets the reported version, normally the value of
// -X main.version at link time. When empty the VCS revision from
// build info is used.
func NewRootCommand(injectedVersion string) *cobra.Command {
	var opts options

	cmd := &cobra.Command{
		Use:     "git-bump",
		Short:   "Bump the latest semver tag",
		Version: resolveVersion(injectedVersion, buildSettings()),
		RunE:    opts.runE,
	}

	cmd.Flags().BoolVar(&opts.major, "major", false, "Bump the major version")
	cmd.Flags().BoolVar(&opts.minor, "minor", false, "Bump the minor version")
	cmd.Flags().BoolVar(&opts.patch, "patch", false, "Bump the patch version")
	cmd.Flags().StringVar(&opts.pattern, "pattern", "", "Only consider tags matching `git tag -l <pattern>`")
	cmd.Flags().BoolVar(&opts.noPush, "no-push", false, "Do not push the new tag")
	cmd.Flags().BoolVar(&opts.release, "release", false, "Create a draft GitHub release for the new tag")

	return cmd
}
