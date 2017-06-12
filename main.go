package main

import (
	"fmt"
	"os"

	orgstats "github.com/caarlos0/org-stats/orgstats"
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
			Name:  "org, o",
			Usage: "GitHub organization to scan",
		},
		cli.StringSliceFlag{
			Name:  "blacklist, b",
			Usage: "Blacklist repos and/or users",
		},
		cli.IntFlag{
			Name:  "top",
			Usage: "How many users to show",
			Value: 3,
		},
	}
	app.Action = func(c *cli.Context) error {
		var token = c.String("token")
		var org = c.String("org")
		var blacklist = c.StringSlice("blacklist")
		var top = c.Int("top")
		if token == "" {
			return cli.NewExitError("missing github api token", 1)
		}
		if org == "" {
			return cli.NewExitError("missing organization name", 1)
		}
		var spin = spin.New("  \033[36m%s Gathering data for '" + org + "'...\033[m")
		spin.Start()
		allStats, err := orgstats.Gather(token, org, blacklist)
		spin.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		printHighlights(allStats, top)
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func printHighlights(s orgstats.Stats, top int) {
	data := []struct {
		stats  []orgstats.StatPair
		trophy string
		kind   string
	}{
		{
			stats:  orgstats.Sort(s, orgstats.ExtractCommits),
			trophy: "Commit",
			kind:   "commits",
		}, {
			stats:  orgstats.Sort(s, orgstats.ExtractAdditions),
			trophy: "Lines Added",
			kind:   "lines added",
		}, {
			stats:  orgstats.Sort(s, orgstats.ExtractDeletions),
			trophy: "Housekeeper",
			kind:   "lines removed",
		},
	}
	for _, d := range data {
		fmt.Printf("\033[1m%s champions are:\033[0m\n", d.trophy)
		var j = top
		if len(d.stats) < j {
			j = len(d.stats)
		}
		for i := 0; i < j; i++ {
			fmt.Printf(
				"%s %s with %d %s!\n",
				emojiForPos(i),
				d.stats[i].Key,
				d.stats[i].Value,
				d.kind,
			)
		}
		fmt.Printf("\n")
	}
}

func emojiForPos(pos int) string {
	var emojis = []string{"\U0001f3c6", "\U0001f948", "\U0001f949"}
	if pos < len(emojis) {
		return emojis[pos]
	}
	return " "
}
