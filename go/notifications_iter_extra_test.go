package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

// NotificationService.listIter pages internally and yields per-thread. The
// generated tests do not cover its error-yield, early-break or multi-page
// loop, so exercise those here via the Iter entrypoint.

func TestNotificationService_Iter_MultiPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("X-Total-Count", "60")
		var items []types.NotificationThread
		if page == "2" {
			items = make([]types.NotificationThread, 10)
		} else {
			items = make([]types.NotificationThread, 50)
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for _, err := range f.Notifications.Iter(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 60 {
		t.Fatalf("got %d threads across pages, want 60", count)
	}
}

func TestNotificationService_Iter_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Notifications.Iter(context.Background()) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from Iter")
	}
}

func TestNotificationService_Iter_EarlyBreak_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode([]types.NotificationThread{{ID: 1}, {ID: 2}, {ID: 3}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Notifications.Iter(context.Background()) {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}
