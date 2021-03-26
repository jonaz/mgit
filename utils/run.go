package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

// Run runs a command.
func Run(head string, parts ...string) (string, error) {
	var err error

	head, err = exec.LookPath(head)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(head, parts...) // #nosec
	cmd.Env = os.Environ()

	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return stdout.String(), errors.New(err.Error() + stderr.String())
	}
	return stdout.String(), nil
}

// RunInteractive runs a command while attaching to stdout and stderr.
func RunInteractive(head string, parts ...string) error {
	head, err := exec.LookPath(head)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	cmd := exec.Command(head, parts...) // #nosec
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
