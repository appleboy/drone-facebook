package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

// Version set at compile-time
var Version string

func main() {
	year := fmt.Sprintf("%v", time.Now().Year())
	app := cli.NewApp()
	app.Name = "facebook plugin"
	app.Usage = "facebook plugin"
	app.Copyright = "Copyright (c) " + year + " Bo-Yi Wu"
	app.Authors = []cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "page.token",
			Usage:  "Token is the access token of the Facebook page to send messages from.",
			EnvVar: "PLUGIN_FB_PAGE_TOKEN,FB_PAGE_TOKEN,PAGE_TOKEN,INPUT_FB_PAGE_TOKEN",
		},
		cli.StringFlag{
			Name:   "verify.token",
			Usage:  "The token used to verify facebook",
			EnvVar: "PLUGIN_FB_VERIFY_TOKEN,FB_VERIFY_TOKEN,VERIFY_TOKEN,INPUT_FB_VERIFY_TOKEN",
		},
		cli.BoolFlag{
			Name:   "verify",
			Usage:  "verifying webhooks on the Facebook Developer Portal",
			EnvVar: "PLUGIN_VERIFY,FB_VERIFY,INPUT_VERIFY",
		},
		cli.StringSliceFlag{
			Name:   "to",
			Usage:  "send message to user",
			EnvVar: "PLUGIN_TO,FB_TO,INPUT_TO",
		},
		cli.StringSliceFlag{
			Name:   "message",
			Usage:  "send facebook message",
			EnvVar: "PLUGIN_MESSAGE,FB_MESSAGE,INPUT_MESSAGE",
		},
		cli.StringSliceFlag{
			Name:   "image",
			Usage:  "image message",
			EnvVar: "PLUGIN_IMAGES,IMAGES,INPUT_PLUGIN_IMAGES",
		},
		cli.StringSliceFlag{
			Name:   "audio",
			Usage:  "audio message",
			EnvVar: "PLUGIN_AUDIOS,AUDIOS,INPUT_AUDIOS",
		},
		cli.StringSliceFlag{
			Name:   "video",
			Usage:  "video message",
			EnvVar: "PLUGIN_VIDEOS,VIDEOS,INPUT_VIDEOS",
		},
		cli.StringSliceFlag{
			Name:   "file",
			Usage:  "file message",
			EnvVar: "PLUGIN_FILES,FILES,INPUT_FILES",
		},
		cli.BoolFlag{
			Name:   "match.email",
			Usage:  "send message when only match email",
			EnvVar: "PLUGIN_ONLY_MATCH_EMAIL,INPUT_ONLY_MATCH_EMAIL",
		},
		cli.StringFlag{
			Name:   "repo",
			Usage:  "repository owner and repository name",
			EnvVar: "DRONE_REPO,GITHUB_REPOSITORY",
		},
		cli.StringFlag{
			Name:   "repo.namespace",
			Usage:  "repository namespace",
			EnvVar: "DRONE_REPO_OWNER,DRONE_REPO_NAMESPACE,GITHUB_ACTOR",
		},
		cli.StringFlag{
			Name:   "repo.name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "commit.sha",
			Usage:  "git commit sha",
			EnvVar: "DRONE_COMMIT_SHA,GITHUB_SHA",
		},
		cli.StringFlag{
			Name:   "commit.ref",
			Usage:  "git commit ref",
			EnvVar: "DRONE_COMMIT_REF,GITHUB_REF",
		},
		cli.StringFlag{
			Name:   "commit.branch",
			Value:  "master",
			Usage:  "git commit branch",
			EnvVar: "DRONE_COMMIT_BRANCH",
		},
		cli.StringFlag{
			Name:   "commit.link",
			Usage:  "git commit link",
			EnvVar: "DRONE_COMMIT_LINK",
		},
		cli.StringFlag{
			Name:   "commit.author",
			Usage:  "git author name",
			EnvVar: "DRONE_COMMIT_AUTHOR",
		},
		cli.StringFlag{
			Name:   "commit.author.email",
			Usage:  "git author email",
			EnvVar: "DRONE_COMMIT_AUTHOR_EMAIL",
		},
		cli.StringFlag{
			Name:   "commit.message",
			Usage:  "commit message",
			EnvVar: "DRONE_COMMIT_MESSAGE",
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
			Name:   "build.tag",
			Usage:  "build tag",
			EnvVar: "DRONE_TAG",
		},
		cli.Float64Flag{
			Name:   "job.started",
			Usage:  "job started",
			EnvVar: "DRONE_JOB_STARTED",
		},
		cli.Float64Flag{
			Name:   "job.finished",
			Usage:  "job finished",
			EnvVar: "DRONE_JOB_FINISHED",
		},
		cli.StringFlag{
			Name:   "env-file",
			Usage:  "source env file",
			EnvVar: "ENV_FILE",
		},
		cli.IntFlag{
			Name:   "port, P",
			Usage:  "webhook port",
			EnvVar: "FACEBOOK_WEBHOOK_PORT",
			Value:  8088,
		},
		cli.BoolFlag{
			Name:   "autotls",
			Usage:  "Auto tls mode",
			EnvVar: "PLUGIN_AUTOTLS,AUTOTLS",
		},
		cli.StringSliceFlag{
			Name:   "host",
			Usage:  "Auto tls host name",
			EnvVar: "PLUGIN_HOSTNAME,HOSTNAME",
		},
		cli.BoolFlag{
			Name:   "github",
			Usage:  "Boolean value, indicates the runtime environment is GitHub Action.",
			EnvVar: "PLUGIN_GITHUB,GITHUB",
		},
		cli.StringFlag{
			Name:   "github.workflow",
			Usage:  "The name of the workflow.",
			EnvVar: "GITHUB_WORKFLOW",
		},
		cli.StringFlag{
			Name:   "github.action",
			Usage:  "The name of the action.",
			EnvVar: "GITHUB_ACTION",
		},
		cli.StringFlag{
			Name:   "github.event.name",
			Usage:  "The webhook name of the event that triggered the workflow.",
			EnvVar: "GITHUB_EVENT_NAME",
		},
		cli.StringFlag{
			Name:   "github.event.path",
			Usage:  "The path to a file that contains the payload of the event that triggered the workflow. Value: /github/workflow/event.json.",
			EnvVar: "GITHUB_EVENT_PATH",
		},
		cli.StringFlag{
			Name:   "github.workspace",
			Usage:  "The GitHub workspace path. Value: /github/workspace.",
			EnvVar: "GITHUB_WORKSPACE",
		},
		cli.StringFlag{
			Name:   "app.secret",
			Usage:  "The app secret from the facebook developer portal",
			EnvVar: "PLUGIN_APP_SECRET,APP_SECRET",
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		GitHub: GitHub{
			Workflow:  c.String("github.workflow"),
			Workspace: c.String("github.workspace"),
			Action:    c.String("github.action"),
			EventName: c.String("github.event.name"),
			EventPath: c.String("github.event.path"),
		},
		Repo: Repo{
			FullName:  c.String("repo"),
			Namespace: c.String("repo.namespace"),
			Name:      c.String("repo.name"),
		},
		Commit: Commit{
			Sha:     c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Link:    c.String("commit.link"),
			Author:  c.String("commit.author"),
			Email:   c.String("commit.author.email"),
			Message: c.String("commit.message"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Started:  c.Float64("job.started"),
			Finished: c.Float64("job.finished"),
			PR:       c.String("pull.request"),
		},
		Config: Config{
			PageToken:   c.String("page.token"),
			VerifyToken: c.String("verify.token"),
			Verify:      c.Bool("verify"),
			MatchEmail:  c.Bool("match.email"),
			To:          c.StringSlice("to"),
			Message:     c.StringSlice("message"),
			Image:       c.StringSlice("image"),
			Audio:       c.StringSlice("audio"),
			Video:       c.StringSlice("video"),
			File:        c.StringSlice("file"),
			Port:        c.Int("port"),
			AutoTLS:     c.Bool("autotls"),
			Host:        c.StringSlice("host"),
			AppSecret:   c.String("app.secret"),
			GitHub:      c.Bool("github"),
		},
	}

	command := c.Args().Get(0)

	if command == "webhook" {
		return plugin.Webhook()
	}

	return plugin.Exec()
}
