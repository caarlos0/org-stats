package highlights

import (
	"fmt"
	"io"

	"github.com/caarlos0/org-stats/orgstats"
	"github.com/charmbracelet/lipgloss"
)

func Write(w io.Writer, s orgstats.Stats, top int, includeReviews bool) error {
	data := []statHighlight{
		{
			stats:  orgstats.Sort(s, orgstats.ExtractCommits),
			trophy: "Commits",
			kind:   "commits",
		}, {
			stats:  orgstats.Sort(s, orgstats.ExtractAdditions),
			trophy: "Lines Added",
			kind:   "lines added",
		}, {
			stats:  orgstats.Sort(s, orgstats.ExtractDeletions),
			trophy: "Housekeeper",
			kind:   "lines removed",
		},
	}

	if includeReviews {
		data = append(data, statHighlight{
			stats:  orgstats.Sort(s, orgstats.Reviews),
			trophy: "Pull Requests Reviewed",
			kind:   "pull requests reviewed",
		})
	}

	var headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.AdaptiveColor{
			Dark:  "#BD7EFC",
			Light: "#7D56F4",
		}).
		MarginTop(1).
		Underline(true)

	var bodyStyle = lipgloss.NewStyle().
		MarginLeft(2)

	// TODO: handle no results for a given topic
	for _, d := range data {
		if _, err := fmt.Fprintln(
			w,
			headerStyle.Render(d.trophy+" champions are:"),
		); err != nil {
			return err
		}
		j := top
		if len(d.stats) < j {
			j = len(d.stats)
		}
		for i := 0; i < j; i++ {
			if _, err := fmt.Fprintln(w,
				bodyStyle.Render(
					fmt.Sprintf(
						"%s %s with %d %s!",
						emojiForPos(i),
						d.stats[i].Key,
						d.stats[i].Value,
						d.kind,
					),
				),
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func emojiForPos(pos int) string {
	emojis := []string{"\U0001f3c6", "\U0001f948", "\U0001f949"}
	if pos < len(emojis) {
		return emojis[pos]
	}
	return " "
}

type statHighlight struct {
	stats  []orgstats.StatPair
	trophy string
	kind   string
}
