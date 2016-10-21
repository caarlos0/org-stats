package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"gopkg.in/kyokomi/emoji.v1"
)

type Stat struct {
	Additions, Deletions, Commits int
}

type Stats struct {
	stats map[string]Stat
}

func NewStats() Stats {
	return Stats{make(map[string]Stat)}
}

func (s Stats) Add(cs *github.ContributorStats) {
	stat := s.stats[*cs.Author.Login]
	var adds int
	var rms int
	var commits int
	for _, week := range cs.Weeks {
		adds += *week.Additions
		rms += *week.Deletions
		commits += *week.Commits
	}
	stat.Additions += adds
	stat.Deletions += rms
	stat.Commits += commits
	s.stats[*cs.Author.Login] = stat
}

func (s Stats) Print() {
	for name, stat := range s.stats {
		fmt.Printf(
			"%s,%d,%d,%d\n",
			name,
			stat.Commits,
			stat.Additions,
			stat.Deletions,
		)
	}
}

func (s Stats) PrintHighlights() {
	commits := ""
	adds := ""
	dels := ""
	for name, stat := range s.stats {
		if stat.Commits > s.stats[commits].Commits {
			commits = name
		}
		if stat.Additions > s.stats[adds].Additions {
			adds = name
		}
		if stat.Deletions > s.stats[dels].Deletions {
			dels = name
		}
	}
	emoji.Printf(
		":trophy: Commit Champion is %s with %d commits!\n",
		commits,
		s.stats[commits].Commits,
	)
	emoji.Printf(
		":trophy: Lines Added Champion is %s with %d lines added!\n",
		adds,
		s.stats[adds].Additions,
	)
	emoji.Printf(
		":trophy: Housekeeper Champion is %s with %d lines removed!\n",
		dels,
		s.stats[dels].Deletions,
	)
}

func repos(org string, client *github.Client) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}

	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(org, opt)
		if err != nil {
			return allRepos, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}
	return allRepos, nil
}

func getStats(org, repo string, client *github.Client) []*github.ContributorStats {
	stats, res, err := client.Repositories.ListContributorsStats(org, repo)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			time.Sleep(time.Duration(15) * time.Second)
			return getStats(org, repo, client)
		}
		fmt.Println("Got error", res.Status)
	}
	return stats
}

func main() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	org := os.Args[1]
	s := spinner.New([]string{"⦾", "⦿"}, 80*time.Millisecond)
	s.Suffix = " Gathering data for '" + org + "'..."
	s.Start()
	allRepos, err := repos(org, client)
	if err != nil {
		log.Fatalln(err)
	}

	allStats := NewStats()
	for _, repo := range allRepos {
		for _, cs := range getStats(org, *repo.Name, client) {
			allStats.Add(cs)
		}
	}
	fmt.Printf("\n\n\n")
	s.Stop()
	allStats.PrintHighlights()
}
