// Package planner defines the boundary to the LLM that writes itineraries.
package planner

import "context"

// Budget is the traveller's total budget in the smallest unit of Currency.
type Budget struct {
	Amount   int64
	Currency string
}

// Request is what the API asks the planner for.
type Request struct {
	Destination string
	// Dates has one YYYY-MM-DD entry per trip day, in order.
	Dates     []string
	Travelers int
	Budget    *Budget
	Interests []string
}

// Activity is one planned item in a day.
type Activity struct {
	Time        string // HH:MM, 24h
	Title       string
	Description *string
}

// Planner returns one activity list per entry in Request.Dates.
type Planner interface {
	Plan(ctx context.Context, req Request) ([][]Activity, error)
}

// Func adapts a plain function to the Planner interface.
type Func func(ctx context.Context, req Request) ([][]Activity, error)

// Plan implements Planner.
func (f Func) Plan(ctx context.Context, req Request) ([][]Activity, error) { return f(ctx, req) }
