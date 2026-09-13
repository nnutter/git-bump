package cmd

import (
	"runtime/debug"
)

// resolveVersion reports the version to display.
//
// The link-time value wins when set (go build -ldflags
// "-X main.version=<version>"). Otherwise it defaults to the short
// VCS revision stamped into the binary at build time, suffixed with
// -dirty when the worktree was modified, or "dev" when no revision
// is available.
func resolveVersion(injected string, settings []debug.BuildSetting) string {
	if injected != "" {
		return injected
	}
	var revision, modified string
	for _, setting := range settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.modified":
			modified = setting.Value
		}
	}
	if revision == "" {
		return "dev"
	}
	short := revision
	if len(short) > 12 {
		short = short[:12]
	}
	if modified == "true" {
		short += "-dirty"
	}
	return short
}

// buildSettings reports the current binary's build settings, or nil
// when build info is unavailable.
func buildSettings() []debug.BuildSetting {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}
	return info.Settings
}
