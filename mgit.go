package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonaz/mgit/providers"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true})

	app := &cli.App{
		Name:  "mgit",
		Usage: "manage multiple git repos at the same time",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bitbucket-url",
				Value:   "",
				Usage:   "regexp kubectl context name",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:  "dir",
				Value: ".",
				Usage: "temporary working directory for all the git repos",
			},
			&cli.StringFlag{
				Name:  "loglevel",
				Value: "info",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "clone",
				Usage: "clone all repos and make sure they are up to date",
				Action: func(c *cli.Context) error {
					provider, err := providers.GetProvider(c)
					if err != nil {
						return err
					}
					return provider.Clone(c.String("whitelist"), c.String("has-file"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "has-file",
						Usage: "only clone repo which has file",
					},

					&cli.StringFlag{
						Name:  "whitelist",
						Usage: "only clone repos in comma separated list",
					},
				},
			},
			{
				Name:    "pull-request",
				Aliases: []string{"pr"},
				Usage:   "multip PR open",
				Action: func(c *cli.Context) error {
					provider, err := providers.GetProvider(c)
					if err != nil {
						return err
					}
					return provider.PR()
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "has-file",
						Usage: "only clone repo which has file",
					},

					&cli.StringFlag{
						Name:  "whitelist",
						Usage: "only clone repos in comma separated list",
					},
				},
			},
			{
				Name:  "git",
				Usage: "run git commands in multiple repos",
				Action: func(c *cli.Context) error {
					provider, err := providers.GetProvider(c)
					if err != nil {
						return err
					}
					return provider.Git(c.Args().Slice())
				},
			},

			{
				Name:  "replace",
				Usage: "replace text in multiple repos",
				Action: func(c *cli.Context) error {
					provider, err := providers.GetProvider(c)
					if err != nil {
						return err
					}
					return provider.Replace(c.String("regexp"), c.String("with"), c.String("file-regexp"))
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "with",
						Usage: "content to replace that line with",
					},
					&cli.StringFlag{
						Name:  "regexp",
						Usage: "regexp to find a line in file",
					},
					&cli.StringFlag{
						Name:  "file-regexp",
						Usage: "regexp to filter files",
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
