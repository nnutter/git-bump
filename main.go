// Package main is the entry point for git-bump.
package main

import (
	"os"

	"github.com/nnutter/git-bump/cmd"
)

// version is the release version, set at link time:
//
//	go build -ldflags "-X main.version=<version>" .
//
// When empty, the command defaults to the VCS revision from build
// info.
var version string

func main() {
	if err := cmd.Execute(version); err != nil {
		os.Exit(1)
	}
}
