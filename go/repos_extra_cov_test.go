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

func TestRepoService_IterSearchRepositories_MultiPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("X-Total-Count", "60")
		var data []*types.Repository
		if page == "2" {
			data = make([]*types.Repository, 10)
			for i := range data {
				data[i] = &types.Repository{}
			}
		} else {
			data = make([]*types.Repository, 50)
			for i := range data {
				data[i] = &types.Repository{}
			}
		}
		json.NewEncoder(w).Encode(types.SearchResults{Data: data, OK: true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for _, err := range f.Repos.IterSearchRepositories(context.Background(), "go") {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 60 {
		t.Fatalf("got %d repos across pages, want 60", count)
	}
}

func TestRepoService_IterSearchRepositories_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Repos.IterSearchRepositories(context.Background(), "go") {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from IterSearchRepositories")
	}
}

func TestRepoService_IterSearchRepositories_EarlyBreak_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode(types.SearchResults{Data: []*types.Repository{{}, {}, {}}, OK: true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Repos.IterSearchRepositories(context.Background(), "go") {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
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
