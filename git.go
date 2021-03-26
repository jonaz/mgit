package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jonaz/mgit/config"
	"github.com/jonaz/mgit/git"
	"github.com/jonaz/mgit/service"
	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func multiClone(c *cli.Context) error {
	username, password, err := utils.Credentials()
	if err != nil {
		return err
	}

	config := config.Config{
		BitbucketURL:      c.String("bitbucket-url"),
		BitbucketUser:     username,
		BitbucketPassword: password,
	}
	bit := service.NewBitbucket(config)

	projects, err := bit.ListProjects()
	if err != nil {
		return err
	}

	for _, project := range projects.Values {
		repos, err := bit.ListRepos(project.Key)
		if err != nil {
			logrus.Error(err)
			continue
		}

		for _, repo := range repos.Values {
			if c.String("whitelist") != "" {
				s := strings.Split(c.String("whitelist"), ",")
				if !inSlice(s, repo.Name) {
					logrus.Debugf("skipping repo %s", repo.Slug)
					continue
				}

				logrus.Infof("cloning repo %s", repo.Name)
				git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(c.String("dir")))

				continue
			}

			if c.String("has-file") != "" {
				files, err := bit.ListFiles(project.Key, repo.Slug)
				if err != nil {
					logrus.Error(err)
					continue
				}
				if inSlice(files.Values, c.String("has-fil")) {
					logrus.Infof("cloning repo %s", repo.Name)
					git.Clone(repo.Links.Clone.GetSSH(), repo.RepoPath(c.String("dir")))
				}
				continue
			}
		}
	}

	return nil
}

func inSlice(files []string, filename string) bool {
	for _, file := range files {
		if file == filename {
			return true
		}
	}
	return false
}
func multiPR(c *cli.Context) error {
	return inEachRepo(c, func(path string) error {
		logrus.Info(path)
		OpenBitbucketPR(c.String("bitbucket-url"), "", path)
		return nil
	})
}
func multiGit(c *cli.Context) error {
	return inEachRepo(c, func(path string) error {
		logrus.Info(path)
		args := []string{"-C", path}
		args = append(args, c.Args().Slice()...)
		return utils.RunInteractive("git", args...)
	})
}
func multiReplace(c *cli.Context) error {
	if c.String("with") == "" {
		return fmt.Errorf("missing --with flag to replace wih")
	}
	if c.String("regexp") == "" {
		return fmt.Errorf("missing --regexp flag to find what to replace")
	}

	reg, _ := regexp.Compile(c.String("regexp"))
	fileReg, _ := regexp.Compile(c.String("file-regexp"))
	return inEachRepo(c, func(path string) error {
		logrus.Info(path)

		return filepath.Walk(c.String("dir"), func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() && info.Name() == ".git" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				return nil
			}

			if c.String("file-regexp") != "" {
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

			newContent := reg.ReplaceAll(read, []byte(c.String("with")))
			err = ioutil.WriteFile(path, []byte(newContent), info.Mode())
			if err != nil {
				return err
			}
			return nil
		})

	})
}

func inEachRepo(c *cli.Context, fn func(path string) error) error {
	files, err := ioutil.ReadDir(c.String("dir"))
	if err != nil {
		return err
	}
	for _, v := range files {
		if !v.IsDir() {
			continue // skip non dirs
		}
		err := fn(filepath.Join(c.String("dir"), v.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenBitbucketPR(bitbucketURL, targetBranch, gitDir string) error {

	if targetBranch == "" {
		targetBranch = "master"
		_, err := utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "master")
		if err != nil && err.Error() == "exit status 1" {
			targetBranch = "develop"
			_, err = utils.Run("git", "-C", gitDir, "rev-parse", "-q", "--verify", "develop")
			if err != nil && err.Error() == "exit status 1" {
				return err
			}
		}
	}

	origin, err := utils.Run("git", "-C", gitDir, "config", "--get", "remote.origin.url")
	if err != nil {
		return err
	}

	u, err := url.Parse(strings.TrimSpace(origin))
	if err != nil {
		return err
	}

	sourceBranch, err := utils.Run("git", "-C", gitDir, "symbolic-ref", "HEAD")
	if err != nil {
		return err
	}

	paths := strings.Split(u.Path, "/")
	burl := path.Join(bitbucketURL, "projects/%s/repos/%s/pull-requests?create&targetBranch=%s&sourceBranch=%s")
	prURL := fmt.Sprintf(burl, strings.TrimSpace(paths[1]), strings.TrimSuffix(strings.TrimSpace(paths[2]), ".git"), "refs/heads/"+targetBranch, strings.TrimSpace(sourceBranch))
	fmt.Println("Opening in browser:", prURL)
	return utils.OpenBrowser(prURL)
}
