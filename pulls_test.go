package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestPullService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PullRequest{
			{ID: 1, Title: "add feature"},
			{ID: 2, Title: "fix bug"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Pulls.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].Title != "add feature" {
		t.Errorf("got title=%q, want %q", result.Items[0].Title, "add feature")
	}
}

func TestPullService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.PullRequest{ID: 1, Title: "add feature", Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Title != "add feature" {
		t.Errorf("got title=%q", pr.Title)
	}
}

func TestPullService_Create_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreatePullRequestOption
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.PullRequest{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.Create(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.CreatePullRequestOption{
		Title: "new pull request",
		Head:  "feature-branch",
		Base:  "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Title != "new pull request" {
		t.Errorf("got title=%q", pr.Title)
	}
}

func TestPullService_ListReviewers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/reviewers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{
			{UserName: "alice"},
			{UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	reviewers, err := f.Pulls.ListReviewers(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(reviewers) != 2 || reviewers[0].UserName != "alice" || reviewers[1].UserName != "bob" {
		t.Fatalf("got %#v", reviewers)
	}
}

func TestPullService_ListFiles_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/files" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.ChangedFile{
			{Filename: "README.md", Status: "modified", Additions: 2, Deletions: 1},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	files, err := f.Pulls.ListFiles(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d files, want 1", len(files))
	}
	if files[0].Filename != "README.md" || files[0].Status != "modified" {
		t.Fatalf("got %#v", files[0])
	}
}

func TestPullService_IterFiles_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/files" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		switch requests {
		case 1:
			if got := r.URL.Query().Get("page"); got != "1" {
				t.Errorf("got page=%q, want %q", got, "1")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.ChangedFile{{Filename: "README.md", Status: "modified"}})
		case 2:
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Errorf("got page=%q, want %q", got, "2")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.ChangedFile{{Filename: "docs/guide.md", Status: "added"}})
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for file, err := range f.Pulls.IterFiles(context.Background(), "core", "go-forge", 7) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, file.Filename)
	}
	if len(got) != 2 || got[0] != "README.md" || got[1] != "docs/guide.md" {
		t.Fatalf("got %#v", got)
	}
}

func TestPullService_IterReviewers_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/reviewers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		switch requests {
		case 1:
			if got := r.URL.Query().Get("page"); got != "1" {
				t.Errorf("got page=%q, want %q", got, "1")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}})
		case 2:
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Errorf("got page=%q, want %q", got, "2")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.User{{UserName: "bob"}})
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for reviewer, err := range f.Pulls.IterReviewers(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, reviewer.UserName)
	}
	if len(got) != 2 || got[0] != "alice" || got[1] != "bob" {
		t.Fatalf("got %#v", got)
	}
}

func TestPullService_RequestReviewers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/requested_reviewers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.PullReviewRequestOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body.Reviewers) != 2 || body.Reviewers[0] != "alice" || body.Reviewers[1] != "bob" {
			t.Fatalf("got reviewers %#v", body.Reviewers)
		}
		if len(body.TeamReviewers) != 1 || body.TeamReviewers[0] != "platform" {
			t.Fatalf("got team reviewers %#v", body.TeamReviewers)
		}
		json.NewEncoder(w).Encode([]types.PullReview{
			{ID: 101, Body: "requested"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	reviews, err := f.Pulls.RequestReviewers(context.Background(), "core", "go-forge", 7, &types.PullReviewRequestOptions{
		Reviewers:     []string{"alice", "bob"},
		TeamReviewers: []string{"platform"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 || reviews[0].ID != 101 || reviews[0].Body != "requested" {
		t.Fatalf("got %#v", reviews)
	}
}

func TestPullService_CancelReviewRequests_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/requested_reviewers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.PullReviewRequestOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if len(body.Reviewers) != 1 || body.Reviewers[0] != "alice" {
			t.Fatalf("got reviewers %#v", body.Reviewers)
		}
		if len(body.TeamReviewers) != 1 || body.TeamReviewers[0] != "platform" {
			t.Fatalf("got team reviewers %#v", body.TeamReviewers)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Pulls.CancelReviewRequests(context.Background(), "core", "go-forge", 7, &types.PullReviewRequestOptions{
		Reviewers:     []string{"alice"},
		TeamReviewers: []string{"platform"},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestPullService_Merge_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/merge" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["Do"] != "merge" {
			t.Errorf("got Do=%q, want %q", body["Do"], "merge")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Pulls.Merge(context.Background(), "core", "go-forge", 7, "merge")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPullService_Merge_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "already merged"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Pulls.Merge(context.Background(), "core", "go-forge", 7, "merge"); !IsConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestPullService_CancelScheduledAutoMerge_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/merge" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Pulls.CancelScheduledAutoMerge(context.Background(), "core", "go-forge", 7); err != nil {
		t.Fatal(err)
	}
}
