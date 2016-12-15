package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/org-stats/internal/stats"
	"github.com/caarlos0/spin"
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
		s := spin.New("  \033[36m%s Gathering data for '" + org + "'...\033[m")
		s.Set(spin.Spin10)
		s.Start()
		allStats, err := stats.Gather(token, org)
		s.Stop()
		if err != nil {
			return err
		}
		printHighlights(allStats)
		return nil
	}
	app.Run(os.Args)
}

func printHighlights(s stats.Stats) {
	data := []struct {
		stats  []stats.StatPair
		trophy string
		kind   string
	}{
		{
			stats:  stats.Sort(s, stats.ExtractCommits),
			trophy: "Commit",
			kind:   "commits",
		}, {
			stats:  stats.Sort(s, stats.ExtractAdditions),
			trophy: "Lines Added",
			kind:   "lines added",
		}, {
			stats:  stats.Sort(s, stats.ExtractDeletions),
			trophy: "Housekeeper",
			kind:   "lines removed",
		},
	}
	var emojis = []string{"\U0001f3c6", "\U0001f948", "\U0001f949"}
	for _, d := range data {
		fmt.Printf("\033[1m%s champions are:\033[0m\n", d.trophy)
		for i := 0; i < 3; i++ {
			fmt.Printf(
				"%s %s with %d %s!\n",
				emojis[i],
				d.stats[i].Key,
				d.stats[i].Value,
				d.kind,
			)
		}
		fmt.Printf("\n")
	}
}
