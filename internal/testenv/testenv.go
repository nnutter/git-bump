// Package testenv sterilizes the test process environment from Git.
//
// Commands like git rebase --exec run tests with GIT_* variables set
// (GIT_DIR, GIT_WORK_TREE, GIT_INDEX_FILE, ...). Any test that creates
// or opens repositories must call Sterilize first so those leaked
// variables cannot change repository discovery.
package testenv

import (
	"os"
	"strings"
	"testing"
)

// Sterilize unsets every GIT_* environment variable for the duration
// of the test, restoring the original values on cleanup.
func Sterilize(t *testing.T) {
	t.Helper()

	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if !strings.HasPrefix(key, "GIT_") {
			continue
		}
		value, ok := os.LookupEnv(key)
		if !ok {
			continue
		}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("unset %s: %v", key, err)
		}
		t.Cleanup(func() {
			if err := os.Setenv(key, value); err != nil {
				t.Errorf("restore %s: %v", key, err)
			}
		})
	}
}
