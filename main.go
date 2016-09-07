package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "slack plugin"
	app.Usage = "slack plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "webhook",
			Usage:  "slack webhook url",
			EnvVar: "SLACK_WEBHOOK,PLUGIN_WEBHOOK",
		},
		cli.StringFlag{
			Name:   "channel",
			Usage:  "slack channel",
			EnvVar: "PLUGIN_CHANNEL",
		},
		cli.StringFlag{
			Name:   "recipient",
			Usage:  "slack recipient",
			EnvVar: "PLUGIN_RECIPIENT",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "slack username",
			Value:  "drone",
			EnvVar: "PLUGIN_USERNAME",
		},
		cli.StringFlag{
			Name:   "template",
			Usage:  "slack template",
			EnvVar: "PLUGIN_TEMPLATE",
		},
		cli.StringFlag{
			Name:   "image",
			Usage:  "slack image url",
			EnvVar: "PLUGIN_IMAGE_URL",
		},
		cli.StringFlag{
			Name:   "icon.url",
			Usage:  "slack icon url",
			EnvVar: "PLUGIN_ICON_URL",
		},
		cli.StringFlag{
			Name:   "icon.emoji",
			Usage:  "slack emoji url",
			EnvVar: "PLUGIN_ICON_EMOJI",
		},
		cli.StringFlag{
			Name:   "repo.owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "build.event",
			Value:  "push",
			Usage:  "build event",
			EnvVar: "DRONE_BUILD_EVENT",
		},
		cli.IntFlag{
			Name:   "build.number",
			Usage:  "build number",
			EnvVar: "DRONE_BUILD_NUMBER",
		},
		cli.StringFlag{
			Name:   "build.status",
			Usage:  "build status",
			Value:  "success",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "build.link",
			Usage:  "build link",
			EnvVar: "DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
		cli.BoolTFlag{
			Name:  "letter",
			Usage: "Enable / disable .Pipeline-Letter proccesing",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repo: Repo{
			Owner: c.String("repo.owner"),
			Name:  c.String("repo.name"),
		},
		Build: Build{
			Number: c.Int("build.number"),
			Event:  c.String("build.event"),
			Status: c.String("build.status"),
			Commit: c.String("commit.sha"),
			Branch: c.String("commit.branch"),
			Author: c.String("commit.author"),
			Link:   c.String("build.link"),
		},
		Config: Config{
			Webhook:   c.StringSlice("webhook"),
			Channel:   c.String("channel"),
			Recipient: c.String("recipient"),
			Username:  c.String("username"),
			Template:  c.String("template"),
			ImageURL:  c.String("image"),
			IconURL:   c.String("icon.url"),
			IconEmoji: c.String("icon.emoji"),
			Letter:    c.Bool("letter"),
		},
	}

	return plugin.Exec()
}
