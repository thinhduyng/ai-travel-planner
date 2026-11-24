// Contract: docs/contract/create-itinerary.md. Test names cite the § they check.
package contract

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/thinhduyng/ai-travel-planner/api/internal/planner"
)

func TestSuccess(t *testing.T) {
	base := start(t, nil)

	t.Run("§2.1 returns 201 with Location /v1/itineraries/{id}", func(t *testing.T) {
		r := post(t, base, validBody())
		expectStatus(t, r, http.StatusCreated)
		if got, want := r.Header.Get("Location"), fmt.Sprintf("/v1/itineraries/%v", r.Body["id"]); got != want {
			t.Fatalf("Location = %q, want %q", got, want)
		}
	})

	t.Run("§2.2 body has exactly the documented fields", func(t *testing.T) {
		r := post(t, base, validBody())
		var keys []string
		for k := range r.Body {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		want := []string{"budget", "createdAt", "days", "destination", "endDate", "id", "interests", "startDate", "travelers"}
		if !reflect.DeepEqual(keys, want) {
			t.Fatalf("keys = %v, want %v", keys, want)
		}
	})

	t.Run("§2.3 id is a UUID v4", func(t *testing.T) {
		r := post(t, base, validBody())
		if id, _ := r.Body["id"].(string); !uuidV4.MatchString(id) {
			t.Fatalf("id = %q is not a UUID v4", id)
		}
	})

	t.Run("§2.4 echoes trimmed destination and the given fields", func(t *testing.T) {
		r := post(t, base, with("destination", "  Hoi An  "))
		expectStatus(t, r, http.StatusCreated)
		if r.Body["destination"] != "Hoi An" {
			t.Fatalf("destination = %q", r.Body["destination"])
		}
		if r.Body["startDate"] != "2026-12-20" || r.Body["endDate"] != "2026-12-22" || r.Body["travelers"] != float64(2) {
			t.Fatalf("dates/travelers not echoed: %s", r.Raw)
		}
		wantBudget := map[string]any{"amount": float64(5000000), "currency": "VND"}
		if !reflect.DeepEqual(r.Body["budget"], wantBudget) {
			t.Fatalf("budget = %v", r.Body["budget"])
		}
		if !reflect.DeepEqual(r.Body["interests"], []any{"food", "nature"}) {
			t.Fatalf("interests = %v", r.Body["interests"])
		}
	})

	t.Run("§2.4 applies defaults: travelers=1, interests=[], budget=null", func(t *testing.T) {
		r := post(t, base, map[string]any{"destination": "Hue", "startDate": "2026-12-20", "endDate": "2026-12-20"})
		expectStatus(t, r, http.StatusCreated)
		if r.Body["travelers"] != float64(1) {
			t.Fatalf("travelers = %v, want 1", r.Body["travelers"])
		}
		if !reflect.DeepEqual(r.Body["interests"], []any{}) {
			t.Fatalf("interests = %v, want []", r.Body["interests"])
		}
		if v, ok := r.Body["budget"]; !ok || v != nil {
			t.Fatalf("budget = %v (present=%v), want null", v, ok)
		}
	})

	t.Run("§2.5 one day per trip day, numbered from 1, with consecutive dates", func(t *testing.T) {
		d := days(t, post(t, base, validBody()))
		want := [][2]any{{float64(1), "2026-12-20"}, {float64(2), "2026-12-21"}, {float64(3), "2026-12-22"}}
		if len(d) != len(want) {
			t.Fatalf("got %d days, want 3", len(d))
		}
		for i, w := range want {
			if d[i]["day"] != w[0] || d[i]["date"] != w[1] {
				t.Fatalf("day %d = %v/%v, want %v/%v", i, d[i]["day"], d[i]["date"], w[0], w[1])
			}
		}
	})

	t.Run("§2.5 startDate == endDate is a one-day trip", func(t *testing.T) {
		r := post(t, base, with("endDate", "2026-12-20"))
		expectStatus(t, r, http.StatusCreated)
		if n := len(days(t, r)); n != 1 {
			t.Fatalf("got %d days, want 1", n)
		}
	})

	t.Run("§2.5 dates cross month and year boundaries correctly", func(t *testing.T) {
		d := days(t, post(t, base, with("startDate", "2026-12-31", "endDate", "2027-01-01")))
		if len(d) != 2 || d[0]["date"] != "2026-12-31" || d[1]["date"] != "2027-01-01" {
			t.Fatalf("days = %v", d)
		}
	})

	t.Run("§2.6 each day has >=1 activity with HH:MM time and non-empty title", func(t *testing.T) {
		for _, day := range days(t, post(t, base, validBody())) {
			acts, _ := day["activities"].([]any)
			if len(acts) == 0 {
				t.Fatalf("day %v has no activities", day["day"])
			}
			for _, a := range acts {
				act := a.(map[string]any)
				if tm, _ := act["time"].(string); !hhmm.MatchString(tm) {
					t.Fatalf("time = %v", act["time"])
				}
				if title, _ := act["title"].(string); title == "" {
					t.Fatalf("empty title: %v", act)
				}
				if d, ok := act["description"]; ok {
					if _, isStr := d.(string); !isStr {
						t.Fatalf("description is not a string: %v", d)
					}
				}
			}
		}
	})

	t.Run("§2.7 activities within a day are sorted by time", func(t *testing.T) {
		for _, day := range days(t, post(t, base, validBody())) {
			var times []string
			for _, a := range day["activities"].([]any) {
				times = append(times, a.(map[string]any)["time"].(string))
			}
			if !sort.StringsAreSorted(times) {
				t.Fatalf("day %v times not sorted: %v", day["day"], times)
			}
		}
	})

	t.Run("§2.8 createdAt is an ISO 8601 UTC timestamp", func(t *testing.T) {
		r := post(t, base, validBody())
		if s, _ := r.Body["createdAt"].(string); !isoUTC.MatchString(s) {
			t.Fatalf("createdAt = %q", s)
		}
	})
}

func TestValidation(t *testing.T) {
	base := start(t, nil)

	cases := []struct {
		name  string
		body  any
		field string
	}{
		{"§1.1 destination is required", without("destination"), "destination"},
		{"§1.1 destination of only whitespace is rejected", with("destination", "   "), "destination"},
		{"§1.1 destination longer than 100 characters is rejected", with("destination", strings.Repeat("x", 101)), "destination"},
		{"§1.2 dates must be YYYY-MM-DD", with("startDate", "20/12/2026"), "startDate"},
		{"§1.2 dates must be real calendar dates", with("startDate", "2026-02-30", "endDate", "2026-03-01"), "startDate"},
		{"§1.3 endDate before startDate is rejected", with("startDate", "2026-12-22", "endDate", "2026-12-20"), "endDate"},
		{"§1.4 a 15-day trip is rejected", with("startDate", "2026-12-01", "endDate", "2026-12-15"), "endDate"},
		{"§1.5 travelers=0 is rejected", with("travelers", 0), "travelers"},
		{"§1.5 travelers=21 is rejected", with("travelers", 21), "travelers"},
		{"§1.5 travelers=2.5 is rejected", with("travelers", 2.5), "travelers"},
		{`§1.5 travelers="2" is rejected`, with("travelers", "2"), "travelers"},
		{"§1.6 budget amount 0 is rejected", with("budget", map[string]any{"amount": 0, "currency": "VND"}), "budget.amount"},
		{"§1.6 budget amount 10.5 is rejected", with("budget", map[string]any{"amount": 10.5, "currency": "VND"}), "budget.amount"},
		{"§1.6 lower-case currency is rejected", with("budget", map[string]any{"amount": 100, "currency": "vnd"}), "budget.currency"},
		{"§1.6 4-letter currency is rejected", with("budget", map[string]any{"amount": 100, "currency": "DONG"}), "budget.currency"},
		{"§1.7 more than 10 interests is rejected", with("interests", []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}), "interests"},
		{"§1.7 an interest longer than 30 characters is rejected", with("interests", []string{strings.Repeat("x", 31)}), "interests.0"},
		{"§1.8 unknown top-level fields are rejected", with("hotel", "Ritz"), "hotel"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expectValidationError(t, post(t, base, c.body), c.field)
		})
	}

	t.Run("§1.4 a 14-day trip is accepted", func(t *testing.T) {
		r := post(t, base, with("startDate", "2026-12-01", "endDate", "2026-12-14"))
		expectStatus(t, r, http.StatusCreated)
		if n := len(days(t, r)); n != 14 {
			t.Fatalf("got %d days, want 14", n)
		}
	})

	t.Run("§3.2 details has one entry per broken field", func(t *testing.T) {
		fields := expectError(t, post(t, base, with("destination", "", "travelers", 0)), http.StatusBadRequest, "VALIDATION_ERROR")
		sort.Strings(fields)
		if !reflect.DeepEqual(fields, []string{"destination", "travelers"}) {
			t.Fatalf("fields = %v, want [destination travelers]", fields)
		}
	})
}

func TestErrors(t *testing.T) {
	base := start(t, nil)

	t.Run("§3.3 malformed JSON gives 400 MALFORMED_JSON", func(t *testing.T) {
		expectError(t, post(t, base, `{"destination": `), http.StatusBadRequest, "MALFORMED_JSON")
	})

	t.Run("§3.4 non-JSON Content-Type gives 415 UNSUPPORTED_MEDIA_TYPE", func(t *testing.T) {
		r := post(t, base, "destination=Hue", "application/x-www-form-urlencoded")
		expectError(t, r, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE")
	})
}

func TestPlannerFailures(t *testing.T) {
	leak := regexp.MustCompile(`(?i)anthropic|sk-ant|overloaded|goroutine|\.go:\d+`)

	t.Run("§3.5 planner error gives 502 PLANNER_UNAVAILABLE, §3.6 without leaking it", func(t *testing.T) {
		base := start(t, planner.Func(func(context.Context, planner.Request) ([][]planner.Activity, error) {
			return nil, errors.New("anthropic: 529 overloaded sk-ant-secret")
		}))
		r := post(t, base, validBody())
		expectError(t, r, http.StatusBadGateway, "PLANNER_UNAVAILABLE")
		if leak.MatchString(r.Raw) {
			t.Fatalf("response leaks provider details: %s", r.Raw)
		}
	})

	t.Run("§3.5 wrong number of days gives 502 and no partial itinerary", func(t *testing.T) {
		base := start(t, planner.Func(func(context.Context, planner.Request) ([][]planner.Activity, error) {
			return [][]planner.Activity{{{Time: "09:00", Title: "Only one"}}}, nil
		}))
		r := post(t, base, validBody())
		expectError(t, r, http.StatusBadGateway, "PLANNER_UNAVAILABLE")
		if _, ok := r.Body["days"]; ok {
			t.Fatalf("partial itinerary returned: %s", r.Raw)
		}
	})

	t.Run("§3.5 invalid activity time gives 502", func(t *testing.T) {
		base := start(t, planner.Func(func(_ context.Context, req planner.Request) ([][]planner.Activity, error) {
			out := make([][]planner.Activity, len(req.Dates))
			for i := range out {
				out[i] = []planner.Activity{{Time: "9am", Title: "Breakfast"}}
			}
			return out, nil
		}))
		expectError(t, post(t, base, validBody()), http.StatusBadGateway, "PLANNER_UNAVAILABLE")
	})
}
