// Package main is the entry point for git-bump.
package main

import (
	"context"
	"os"

	"charm.land/fang/v2"

	"github.com/nnutter/git-bump/internal/cmd"
)

// version is the release version, set at link time:
//
//	go build -ldflags "-X main.version=<version>" .
//
// When empty, the command defaults to the VCS revision from build
// info.
var version string

func main() {
	root := cmd.NewRootCommand(version)
	if err := fang.Execute(context.Background(), root, fang.WithVersion(root.Version)); err != nil {
		os.Exit(1)
	}
}
