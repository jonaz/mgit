package providers

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var ErrNoProviderFound = fmt.Errorf("no provider found")

type Provider interface {
	Clone(whitelist []string, hasFile string) error
	Git(args []string) error
	PR() error
	Replace(regexp, with, fileRegexp string) error
}

func GetProvider(c *cli.Context) (Provider, error) {
	if c.String("bitbucket-url") != "" {
		return NewBitbucket(
			c.String("dir"),
			c.String("bitbucket-url"),
		), nil
	}

	return nil, ErrNoProviderFound
}

type DefaultProvider struct {
	Dir string
}

func (d *DefaultProvider) Git(gitArgs []string) error {
	return utils.InEachRepo(d.Dir, func(path string) error {
		logrus.Info(path)
		args := []string{"-C", path}
		args = append(args, gitArgs...)
		return utils.RunInteractive("git", args...)
	})
}

func (d *DefaultProvider) Replace(regex, with, fileRegex string) error {
	if with == "" {
		return fmt.Errorf("missing --with flag to replace with")
	}
	if regex == "" {
		return fmt.Errorf("missing --regexp flag to find what to replace")
	}

	reg, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	fileReg, err := regexp.Compile(fileRegex)
	if err != nil {
		return err
	}

	return utils.InEachRepo(d.Dir, func(path string) error {
		logrus.Infof("scanning repo for replace: %s", path)

		return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			if fileRegex != "" {
				if !fileReg.MatchString(info.Name()) {
					return nil
				}
			}

			logrus.Debugf("checking path %s for matching regexp", path)
			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			if !reg.Match(read) {
				return nil
			}

			logrus.Infof("found file to replace in: %s", path)
			newContent := reg.ReplaceAll(read, []byte(with))
			err = ioutil.WriteFile(path, newContent, info.Mode())
			if err != nil {
				return err
			}
			return nil
		})
	})
}
