package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

// GetUserByID validates the id, searches, and matches by exact id. Cover the
// invalid-id, search-error, found and not-found-after-search branches.

func TestUserService_GetUserByID_InvalidID_Bad(t *testing.T) {
	f := NewForge("http://localhost", "tok")
	_, err := f.Users.GetUserByID(context.Background(), 0)
	if err == nil || !IsNotFound(err) {
		t.Fatalf("id < 1 should yield a not-found APIError, got %v", err)
	}
}

func TestUserService_GetUserByID_SearchError_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Users.GetUserByID(context.Background(), 5); err == nil {
		t.Fatal("expected the search error to propagate")
	}
}

func TestUserService_GetUserByID_Found_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode(map[string]any{"data": []types.User{{ID: 5}}, "ok": true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.GetUserByID(context.Background(), 5)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != 5 {
		t.Fatalf("got id %d, want 5", user.ID)
	}
}

func TestUserService_GetUserByID_NotFoundAfterSearch_Bad(t *testing.T) {
	// The search returns a different user, so no exact id match is found.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode(map[string]any{"data": []types.User{{ID: 99}}, "ok": true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Users.GetUserByID(context.Background(), 5)
	if err == nil || !IsNotFound(err) {
		t.Fatalf("no exact match should yield a not-found APIError, got %v", err)
	}
}

// Representative early-break coverage for the custom-yield iterators: breaking
// after the first item must stop iteration (the `if !yield { return }` guard).

func TestCommitService_IterStatuses_EarlyBreak_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]types.CommitStatus{{ID: 1}, {ID: 2}, {ID: 3}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Commits.IterStatuses(context.Background(), "core", "go-forge", "main") {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}

func TestRepoService_IterTopics_EarlyBreak_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(types.TopicName{TopicNames: []string{"go", "forge", "client"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Repos.IterTopics(context.Background(), "core", "go-forge") {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}
