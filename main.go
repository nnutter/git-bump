// Package main is the entry point for git-bump.
package main

import (
	"os"

	"github.com/nnutter/git-bump/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
