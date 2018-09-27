package orgstats

import (
	"context"
	"net/url"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func newClient(ctx context.Context, token, baseURL string) (*github.Client, error) {
	var ts = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	var client = github.NewClient(oauth2.NewClient(ctx, ts))

	if baseURL != "" {
		api, err := url.Parse(baseURL)
		if err != nil {
			return client, err
		}
		client.BaseURL = api
	}

	return client, nil
}
