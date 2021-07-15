package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func TestBuildBlacklists(t *testing.T) {
	users, repos := buildBlacklists([]string{
		"user:foo",
		"repo:bar",
		"something else",
		"yada:yada",
	})

	is := is.New(t)
	is.Equal(users, []string{"foo", "something else", "yada:yada"})
	is.Equal(repos, []string{"bar", "something else", "yada:yada"})
}
