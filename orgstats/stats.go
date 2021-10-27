package orgstats

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/go-github/v39/github"
)

// Stat represents an user adds, rms and commits count
type Stat struct {
	Additions, Deletions, Commits, Reviews int
}

// Stats contains the user->Stat mapping
type Stats struct {
	data  map[string]Stat
	since time.Time
}

func (s Stats) Logins() []string {
	logins := make([]string, 0, len(s.data))
	for login := range s.data {
		logins = append(logins, login)
	}
	return logins
}

func (s Stats) For(login string) Stat {
	return s.data[login]
}

// NewStats return a new Stats map
func NewStats(since time.Time) Stats {
	return Stats{
		data:  make(map[string]Stat),
		since: since,
	}
}

// Gather a given organization's stats
func Gather(
	ctx context.Context,
	client *github.Client,
	org string,
	userBlacklist, repoBlacklist []string,
	since time.Time,
	includeReviewStats bool,
) (Stats, error) {

	allStats := NewStats(since)
	if err := gatherLineStats(
		ctx,
		client,
		org,
		userBlacklist,
		repoBlacklist,
		&allStats,
	); err != nil {
		return Stats{}, err
	}

	if !includeReviewStats {
		return allStats, nil
	}

	for user := range allStats.data {
		if err := gatherReviewStats(
			ctx,
			client,
			org,
			user,
			userBlacklist,
			repoBlacklist,
			&allStats,
			since,
		); err != nil {
			return Stats{}, err
		}
	}

	return allStats, nil
}

func gatherReviewStats(
	ctx context.Context,
	client *github.Client,
	org, user string,
	userBlacklist, repoBlacklist []string,
	allStats *Stats,
	since time.Time,
) error {
	ts := since.Format("2006-01-02")
	// review:approved, review:changes_requested
	reviewed, err := search(ctx, client, fmt.Sprintf("user:%s is:pr reviewed-by:%s created:>%s", org, user, ts))
	if err != nil {
		return err
	}
	allStats.addReviewStats(user, reviewed)
	return nil
}

func search(
	ctx context.Context,
	client *github.Client,
	query string,
) (int, error) {
	log.Printf("searching '%s'", query)
	result, _, err := client.Search.Issues(ctx, query, &github.SearchOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	if rateErr, ok := err.(*github.RateLimitError); ok {
		handleRateLimit(rateErr)
		return search(ctx, client, query)
	}
	if _, ok := err.(*github.AcceptedError); ok {
		return search(ctx, client, query)
	}
	if err != nil {
		return 0, fmt.Errorf("failed to search: %s: %w", query, err)
	}
	return *result.Total, nil
}

func gatherLineStats(
	ctx context.Context,
	client *github.Client,
	org string,
	userBlacklist, repoBlacklist []string,
	allStats *Stats,
) error {
	allRepos, err := repos(ctx, client, org)
	if err != nil {
		return err
	}

	for _, repo := range allRepos {
		if isBlacklisted(repoBlacklist, repo.GetName()) {
			continue
		}
		stats, serr := getStats(ctx, client, org, *repo.Name)
		if serr != nil {
			return serr
		}
		for _, cs := range stats {
			if isBlacklisted(userBlacklist, cs.Author.GetLogin()) {
				continue
			}
			allStats.add(cs)
		}
	}
	return err
}

func isBlacklisted(blacklist []string, s string) bool {
	for _, b := range blacklist {
		if strings.EqualFold(s, b) {
			return true
		}
	}
	return false
}

func (s *Stats) addReviewStats(user string, reviewed int) {
	stat := s.data[user]
	stat.Reviews += reviewed
	s.data[user] = stat
}

func (s *Stats) add(cs *github.ContributorStats) {
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
	if stat.Additions+stat.Deletions+stat.Commits == 0 && !s.since.IsZero() {
		// ignore users with no activity when running with a since time
		return
	}
	s.data[*cs.Author.Login] = stat
}

func repos(ctx context.Context, client *github.Client, org string) ([]*github.Repository, error) {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if rateErr, ok := err.(*github.RateLimitError); ok {
			handleRateLimit(rateErr)
			continue
		}
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
			handleRateLimit(rateErr)
			return getStats(ctx, client, org, repo)
		}
		if _, ok := err.(*github.AcceptedError); ok {
			return getStats(ctx, client, org, repo)
		}
	}
	return stats, err
}

func handleRateLimit(err *github.RateLimitError) {
	s := err.Rate.Reset.UTC().Sub(time.Now().UTC())
	if s < 0 {
		s = 5 * time.Second
	}
	log.Printf("hit rate limit, waiting %v", s)
	time.Sleep(s)
}
