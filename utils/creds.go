package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func Credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, "Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Fprint(os.Stderr, "Enter Password: ")

	// golangci-lint says unnecessary conversion (unconvert) on this line. This is not unnecessary. Only works on windows if its this ways.
	pw, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println()

	return strings.TrimSpace(username), strings.TrimSpace(string(pw)), nil
}
