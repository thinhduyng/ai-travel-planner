package contract

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/thinhduyng/ai-travel-planner/api/internal/planner"
	"github.com/thinhduyng/ai-travel-planner/api/internal/server"
)

var (
	uuidV4 = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	isoUTC = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(\.\d+)?Z$`)
	hhmm   = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)
)

// start runs the real server on a random port: contract tests go over HTTP.
func start(t *testing.T, p planner.Planner) string {
	t.Helper()
	srv := httptest.NewServer(server.New(server.Options{Planner: p}))
	t.Cleanup(srv.Close)
	return srv.URL
}

func validBody() map[string]any {
	return map[string]any{
		"destination": "Da Lat, Vietnam",
		"startDate":   "2026-12-20",
		"endDate":     "2026-12-22",
		"travelers":   2,
		"budget":      map[string]any{"amount": 5000000, "currency": "VND"},
		"interests":   []any{"food", "nature"},
	}
}

func with(kv ...any) map[string]any {
	b := validBody()
	for i := 0; i < len(kv); i += 2 {
		b[kv[i].(string)] = kv[i+1]
	}
	return b
}

func without(key string) map[string]any {
	b := validBody()
	delete(b, key)
	return b
}

type response struct {
	Status int
	Header http.Header
	Raw    string
	Body   map[string]any
}

func post(t *testing.T, base string, body any, contentType ...string) response {
	t.Helper()
	var raw string
	if s, ok := body.(string); ok {
		raw = s
	} else {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		raw = string(b)
	}
	ct := "application/json"
	if len(contentType) > 0 {
		ct = contentType[0]
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, base+"/v1/itineraries", strings.NewReader(raw))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", ct)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()
	data, _ := io.ReadAll(res.Body)
	out := response{Status: res.StatusCode, Header: res.Header, Raw: string(data)}
	_ = json.Unmarshal(data, &out.Body)
	return out
}

func days(t *testing.T, r response) []map[string]any {
	t.Helper()
	list, ok := r.Body["days"].([]any)
	if !ok {
		t.Fatalf("days is not an array: %s", r.Raw)
	}
	out := make([]map[string]any, len(list))
	for i, d := range list {
		out[i] = d.(map[string]any)
	}
	return out
}

func expectStatus(t *testing.T, r response, want int) {
	t.Helper()
	if r.Status != want {
		t.Fatalf("status = %d, want %d; body: %s", r.Status, want, r.Raw)
	}
}

// expectError checks the §3.1 shape and code.
func expectError(t *testing.T, r response, status int, code string) []string {
	t.Helper()
	expectStatus(t, r, status)
	e, ok := r.Body["error"].(map[string]any)
	if !ok || len(r.Body) != 1 {
		t.Fatalf("body is not {error:{...}}: %s", r.Raw)
	}
	if e["code"] != code {
		t.Fatalf("error.code = %v, want %s", e["code"], code)
	}
	if _, ok := e["message"].(string); !ok {
		t.Fatalf("error.message is not a string: %s", r.Raw)
	}
	details, ok := e["details"].([]any)
	if !ok {
		t.Fatalf("error.details is not an array: %s", r.Raw)
	}
	var fields []string
	for _, d := range details {
		fields = append(fields, d.(map[string]any)["field"].(string))
	}
	return fields
}

func expectValidationError(t *testing.T, r response, field string) {
	t.Helper()
	fields := expectError(t, r, http.StatusBadRequest, "VALIDATION_ERROR")
	for _, f := range fields {
		if f == field {
			return
		}
	}
	t.Fatalf("details fields = %v, want to contain %q", fields, field)
}
