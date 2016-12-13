package stats

import (
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Stat struct {
	Additions, Deletions, Commits int
}

type Stats struct {
	Stats map[string]Stat
}

func NewStats() Stats {
	return Stats{make(map[string]Stat)}
}

func (s Stats) Add(cs *github.ContributorStats) {
	stat := s.Stats[*cs.Author.Login]
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
	s.Stats[*cs.Author.Login] = stat
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

func getStats(org, repo string, client *github.Client) ([]*github.ContributorStats, error) {
	stats, _, err := client.Repositories.ListContributorsStats(org, repo)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			time.Sleep(time.Duration(15) * time.Second)
			return getStats(org, repo, client)
		}
		if _, ok := err.(*github.AcceptedError); ok {
			return getStats(org, repo, client)
		}
	}
	return stats, err
}

func Gather(token, org string) (Stats, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	allStats := NewStats()

	allRepos, err := repos(org, client)
	if err != nil {
		return allStats, err
	}

	for _, repo := range allRepos {
		stats, err := getStats(org, *repo.Name, client)
		if err != nil {
			return allStats, err
		}
		for _, cs := range stats {
			allStats.Add(cs)
		}
	}
	return allStats, err
}
