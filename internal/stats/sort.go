package stats

import "sort"

type Extract func(st Stat) int

var ExtractCommits = func(st Stat) int {
	return st.Commits
}

var ExtractAdditions = func(st Stat) int {
	return st.Additions
}

var ExtractDeletions = func(st Stat) int {
	return st.Deletions
}

func Sort(s Stats, extract Extract) []StatPair {
	var result statPairList
	for key, value := range s.Stats {
		result = append(result, StatPair{Key: key, Value: extract(value)})
	}
	sort.Sort(sort.Reverse(result))
	return result
}

type StatPair struct {
	Key   string
	Value int
}

type statPairList []StatPair

func (b statPairList) Len() int {
	return len(b)
}

func (b statPairList) Less(i, j int) bool {
	return b[i].Value < b[j].Value
}

func (b statPairList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
