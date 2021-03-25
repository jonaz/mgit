package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jonaz/mgit/config"
	"github.com/jonaz/mgit/service"
	"github.com/jonaz/mgit/utils"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mkubectl",
		Usage: "run kubectl command in multiple contexts",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bitbucket-url",
				Value:   "",
				Usage:   "regexp kubectl context name",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "namespace",
				Value:   "",
				Usage:   "kubectl namespace",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "clone",
				Usage:  "clone all repos and make sure they are up to date",
				Action: clone,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "has-file",
						Usage: "only clone repo which has file",
					},
				},
			},
		},
	}
	app.Before = globalBefore

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func clone(c *cli.Context) error {
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
			fmt.Println(repo.Links.Clone)
			files, err := bit.ListFiles(project.Key, repo.Slug)
			if err != nil {
				logrus.Error(err)
				continue
			}
			for _, file := range files.Values {
				fmt.Println(file)
			}
		}
	}

	return nil
}

func globalBefore(c *cli.Context) error {
	lvl, err := logrus.ParseLevel(c.String("loglevel"))
	if err != nil {
		return err
	}
	if lvl != logrus.InfoLevel {
		fmt.Fprintf(os.Stderr, "using loglevel: %s\n", lvl.String())
	}
	logrus.SetLevel(lvl)
	return nil
}
