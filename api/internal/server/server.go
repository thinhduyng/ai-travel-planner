// Package server wires the HTTP routes.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/thinhduyng/ai-travel-planner/api/internal/api"
	"github.com/thinhduyng/ai-travel-planner/api/internal/itinerary"
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
	// Routes for the contract endpoints come from the OpenAPI spec (internal/api).
	api.HandlerWithOptions(&itinerary.Handler{Planner: opts.Planner}, api.StdHTTPServerOptions{BaseRouter: mux})
	return recoverer(mux)
}

// recoverer turns panics into a generic 500 with no internals (contract §3.6).
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				slog.ErrorContext(r.Context(), "panic", "value", v)
				writeJSON(w, http.StatusInternalServerError, map[string]any{
					"error": map[string]any{"code": "INTERNAL_ERROR", "message": "Something went wrong", "details": []any{}},
				})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
