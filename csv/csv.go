package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strconv"

	"github.com/caarlos0/org-stats/orgstats"
)

func Write(w io.Writer, s orgstats.Stats, includeReviews bool) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	headers := []string{"login", "commits", "lines-added", "lines-removed"}
	if includeReviews {
		headers = append(headers, "reviews")
	}
	if err := cw.Write(headers); err != nil {
		return fmt.Errorf("failed to write csv: %w", err)
	}

	logins := s.Logins()
	sort.Strings(logins)

	for _, login := range logins {
		stat := s.For(login)
		record := []string{
			login,
			strconv.Itoa(stat.Commits),
			strconv.Itoa(stat.Additions),
			strconv.Itoa(stat.Deletions),
		}
		if includeReviews {
			record = append(record, strconv.Itoa(stat.Reviews))
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("failed to write csv: %w", err)
		}
	}

	return cw.Error()
}
