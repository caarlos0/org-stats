package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/duration"
	"github.com/caarlos0/org-stats/cmd/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	token          string
	organization   string
	githubURL      string
	since          string
	csvPath        string
	blacklist      []string
	top            int
	includeReviews bool
	excludeForks   bool
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&token, "token", "", "github api token (default $GITHUB_TOKEN)")
	_ = rootCmd.MarkFlagRequired(token)

	rootCmd.Flags().StringVarP(&organization, "org", "o", "", "github organization to scan")
	_ = rootCmd.MarkFlagRequired("org")

	rootCmd.Flags().StringSliceVarP(&blacklist, "blacklist", "b", []string{}, "blacklist repos and/or users")
	rootCmd.Flags().IntVar(&top, "top", 3, "how many users to show")
	rootCmd.Flags().StringVar(&githubURL, "github-url", "", "custom github base url (if using github enterprise)")
	rootCmd.Flags().StringVar(&since, "since", "0s", "time to look back to gather info (0s means everything)")
	rootCmd.Flags().BoolVar(&includeReviews, "include-reviews", false, "include pull request reviews in the stats")
	rootCmd.Flags().BoolVar(&excludeForks, "exclude-forks", false, "exclude forked repositories from the stats")
	rootCmd.Flags().StringVar(&csvPath, "csv-path", "", "path to write a csv file with all data collected")

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.AddCommand(versionCmd, docsCmd, manCmd)
}

var rootCmd = &cobra.Command{
	Use:   "org-stats",
	Short: "Get the contributor stats summary from all repos of any given organization",
	Long: `org-stats can be used to get an overall sense of your org's contributors.

It uses the GitHub API to grab the repositories in the given organization.
Then, iterating one by one, it gets statistics of lines added, removed and number of commits of contributors.
After that, if opted in, it does several searches to get the number of pull requests reviewed by each of the previously find contributors.
Finally, it prints a rank by each category.


Important notes:
* GitHub's API rate limits for unauthenticated requests have been lowered significantly in the recent past. Using the ` + "`--token`" + ` option for compiling stats will speed up gathering of data considerably, since for authenticated requests it will be less likely that rate-limiting timelocks have to be awaited.
* The ` + "`--since`" + ` filter does not work "that well" because GitHub summarizes thedata by week, so the data is not as granular as it should be.
* The ` + "`--include-reviews`" + ` only grabs reviews from users that had contributions on the previous step.
* In the ` + "`--blacklist`" + ` option, 'foo' blacklists both the 'foo' user and 'foo' repo, while 'user:foo' blacklists only the user and 'repo:foo' only the repository.
* The ` + "`--since`" + ` option accepts all the regular time. Accepts any duration Go standard library accepts, plus a few more: 1y (365d), 1mo (30d), 1w (7d) and 1d (24h).
* The ` + "`--token`" + ` token permissions need to include 'repo - Full control of private repositories'. Required only if you need to fetch data from private repositories in your organization.
}`,
	PreRun: func(*cobra.Command, []string) {
		if token == "" {
			token = os.Getenv("GITHUB_TOKEN")
		}
	},
	RunE: func(*cobra.Command, []string) error {
		ctx := context.Background()
		client, err := newClient(ctx, token, githubURL)
		if err != nil {
			return err
		}

		sinceD, err := duration.Parse(since)
		if err != nil {
			return fmt.Errorf("invalid --since duration: '%s'", since)
		}

		userBlacklist, repoBlacklist := buildBlacklists(blacklist)

		csv := io.Discard
		if csvPath != "" {
			if err := os.MkdirAll(filepath.Dir(csvPath), 0o755); err != nil {
				return fmt.Errorf("failed to create csv file: %w", err)
			}
			f, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return fmt.Errorf("failed to create csv file: %w", err)
			}
			defer f.Close()
			csv = f
		}

		f, err := tea.LogToFile(filepath.Join(os.TempDir(), "org-stats.log"), "org-stats")
		if err != nil {
			return err
		}
		defer f.Close()

		sinceT := time.Time{}
		if sinceD > 0 {
			sinceT = time.Now().UTC().Add(-1 * time.Duration(sinceD))
		}

		p := tea.NewProgram(ui.NewInitialModel(
			client,
			organization,
			userBlacklist,
			repoBlacklist,
			sinceT,
			top,
			includeReviews,
			excludeForks,
			csv,
		))
		_, err = p.Run()
		return err
	},
}
