package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// GitRoot find the root of the git dir and returns the absolute path.
func Root() (string, error) {
	tmp, err := exec.Command("git", "rev-parse", "--show-toplevel").Output() // #nosec
	if err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("stderr: %s", err.Error())
		}
		return "", err
	}
	return strings.TrimSpace(string(tmp)), nil
}
