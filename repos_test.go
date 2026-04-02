package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"dappco.re/go/core/forge/types"
)

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
