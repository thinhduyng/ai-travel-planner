// Package server wires the HTTP routes.
package server

import (
	"encoding/json"
	"net/http"

	"github.com/thinhduyng/ai-travel-planner/api/internal/planner"
)

// Options configures the server. A nil Planner means planner.Fake.
type Options struct {
	Planner planner.Planner
}

// New returns the API's HTTP handler.
func New(opts Options) http.Handler {
	if opts.Planner == nil {
		opts.Planner = planner.Fake{}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
