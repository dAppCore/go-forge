package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
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

func TestPullService_ListFiltered_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		want := map[string]string{
			"state":     "open",
			"sort":      "priority",
			"milestone": "7",
			"poster":    "alice",
			"page":      "1",
			"limit":     "50",
		}
		for key, wantValue := range want {
			if got := r.URL.Query().Get(key); got != wantValue {
				t.Errorf("got %s=%q, want %q", key, got, wantValue)
			}
		}
		if got := r.URL.Query()["labels"]; !reflect.DeepEqual(got, []string{"1", "2"}) {
			t.Errorf("got labels=%v, want %v", got, []string{"1", "2"})
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PullRequest{{ID: 1, Title: "add feature"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	prs, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge", PullListOptions{
		State:     "open",
		Sort:      "priority",
		Milestone: 7,
		Labels:    []int64{1, 2},
		Poster:    "alice",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 1 || prs[0].Title != "add feature" {
		t.Fatalf("got %#v", prs)
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

func TestPullService_GetByBaseHead_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/main/feature" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.PullRequest{Index: 7, Title: "Add feature"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.GetByBaseHead(context.Background(), "core", "go-forge", "main", "feature")
	if err != nil {
		t.Fatal(err)
	}
	if pr.Index != 7 || pr.Title != "Add feature" {
		t.Fatalf("got %+v", pr)
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

func TestPullService_MergePullRequest_CompatMergeStyle_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/merge" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if got := body["Do"]; got != "squash" {
			t.Fatalf("got Do=%v, want squash", got)
		}
		if got := body["MergeMessageField"]; got != "PR: Add feature" {
			t.Fatalf("got MergeMessageField=%v, want %q", got, "PR: Add feature")
		}
		if _, ok := body["MergeStyle"]; ok {
			t.Fatalf("did not expect MergeStyle in request body: %#v", body)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Pulls.MergePullRequest(context.Background(), "core", "go-forge", 7, &types.MergePullRequestOption{
		MergeMessageField: "PR: Add feature",
		MergeStyle:        "squash",
	})
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

func TestPulls_PullListOptions_String_Good(t *core.T) {
	subject := (*PullListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullListOptions_String_Bad(t *core.T) {
	subject := (*PullListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullListOptions_String_Ugly(t *core.T) {
	subject := (*PullListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullListOptions_GoString_Good(t *core.T) {
	subject := (*PullListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullListOptions_GoString_Bad(t *core.T) {
	subject := (*PullListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullListOptions_GoString_Ugly(t *core.T) {
	subject := (*PullListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequestsPage_Good(t *core.T) {
	subject := (*PullService).ListPullRequestsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequestsPage_Bad(t *core.T) {
	subject := (*PullService).ListPullRequestsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequestsPage_Ugly(t *core.T) {
	subject := (*PullService).ListPullRequestsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequests_Good(t *core.T) {
	subject := (*PullService).ListPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequests_Bad(t *core.T) {
	subject := (*PullService).ListPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullRequests_Ugly(t *core.T) {
	subject := (*PullService).ListPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullRequests_Good(t *core.T) {
	subject := (*PullService).IterPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullRequests_Bad(t *core.T) {
	subject := (*PullService).IterPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullRequests_Ugly(t *core.T) {
	subject := (*PullService).IterPullRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreatePullRequest_Good(t *core.T) {
	subject := (*PullService).CreatePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreatePullRequest_Bad(t *core.T) {
	subject := (*PullService).CreatePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreatePullRequest_Ugly(t *core.T) {
	subject := (*PullService).CreatePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullRequest_Good(t *core.T) {
	subject := (*PullService).GetPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullRequest_Bad(t *core.T) {
	subject := (*PullService).GetPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullRequest_Ugly(t *core.T) {
	subject := (*PullService).GetPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_EditPullRequest_Good(t *core.T) {
	subject := (*PullService).EditPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_EditPullRequest_Bad(t *core.T) {
	subject := (*PullService).EditPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_EditPullRequest_Ugly(t *core.T) {
	subject := (*PullService).EditPullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullRequest_Good(t *core.T) {
	subject := (*PullService).DeletePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullRequest_Bad(t *core.T) {
	subject := (*PullService).DeletePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullRequest_Ugly(t *core.T) {
	subject := (*PullService).DeletePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Merge_Good(t *core.T) {
	subject := (*PullService).Merge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Merge_Bad(t *core.T) {
	subject := (*PullService).Merge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Merge_Ugly(t *core.T) {
	subject := (*PullService).Merge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_MergePullRequest_Good(t *core.T) {
	subject := (*PullService).MergePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_MergePullRequest_Bad(t *core.T) {
	subject := (*PullService).MergePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_MergePullRequest_Ugly(t *core.T) {
	subject := (*PullService).MergePullRequest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelScheduledAutoMerge_Good(t *core.T) {
	subject := (*PullService).CancelScheduledAutoMerge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelScheduledAutoMerge_Bad(t *core.T) {
	subject := (*PullService).CancelScheduledAutoMerge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelScheduledAutoMerge_Ugly(t *core.T) {
	subject := (*PullService).CancelScheduledAutoMerge
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Update_Good(t *core.T) {
	subject := (*PullService).Update
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Update_Bad(t *core.T) {
	subject := (*PullService).Update
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_Update_Ugly(t *core.T) {
	subject := (*PullService).Update
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetDiffOrPatch_Good(t *core.T) {
	subject := (*PullService).GetDiffOrPatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetDiffOrPatch_Bad(t *core.T) {
	subject := (*PullService).GetDiffOrPatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetDiffOrPatch_Ugly(t *core.T) {
	subject := (*PullService).GetDiffOrPatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListCommits_Good(t *core.T) {
	subject := (*PullService).ListCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListCommits_Bad(t *core.T) {
	subject := (*PullService).ListCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListCommits_Ugly(t *core.T) {
	subject := (*PullService).ListCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterCommits_Good(t *core.T) {
	subject := (*PullService).IterCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterCommits_Bad(t *core.T) {
	subject := (*PullService).IterCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterCommits_Ugly(t *core.T) {
	subject := (*PullService).IterCommits
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviews_Good(t *core.T) {
	subject := (*PullService).ListReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviews_Bad(t *core.T) {
	subject := (*PullService).ListReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviews_Ugly(t *core.T) {
	subject := (*PullService).ListReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullReviews_Good(t *core.T) {
	subject := (*PullService).ListPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullReviews_Bad(t *core.T) {
	subject := (*PullService).ListPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListPullReviews_Ugly(t *core.T) {
	subject := (*PullService).ListPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviews_Good(t *core.T) {
	subject := (*PullService).IterReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviews_Bad(t *core.T) {
	subject := (*PullService).IterReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviews_Ugly(t *core.T) {
	subject := (*PullService).IterReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullReviews_Good(t *core.T) {
	subject := (*PullService).IterPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullReviews_Bad(t *core.T) {
	subject := (*PullService).IterPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterPullReviews_Ugly(t *core.T) {
	subject := (*PullService).IterPullReviews
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListFiles_Good(t *core.T) {
	subject := (*PullService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListFiles_Bad(t *core.T) {
	subject := (*PullService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListFiles_Ugly(t *core.T) {
	subject := (*PullService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterFiles_Good(t *core.T) {
	subject := (*PullService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterFiles_Bad(t *core.T) {
	subject := (*PullService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterFiles_Ugly(t *core.T) {
	subject := (*PullService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetByBaseHead_Good(t *core.T) {
	subject := (*PullService).GetByBaseHead
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetByBaseHead_Bad(t *core.T) {
	subject := (*PullService).GetByBaseHead
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetByBaseHead_Ugly(t *core.T) {
	subject := (*PullService).GetByBaseHead
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewers_Good(t *core.T) {
	subject := (*PullService).ListReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewers_Bad(t *core.T) {
	subject := (*PullService).ListReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewers_Ugly(t *core.T) {
	subject := (*PullService).ListReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewers_Good(t *core.T) {
	subject := (*PullService).IterReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewers_Bad(t *core.T) {
	subject := (*PullService).IterReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewers_Ugly(t *core.T) {
	subject := (*PullService).IterReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_RequestReviewers_Good(t *core.T) {
	subject := (*PullService).RequestReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_RequestReviewers_Bad(t *core.T) {
	subject := (*PullService).RequestReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_RequestReviewers_Ugly(t *core.T) {
	subject := (*PullService).RequestReviewers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelReviewRequests_Good(t *core.T) {
	subject := (*PullService).CancelReviewRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelReviewRequests_Bad(t *core.T) {
	subject := (*PullService).CancelReviewRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CancelReviewRequests_Ugly(t *core.T) {
	subject := (*PullService).CancelReviewRequests
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_SubmitReview_Good(t *core.T) {
	subject := (*PullService).SubmitReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_SubmitReview_Bad(t *core.T) {
	subject := (*PullService).SubmitReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_SubmitReview_Ugly(t *core.T) {
	subject := (*PullService).SubmitReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReview_Good(t *core.T) {
	subject := (*PullService).GetReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReview_Bad(t *core.T) {
	subject := (*PullService).GetReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReview_Ugly(t *core.T) {
	subject := (*PullService).GetReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullReview_Good(t *core.T) {
	subject := (*PullService).GetPullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullReview_Bad(t *core.T) {
	subject := (*PullService).GetPullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetPullReview_Ugly(t *core.T) {
	subject := (*PullService).GetPullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReview_Good(t *core.T) {
	subject := (*PullService).DeleteReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReview_Bad(t *core.T) {
	subject := (*PullService).DeleteReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReview_Ugly(t *core.T) {
	subject := (*PullService).DeleteReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullReview_Good(t *core.T) {
	subject := (*PullService).DeletePullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullReview_Bad(t *core.T) {
	subject := (*PullService).DeletePullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeletePullReview_Ugly(t *core.T) {
	subject := (*PullService).DeletePullReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewComments_Good(t *core.T) {
	subject := (*PullService).ListReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewComments_Bad(t *core.T) {
	subject := (*PullService).ListReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_ListReviewComments_Ugly(t *core.T) {
	subject := (*PullService).ListReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewComments_Good(t *core.T) {
	subject := (*PullService).IterReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewComments_Bad(t *core.T) {
	subject := (*PullService).IterReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_IterReviewComments_Ugly(t *core.T) {
	subject := (*PullService).IterReviewComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReviewComment_Good(t *core.T) {
	subject := (*PullService).GetReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReviewComment_Bad(t *core.T) {
	subject := (*PullService).GetReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_GetReviewComment_Ugly(t *core.T) {
	subject := (*PullService).GetReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreateReviewComment_Good(t *core.T) {
	subject := (*PullService).CreateReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreateReviewComment_Bad(t *core.T) {
	subject := (*PullService).CreateReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_CreateReviewComment_Ugly(t *core.T) {
	subject := (*PullService).CreateReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReviewComment_Good(t *core.T) {
	subject := (*PullService).DeleteReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReviewComment_Bad(t *core.T) {
	subject := (*PullService).DeleteReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DeleteReviewComment_Ugly(t *core.T) {
	subject := (*PullService).DeleteReviewComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DismissReview_Good(t *core.T) {
	subject := (*PullService).DismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DismissReview_Bad(t *core.T) {
	subject := (*PullService).DismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_DismissReview_Ugly(t *core.T) {
	subject := (*PullService).DismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_UndismissReview_Good(t *core.T) {
	subject := (*PullService).UndismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_UndismissReview_Bad(t *core.T) {
	subject := (*PullService).UndismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPulls_PullService_UndismissReview_Ugly(t *core.T) {
	subject := (*PullService).UndismissReview
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
