package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"dappco.re/go/core/forge/types"
)

func TestIssueService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Issue{
			{ID: 1, Title: "bug report"},
			{ID: 2, Title: "feature request"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Issues.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].Title != "bug report" {
		t.Errorf("got title=%q, want %q", result.Items[0].Title, "bug report")
	}
}

func TestIssueService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Title: "bug report", Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "bug report" {
		t.Errorf("got title=%q", issue.Title)
	}
}

func TestIssueService_Create_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreateIssueOption
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.Create(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.CreateIssueOption{
		Title: "new issue",
		Body:  "description here",
	})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "new issue" {
		t.Errorf("got title=%q", issue.Title)
	}
}

func TestIssueService_Update_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditIssueOption
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.Update(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"}, &types.EditIssueOption{
		Title: "updated issue",
	})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "updated issue" {
		t.Errorf("got title=%q", issue.Title)
	}
}

func TestIssueService_Delete_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.Delete(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"}); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_SearchIssuesPage_Good(t *testing.T) {
	since := time.Date(2026, time.March, 1, 12, 30, 0, 0, time.UTC)
	before := time.Date(2026, time.March, 2, 12, 30, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/issues/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		want := map[string]string{
			"state":            "open",
			"labels":           "bug,help wanted",
			"milestones":       "v1.0",
			"q":                "panic",
			"priority_repo_id": "42",
			"type":             "issues",
			"since":            since.Format(time.RFC3339),
			"before":           before.Format(time.RFC3339),
			"assigned":         "true",
			"created":          "true",
			"mentioned":        "true",
			"review_requested": "true",
			"reviewed":         "true",
			"owner":            "core",
			"team":             "platform",
			"page":             "2",
			"limit":            "25",
		}
		for key, wantValue := range want {
			if got := r.URL.Query().Get(key); got != wantValue {
				t.Errorf("got %s=%q, want %q", key, got, wantValue)
			}
		}
		w.Header().Set("X-Total-Count", "100")
		json.NewEncoder(w).Encode([]types.Issue{
			{ID: 1, Title: "panic in parser"},
			{ID: 2, Title: "panic in generator"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Issues.SearchIssuesPage(context.Background(), SearchIssuesOptions{
		State:           "open",
		Labels:          "bug,help wanted",
		Milestones:      "v1.0",
		Query:           "panic",
		PriorityRepoID:  42,
		Type:            "issues",
		Since:           &since,
		Before:          &before,
		Assigned:        true,
		Created:         true,
		Mentioned:       true,
		ReviewRequested: true,
		Reviewed:        true,
		Owner:           "core",
		Team:            "platform",
	}, ListOptions{Page: 2, Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(page.Items), 2; got != want {
		t.Fatalf("got %d items, want %d", got, want)
	}
	if !page.HasMore {
		t.Fatalf("expected HasMore to be true")
	}
	if page.TotalCount != 100 {
		t.Fatalf("got total count %d, want 100", page.TotalCount)
	}
	if page.Items[0].Title != "panic in parser" {
		t.Fatalf("got first title %q", page.Items[0].Title)
	}
}

func TestIssueService_SearchIssues_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/issues/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "panic" {
			t.Errorf("got q=%q, want %q", got, "panic")
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Issue{
			{ID: 1, Title: "panic in parser"},
			{ID: 2, Title: "panic in generator"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.SearchIssues(context.Background(), SearchIssuesOptions{Query: "panic"})
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(issues), 2; got != want {
		t.Fatalf("got %d items, want %d", got, want)
	}
	if issues[1].Title != "panic in generator" {
		t.Fatalf("got second title %q", issues[1].Title)
	}
}

func TestIssueService_IterSearchIssues_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/issues/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "panic" {
			t.Errorf("got q=%q, want %q", got, "panic")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 1, Title: "panic in parser"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.Issue
	for issue, err := range f.Issues.IterSearchIssues(context.Background(), SearchIssuesOptions{Query: "panic"}) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, issue)
	}
	if got, want := len(seen), 1; got != want {
		t.Fatalf("got %d items, want %d", got, want)
	}
	if seen[0].Title != "panic in parser" {
		t.Fatalf("got title %q", seen[0].Title)
	}
}

func TestIssueService_CreateComment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/comments" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreateIssueCommentOption
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Comment{ID: 7, Body: body.Body})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	comment, err := f.Issues.CreateComment(context.Background(), "core", "go-forge", 1, "first!")
	if err != nil {
		t.Fatal(err)
	}
	if comment.Body != "first!" {
		t.Errorf("got body=%q", comment.Body)
	}
}

func TestIssueService_ListTimeline_Good(t *testing.T) {
	since := time.Date(2026, time.March, 1, 12, 30, 0, 0, time.UTC)
	before := time.Date(2026, time.March, 2, 12, 30, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/timeline" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("since"); got != since.Format(time.RFC3339) {
			t.Errorf("got since=%q, want %q", got, since.Format(time.RFC3339))
		}
		if got := r.URL.Query().Get("before"); got != before.Format(time.RFC3339) {
			t.Errorf("got before=%q, want %q", got, before.Format(time.RFC3339))
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.TimelineComment{
			{ID: 11, Type: "comment", Body: "first"},
			{ID: 12, Type: "state_change", Body: "second"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	events, err := f.Issues.ListTimeline(context.Background(), "core", "go-forge", 1, &since, &before)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(events, []types.TimelineComment{
		{ID: 11, Type: "comment", Body: "first"},
		{ID: 12, Type: "state_change", Body: "second"},
	}) {
		t.Fatalf("got %#v", events)
	}
}

func TestIssueService_IterTimeline_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/timeline" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.TimelineComment{{ID: 11, Type: "comment", Body: "first"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.TimelineComment
	for event, err := range f.Issues.IterTimeline(context.Background(), "core", "go-forge", 1, nil, nil) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, event)
	}
	if !reflect.DeepEqual(seen, []types.TimelineComment{{ID: 11, Type: "comment", Body: "first"}}) {
		t.Fatalf("got %#v", seen)
	}
}

func TestIssueService_ListSubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{
			{ID: 1, UserName: "alice"},
			{ID: 2, UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Issues.ListSubscriptions(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(users, []types.User{{ID: 1, UserName: "alice"}, {ID: 2, UserName: "bob"}}) {
		t.Fatalf("got %#v", users)
	}
}

func TestIssueService_IterSubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.User{{ID: 1, UserName: "alice"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.User
	for user, err := range f.Issues.IterSubscriptions(context.Background(), "core", "go-forge", 1) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, user)
	}
	if !reflect.DeepEqual(seen, []types.User{{ID: 1, UserName: "alice"}}) {
		t.Fatalf("got %#v", seen)
	}
}

func TestIssueService_CheckSubscription_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/subscriptions/check" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.WatchInfo{Subscribed: true, Ignored: false})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Issues.CheckSubscription(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Subscribed || result.Ignored {
		t.Fatalf("got %#v", result)
	}
}

func TestIssueService_SubscribeUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/subscriptions/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.SubscribeUser(context.Background(), "core", "go-forge", 1, "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_UnsubscribeUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/subscriptions/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.UnsubscribeUser(context.Background(), "core", "go-forge", 1, "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_ListDependencies_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/dependencies" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 11, Index: 11, Title: "blocking issue"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListDependencies(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(issues, []types.Issue{{ID: 11, Index: 11, Title: "blocking issue"}}) {
		t.Fatalf("got %#v", issues)
	}
}

func TestIssueService_AddDependency_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/dependencies" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.IssueMeta
		json.NewDecoder(r.Body).Decode(&body)
		if body.Owner != "core" || body.Name != "go-forge" || body.Index != 2 {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Issue{ID: 11, Index: 11, Title: "blocking issue"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.AddDependency(context.Background(), "core", "go-forge", 1, types.IssueMeta{Owner: "core", Name: "go-forge", Index: 2}); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_RemoveDependency_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/dependencies" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.IssueMeta
		json.NewDecoder(r.Body).Decode(&body)
		if body.Owner != "core" || body.Name != "go-forge" || body.Index != 2 {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.Issue{ID: 11, Index: 11, Title: "blocking issue"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.RemoveDependency(context.Background(), "core", "go-forge", 1, types.IssueMeta{Owner: "core", Name: "go-forge", Index: 2}); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_ListBlocks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/blocks" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 22, Index: 22, Title: "blocked issue"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListBlocks(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(issues, []types.Issue{{ID: 22, Index: 22, Title: "blocked issue"}}) {
		t.Fatalf("got %#v", issues)
	}
}

func TestIssueService_AddBlock_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/blocks" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.IssueMeta
		json.NewDecoder(r.Body).Decode(&body)
		if body.Owner != "core" || body.Name != "go-forge" || body.Index != 3 {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Issue{ID: 22, Index: 22, Title: "blocked issue"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.AddBlock(context.Background(), "core", "go-forge", 1, types.IssueMeta{Owner: "core", Name: "go-forge", Index: 3}); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_RemoveBlock_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/blocks" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.IssueMeta
		json.NewDecoder(r.Body).Decode(&body)
		if body.Owner != "core" || body.Name != "go-forge" || body.Index != 3 {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.Issue{ID: 22, Index: 22, Title: "blocked issue"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.RemoveBlock(context.Background(), "core", "go-forge", 1, types.IssueMeta{Owner: "core", Name: "go-forge", Index: 3}); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_Pin_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/pin" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Issues.Pin(context.Background(), "core", "go-forge", 42)
	if err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_MovePin_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/pin/3" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.MovePin(context.Background(), "core", "go-forge", 42, 3); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_DeleteStopwatch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/stopwatch/delete" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.DeleteStopwatch(context.Background(), "core", "go-forge", 42); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_ListTimes_Good(t *testing.T) {
	since := time.Date(2026, time.March, 3, 9, 15, 0, 0, time.UTC)
	before := time.Date(2026, time.March, 4, 9, 15, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/times" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("user"); got != "alice" {
			t.Errorf("got user=%q, want %q", got, "alice")
		}
		if got := r.URL.Query().Get("since"); got != since.Format(time.RFC3339) {
			t.Errorf("got since=%q, want %q", got, since.Format(time.RFC3339))
		}
		if got := r.URL.Query().Get("before"); got != before.Format(time.RFC3339) {
			t.Errorf("got before=%q, want %q", got, before.Format(time.RFC3339))
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.TrackedTime{
			{ID: 11, Time: 30, UserName: "alice"},
			{ID: 12, Time: 90, UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	times, err := f.Issues.ListTimes(context.Background(), "core", "go-forge", 42, "alice", &since, &before)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(times, []types.TrackedTime{
		{ID: 11, Time: 30, UserName: "alice"},
		{ID: 12, Time: 90, UserName: "bob"},
	}) {
		t.Fatalf("got %#v", times)
	}
}

func TestIssueService_AddTime_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/times" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.AddTimeOption
		json.NewDecoder(r.Body).Decode(&body)
		if body.Time != 180 || body.User != "alice" {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.TrackedTime{ID: 99, Time: body.Time, UserName: body.User})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	got, err := f.Issues.AddTime(context.Background(), "core", "go-forge", 42, &types.AddTimeOption{
		Time: 180,
		User: "alice",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 99 || got.Time != 180 || got.UserName != "alice" {
		t.Fatalf("got %#v", got)
	}
}

func TestIssueService_ResetTime_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/times" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.ResetTime(context.Background(), "core", "go-forge", 42); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_List_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Issues.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList); err == nil {
		t.Fatal("expected error")
	}
}

func TestIssueService_ListIgnoresIndexParam_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "0")
		json.NewEncoder(w).Encode([]types.Issue{})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Issues.List(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "99"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 0 {
		t.Errorf("got %d items, want 0", len(result.Items))
	}
}
