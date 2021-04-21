package providers

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/jonaz/mgit/git"
	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var ErrNoProviderFound = fmt.Errorf("no provider found")

type Provider interface {
	Clone(whitelist []string, hasFile string) error
	Git(args []string) error
	PR() error
	Replace(regexp, with, fileRegexp, pathRegex string, contentRegex []string) error
	ShouldProcessRepo(path string) (bool, error)
	WorkDir() string
}

func GetProvider(c *cli.Context) (Provider, error) {
	if c.String("bitbucket-url") != "" {
		return NewBitbucket(
			c.String("dir"),
			c.String("bitbucket-url"),
		), nil
	}

	return &DefaultProvider{Dir: c.String("dir")}, nil
}

type DefaultProvider struct {
	Dir              string
	RepoURLWhitelist []string
}

func (d *DefaultProvider) WorkDir() string {
	return d.Dir
}

func (d *DefaultProvider) Clone(whitelist []string, hasFile string) error {
	return fmt.Errorf("clone is not defined in this provider. Set --<provider>-url flag")
}

func (d *DefaultProvider) PR() error {
	return fmt.Errorf("prd is not defined in this provider. Set --<provider>-url flag")
}

func (d *DefaultProvider) ShouldProcessRepo(path string) (bool, error) {
	if len(d.RepoURLWhitelist) == 0 {
		return true, nil
	}
	r := git.NewRepo(path)
	origin, err := r.RemoteURL()
	if err != nil {
		return false, err
	}
	return utils.InSlice(d.RepoURLWhitelist, origin), nil
}

func (d *DefaultProvider) Git(gitArgs []string) error {
	return utils.InEachRepo(d.Dir, func(path string) error {
		if ok, err := d.ShouldProcessRepo(path); !ok {
			if err != nil {
				return err
			}
			return nil
		}
		logrus.Info(path)
		args := []string{"-C", path}
		args = append(args, gitArgs...)
		return utils.RunInteractive("git", args...)
	})
}

var (
	ErrMissingWithFlag   = fmt.Errorf("missing --with flag to replace with")
	ErrMissingRegexpFlag = fmt.Errorf("missing --regexp flag to find what to replace")
)

func (d *DefaultProvider) Replace(regex, with, fileRegex, pathRegex string, contentRegex []string) error {
	if with == "" {
		return ErrMissingWithFlag
	}
	if regex == "" {
		return ErrMissingRegexpFlag
	}

	reg, err := regexp.Compile(regex)
	if err != nil {
		return err
	}

	fileReg, err := regexp.Compile(fileRegex)
	if err != nil {
		return err
	}
	pathReg, err := regexp.Compile(pathRegex)
	if err != nil {
		return err
	}

	contentReg := make([]*regexp.Regexp, len(contentRegex))
	for i, s := range contentRegex {
		reg, err := regexp.Compile(s)
		if err != nil {
			return err
		}
		contentReg[i] = reg
	}

	return utils.InEachRepo(d.Dir, func(path string) error {
		if ok, err := d.ShouldProcessRepo(path); !ok {
			if err != nil {
				return err
			}
			return nil
		}
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
			if pathRegex != "" {
				if !pathReg.MatchString(path) {
					return nil
				}
			}

			logrus.Debugf("checking path %s for matching regexp", path)
			read, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			for _, reg := range contentReg {
				if !reg.Match(read) {
					return nil
				}
			}

			if !reg.Match(read) {
				return nil
			}

			logrus.Infof("found file to replace in: %s", path)
			newContent := reg.ReplaceAll(read, []byte(with))
			err = ioutil.WriteFile(path, newContent, info.Mode())
			if err != nil {
				return fmt.Errorf("error saving file after replace: %w", err)
			}
			return nil
		})
	})
}
