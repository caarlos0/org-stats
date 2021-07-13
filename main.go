package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/duration"
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
			Usage: "Blacklist repos and/or users. E.g. 'foo' blacklists both the 'foo' user and 'foo' repo, 'user:foo' blacklists only the user and `repo:foo` only the repo.",
		},
		cli.IntFlag{
			Name:  "top",
			Usage: "How many users to show",
			Value: 3,
		},
		cli.StringFlag{
			Name:  "github-url",
			Usage: "Custom GitHub URL (for GitHub Enterprise for example)",
		},
		cli.StringFlag{
			Name:  "since",
			Usage: "Time to look back to gather info (0s means everything). Examples: e.g. 2y, 1mo, 1w, 10d, 20h, 15m, 25s, 10ms, etc. Note that GitHub data comes summarized by week, so this is not",
			Value: "0s",
		},
	}
	app.Action = func(c *cli.Context) error {
		token := c.String("token")
		org := c.String("org")
		blacklist := c.StringSlice("blacklist")
		top := c.Int("top")
		if token == "" {
			return cli.NewExitError("missing github api token", 1)
		}
		if org == "" {
			return cli.NewExitError("missing organization name", 1)
		}
		spin := spin.New("  \033[36m%s Gathering data for '" + org + "'...\033[m")
		spin.Start()

		since, err := duration.Parse(c.String("since"))
		if err != nil {
			return cli.NewExitError("invalid --since duration", 1)
		}

		var userBlacklist []string
		var repoBlacklist []string

		for _, b := range blacklist {
			if strings.HasPrefix(b, "user:") {
				userBlacklist = append(userBlacklist, strings.TrimPrefix(b, "user:"))
			} else if strings.HasPrefix(b, "repo:") {
				repoBlacklist = append(repoBlacklist, strings.TrimPrefix(b, "repo:"))
			} else {
				userBlacklist = append(userBlacklist, b)
				repoBlacklist = append(repoBlacklist, b)
			}
		}

		allStats, err := orgstats.Gather(
			token,
			org,
			userBlacklist,
			repoBlacklist,
			c.String("github-url"),
			time.Now().UTC().Add(-1*time.Duration(since)),
		)
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
		j := top
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
	emojis := []string{"\U0001f3c6", "\U0001f948", "\U0001f949"}
	if pos < len(emojis) {
		return emojis[pos]
	}
	return " "
}
