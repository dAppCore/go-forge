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

// repoTimeQuery feeds the repository tracked-time endpoints. Cover its
// no-filter, populated and empty-filter shapes, plus the named-user branch of
// ListUserRepos / IterUserRepos that the generated tests leave unexercised.

func TestRepoTimeQuery_NoFilters_Bad(t *testing.T) {
	if got := repoTimeQuery(); got != nil {
		t.Errorf("no filters should yield nil, got %v", got)
	}
}

func TestRepoTimeQuery_EmptyFilter_Bad(t *testing.T) {
	// A present-but-empty filter produces no query keys, so the result is nil.
	if got := repoTimeQuery(RepoTimeListOptions{}); got != nil {
		t.Errorf("empty filter should yield nil, got %v", got)
	}
}

func TestRepoTimeQuery_Populated_Good(t *testing.T) {
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	q := repoTimeQuery(RepoTimeListOptions{User: "alice", Since: &since})
	if q["user"] != "alice" {
		t.Errorf("got user=%q", q["user"])
	}
	if q["since"] != since.Format(time.RFC3339) {
		t.Errorf("got since=%q", q["since"])
	}
}

func TestRepoService_ListTimes_Filtered_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("user"); got != "alice" {
			t.Errorf("got user=%q, want alice", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.TrackedTime{{ID: 1}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	times, err := f.Repos.ListTimes(context.Background(), "core", "go-forge", RepoTimeListOptions{User: "alice"})
	if err != nil {
		t.Fatal(err)
	}
	if len(times) != 1 {
		t.Fatalf("got %d times", len(times))
	}
}

func TestRepoService_ListUserRepos_NamedUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/users/alice/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "r1"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Repos.ListUserRepos(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 || repos[0].Name != "r1" {
		t.Fatalf("got %#v", repos)
	}
}

func TestRepoService_IterUserRepos_NamedUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/users/bob/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "r2"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for repo, err := range f.Repos.IterUserRepos(context.Background(), "bob") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, repo.Name)
	}
	if len(got) != 1 || got[0] != "r2" {
		t.Fatalf("got %v", got)
	}
}
