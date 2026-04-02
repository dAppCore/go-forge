package forge

import (
	"bytes"
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"dappco.re/go/core/forge/types"
)

func TestRepoService_ListActivityFeeds_Good(t *testing.T) {
	date := time.Date(2026, time.April, 2, 15, 4, 5, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/activities/feeds" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("date"); got != "2026-04-02" {
			t.Errorf("wrong date: %s", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Activity{{
			ID:      7,
			OpType:  "create_repo",
			Content: "created repository",
		}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	activities, err := f.Repos.ListActivityFeeds(context.Background(), "core", "go-forge", ActivityFeedListOptions{Date: &date})
	if err != nil {
		t.Fatal(err)
	}
	if len(activities) != 1 || activities[0].ID != 7 || activities[0].OpType != "create_repo" {
		t.Fatalf("got %#v", activities)
	}
}

func TestRepoService_GetByID_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repositories/42" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Repository{
			ID:       42,
			Name:     "go-forge",
			FullName: "core/go-forge",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.GetByID(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if repo.ID != 42 || repo.Name != "go-forge" || repo.FullName != "core/go-forge" {
		t.Fatalf("got %#v", repo)
	}
}

func TestRepoService_GetRunnerRegistrationToken_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/runners/registration-token" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("token", "runner-token")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	token, err := f.Repos.GetRunnerRegistrationToken(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if token != "runner-token" {
		t.Fatalf("got token=%q, want %q", token, "runner-token")
	}
}

func TestRepoService_Migrate_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/migrate" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var opts types.MigrateRepoOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if opts.CloneAddr != "https://example.com/source.git" || opts.RepoName != "go-forge" || opts.RepoOwner != "core" {
			t.Fatalf("got %#v", opts)
		}
		json.NewEncoder(w).Encode(types.Repository{
			ID:       99,
			Name:     opts.RepoName,
			FullName: opts.RepoOwner + "/" + opts.RepoName,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.Migrate(context.Background(), &types.MigrateRepoOptions{
		CloneAddr: "https://example.com/source.git",
		RepoName:  "go-forge",
		RepoOwner: "core",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.ID != 99 || repo.Name != "go-forge" || repo.FullName != "core/go-forge" {
		t.Fatalf("got %#v", repo)
	}
}

func TestRepoService_ListTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.TopicName{TopicNames: []string{"go", "forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	topics, err := f.Repos.ListTopics(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(topics, []string{"go", "forge"}) {
		t.Fatalf("got %#v", topics)
	}
}

func TestRepoService_IterTopics_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.TopicName{TopicNames: []string{"go", "forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for topic, err := range f.Repos.IterTopics(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, topic)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if !reflect.DeepEqual(got, []string{"go", "forge"}) {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_SearchTopics_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/topics/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "go" {
			t.Errorf("wrong query: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("wrong page: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("wrong limit: %s", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.TopicResponse{
			{Name: "go", RepoCount: 10},
			{Name: "forge", RepoCount: 4},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	topics, err := f.Repos.SearchTopics(context.Background(), "go")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(topics) != 2 || topics[0].Name != "go" || topics[1].RepoCount != 4 {
		t.Fatalf("got %#v", topics)
	}
}

func TestRepoService_IterSearchTopics_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/topics/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "go" {
			t.Errorf("wrong query: %s", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.TopicResponse{{Name: "go", RepoCount: 10}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.TopicResponse
	for topic, err := range f.Repos.IterSearchTopics(context.Background(), "go") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, topic)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(got) != 1 || got[0].Name != "go" || got[0].RepoCount != 10 {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_SearchRepositoriesPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "go" {
			t.Errorf("wrong query: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("wrong page: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "2" {
			t.Errorf("wrong limit: %s", got)
		}
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode(types.SearchResults{
			Data: []*types.Repository{
				{Name: "go-forge"},
				{Name: "go-core"},
			},
			OK: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Repos.SearchRepositoriesPage(context.Background(), "go", ListOptions{Page: 1, Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 1 || page.TotalCount != 3 || !page.HasMore {
		t.Fatalf("got %#v", page)
	}
	if !reflect.DeepEqual(page.Items, []types.Repository{{Name: "go-forge"}, {Name: "go-core"}}) {
		t.Fatalf("got %#v", page.Items)
	}
}

func TestRepoService_SearchRepositories_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "go" {
			t.Errorf("wrong query: %s", got)
		}
		switch r.URL.Query().Get("page") {
		case "1":
			w.Header().Set("X-Total-Count", "3")
			json.NewEncoder(w).Encode(types.SearchResults{
				Data: []*types.Repository{
					{Name: "go-forge"},
					{Name: "go-core"},
				},
				OK: true,
			})
		case "2":
			w.Header().Set("X-Total-Count", "3")
			json.NewEncoder(w).Encode(types.SearchResults{
				Data: []*types.Repository{{Name: "go-utils"}},
				OK:   true,
			})
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Repos.SearchRepositories(context.Background(), "go")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 2 {
		t.Fatalf("expected 2 requests, got %d", requests)
	}
	if !reflect.DeepEqual(repos, []types.Repository{{Name: "go-forge"}, {Name: "go-core"}, {Name: "go-utils"}}) {
		t.Fatalf("got %#v", repos)
	}
}

func TestRepoService_IterSearchRepositories_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode(types.SearchResults{
			Data: []*types.Repository{{Name: "go-forge"}},
			OK:   true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.Repository
	for repo, err := range f.Repos.IterSearchRepositories(context.Background(), "go") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, repo)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if !reflect.DeepEqual(got, []types.Repository{{Name: "go-forge"}}) {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_UpdateTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.RepoTopicOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(body.Topics, []string{"go", "forge"}) {
			t.Fatalf("got %#v", body.Topics)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.UpdateTopics(context.Background(), "core", "go-forge", []string{"go", "forge"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_AddTopic_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/api/v1/repos/core/go-forge/topics/release%20candidate" {
			t.Errorf("wrong path: %s", r.URL.EscapedPath())
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.AddTopic(context.Background(), "core", "go-forge", "release candidate"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteTopic_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/api/v1/repos/core/go-forge/topics/release%20candidate" {
			t.Errorf("wrong path: %s", r.URL.EscapedPath())
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteTopic(context.Background(), "core", "go-forge", "release candidate"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_GetTag_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tags/v1.0.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Tag{
			Name:    "v1.0.0",
			Message: "Release 1.0.0",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tag, err := f.Repos.GetTag(context.Background(), "core", "go-forge", "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if tag.Name != "v1.0.0" || tag.Message != "Release 1.0.0" {
		t.Fatalf("got %#v", tag)
	}
}

func TestRepoService_DeleteTag_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tags/v1.0.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteTag(context.Background(), "core", "go-forge", "v1.0.0"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_GetLanguages_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/languages" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]int64{
			"go":    1200,
			"shell": 300,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	languages, err := f.Repos.GetLanguages(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(languages, map[string]int64{"go": 1200, "shell": 300}) {
		t.Fatalf("got %#v", languages)
	}
}

func TestRepoService_GetRawFileOrLFS_Good(t *testing.T) {
	want := []byte("lfs-pointer-or-content")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/media/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("ref"); got != "main" {
			t.Errorf("wrong ref: %s", got)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(want)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	got, err := f.Repos.GetRawFileOrLFS(context.Background(), "core", "go-forge", "README.md", "main")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRepoService_GetEditorConfig_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/editorconfig/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("ref"); got != "main" {
			t.Errorf("wrong ref: %s", got)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.GetEditorConfig(context.Background(), "core", "go-forge", "README.md", "main"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_GetEditorConfig_Bad_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.GetEditorConfig(context.Background(), "core", "go-forge", "README.md", "main"); !IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRepoService_ApplyDiffPatch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/diffpatch" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.UpdateFileOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.SHA != "abc123" || body.Message != "apply patch" || body.ContentBase64 != "ZGlmZiBjb250ZW50" {
			t.Fatalf("got %#v", body)
		}
		json.NewEncoder(w).Encode(types.FileResponse{
			Commit: &types.FileCommitResponse{SHA: "commit-1"},
			Content: &types.ContentsResponse{
				Path: "README.md",
				SHA:  "file-1",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	resp, err := f.Repos.ApplyDiffPatch(context.Background(), "core", "go-forge", &types.UpdateFileOptions{
		SHA:           "abc123",
		Message:       "apply patch",
		ContentBase64: "ZGlmZiBjb250ZW50",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Commit == nil || resp.Commit.SHA != "commit-1" {
		t.Fatalf("got %#v", resp)
	}
	if resp.Content == nil || resp.Content.Path != "README.md" || resp.Content.SHA != "file-1" {
		t.Fatalf("got %#v", resp)
	}
}

func TestRepoService_ListForks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/forks" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Repository{
			{ID: 11, Name: "go-forge-fork", FullName: "alice/go-forge-fork"},
			{ID: 12, Name: "go-forge-fork-2", FullName: "bob/go-forge-fork-2"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	forks, err := f.Repos.ListForks(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(forks) != 2 || forks[0].FullName != "alice/go-forge-fork" || forks[1].FullName != "bob/go-forge-fork-2" {
		t.Fatalf("got %#v", forks)
	}
}

func TestRepoService_ListTagProtections_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tag_protections" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.TagProtection{
			{ID: 1, NamePattern: "v*"},
			{ID: 2, NamePattern: "release-*"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tagProtections, err := f.Repos.ListTagProtections(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(tagProtections) != 2 || tagProtections[0].ID != 1 || tagProtections[1].NamePattern != "release-*" {
		t.Fatalf("got %#v", tagProtections)
	}
}

func TestRepoService_GetTagProtection_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tag_protections/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.TagProtection{
			ID:          7,
			NamePattern: "v*",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tagProtection, err := f.Repos.GetTagProtection(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if tagProtection.ID != 7 || tagProtection.NamePattern != "v*" {
		t.Fatalf("got %#v", tagProtection)
	}
}

func TestRepoService_CreateTagProtection_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tag_protections" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreateTagProtectionOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.NamePattern != "v*" || !reflect.DeepEqual(body.WhitelistTeams, []string{"release-team"}) || !reflect.DeepEqual(body.WhitelistUsernames, []string{"alice"}) {
			t.Fatalf("got %#v", body)
		}
		json.NewEncoder(w).Encode(types.TagProtection{
			ID:                 9,
			NamePattern:        body.NamePattern,
			WhitelistTeams:     body.WhitelistTeams,
			WhitelistUsernames: body.WhitelistUsernames,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tagProtection, err := f.Repos.CreateTagProtection(context.Background(), "core", "go-forge", &types.CreateTagProtectionOption{
		NamePattern:        "v*",
		WhitelistTeams:     []string{"release-team"},
		WhitelistUsernames: []string{"alice"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if tagProtection.ID != 9 || tagProtection.NamePattern != "v*" {
		t.Fatalf("got %#v", tagProtection)
	}
}

func TestRepoService_EditTagProtection_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tag_protections/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditTagProtectionOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.NamePattern != "release-*" || !reflect.DeepEqual(body.WhitelistTeams, []string{"release-team"}) {
			t.Fatalf("got %#v", body)
		}
		json.NewEncoder(w).Encode(types.TagProtection{
			ID:             7,
			NamePattern:    body.NamePattern,
			WhitelistTeams: body.WhitelistTeams,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tagProtection, err := f.Repos.EditTagProtection(context.Background(), "core", "go-forge", 7, &types.EditTagProtectionOption{
		NamePattern:    "release-*",
		WhitelistTeams: []string{"release-team"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if tagProtection.ID != 7 || tagProtection.NamePattern != "release-*" {
		t.Fatalf("got %#v", tagProtection)
	}
}

func TestRepoService_DeleteTagProtection_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tag_protections/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteTagProtection(context.Background(), "core", "go-forge", 7); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_ListKeysWithFilters_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("key_id"); got != "7" {
			t.Errorf("got key_id=%q, want %q", got, "7")
		}
		if got := r.URL.Query().Get("fingerprint"); got != "aa:bb:cc" {
			t.Errorf("got fingerprint=%q, want %q", got, "aa:bb:cc")
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("got limit=%q, want %q", got, "50")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.DeployKey{{ID: 7, Title: "deploy"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	keys, err := f.Repos.ListKeys(context.Background(), "core", "go-forge", RepoKeyListOptions{
		KeyID:       7,
		Fingerprint: "aa:bb:cc",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 {
		t.Fatalf("got %d keys, want 1", len(keys))
	}
	if keys[0].ID != 7 || keys[0].Title != "deploy" {
		t.Fatalf("got %#v", keys[0])
	}
}

func TestRepoService_GetKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/keys/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.DeployKey{ID: 7, Title: "deploy"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Repos.GetKey(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 7 || key.Title != "deploy" {
		t.Fatalf("got %#v", key)
	}
}

func TestRepoService_CreateKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateKeyOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Title != "deploy" || opts.Key != "ssh-ed25519 AAAA..." || !opts.ReadOnly {
			t.Fatalf("got %#v", opts)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.DeployKey{ID: 9, Title: opts.Title, Key: opts.Key, ReadOnly: opts.ReadOnly})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Repos.CreateKey(context.Background(), "core", "go-forge", &types.CreateKeyOption{
		Title:    "deploy",
		Key:      "ssh-ed25519 AAAA...",
		ReadOnly: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 9 || key.Title != "deploy" || !key.ReadOnly {
		t.Fatalf("got %#v", key)
	}
}

func TestRepoService_DeleteKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/keys/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteKey(context.Background(), "core", "go-forge", 7); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteTag_Bad_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tags/missing" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Repos.DeleteTag(context.Background(), "core", "go-forge", "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Fatalf("got %v, want not found", err)
	}
}

func TestRepoService_ListIssueTemplates_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issue_templates" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.IssueTemplate{
			{
				Name:    "bug report",
				Title:   "Bug report",
				Content: "Describe the problem",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	templates, err := f.Repos.ListIssueTemplates(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(templates) != 1 {
		t.Fatalf("got %d templates, want 1", len(templates))
	}
	if templates[0].Name != "bug report" || templates[0].Title != "Bug report" {
		t.Fatalf("got %#v", templates[0])
	}
}

func TestRepoService_GetIssueConfig_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issue_config" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.IssueConfig{
			BlankIssuesEnabled: true,
			ContactLinks: []*types.IssueConfigContactLink{{
				Name:  "Security",
				URL:   "https://example.com/security",
				About: "Report a vulnerability",
			}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	cfg, err := f.Repos.GetIssueConfig(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.BlankIssuesEnabled {
		t.Fatalf("expected blank issues to be enabled, got %#v", cfg)
	}
	if len(cfg.ContactLinks) != 1 || cfg.ContactLinks[0].Name != "Security" {
		t.Fatalf("got %#v", cfg.ContactLinks)
	}
}

func TestRepoService_ValidateIssueConfig_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issue_config/validate" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.IssueConfigValidation{
			Valid:   false,
			Message: "invalid contact link URL",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.ValidateIssueConfig(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if result.Valid || result.Message != "invalid contact link URL" {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_GetNewPinAllowed_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/new_pin_allowed" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.NewIssuePinsAllowed{
			Issues:       true,
			PullRequests: false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.GetNewPinAllowed(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Issues || result.PullRequests {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_UpdateAvatar_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/avatar" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.UpdateRepoAvatarOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Image != "iVBORw0KGgoAAAANSUhEUg==" {
			t.Fatalf("got image=%q", body.Image)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.UpdateAvatar(context.Background(), "core", "go-forge", &types.UpdateRepoAvatarOption{
		Image: "iVBORw0KGgoAAAANSUhEUg==",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteAvatar_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/avatar" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteAvatar(context.Background(), "core", "go-forge"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_ListPushMirrors_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/push_mirrors" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PushMirror{
			{RemoteName: "mirror-a"},
			{RemoteName: "mirror-b"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	mirrors, err := f.Repos.ListPushMirrors(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(mirrors) != 2 || mirrors[0].RemoteName != "mirror-a" || mirrors[1].RemoteName != "mirror-b" {
		t.Fatalf("got %#v", mirrors)
	}
}

func TestRepoService_GetPushMirror_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/push_mirrors/mirror-a" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.PushMirror{
			RemoteName:    "mirror-a",
			RemoteAddress: "ssh://git@example.com/core/go-forge.git",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	mirror, err := f.Repos.GetPushMirror(context.Background(), "core", "go-forge", "mirror-a")
	if err != nil {
		t.Fatal(err)
	}
	if mirror.RemoteName != "mirror-a" {
		t.Fatalf("got remote_name=%q", mirror.RemoteName)
	}
}

func TestRepoService_CreatePushMirror_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/push_mirrors" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreatePushMirrorOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.RemoteAddress != "ssh://git@example.com/core/go-forge.git" || !body.SyncOnCommit {
			t.Fatalf("got %#v", body)
		}
		json.NewEncoder(w).Encode(types.PushMirror{
			RemoteName:    "mirror-a",
			RemoteAddress: body.RemoteAddress,
			SyncOnCommit:  body.SyncOnCommit,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	mirror, err := f.Repos.CreatePushMirror(context.Background(), "core", "go-forge", &types.CreatePushMirrorOption{
		RemoteAddress:  "ssh://git@example.com/core/go-forge.git",
		RemoteUsername: "git",
		RemotePassword: "secret",
		Interval:       "1h",
		SyncOnCommit:   true,
		UseSSH:         true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if mirror.RemoteName != "mirror-a" || !mirror.SyncOnCommit {
		t.Fatalf("got %#v", mirror)
	}
}

func TestRepoService_DeletePushMirror_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/push_mirrors/mirror-a" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeletePushMirror(context.Background(), "core", "go-forge", "mirror-a"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_SyncPushMirrors_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/push_mirrors-sync" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.SyncPushMirrors(context.Background(), "core", "go-forge"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_GetSubscription_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.WatchInfo{
			Subscribed: true,
			Ignored:    false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.GetSubscription(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Subscribed || result.Ignored {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_ListPinnedPullRequests_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/pinned" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PullRequest{
			{ID: 7, Title: "pin me"},
			{ID: 8, Title: "keep me"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pulls, err := f.Repos.ListPinnedPullRequests(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(pulls), 2; got != want {
		t.Fatalf("got %d pull requests, want %d", got, want)
	}
	if pulls[0].Title != "pin me" {
		t.Fatalf("got first title %q", pulls[0].Title)
	}
}

func TestRepoService_IterPinnedPullRequests_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/pinned" {
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
			json.NewEncoder(w).Encode([]types.PullRequest{{ID: 7, Title: "pin me"}})
		case 2:
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Errorf("got page=%q, want %q", got, "2")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.PullRequest{{ID: 8, Title: "keep me"}})
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for pr, err := range f.Repos.IterPinnedPullRequests(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, pr.Title)
	}
	if len(got) != 2 || got[0] != "pin me" || got[1] != "keep me" {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_ListStargazers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/stargazers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}, {UserName: "bob"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListStargazers(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 || users[0].UserName != "alice" || users[1].UserName != "bob" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_ListSubscribers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscribers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "charlie"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListSubscribers(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].UserName != "charlie" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_Compare_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/compare/main...feature" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Compare{
			TotalCommits: 2,
			Commits: []*types.Commit{
				{SHA: "abc123"},
				{SHA: "def456"},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	got, err := f.Repos.Compare(context.Background(), "core", "go-forge", "main...feature")
	if err != nil {
		t.Fatal(err)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if got.TotalCommits != 2 || len(got.Commits) != 2 || got.Commits[0].SHA != "abc123" || got.Commits[1].SHA != "def456" {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_GetSigningKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/signing-key.gpg" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("-----BEGIN PGP PUBLIC KEY BLOCK-----\n..."))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Repos.GetSigningKey(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	want := "-----BEGIN PGP PUBLIC KEY BLOCK-----\n..."
	if key != want {
		t.Fatalf("got %q, want %q", key, want)
	}
}

func TestRepoService_ListFlags_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/flags" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]string{"alpha", "beta"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	flags, err := f.Repos.ListFlags(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(flags, []string{"alpha", "beta"}) {
		t.Fatalf("got %#v", flags)
	}
}

func TestRepoService_IterFlags_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/flags" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("got limit=%q, want %q", got, "50")
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]string{"alpha", "beta"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for flag, err := range f.Repos.IterFlags(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, flag)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if !reflect.DeepEqual(got, []string{"alpha", "beta"}) {
		t.Fatalf("got %#v", got)
	}
}

func TestRepoService_ReplaceFlags_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/flags" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.ReplaceFlagsOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(body.Flags, []string{"alpha", "beta"}) {
			t.Fatalf("got %#v", body.Flags)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.ReplaceFlags(context.Background(), "core", "go-forge", &types.ReplaceFlagsOption{Flags: []string{"alpha", "beta"}}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteFlags_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/flags" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteFlags(context.Background(), "core", "go-forge"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_ListAssignees_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/assignees" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}, {UserName: "bob"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListAssignees(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 || users[0].UserName != "alice" || users[1].UserName != "bob" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_IterAssignees_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/assignees" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}, {UserName: "bob"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var names []string
	for user, err := range f.Repos.IterAssignees(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, user.UserName)
	}
	if len(names) != 2 || names[0] != "alice" || names[1] != "bob" {
		t.Fatalf("got %#v", names)
	}
}

func TestRepoService_ListCollaborators_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListCollaborators(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].UserName != "alice" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_AddCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.AddCollaboratorOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Permission != "write" {
			t.Fatalf("got permission=%q, want %q", body.Permission, "write")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.AddCollaborator(context.Background(), "core", "go-forge", "alice", &types.AddCollaboratorOption{Permission: "write"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteCollaborator(context.Background(), "core", "go-forge", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_CheckCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	ok, err := f.Repos.CheckCollaborator(context.Background(), "core", "go-forge", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected collaborator check to return true")
	}
}

func TestRepoService_GetCollaboratorPermission_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice/permission" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.RepoCollaboratorPermission{
			Permission: "write",
			RoleName:   "collaborator",
			User:       &types.User{UserName: "alice"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	perm, err := f.Repos.GetCollaboratorPermission(context.Background(), "core", "go-forge", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if perm.Permission != "write" || perm.User == nil || perm.User.UserName != "alice" {
		t.Fatalf("got %#v", perm)
	}

	perm, err = f.Repos.GetRepoPermissions(context.Background(), "core", "go-forge", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if perm.Permission != "write" || perm.User == nil || perm.User.UserName != "alice" {
		t.Fatalf("got %#v", perm)
	}
}

func TestRepoService_ListRepoTeams_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/teams" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.Team{{ID: 7, Name: "platform"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	teams, err := f.Repos.ListRepoTeams(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(teams) != 1 || teams[0].ID != 7 || teams[0].Name != "platform" {
		t.Fatalf("got %#v", teams)
	}
}

func TestRepoService_GetRepoTeam_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/teams/platform" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Team{ID: 7, Name: "platform"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	team, err := f.Repos.GetRepoTeam(context.Background(), "core", "go-forge", "platform")
	if err != nil {
		t.Fatal(err)
	}
	if team.ID != 7 || team.Name != "platform" {
		t.Fatalf("got %#v", team)
	}
}

func TestRepoService_AddRepoTeam_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/teams/platform" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.AddRepoTeam(context.Background(), "core", "go-forge", "platform"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteRepoTeam_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/teams/platform" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteRepoTeam(context.Background(), "core", "go-forge", "platform"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_Watch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.WatchInfo{
			Subscribed: true,
			Ignored:    false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.Watch(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Subscribed || result.Ignored {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_Unwatch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.Unwatch(context.Background(), "core", "go-forge"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_ForkWithOptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/forks" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var opts types.CreateForkOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if opts.Name != "go-forge-fork" || opts.Organization != "core-team" {
			t.Fatalf("got %#v", opts)
		}
		json.NewEncoder(w).Encode(types.Repository{
			Name:     opts.Name,
			FullName: opts.Organization + "/" + opts.Name,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.ForkWithOptions(context.Background(), "core", "go-forge", &types.CreateForkOption{
		Name:         "go-forge-fork",
		Organization: "core-team",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "go-forge-fork" || repo.FullName != "core-team/go-forge-fork" {
		t.Fatalf("got %#v", repo)
	}
}

func TestRepoService_PathParamsAreEscaped_Good(t *testing.T) {
	owner := "acme org"
	repo := "my/repo"
	org := "team alpha"

	t.Run("ListOrgRepos", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() != "/api/v1/orgs/team%20alpha/repos" {
				t.Errorf("got path %q, want %q", r.URL.EscapedPath(), "/api/v1/orgs/team%20alpha/repos")
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode([]types.Repository{})
		}))
		defer srv.Close()

		f := NewForge(srv.URL, "tok")
		_, err := f.Repos.ListOrgRepos(context.Background(), org)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Fork", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			want := "/api/v1/repos/acme%20org/my%2Frepo/forks"
			if r.URL.EscapedPath() != want {
				t.Errorf("got path %q, want %q", r.URL.EscapedPath(), want)
				http.NotFound(w, r)
				return
			}
			var opts types.CreateForkOption
			if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if opts.Organization != "" {
				t.Fatalf("got organisation %q, want empty", opts.Organization)
			}
			json.NewEncoder(w).Encode(types.Repository{Name: repo})
		}))
		defer srv.Close()

		f := NewForge(srv.URL, "tok")
		if _, err := f.Repos.Fork(context.Background(), owner, repo, ""); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Generate", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			want := "/api/v1/repos/acme%20org/template%2Frepo/generate"
			if r.URL.EscapedPath() != want {
				t.Errorf("got path %q, want %q", r.URL.EscapedPath(), want)
				http.NotFound(w, r)
				return
			}
			if r.Method != http.MethodPost {
				t.Errorf("expected POST, got %s", r.Method)
			}
			var body types.GenerateRepoOption
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %v", err)
			}
			if body.Owner != "acme org" || body.Name != "generated repo" || !body.Private || !body.Topics {
				t.Fatalf("got %#v", body)
			}
			json.NewEncoder(w).Encode(types.Repository{Name: body.Name, FullName: "acme org/" + body.Name})
		}))
		defer srv.Close()

		f := NewForge(srv.URL, "tok")
		repo, err := f.Repos.Generate(context.Background(), owner, "template/repo", &types.GenerateRepoOption{
			Owner:   "acme org",
			Name:    "generated repo",
			Private: true,
			Topics:  true,
		})
		if err != nil {
			t.Fatal(err)
		}
		if repo.Name != "generated repo" || repo.FullName != "acme org/generated repo" {
			t.Fatalf("got %#v", repo)
		}
	})
}
