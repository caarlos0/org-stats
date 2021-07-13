package orgstats

import "sort"

// Extract is a function that converts a multiple stat into a single stat
type Extract func(st Stat) int

// ExtractCommits extract the commit section of the given stat
var ExtractCommits = func(st Stat) int {
	return st.Commits
}

// ExtractAdditions extract the adds section of the given stat
var ExtractAdditions = func(st Stat) int {
	return st.Additions
}

// ExtractDeletions extract the rms section of the given stat
var ExtractDeletions = func(st Stat) int {
	return st.Deletions
}

func Sort(s Stats, extract Extract) []StatPair {
	var result []StatPair
	for key, value := range s.data {
		result = append(result, StatPair{Key: key, Value: extract(value)})
	}
	sort.Slice(result, func(i int, j int) bool {
		return result[i].Value > result[j].Value
	})
	return result
}

type StatPair struct {
	Key   string
	Value int
}
