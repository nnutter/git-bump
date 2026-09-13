package gittags

import (
	"fmt"
	"os/exec"
)

// OpenURL opens url in the browser, preferring the open command
// (macOS, or anywhere it exists) with a fallback to xdg-open (Linux,
// if it exists). When neither opener exists it is a silent no-op so
// headless environments still succeed.
func OpenURL(url string) error {
	opener, err := exec.LookPath("open")
	if err != nil {
		opener, err = exec.LookPath("xdg-open")
		if err != nil {
			return nil
		}
	}
	if out, err := exec.Command(opener, url).CombinedOutput(); err != nil {
		return fmt.Errorf("open %q: %w: %s", url, err, out)
	}
	return nil
}
