package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/org-stats/internal/stats"
	"github.com/kyokomi/emoji"
	"github.com/tj/go-spin"
	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "org-stats"
	app.Version = version
	app.Author = "Carlos Alexandro Becker (caarlos0@gmail.com)"
	app.Usage = "Get the contributor stats summary from all repos of any given organization"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			EnvVar: "GITHUB_TOKEN",
			Name:   "token",
			Usage:  "Your GitHub token",
		},
		cli.StringFlag{
			Name:  "org",
			Usage: "GitHub organization to scan.",
		},
	}
	app.Action = func(c *cli.Context) error {
		token := c.String("token")
		org := c.String("org")
		if token == "" || org == "" {
			return cli.ShowAppHelp(c)
		}
		active := true
		s := spin.New()
		s.Set(`⦾⦿`)
		go func() {
			for {
				if !active {
					return
				}
				fmt.Printf(
					"\r  \033[36m%s Gathering data for '%s'...\033[m",
					s.Next(),
					org,
				)
				time.Sleep(100 * time.Millisecond)
			}
		}()
		allStats, err := stats.Gather(token, org)
		active = false
		fmt.Printf("\r")
		if err != nil {
			return err
		}
		printHighlights(allStats)
		return nil
	}
	app.Run(os.Args)
}

func printHighlights(s stats.Stats) {
	var commits, adds, dels string
	for name, stat := range s.Stats {
		if stat.Commits > s.Stats[commits].Commits {
			commits = name
		}
		if stat.Additions > s.Stats[adds].Additions {
			adds = name
		}
		if stat.Deletions > s.Stats[dels].Deletions {
			dels = name
		}
	}
	emoji.Printf(
		":trophy: Commit Champion is %s with %d commits!\n",
		commits,
		s.Stats[commits].Commits,
	)
	emoji.Printf(
		":trophy: Lines Added Champion is %s with %d lines added!\n",
		adds,
		s.Stats[adds].Additions,
	)
	emoji.Printf(
		":trophy: Housekeeper Champion is %s with %d lines removed!\n",
		dels,
		s.Stats[dels].Deletions,
	)
}
