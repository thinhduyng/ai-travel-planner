package itinerary

import (
	"regexp"
	"time"
)

const dateLayout = "2006-01-02"

var dateRE = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

// parseDate accepts only real calendar dates in YYYY-MM-DD form (rejects 2026-02-30).
func parseDate(s string) (time.Time, bool) {
	if !dateRE.MatchString(s) {
		return time.Time{}, false
	}
	t, err := time.Parse(dateLayout, s)
	return t, err == nil
}

// tripDates returns every date from start to end inclusive.
func tripDates(start, end time.Time) []string {
	n := tripLength(start, end)
	if n < 1 {
		return nil
	}
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = start.AddDate(0, 0, i).Format(dateLayout)
	}
	return out
}

// tripLength is inclusive: same start and end is 1 day.
func tripLength(start, end time.Time) int {
	return int(end.Sub(start).Hours()/24) + 1
}
