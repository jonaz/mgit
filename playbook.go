package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jonaz/mgit/git"
	"github.com/jonaz/mgit/models"
	"github.com/jonaz/mgit/providers"
	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func generatePlaybook(c *cli.Context) error {
	p := &models.Playbook{
		Tasks: []models.Task{
			models.Task{
				Actions: []models.Action{
					{
						Command:    "",
						Regexp:     "",
						With:       "",
						FileRegexp: "",
					},
				},
			},
		},
	}
	err := utils.InEachRepo(c.String("dir"), func(path string) error {
		r := git.NewRepo(path)
		origin, err := r.RemoteURL()
		if err != nil {
			return err
		}
		p.Tasks[0].Repos = append(p.Tasks[0].Repos, origin)

		return nil
	})
	if err != nil {
		return err
	}

	return yaml.NewEncoder(os.Stdout).Encode(p)
}

func runPlaybook(c *cli.Context) error {
	file, err := os.Open(c.Args().First())
	if err != nil {
		return err
	}
	defer file.Close()

	p := &models.Playbook{}
	err = yaml.NewDecoder(file).Decode(&p)
	if err != nil {
		return err
	}

	for _, task := range p.Tasks {
		provider := &providers.DefaultProvider{Dir: c.String("dir")}

		for _, repoURL := range task.Repos {
			dir := repoDir(c.String("dir"), repoURL)
			logrus.Debugf("will clone into %s", dir)
			if _, err := os.Stat(dir); errors.Is(err, os.ErrNotExist) {
				_, err := git.Clone(repoURL, dir)
				if err != nil {
					logrus.Error(err)
					continue
				}
			} else {
				r := git.NewRepo(dir)
				err := r.Pull()
				if err != nil {
					logrus.Error(err)
				}
			}

			provider.RepoURLWhitelist = append(provider.RepoURLWhitelist, repoURL)
		}

		for _, a := range task.Actions {
			err := runAction(provider, a)
			if err != nil {
				return err
			}
		}

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			return repo.Checkout(task.TargetBranch)
		})
		if err != nil {
			return err
		}

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			return repo.CommitAndPush(task.CommitMessage, "")
		})
		if err != nil {
			return err
		}

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			repoURL, err := repo.RemoteURL()
			if err != nil {
				return err
			}
			return repo.Push(repoURL, c.Bool("force"))
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func eachRepoInPlay(provider providers.Provider, cb func(repo git.Repo) error) error {
	return utils.InEachRepo(provider.WorkDir(), func(path string) error {
		ok, err := provider.ShouldProcessRepo(path)
		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		r := git.NewRepo(path)

		return cb(r)
	})
}

func repoDir(workDir, repoURL string) string {
	tmp := strings.Split(repoURL, "/")
	dir := filepath.Join(workDir, tmp[len(tmp)-2]+"_"+tmp[len(tmp)-1])
	return strings.TrimSuffix(dir, ".git")
}

func runAction(provider providers.Provider, action models.Action) error {
	if action.Regexp != "" && action.Command == "" {
		err := provider.Replace(action.Regexp, action.With, action.FileRegexp)
		if err != nil {
			return err
		}
	}
	if action.Regexp == "" && action.Command != "" {
		err := eachRepoInPlay(provider, func(repo git.Repo) error {
			logrus.Infof("running command: %s", action.Command)
			args := strings.Fields(action.Command)
			cmd := exec.Command(args[0], args[1:]...) // #nosec
			cmd.Dir = repo.WorkDir()
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			return cmd.Run()
		})
		if err != nil {
			return err
		}
	}
	return nil
}
