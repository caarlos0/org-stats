package cmd

import "strings"

func buildBlacklists(blacklist []string) ([]string, []string) {
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
	return userBlacklist, repoBlacklist
}
