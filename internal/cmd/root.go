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
	noOpen  bool
}

func (o *options) runE(cmd *cobra.Command, _ []string) error {
	var kind bump.Kind
	count := 0
	for _, flag := range []struct {
		set  bool
		kind bump.Kind
	}{
		{o.major, bump.Major{}},
		{o.minor, bump.Minor{}},
		{o.patch, bump.Patch{}},
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
	latest, err := gittags.Latest(repo, o.pattern)
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
	if !o.noPush {
		if err := gittags.Push(repo, next); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), next); err != nil {
		return err
	}
	if o.release {
		url, err := gittags.CreateDraftRelease(repo, next)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), url); err != nil {
			return err
		}
		if !o.noOpen {
			if err := gittags.OpenURL(url); err != nil {
				if _, err := fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not open browser: %v\n", err); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return nil
}
