package cmd

import (
	"context"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

func newClient(ctx context.Context, token, baseURL string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	if baseURL == "" {
		if token == "" {
			return github.NewClient(nil), nil
		} else {
			return github.NewClient(oauth2.NewClient(ctx, ts)), nil
		}
	}

	return github.NewEnterpriseClient(baseURL, "", oauth2.NewClient(ctx, ts))
}
