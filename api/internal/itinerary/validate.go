package itinerary

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/thinhduyng/ai-travel-planner/api/internal/api"
)

// MaxTripDays is contract §1.4.
const MaxTripDays = 14

var currencyRE = regexp.MustCompile(`^[A-Z]{3}$`)

// input is a request after validation, with defaults applied.
type input struct {
	Destination string
	Start, End  time.Time
	Travelers   int
	Budget      *api.Budget
	Interests   []string
}

// problems collects at most one issue per field (contract §3.2), in order.
type problems struct{ list []api.ErrorDetail }

func (p *problems) add(field, issue string) {
	for _, d := range p.list {
		if d.Field == field {
			return
		}
	}
	p.list = append(p.list, api.ErrorDetail{Field: field, Issue: issue})
}

func validate(req api.CreateItineraryRequest) (input, []api.ErrorDetail) {
	var p problems
	in := input{Travelers: 1, Interests: []string{}, Budget: req.Budget}

	// §1.1
	in.Destination = strings.TrimSpace(req.Destination)
	switch n := utf8.RuneCountInString(in.Destination); {
	case n == 0:
		p.add("destination", "is required and must not be empty")
	case n > 100:
		p.add("destination", "must be at most 100 characters")
	}

	// §1.2
	start, okStart := parseDate(req.StartDate)
	if !okStart {
		p.add("startDate", "must be a real date in YYYY-MM-DD format")
	}
	end, okEnd := parseDate(req.EndDate)
	if !okEnd {
		p.add("endDate", "must be a real date in YYYY-MM-DD format")
	}
	// §1.3, §1.4
	if okStart && okEnd {
		switch n := tripLength(start, end); {
		case n < 1:
			p.add("endDate", "must not be before startDate")
		case n > MaxTripDays:
			p.add("endDate", fmt.Sprintf("trip must be at most %d days", MaxTripDays))
		}
	}
	in.Start, in.End = start, end

	// §1.5
	if req.Travelers != nil {
		if *req.Travelers < 1 || *req.Travelers > 20 {
			p.add("travelers", "must be an integer from 1 to 20")
		}
		in.Travelers = *req.Travelers
	}

	// §1.6
	if b := req.Budget; b != nil {
		if b.Amount < 1 {
			p.add("budget.amount", "must be a positive integer")
		}
		if !currencyRE.MatchString(b.Currency) {
			p.add("budget.currency", "must be 3 upper-case letters (ISO 4217)")
		}
	}

	// §1.7
	if req.Interests != nil {
		if len(*req.Interests) > 10 {
			p.add("interests", "must have at most 10 items")
		}
		for i, s := range *req.Interests {
			if n := utf8.RuneCountInString(s); n < 1 || n > 30 {
				p.add(fmt.Sprintf("interests.%d", i), "must be 1-30 characters")
			}
		}
		in.Interests = *req.Interests
	}

	return in, p.list
}
