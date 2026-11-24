// Package itinerary implements POST /v1/itineraries (docs/contract/create-itinerary.md).
package itinerary

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/thinhduyng/ai-travel-planner/api/internal/api"
	"github.com/thinhduyng/ai-travel-planner/api/internal/planner"
)

const maxBodyBytes = 64 << 10

var timeRE = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

// Handler implements the generated api.ServerInterface.
type Handler struct {
	Planner planner.Planner
	Now     func() time.Time
}

var _ api.ServerInterface = (*Handler)(nil)

// CreateItinerary handles POST /v1/itineraries.
func (h *Handler) CreateItinerary(w http.ResponseWriter, r *http.Request) {
	// §3.4
	if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
		WriteError(w, http.StatusUnsupportedMediaType, api.UNSUPPORTEDMEDIATYPE, "Content-Type must be application/json", nil)
		return
	}

	req, details, malformed := decode(r)
	if malformed {
		WriteError(w, http.StatusBadRequest, api.MALFORMEDJSON, "Request body is not valid JSON", nil) // §3.3
		return
	}
	if len(details) > 0 {
		WriteError(w, http.StatusBadRequest, api.VALIDATIONERROR, "Request body is invalid", details) // §3.2
		return
	}
	in, details := validate(req)
	if len(details) > 0 {
		WriteError(w, http.StatusBadRequest, api.VALIDATIONERROR, "Request body is invalid", details)
		return
	}

	dates := tripDates(in.Start, in.End)
	var budget *planner.Budget
	if in.Budget != nil {
		budget = &planner.Budget{Amount: in.Budget.Amount, Currency: in.Budget.Currency}
	}
	plan, err := h.Planner.Plan(r.Context(), planner.Request{
		Destination: in.Destination,
		Dates:       dates,
		Travelers:   in.Travelers,
		Budget:      budget,
		Interests:   in.Interests,
	})
	if err != nil || !usable(plan, len(dates)) {
		// §3.5 / §3.6: log the cause, never send it.
		slog.ErrorContext(r.Context(), "planner failed", "err", err, "days", len(plan), "want", len(dates))
		WriteError(w, http.StatusBadGateway, api.PLANNERUNAVAILABLE, "Could not generate an itinerary, try again later", nil)
		return
	}

	now := time.Now
	if h.Now != nil {
		now = h.Now
	}
	out := api.Itinerary{
		Id:          newUUIDv4(),
		Destination: in.Destination,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Travelers:   in.Travelers,
		Budget:      in.Budget,
		Interests:   in.Interests,
		Days:        make([]api.Day, len(dates)),
		CreatedAt:   now().UTC(),
	}
	for i, date := range dates {
		acts := make([]api.Activity, len(plan[i]))
		for j, a := range plan[i] {
			acts[j] = api.Activity{Time: a.Time, Title: a.Title, Description: a.Description}
		}
		sort.SliceStable(acts, func(x, y int) bool { return acts[x].Time < acts[y].Time }) // §2.7
		out.Days[i] = api.Day{Day: i + 1, Date: date, Activities: acts}
	}

	w.Header().Set("Location", "/v1/itineraries/"+out.Id) // §2.1
	writeJSON(w, http.StatusOK, out)
}

// decode reads the body strictly. It returns field-level details for type
// mismatches and unknown fields, or malformed=true for broken JSON.
func decode(r *http.Request) (req api.CreateItineraryRequest, details []api.ErrorDetail, malformed bool) {
	dec := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes))
	dec.DisallowUnknownFields() // §1.8
	err := dec.Decode(&req)
	if err == nil {
		return req, nil, false
	}
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			return req, nil, true // body is not an object at all
		}
		return req, []api.ErrorDetail{{Field: field, Issue: "must be of type " + typeErr.Type.String()}}, false
	}
	if f, ok := strings.CutPrefix(err.Error(), "json: unknown field "); ok {
		return req, []api.ErrorDetail{{Field: strings.Trim(f, `"`), Issue: "unknown field"}}, false
	}
	return req, nil, true
}

// usable checks the planner's output before any of it reaches a client (§3.5).
func usable(plan [][]planner.Activity, days int) bool {
	if len(plan) != days {
		return false
	}
	for _, day := range plan {
		if len(day) == 0 {
			return false
		}
		for _, a := range day {
			if !timeRE.MatchString(a.Time) || strings.TrimSpace(a.Title) == "" {
				return false
			}
		}
	}
	return true
}

func newUUIDv4() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// WriteError writes the contract's error shape (§3.1). details is never null.
func WriteError(w http.ResponseWriter, status int, code api.ErrorErrorCode, msg string, details []api.ErrorDetail) {
	var body api.Error
	body.Error.Code = code
	body.Error.Message = msg
	body.Error.Details = details
	if body.Error.Details == nil {
		body.Error.Details = []api.ErrorDetail{}
	}
	writeJSON(w, status, body)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
