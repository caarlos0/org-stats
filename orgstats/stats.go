package orgstats

import (
	"context"
	"strings"
	"time"

	"github.com/google/go-github/github"
)

// Stat represents an user adds, rms and commits count
type Stat struct {
	Additions, Deletions, Commits int
}

// Stats contains the user->Stat mapping
type Stats map[string]Stat

// NewStats return a new Stats map
func NewStats() Stats {
	return make(map[string]Stat)
}

// Gather a given organization's stats
func Gather(token, org string, blacklist []string, url string, weeks int32) (Stats, error) {
	var ctx = context.Background()
	var allStats = NewStats()
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
			allStats.add(cs, weeks)
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

func (s Stats) add(cs *github.ContributorStats, weeks int32) {
	if cs.Author == nil {
		return
	}
	stat := s[*cs.Author.Login]
	var adds int
	var rms int
	var commits int
	t1 := time.Now()
	lastWeeksDuration := time.Hour * 24 * 7 * time.Duration(weeks)
	for _, week := range cs.Weeks {
		diff := t1.Sub(week.Week.Time)
		if weeks == -1 || diff <= lastWeeksDuration {
			adds += *week.Additions
			rms += *week.Deletions
			commits += *week.Commits
		}
	}
	stat.Additions += adds
	stat.Deletions += rms
	stat.Commits += commits
	s[*cs.Author.Login] = stat
}

func repos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	var opt = &github.RepositoryListByOrgOptions{
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
		if _, ok := err.(*github.RateLimitError); ok {
			time.Sleep(time.Duration(15) * time.Second)
			return getStats(ctx, client, org, repo)
		}
		if _, ok := err.(*github.AcceptedError); ok {
			return getStats(ctx, client, org, repo)
		}
	}
	return stats, err
}
