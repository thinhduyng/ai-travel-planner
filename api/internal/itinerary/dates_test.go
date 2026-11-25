package itinerary

import (
	"reflect"
	"testing"
)

func TestParseDate(t *testing.T) {
	cases := map[string]bool{
		"2026-12-20": true,
		"2028-02-29": true, // leap year
		"2026-02-29": false,
		"2026-02-30": false,
		"2026-13-01": false,
		"2026-1-01":  false,
		"20/12/2026": false,
		"":           false,
	}
	for in, want := range cases {
		if _, ok := parseDate(in); ok != want {
			t.Errorf("parseDate(%q) ok = %v, want %v", in, ok, want)
		}
	}
}

func TestTripLength(t *testing.T) {
	cases := []struct {
		start, end string
		want       int
	}{
		{"2026-12-20", "2026-12-20", 1},
		{"2026-12-20", "2026-12-22", 3},
		{"2026-12-31", "2027-01-01", 2},
		{"2026-12-22", "2026-12-20", -1},
	}
	for _, c := range cases {
		s, _ := parseDate(c.start)
		e, _ := parseDate(c.end)
		if got := tripLength(s, e); got != c.want {
			t.Errorf("tripLength(%s, %s) = %d, want %d", c.start, c.end, got, c.want)
		}
	}
}

func TestTripDates(t *testing.T) {
	cases := []struct {
		start, end string
		want       []string
	}{
		{"2026-12-20", "2026-12-20", []string{"2026-12-20"}},
		{"2026-12-30", "2027-01-01", []string{"2026-12-30", "2026-12-31", "2027-01-01"}},
		{"2028-02-28", "2028-03-01", []string{"2028-02-28", "2028-02-29", "2028-03-01"}},
	}
	for _, c := range cases {
		s, _ := parseDate(c.start)
		e, _ := parseDate(c.end)
		if got := tripDates(s, e); !reflect.DeepEqual(got, c.want) {
			t.Errorf("tripDates(%s, %s) = %q, want %q", c.start, c.end, got, c.want)
		}
	}
}
