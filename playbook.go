package main

import (
	"errors"
	"os"
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
				Replace: []models.Replace{
					{
						Regexp:     "",
						With:       "",
						FileRegexp: "",
					},
				},
			},
		}}
	err := utils.InEachRepo(c.String("dir"), func(path string) error {
		origin, err := utils.Run("git", "-C", path, "config", "--get", "remote.origin.url")
		if err != nil {
			return err
		}
		origin = strings.TrimSpace(origin)

		p.Tasks[0].Repos = append(p.Tasks[0].Repos, origin)

		return nil
	})
	if err != nil {
		return err
	}

	return yaml.NewEncoder(os.Stdout).Encode(p)
}

func playbook(c *cli.Context) error {
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

	provider, err := providers.GetProvider(c)
	if err != nil {
		return err
	}

	for _, task := range p.Tasks {
		for _, repoURL := range task.Repos {
			tmp := strings.Split(repoURL, "/")
			dir := filepath.Join(c.String("dir"), tmp[len(tmp)-2]+"_"+tmp[len(tmp)-1])
			dir = strings.TrimSuffix(dir, ".git")
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
					continue
				}
			}
		}

		for _, rep := range task.Replace {
			err := provider.Replace(rep.Regexp, rep.With, rep.FileRegexp)
			if err != nil {
				logrus.Error(err)
				continue
			}
		}
	}

	return nil
}
