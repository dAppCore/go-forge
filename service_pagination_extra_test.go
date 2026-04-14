package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestRepoService_ListOrgReposPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Errorf("got limit=%q, want %q", got, "25")
		}
		w.Header().Set("X-Total-Count", "51")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "go-forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Repos.ListOrgReposPage(context.Background(), "core", ListOptions{Page: 2, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 2 || page.TotalCount != 51 || len(page.Items) != 1 || page.Items[0].Name != "go-forge" {
		t.Fatalf("got %#v", page)
	}
}

func TestRepoService_ListUserReposPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/Virgil/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("got limit=%q, want %q", got, "10")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "go-user"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Repos.ListUserReposPage(context.Background(), "Virgil", ListOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 1 || len(page.Items) != 1 || page.Items[0].Name != "go-user" {
		t.Fatalf("got %#v", page)
	}
}

func TestIssueService_ListRepoIssuesPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want %q", got, "open")
		}
		if got := r.URL.Query().Get("page"); got != "3" {
			t.Errorf("got page=%q, want %q", got, "3")
		}
		if got := r.URL.Query().Get("limit"); got != "15" {
			t.Errorf("got limit=%q, want %q", got, "15")
		}
		w.Header().Set("X-Total-Count", "42")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 7, Title: "bug"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Issues.ListRepoIssuesPage(context.Background(), "core", "go-forge", ListOptions{Page: 3, PageSize: 15}, &types.ListIssueOption{State: "open"})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 3 || page.TotalCount != 42 || len(page.Items) != 1 || page.Items[0].Title != "bug" {
		t.Fatalf("got %#v", page)
	}
}

func TestPullService_ListPullRequestsPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want %q", got, "open")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "20" {
			t.Errorf("got limit=%q, want %q", got, "20")
		}
		w.Header().Set("X-Total-Count", "21")
		json.NewEncoder(w).Encode([]types.PullRequest{{ID: 9, Title: "feature"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Pulls.ListPullRequestsPage(context.Background(), "core", "go-forge", ListOptions{Page: 2, PageSize: 20}, &types.ListPullRequestsOption{State: "open"})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 2 || page.TotalCount != 21 || len(page.Items) != 1 || page.Items[0].Title != "feature" {
		t.Fatalf("got %#v", page)
	}
}

func TestCommitService_ListCommitsPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("sha"); got != "main" {
			t.Errorf("got sha=%q, want %q", got, "main")
		}
		if got := r.URL.Query().Get("page"); got != "4" {
			t.Errorf("got page=%q, want %q", got, "4")
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("got limit=%q, want %q", got, "5")
		}
		w.Header().Set("X-Total-Count", "17")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "abc123"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Commits.ListCommitsPage(context.Background(), "core", "go-forge", ListOptions{Page: 4, PageSize: 5}, &types.ListCommitsOption{Sha: "main"})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 4 || page.TotalCount != 17 || len(page.Items) != 1 || page.Items[0].SHA != "abc123" {
		t.Fatalf("got %#v", page)
	}
}

func TestBranchService_ListBranchesPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "30" {
			t.Errorf("got limit=%q, want %q", got, "30")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Branch{{Name: "main"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Branches.ListBranchesPage(context.Background(), "core", "go-forge", ListOptions{Page: 1, PageSize: 30})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 1 || len(page.Items) != 1 || page.Items[0].Name != "main" {
		t.Fatalf("got %#v", page)
	}
}

func TestOrgService_ListOrgTeamsPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/teams" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("got limit=%q, want %q", got, "10")
		}
		w.Header().Set("X-Total-Count", "11")
		json.NewEncoder(w).Encode([]types.Team{{ID: 1, Name: "maintainers"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Orgs.ListOrgTeamsPage(context.Background(), "core", ListOptions{Page: 2, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 2 || page.TotalCount != 11 || len(page.Items) != 1 || page.Items[0].Name != "maintainers" {
		t.Fatalf("got %#v", page)
	}
}

func TestWebhookService_ListRepoHooksPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "5" {
			t.Errorf("got limit=%q, want %q", got, "5")
		}
		w.Header().Set("X-Total-Count", "6")
		json.NewEncoder(w).Encode([]types.Hook{{ID: 1, Type: "forgejo"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Webhooks.ListRepoHooksPage(context.Background(), "core", "go-forge", ListOptions{Page: 1, PageSize: 5})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 6 || len(page.Items) != 1 || page.Items[0].Type != "forgejo" {
		t.Fatalf("got %#v", page)
	}
}

func TestReleaseService_ListReleasesPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("draft"); got != "true" {
			t.Errorf("got draft=%q, want %q", got, "true")
		}
		if got := r.URL.Query().Get("q"); got != "v1" {
			t.Errorf("got q=%q, want %q", got, "v1")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("got limit=%q, want %q", got, "10")
		}
		w.Header().Set("X-Total-Count", "12")
		json.NewEncoder(w).Encode([]types.Release{{ID: 1, TagName: "v1.0.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Releases.ListReleasesPage(context.Background(), "core", "go-forge", ListOptions{Page: 2, PageSize: 10}, ReleaseListOptions{Draft: true, Query: "v1"})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 2 || page.TotalCount != 12 || len(page.Items) != 1 || page.Items[0].TagName != "v1.0.0" {
		t.Fatalf("got %#v", page)
	}
}

func TestMilestoneService_ListMilestonesPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want %q", got, "open")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "10" {
			t.Errorf("got limit=%q, want %q", got, "10")
		}
		w.Header().Set("X-Total-Count", "11")
		json.NewEncoder(w).Encode([]types.Milestone{{ID: 2, Title: "Sprint"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Milestones.ListMilestonesPage(context.Background(), "core", "go-forge", ListOptions{Page: 2, PageSize: 10}, MilestoneListOptions{State: "open"})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 2 || page.TotalCount != 11 || len(page.Items) != 1 || page.Items[0].Title != "Sprint" {
		t.Fatalf("got %#v", page)
	}
}

func TestLabelService_ListLabelsPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "8" {
			t.Errorf("got limit=%q, want %q", got, "8")
		}
		w.Header().Set("X-Total-Count", "8")
		json.NewEncoder(w).Encode([]types.Label{{ID: 1, Name: "bug"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Labels.ListLabelsPage(context.Background(), "core", "go-forge", ListOptions{Page: 1, PageSize: 8})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 8 || len(page.Items) != 1 || page.Items[0].Name != "bug" {
		t.Fatalf("got %#v", page)
	}
}
