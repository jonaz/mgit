package main

import (
	"errors"
	"os"
	"os/exec"

	"github.com/google/shlex"
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

func readFile(fn string) (*models.Playbook, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	p := &models.Playbook{}
	err = yaml.NewDecoder(file).Decode(&p)
	return p, err
}

func runPlaybook(c *cli.Context) error {
	p, err := readFile(c.Args().First())
	if err != nil {
		return err
	}

	for _, task := range p.Tasks {
		task := task
		provider := &providers.DefaultProvider{Dir: c.String("dir")}

		for _, repoURL := range task.Repos {
			dir := utils.RepoDir(c.String("dir"), repoURL)
			logrus.Infof("clone %s into %s", repoURL, dir)
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

		// ignore unchanged repos
		newWhitelist := []string{}
		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			hasChanges, err := repo.HasLocalChanges()
			if err != nil {
				return err
			}

			if hasChanges {
				u, _ := repo.RemoteURL()
				newWhitelist = append(newWhitelist, u)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(newWhitelist) == 0 {
			logrus.Info("found no changes in repos in current play. aborting")
			return nil
		}
		provider.RepoURLWhitelist = newWhitelist

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			logrus.Infof("%s: create branch %s", repo.WorkDir(), task.TargetBranch)
			return repo.Checkout(task.TargetBranch)
		})
		if err != nil {
			return err
		}

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			logrus.Infof("%s: commit %s", repo.WorkDir(), task.TargetBranch)
			return repo.Commit(task.CommitMessage, "")
		})
		if err != nil {
			return err
		}

		err = eachRepoInPlay(provider, func(repo git.Repo) error {
			repoURL, err := repo.RemoteURL()
			if err != nil {
				return err
			}
			logrus.Infof("%s: push %s", repo.WorkDir(), repoURL)
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

func openPR(c *cli.Context) error {
	p, err := readFile(c.Args().First())
	if err != nil {
		return err
	}
	provider, err := providers.GetProvider(c)
	if err != nil {
		return err
	}
	repos := []string{}
	for _, v := range p.Tasks {
		repos = append(repos, v.Repos...)
	}
	return provider.PR(repos)
}

func runAction(provider providers.Provider, action models.Action) error {
	if action.Regexp != "" && action.Command == "" {
		err := provider.Replace(action.Regexp, action.With, action.FileRegexp, action.PathRegexp, action.ContentRegexp)
		if err != nil {
			return err
		}
	}
	if action.Regexp == "" && action.Command != "" {
		if len(action.ContentRegexp) > 0 || action.FileRegexp != "" || action.PathRegexp != "" {
			return provider.CommandEachMatchingFile(action.Command, action.FileRegexp, action.PathRegexp, action.ContentRegexp)
		}

		err := eachRepoInPlay(provider, func(repo git.Repo) error {
			logrus.Infof("%s: running command: %s", repo.WorkDir(), action.Command)
			args, err := shlex.Split(action.Command)
			if err != nil {
				return err
			}
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
