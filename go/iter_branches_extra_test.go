package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dappco.re/go/forge/types"
)

// A cluster of Iter* helpers share a custom-yield shape whose error and
// early-break branches the generated happy-path tests leave uncovered. These
// tests exercise those branches plus the time-filter application in
// IterTimeline.

func TestIssueService_IterTimeline_Filtered_Good(t *testing.T) {
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	before := time.Date(2026, 6, 7, 8, 9, 10, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("since"); got != since.Format(time.RFC3339) {
			t.Errorf("got since=%q", got)
		}
		if got := r.URL.Query().Get("before"); got != before.Format(time.RFC3339) {
			t.Errorf("got before=%q", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.TimelineComment{{ID: 1}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for c, err := range f.Issues.IterTimeline(context.Background(), "core", "go-forge", 1, &since, &before) {
		if err != nil {
			t.Fatal(err)
		}
		if c.ID != 1 {
			t.Errorf("got id %d", c.ID)
		}
		count++
	}
	if count != 1 {
		t.Fatalf("got %d items", count)
	}
}

func TestIssueService_IterTimeline_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Issues.IterTimeline(context.Background(), "core", "go-forge", 1, nil, nil) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from IterTimeline")
	}
}

func TestLabelService_IterLabelTemplates_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Labels.IterLabelTemplates(context.Background()) {
		if err != nil {
			sawErr = true
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from IterLabelTemplates")
	}
}

func TestLabelService_IterLabelTemplates_EarlyBreak_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]string{"Default", "Advanced", "GitHub"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Labels.IterLabelTemplates(context.Background()) {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}
