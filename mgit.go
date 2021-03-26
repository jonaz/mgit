package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: true})

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
				Name:   "clone",
				Usage:  "clone all repos and make sure they are up to date",
				Action: multiClone,
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
				Action:  multiPR,
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
				Name:    "git",
				Aliases: []string{"pr"},
				Usage:   "run git commands in multiple repos",
				Action:  multiGit,
			},

			{
				Name:    "replace",
				Aliases: []string{"pr"},
				Usage:   "replace text in multiple repos",
				Action:  multiReplace,
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
