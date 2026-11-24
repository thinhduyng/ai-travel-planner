package planner

import (
	"context"
	"fmt"
)

// Fake is a deterministic planner for tests and CI: no network, no API key.
type Fake struct{}

// Plan implements Planner. Activities come back out of order on purpose:
// sorting them is the API's job (contract §2.7).
func (Fake) Plan(_ context.Context, req Request) ([][]Activity, error) {
	days := make([][]Activity, len(req.Dates))
	for i, date := range req.Dates {
		desc := "Morning of " + date
		days[i] = []Activity{
			{Time: "14:00", Title: fmt.Sprintf("Explore %s (day %d)", req.Destination, i+1)},
			{Time: "09:00", Title: "Breakfast at a local spot", Description: &desc},
		}
	}
	return days, nil
}
