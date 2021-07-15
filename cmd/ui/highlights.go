package ui

import (
	"bytes"

	"github.com/caarlos0/org-stats/highlights"
	"github.com/caarlos0/org-stats/orgstats"
	tea "github.com/charmbracelet/bubbletea"
)

func NewHighlightsModel(stats orgstats.Stats, top int, includeReviews bool) HighlightsModel {
	return HighlightsModel{
		stats:          stats,
		top:            top,
		includeReviews: includeReviews,
	}
}

type HighlightsModel struct {
	stats          orgstats.Stats
	top            int
	includeReviews bool
}

func (m HighlightsModel) Init() tea.Cmd {
	return tea.Quit
}

func (m HighlightsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil // noop
}

func (m HighlightsModel) View() string {
	var b bytes.Buffer
	_ = highlights.Write(&b, m.stats, m.top, m.includeReviews)
	return b.String()
}
