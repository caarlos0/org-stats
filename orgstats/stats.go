package orgstats

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/v37/github"
)

// Stat represents an user adds, rms and commits count
type Stat struct {
	Additions, Deletions, Commits int
}

// Stats contains the user->Stat mapping
type Stats struct {
	data  map[string]Stat
	since time.Time
}

// NewStats return a new Stats map
func NewStats(since time.Time) Stats {
	return Stats{
		data:  make(map[string]Stat),
		since: since,
	}
}

// Gather a given organization's stats
func Gather(token, org string, blacklist []string, url string, since time.Time) (Stats, error) {
	ctx := context.Background()
	allStats := NewStats(since)
	client, err := newClient(ctx, token, url)
	if err != nil {
		return allStats, err
	}

	allRepos, err := repos(ctx, client, org)
	if err != nil {
		return allStats, err
	}

	for _, repo := range allRepos {
		if isBlacklisted(blacklist, repo.GetName()) {
			continue
		}
		stats, serr := getStats(ctx, client, org, *repo.Name)
		if serr != nil {
			return allStats, serr
		}
		for _, cs := range stats {
			if isBlacklisted(blacklist, cs.Author.GetLogin()) {
				continue
			}
			allStats.add(cs)
		}
	}
	return allStats, err
}

func isBlacklisted(blacklist []string, s string) bool {
	for _, b := range blacklist {
		if strings.EqualFold(s, b) {
			return true
		}
	}
	return false
}

func (s Stats) add(cs *github.ContributorStats) {
	if cs.Author == nil {
		return
	}
	stat := s.data[*cs.Author.Login]
	var adds int
	var rms int
	var commits int
	for _, week := range cs.Weeks {
		if !s.since.IsZero() && week.Week.Time.UTC().Before(s.since) {
			continue
		}
		adds += *week.Additions
		rms += *week.Deletions
		commits += *week.Commits
	}
	stat.Additions += adds
	stat.Deletions += rms
	stat.Commits += commits
	s.data[*cs.Author.Login] = stat
}

func repos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
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

func getStats(ctx context.Context, client *github.Client, org, repo string) ([]*github.ContributorStats, error) {
	stats, _, err := client.Repositories.ListContributorsStats(ctx, org, repo)
	if err != nil {
		if rateErr, ok := err.(*github.RateLimitError); ok {
			time.Sleep(time.Now().UTC().Sub(rateErr.Rate.Reset.Time.UTC()))
			return getStats(ctx, client, org, repo)
		}
		if _, ok := err.(*github.AcceptedError); ok {
			return getStats(ctx, client, org, repo)
		}
	}
	return stats, err
}
