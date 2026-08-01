package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

// IterKeys / IterUserKeys build a query from optional filters before
// delegating to ListIter. The generated tests call them filter-free, leaving
// the filter-application loop and the iterator error path uncovered.

func TestUserService_IterKeys_Filtered_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/user/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("fingerprint"); got != "AB:CD" {
			t.Errorf("got fingerprint=%q, want AB:CD", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PublicKey{{ID: 1}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for key, err := range f.Users.IterKeys(context.Background(), UserKeyListOptions{Fingerprint: "AB:CD"}) {
		if err != nil {
			t.Fatal(err)
		}
		if key.ID != 1 {
			t.Errorf("got key id %d", key.ID)
		}
		count++
	}
	if count != 1 {
		t.Fatalf("got %d keys", count)
	}
}

func TestUserService_IterKeys_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Users.IterKeys(context.Background()) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from IterKeys")
	}
}

func TestUserService_IterUserKeys_Filtered_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/users/alice/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("fingerprint"); got != "EF:01" {
			t.Errorf("got fingerprint=%q, want EF:01", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PublicKey{{ID: 2}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for _, err := range f.Users.IterUserKeys(context.Background(), "alice", UserKeyListOptions{Fingerprint: "EF:01"}) {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 1 {
		t.Fatalf("got %d keys", count)
	}
}

func TestUserService_IterUserKeys_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "nope"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Users.IterUserKeys(context.Background(), "alice") {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from IterUserKeys")
	}
}
