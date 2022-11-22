package utils

import (
	"errors"
	"runtime"
)

// OpenBrowser opens a link in the correct browser depending on OS.
func OpenBrowser(link string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		_, err = Run("xdg-open", link)
	case "darwin":
		_, err = Run("open", link)
	case "windows":
		_, err = Run("rundll32", "url.dll,FileProtocolHandler", link)
	default:
		return errors.New("Unknown operating system, dont know how to open the link in the browser")
	}
	return err
}
