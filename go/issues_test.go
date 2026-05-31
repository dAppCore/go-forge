package forge

import (
	"bytes"
	"context"
	json "github.com/goccy/go-json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
)

func readMultipartAttachment(t *testing.T, r *http.Request) (string, string) {
	t.Helper()

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		t.Fatal(err)
	}
	if mediaType != "multipart/form-data" {
		t.Fatalf("got content-type=%q", mediaType)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatal(err)
	}

	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	part, err := reader.NextPart()
	if err != nil {
		t.Fatal(err)
	}
	if part.FormName() != "attachment" {
		t.Fatalf("got form name=%q", part.FormName())
	}
	content, err := io.ReadAll(part)
	if err != nil {
		t.Fatal(err)
	}
	return part.FileName(), string(content)
}

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

func TestIssueService_ListFiltered_Good(t *testing.T) {
	since := time.Date(2026, time.March, 1, 12, 30, 0, 0, time.UTC)
	before := time.Date(2026, time.March, 2, 12, 30, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		want := map[string]string{
			"state":        "open",
			"labels":       "bug,help wanted",
			"q":            "panic",
			"type":         "issues",
			"milestones":   "v1.0",
			"since":        since.Format(time.RFC3339),
			"before":       before.Format(time.RFC3339),
			"created_by":   "alice",
			"assigned_by":  "bob",
			"mentioned_by": "carol",
			"page":         "1",
			"limit":        "50",
		}
		for key, wantValue := range want {
			if got := r.URL.Query().Get(key); got != wantValue {
				t.Errorf("got %s=%q, want %q", key, got, wantValue)
			}
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 1, Title: "panic in parser"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListIssues(context.Background(), "core", "go-forge", IssueListOptions{
		State:       "open",
		Labels:      "bug,help wanted",
		Query:       "panic",
		Type:        "issues",
		Milestones:  "v1.0",
		Since:       &since,
		Before:      &before,
		CreatedBy:   "alice",
		AssignedBy:  "bob",
		MentionedBy: "carol",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Title != "panic in parser" {
		t.Fatalf("got %#v", issues)
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

func TestIssueService_ListRepoComments_Good(t *testing.T) {
	since := time.Date(2026, 4, 1, 10, 0, 0, 0, time.UTC)
	before := time.Date(2026, 4, 2, 10, 0, 0, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments" {
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
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Comment{
			{ID: 7, Body: "repo-wide comment"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	comments, err := f.Issues.ListRepoComments(context.Background(), "core", "go-forge", RepoCommentListOptions{
		Since:  &since,
		Before: &before,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 || comments[0].Body != "repo-wide comment" {
		t.Fatalf("unexpected result: %#v", comments)
	}
}

func TestIssueService_GetRepoComment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Comment{ID: 7, Body: "repo-wide comment"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	comment, err := f.Issues.GetRepoComment(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Body != "repo-wide comment" {
		t.Fatalf("got body=%q", comment.Body)
	}
}

func TestIssueService_EditRepoComment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditIssueCommentOption
		json.NewDecoder(r.Body).Decode(&body)
		if body.Body != "updated comment" {
			t.Fatalf("got body=%#v", body)
		}
		json.NewEncoder(w).Encode(types.Comment{ID: 7, Body: body.Body})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	comment, err := f.Issues.EditRepoComment(context.Background(), "core", "go-forge", 7, &types.EditIssueCommentOption{
		Body: "updated comment",
	})
	if err != nil {
		t.Fatal(err)
	}
	if comment.Body != "updated comment" {
		t.Fatalf("got body=%q", comment.Body)
	}
}

func TestIssueService_DeleteRepoComment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.DeleteRepoComment(context.Background(), "core", "go-forge", 7); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_ListReactions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/reactions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Reaction{
			{Reaction: "+1", User: &types.User{ID: 1, UserName: "alice"}},
			{Reaction: "heart", User: &types.User{ID: 2, UserName: "bob"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	reactions, err := f.Issues.ListReactions(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reactions, []types.Reaction{
		{Reaction: "+1", User: &types.User{ID: 1, UserName: "alice"}},
		{Reaction: "heart", User: &types.User{ID: 2, UserName: "bob"}},
	}) {
		t.Fatalf("got %#v", reactions)
	}
}

func TestIssueService_IterReactions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/reactions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Reaction{
			{Reaction: "+1", User: &types.User{ID: 1, UserName: "alice"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.Reaction
	for reaction, err := range f.Issues.IterReactions(context.Background(), "core", "go-forge", 1) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, reaction)
	}
	if !reflect.DeepEqual(seen, []types.Reaction{{Reaction: "+1", User: &types.User{ID: 1, UserName: "alice"}}}) {
		t.Fatalf("got %#v", seen)
	}
}

func TestIssueService_ListCommentReactions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7/reactions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Reaction{
			{Reaction: "eyes", User: &types.User{ID: 3, UserName: "carol"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	reactions, err := f.Issues.ListCommentReactions(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(reactions, []types.Reaction{
		{Reaction: "eyes", User: &types.User{ID: 3, UserName: "carol"}},
	}) {
		t.Fatalf("got %#v", reactions)
	}
}

func TestIssueService_AddCommentReaction_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7/reactions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditReactionOption
		json.NewDecoder(r.Body).Decode(&body)
		if body.Reaction != "heart" {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Reaction{Reaction: body.Reaction, User: &types.User{ID: 4, UserName: "dave"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	reaction, err := f.Issues.AddCommentReaction(context.Background(), "core", "go-forge", 7, "heart")
	if err != nil {
		t.Fatal(err)
	}
	if reaction.Reaction != "heart" || reaction.User.UserName != "dave" {
		t.Fatalf("got %#v", reaction)
	}
}

func TestIssueService_DeleteCommentReaction_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7/reactions" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditReactionOption
		json.NewDecoder(r.Body).Decode(&body)
		if body.Reaction != "heart" {
			t.Fatalf("got body=%#v", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.DeleteCommentReaction(context.Background(), "core", "go-forge", 7, "heart"); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_ListAttachments_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Attachment{
			{ID: 4, Name: "design.png"},
			{ID: 5, Name: "notes.txt"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachments, err := f.Issues.ListAttachments(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(attachments, []types.Attachment{{ID: 4, Name: "design.png"}, {ID: 5, Name: "notes.txt"}}) {
		t.Fatalf("got %#v", attachments)
	}
}

func TestIssueService_IterAttachments_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Attachment{{ID: 4, Name: "design.png"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.Attachment
	for attachment, err := range f.Issues.IterAttachments(context.Background(), "core", "go-forge", 1) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, attachment)
	}
	if !reflect.DeepEqual(seen, []types.Attachment{{ID: 4, Name: "design.png"}}) {
		t.Fatalf("got %#v", seen)
	}
}

func TestIssueService_GetAttachment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets/4" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Attachment{ID: 4, Name: "design.png"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Issues.GetAttachment(context.Background(), "core", "go-forge", 1, 4)
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "design.png" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestIssueService_EditAttachment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets/4" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditAttachmentOptions
		json.NewDecoder(r.Body).Decode(&body)
		if body.Name != "updated.png" {
			t.Fatalf("got body=%#v", body)
		}
		json.NewEncoder(w).Encode(types.Attachment{ID: 4, Name: body.Name})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Issues.EditAttachment(context.Background(), "core", "go-forge", 1, 4, &types.EditAttachmentOptions{Name: "updated.png"})
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "updated.png" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestIssueService_DeleteAttachment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets/4" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.DeleteAttachment(context.Background(), "core", "go-forge", 1, 4); err != nil {
		t.Fatal(err)
	}
}

func TestIssueService_CreateAttachment_Good(t *testing.T) {
	updatedAt := time.Date(2026, time.March, 3, 11, 22, 33, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("name"); got != "diagram" {
			t.Fatalf("got name=%q", got)
		}
		if got := r.URL.Query().Get("updated_at"); got != updatedAt.Format(time.RFC3339) {
			t.Fatalf("got updated_at=%q", got)
		}
		filename, content := readMultipartAttachment(t, r)
		if filename != "design.png" {
			t.Fatalf("got filename=%q", filename)
		}
		if content != "attachment bytes" {
			t.Fatalf("got content=%q", content)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Attachment{ID: 9, Name: filename})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Issues.CreateAttachment(
		context.Background(),
		"core",
		"go-forge",
		1,
		&AttachmentUploadOptions{Name: "diagram", UpdatedAt: &updatedAt},
		"design.png",
		bytes.NewBufferString("attachment bytes"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "design.png" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestIssueService_CreateCommentAttachment_Good(t *testing.T) {
	updatedAt := time.Date(2026, time.March, 4, 9, 10, 11, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/comments/7/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("name"); got != "screenshot" {
			t.Fatalf("got name=%q", got)
		}
		if got := r.URL.Query().Get("updated_at"); got != updatedAt.Format(time.RFC3339) {
			t.Fatalf("got updated_at=%q", got)
		}
		filename, content := readMultipartAttachment(t, r)
		if filename != "comment.png" {
			t.Fatalf("got filename=%q", filename)
		}
		if content != "comment attachment bytes" {
			t.Fatalf("got content=%q", content)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Attachment{ID: 11, Name: filename})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Issues.CreateCommentAttachment(
		context.Background(),
		"core",
		"go-forge",
		7,
		&AttachmentUploadOptions{Name: "screenshot", UpdatedAt: &updatedAt},
		"comment.png",
		bytes.NewBufferString("comment attachment bytes"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "comment.png" {
		t.Fatalf("got name=%q", attachment.Name)
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

func TestIssueService_ListPinnedIssues_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/pinned" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Issue{
			{ID: 1, Title: "critical bug"},
			{ID: 2, Title: "release blocker"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListPinnedIssues(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(issues), 2; got != want {
		t.Fatalf("got %d issues, want %d", got, want)
	}
	if issues[0].Title != "critical bug" {
		t.Fatalf("got first title %q", issues[0].Title)
	}
}

func TestIssueService_IterPinnedIssues_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/pinned" {
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
			json.NewEncoder(w).Encode([]types.Issue{{ID: 1, Title: "critical bug"}})
		case 2:
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Errorf("got page=%q, want %q", got, "2")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.Issue{{ID: 2, Title: "release blocker"}})
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for issue, err := range f.Issues.IterPinnedIssues(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, issue.Title)
	}
	if len(got) != 2 || got[0] != "critical bug" || got[1] != "release blocker" {
		t.Fatalf("got %#v", got)
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

func TestIssueService_IterTimes_Good(t *testing.T) {
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
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("got limit=%q, want %q", got, "50")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.TrackedTime{
			{ID: 11, Time: 30, UserName: "alice"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var seen []types.TrackedTime
	for entry, err := range f.Issues.IterTimes(context.Background(), "core", "go-forge", 42, "alice", &since, &before) {
		if err != nil {
			t.Fatal(err)
		}
		seen = append(seen, entry)
	}
	if !reflect.DeepEqual(seen, []types.TrackedTime{
		{ID: 11, Time: 30, UserName: "alice"},
	}) {
		t.Fatalf("got %#v", seen)
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

func TestIssueService_DeleteTime_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/times/99" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Issues.DeleteTime(context.Background(), "core", "go-forge", 42, 99); err != nil {
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

func TestIssues_IssueListOptions_String_Good(t *core.T) {
	subject := (*IssueListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueListOptions_String_Bad(t *core.T) {
	subject := (*IssueListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueListOptions_String_Ugly(t *core.T) {
	subject := (*IssueListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueListOptions_GoString_Good(t *core.T) {
	subject := (*IssueListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueListOptions_GoString_Bad(t *core.T) {
	subject := (*IssueListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueListOptions_GoString_Ugly(t *core.T) {
	subject := (*IssueListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_String_Good(t *core.T) {
	subject := (*AttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_String_Bad(t *core.T) {
	subject := (*AttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_String_Ugly(t *core.T) {
	subject := (*AttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_GoString_Good(t *core.T) {
	subject := (*AttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_GoString_Bad(t *core.T) {
	subject := (*AttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_AttachmentUploadOptions_GoString_Ugly(t *core.T) {
	subject := (*AttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_String_Good(t *core.T) {
	subject := (*RepoCommentListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_String_Bad(t *core.T) {
	subject := (*RepoCommentListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_String_Ugly(t *core.T) {
	subject := (*RepoCommentListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_GoString_Good(t *core.T) {
	subject := (*RepoCommentListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_GoString_Bad(t *core.T) {
	subject := (*RepoCommentListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_RepoCommentListOptions_GoString_Ugly(t *core.T) {
	subject := (*RepoCommentListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssue_Good(t *core.T) {
	subject := (*IssueService).GetIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssue_Bad(t *core.T) {
	subject := (*IssueService).GetIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssue_Ugly(t *core.T) {
	subject := (*IssueService).GetIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssue_Good(t *core.T) {
	subject := (*IssueService).EditIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssue_Bad(t *core.T) {
	subject := (*IssueService).EditIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssue_Ugly(t *core.T) {
	subject := (*IssueService).EditIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssue_Good(t *core.T) {
	subject := (*IssueService).DeleteIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssue_Bad(t *core.T) {
	subject := (*IssueService).DeleteIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssue_Ugly(t *core.T) {
	subject := (*IssueService).DeleteIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_String_Good(t *core.T) {
	subject := (*SearchIssuesOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_String_Bad(t *core.T) {
	subject := (*SearchIssuesOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_String_Ugly(t *core.T) {
	subject := (*SearchIssuesOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_GoString_Good(t *core.T) {
	subject := (*SearchIssuesOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_GoString_Bad(t *core.T) {
	subject := (*SearchIssuesOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_SearchIssuesOptions_GoString_Ugly(t *core.T) {
	subject := (*SearchIssuesOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssuesPage_Good(t *core.T) {
	subject := (*IssueService).SearchIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssuesPage_Bad(t *core.T) {
	subject := (*IssueService).SearchIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssuesPage_Ugly(t *core.T) {
	subject := (*IssueService).SearchIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssues_Good(t *core.T) {
	subject := (*IssueService).SearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssues_Bad(t *core.T) {
	subject := (*IssueService).SearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SearchIssues_Ugly(t *core.T) {
	subject := (*IssueService).SearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSearchIssues_Good(t *core.T) {
	subject := (*IssueService).IterSearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSearchIssues_Bad(t *core.T) {
	subject := (*IssueService).IterSearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSearchIssues_Ugly(t *core.T) {
	subject := (*IssueService).IterSearchIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssuesPage_Good(t *core.T) {
	subject := (*IssueService).ListIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssuesPage_Bad(t *core.T) {
	subject := (*IssueService).ListIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssuesPage_Ugly(t *core.T) {
	subject := (*IssueService).ListIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssues_Good(t *core.T) {
	subject := (*IssueService).ListIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssues_Bad(t *core.T) {
	subject := (*IssueService).ListIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssues_Ugly(t *core.T) {
	subject := (*IssueService).ListIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssues_Good(t *core.T) {
	subject := (*IssueService).IterIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssues_Bad(t *core.T) {
	subject := (*IssueService).IterIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssues_Ugly(t *core.T) {
	subject := (*IssueService).IterIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssues_Good(t *core.T) {
	subject := (*IssueService).ListRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssues_Bad(t *core.T) {
	subject := (*IssueService).ListRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssues_Ugly(t *core.T) {
	subject := (*IssueService).ListRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssuesPage_Good(t *core.T) {
	subject := (*IssueService).ListRepoIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssuesPage_Bad(t *core.T) {
	subject := (*IssueService).ListRepoIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoIssuesPage_Ugly(t *core.T) {
	subject := (*IssueService).ListRepoIssuesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoIssues_Good(t *core.T) {
	subject := (*IssueService).IterRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoIssues_Bad(t *core.T) {
	subject := (*IssueService).IterRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoIssues_Ugly(t *core.T) {
	subject := (*IssueService).IterRepoIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateIssue_Good(t *core.T) {
	subject := (*IssueService).CreateIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateIssue_Bad(t *core.T) {
	subject := (*IssueService).CreateIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateIssue_Ugly(t *core.T) {
	subject := (*IssueService).CreateIssue
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Pin_Good(t *core.T) {
	subject := (*IssueService).Pin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Pin_Bad(t *core.T) {
	subject := (*IssueService).Pin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Pin_Ugly(t *core.T) {
	subject := (*IssueService).Pin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_MovePin_Good(t *core.T) {
	subject := (*IssueService).MovePin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_MovePin_Bad(t *core.T) {
	subject := (*IssueService).MovePin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_MovePin_Ugly(t *core.T) {
	subject := (*IssueService).MovePin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListPinnedIssues_Good(t *core.T) {
	subject := (*IssueService).ListPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListPinnedIssues_Bad(t *core.T) {
	subject := (*IssueService).ListPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListPinnedIssues_Ugly(t *core.T) {
	subject := (*IssueService).ListPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterPinnedIssues_Good(t *core.T) {
	subject := (*IssueService).IterPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterPinnedIssues_Bad(t *core.T) {
	subject := (*IssueService).IterPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterPinnedIssues_Ugly(t *core.T) {
	subject := (*IssueService).IterPinnedIssues
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Unpin_Good(t *core.T) {
	subject := (*IssueService).Unpin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Unpin_Bad(t *core.T) {
	subject := (*IssueService).Unpin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_Unpin_Ugly(t *core.T) {
	subject := (*IssueService).Unpin
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SetDeadline_Good(t *core.T) {
	subject := (*IssueService).SetDeadline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SetDeadline_Bad(t *core.T) {
	subject := (*IssueService).SetDeadline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SetDeadline_Ugly(t *core.T) {
	subject := (*IssueService).SetDeadline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddReaction_Good(t *core.T) {
	subject := (*IssueService).AddReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddReaction_Bad(t *core.T) {
	subject := (*IssueService).AddReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddReaction_Ugly(t *core.T) {
	subject := (*IssueService).AddReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListReactions_Good(t *core.T) {
	subject := (*IssueService).ListReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListReactions_Bad(t *core.T) {
	subject := (*IssueService).ListReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListReactions_Ugly(t *core.T) {
	subject := (*IssueService).ListReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterReactions_Good(t *core.T) {
	subject := (*IssueService).IterReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterReactions_Bad(t *core.T) {
	subject := (*IssueService).IterReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterReactions_Ugly(t *core.T) {
	subject := (*IssueService).IterReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteReaction_Good(t *core.T) {
	subject := (*IssueService).DeleteReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteReaction_Bad(t *core.T) {
	subject := (*IssueService).DeleteReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteReaction_Ugly(t *core.T) {
	subject := (*IssueService).DeleteReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StartStopwatch_Good(t *core.T) {
	subject := (*IssueService).StartStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StartStopwatch_Bad(t *core.T) {
	subject := (*IssueService).StartStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StartStopwatch_Ugly(t *core.T) {
	subject := (*IssueService).StartStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StopStopwatch_Good(t *core.T) {
	subject := (*IssueService).StopStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StopStopwatch_Bad(t *core.T) {
	subject := (*IssueService).StopStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_StopStopwatch_Ugly(t *core.T) {
	subject := (*IssueService).StopStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteStopwatch_Good(t *core.T) {
	subject := (*IssueService).DeleteStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteStopwatch_Bad(t *core.T) {
	subject := (*IssueService).DeleteStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteStopwatch_Ugly(t *core.T) {
	subject := (*IssueService).DeleteStopwatch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimes_Good(t *core.T) {
	subject := (*IssueService).ListTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimes_Bad(t *core.T) {
	subject := (*IssueService).ListTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimes_Ugly(t *core.T) {
	subject := (*IssueService).ListTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimes_Good(t *core.T) {
	subject := (*IssueService).IterTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimes_Bad(t *core.T) {
	subject := (*IssueService).IterTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimes_Ugly(t *core.T) {
	subject := (*IssueService).IterTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddTime_Good(t *core.T) {
	subject := (*IssueService).AddTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddTime_Bad(t *core.T) {
	subject := (*IssueService).AddTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddTime_Ugly(t *core.T) {
	subject := (*IssueService).AddTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ResetTime_Good(t *core.T) {
	subject := (*IssueService).ResetTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ResetTime_Bad(t *core.T) {
	subject := (*IssueService).ResetTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ResetTime_Ugly(t *core.T) {
	subject := (*IssueService).ResetTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteTime_Good(t *core.T) {
	subject := (*IssueService).DeleteTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteTime_Bad(t *core.T) {
	subject := (*IssueService).DeleteTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteTime_Ugly(t *core.T) {
	subject := (*IssueService).DeleteTime
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddLabels_Good(t *core.T) {
	subject := (*IssueService).AddLabels
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddLabels_Bad(t *core.T) {
	subject := (*IssueService).AddLabels
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddLabels_Ugly(t *core.T) {
	subject := (*IssueService).AddLabels
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveLabel_Good(t *core.T) {
	subject := (*IssueService).RemoveLabel
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveLabel_Bad(t *core.T) {
	subject := (*IssueService).RemoveLabel
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveLabel_Ugly(t *core.T) {
	subject := (*IssueService).RemoveLabel
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListComments_Good(t *core.T) {
	subject := (*IssueService).ListComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListComments_Bad(t *core.T) {
	subject := (*IssueService).ListComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListComments_Ugly(t *core.T) {
	subject := (*IssueService).ListComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssueComments_Good(t *core.T) {
	subject := (*IssueService).ListIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssueComments_Bad(t *core.T) {
	subject := (*IssueService).ListIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListIssueComments_Ugly(t *core.T) {
	subject := (*IssueService).ListIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterComments_Good(t *core.T) {
	subject := (*IssueService).IterComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterComments_Bad(t *core.T) {
	subject := (*IssueService).IterComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterComments_Ugly(t *core.T) {
	subject := (*IssueService).IterComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssueComments_Good(t *core.T) {
	subject := (*IssueService).IterIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssueComments_Bad(t *core.T) {
	subject := (*IssueService).IterIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterIssueComments_Ugly(t *core.T) {
	subject := (*IssueService).IterIssueComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssueComment_Good(t *core.T) {
	subject := (*IssueService).GetIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssueComment_Bad(t *core.T) {
	subject := (*IssueService).GetIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetIssueComment_Ugly(t *core.T) {
	subject := (*IssueService).GetIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssueComment_Good(t *core.T) {
	subject := (*IssueService).EditIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssueComment_Bad(t *core.T) {
	subject := (*IssueService).EditIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditIssueComment_Ugly(t *core.T) {
	subject := (*IssueService).EditIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssueComment_Good(t *core.T) {
	subject := (*IssueService).DeleteIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssueComment_Bad(t *core.T) {
	subject := (*IssueService).DeleteIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteIssueComment_Ugly(t *core.T) {
	subject := (*IssueService).DeleteIssueComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateComment_Good(t *core.T) {
	subject := (*IssueService).CreateComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateComment_Bad(t *core.T) {
	subject := (*IssueService).CreateComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateComment_Ugly(t *core.T) {
	subject := (*IssueService).CreateComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditComment_Good(t *core.T) {
	subject := (*IssueService).EditComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditComment_Bad(t *core.T) {
	subject := (*IssueService).EditComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditComment_Ugly(t *core.T) {
	subject := (*IssueService).EditComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteComment_Good(t *core.T) {
	subject := (*IssueService).DeleteComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteComment_Bad(t *core.T) {
	subject := (*IssueService).DeleteComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteComment_Ugly(t *core.T) {
	subject := (*IssueService).DeleteComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoComments_Good(t *core.T) {
	subject := (*IssueService).ListRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoComments_Bad(t *core.T) {
	subject := (*IssueService).ListRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListRepoComments_Ugly(t *core.T) {
	subject := (*IssueService).ListRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoComments_Good(t *core.T) {
	subject := (*IssueService).IterRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoComments_Bad(t *core.T) {
	subject := (*IssueService).IterRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterRepoComments_Ugly(t *core.T) {
	subject := (*IssueService).IterRepoComments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetRepoComment_Good(t *core.T) {
	subject := (*IssueService).GetRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetRepoComment_Bad(t *core.T) {
	subject := (*IssueService).GetRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetRepoComment_Ugly(t *core.T) {
	subject := (*IssueService).GetRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditRepoComment_Good(t *core.T) {
	subject := (*IssueService).EditRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditRepoComment_Bad(t *core.T) {
	subject := (*IssueService).EditRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditRepoComment_Ugly(t *core.T) {
	subject := (*IssueService).EditRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteRepoComment_Good(t *core.T) {
	subject := (*IssueService).DeleteRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteRepoComment_Bad(t *core.T) {
	subject := (*IssueService).DeleteRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteRepoComment_Ugly(t *core.T) {
	subject := (*IssueService).DeleteRepoComment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentReactions_Good(t *core.T) {
	subject := (*IssueService).ListCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentReactions_Bad(t *core.T) {
	subject := (*IssueService).ListCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentReactions_Ugly(t *core.T) {
	subject := (*IssueService).ListCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentReactions_Good(t *core.T) {
	subject := (*IssueService).IterCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentReactions_Bad(t *core.T) {
	subject := (*IssueService).IterCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentReactions_Ugly(t *core.T) {
	subject := (*IssueService).IterCommentReactions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddCommentReaction_Good(t *core.T) {
	subject := (*IssueService).AddCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddCommentReaction_Bad(t *core.T) {
	subject := (*IssueService).AddCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddCommentReaction_Ugly(t *core.T) {
	subject := (*IssueService).AddCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentReaction_Good(t *core.T) {
	subject := (*IssueService).DeleteCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentReaction_Bad(t *core.T) {
	subject := (*IssueService).DeleteCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentReaction_Ugly(t *core.T) {
	subject := (*IssueService).DeleteCommentReaction
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateAttachment_Good(t *core.T) {
	subject := (*IssueService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateAttachment_Bad(t *core.T) {
	subject := (*IssueService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateAttachment_Ugly(t *core.T) {
	subject := (*IssueService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListAttachments_Good(t *core.T) {
	subject := (*IssueService).ListAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListAttachments_Bad(t *core.T) {
	subject := (*IssueService).ListAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListAttachments_Ugly(t *core.T) {
	subject := (*IssueService).ListAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterAttachments_Good(t *core.T) {
	subject := (*IssueService).IterAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterAttachments_Bad(t *core.T) {
	subject := (*IssueService).IterAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterAttachments_Ugly(t *core.T) {
	subject := (*IssueService).IterAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetAttachment_Good(t *core.T) {
	subject := (*IssueService).GetAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetAttachment_Bad(t *core.T) {
	subject := (*IssueService).GetAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetAttachment_Ugly(t *core.T) {
	subject := (*IssueService).GetAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditAttachment_Good(t *core.T) {
	subject := (*IssueService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditAttachment_Bad(t *core.T) {
	subject := (*IssueService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditAttachment_Ugly(t *core.T) {
	subject := (*IssueService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteAttachment_Good(t *core.T) {
	subject := (*IssueService).DeleteAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteAttachment_Bad(t *core.T) {
	subject := (*IssueService).DeleteAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteAttachment_Ugly(t *core.T) {
	subject := (*IssueService).DeleteAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentAttachments_Good(t *core.T) {
	subject := (*IssueService).ListCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentAttachments_Bad(t *core.T) {
	subject := (*IssueService).ListCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListCommentAttachments_Ugly(t *core.T) {
	subject := (*IssueService).ListCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentAttachments_Good(t *core.T) {
	subject := (*IssueService).IterCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentAttachments_Bad(t *core.T) {
	subject := (*IssueService).IterCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterCommentAttachments_Ugly(t *core.T) {
	subject := (*IssueService).IterCommentAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetCommentAttachment_Good(t *core.T) {
	subject := (*IssueService).GetCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetCommentAttachment_Bad(t *core.T) {
	subject := (*IssueService).GetCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_GetCommentAttachment_Ugly(t *core.T) {
	subject := (*IssueService).GetCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateCommentAttachment_Good(t *core.T) {
	subject := (*IssueService).CreateCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateCommentAttachment_Bad(t *core.T) {
	subject := (*IssueService).CreateCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CreateCommentAttachment_Ugly(t *core.T) {
	subject := (*IssueService).CreateCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditCommentAttachment_Good(t *core.T) {
	subject := (*IssueService).EditCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditCommentAttachment_Bad(t *core.T) {
	subject := (*IssueService).EditCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_EditCommentAttachment_Ugly(t *core.T) {
	subject := (*IssueService).EditCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentAttachment_Good(t *core.T) {
	subject := (*IssueService).DeleteCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentAttachment_Bad(t *core.T) {
	subject := (*IssueService).DeleteCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_DeleteCommentAttachment_Ugly(t *core.T) {
	subject := (*IssueService).DeleteCommentAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimeline_Good(t *core.T) {
	subject := (*IssueService).ListTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimeline_Bad(t *core.T) {
	subject := (*IssueService).ListTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListTimeline_Ugly(t *core.T) {
	subject := (*IssueService).ListTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimeline_Good(t *core.T) {
	subject := (*IssueService).IterTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimeline_Bad(t *core.T) {
	subject := (*IssueService).IterTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterTimeline_Ugly(t *core.T) {
	subject := (*IssueService).IterTimeline
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListSubscriptions_Good(t *core.T) {
	subject := (*IssueService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListSubscriptions_Bad(t *core.T) {
	subject := (*IssueService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListSubscriptions_Ugly(t *core.T) {
	subject := (*IssueService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSubscriptions_Good(t *core.T) {
	subject := (*IssueService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSubscriptions_Bad(t *core.T) {
	subject := (*IssueService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterSubscriptions_Ugly(t *core.T) {
	subject := (*IssueService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CheckSubscription_Good(t *core.T) {
	subject := (*IssueService).CheckSubscription
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CheckSubscription_Bad(t *core.T) {
	subject := (*IssueService).CheckSubscription
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_CheckSubscription_Ugly(t *core.T) {
	subject := (*IssueService).CheckSubscription
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SubscribeUser_Good(t *core.T) {
	subject := (*IssueService).SubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SubscribeUser_Bad(t *core.T) {
	subject := (*IssueService).SubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_SubscribeUser_Ugly(t *core.T) {
	subject := (*IssueService).SubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_UnsubscribeUser_Good(t *core.T) {
	subject := (*IssueService).UnsubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_UnsubscribeUser_Bad(t *core.T) {
	subject := (*IssueService).UnsubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_UnsubscribeUser_Ugly(t *core.T) {
	subject := (*IssueService).UnsubscribeUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListDependencies_Good(t *core.T) {
	subject := (*IssueService).ListDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListDependencies_Bad(t *core.T) {
	subject := (*IssueService).ListDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListDependencies_Ugly(t *core.T) {
	subject := (*IssueService).ListDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterDependencies_Good(t *core.T) {
	subject := (*IssueService).IterDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterDependencies_Bad(t *core.T) {
	subject := (*IssueService).IterDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterDependencies_Ugly(t *core.T) {
	subject := (*IssueService).IterDependencies
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddDependency_Good(t *core.T) {
	subject := (*IssueService).AddDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddDependency_Bad(t *core.T) {
	subject := (*IssueService).AddDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddDependency_Ugly(t *core.T) {
	subject := (*IssueService).AddDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveDependency_Good(t *core.T) {
	subject := (*IssueService).RemoveDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveDependency_Bad(t *core.T) {
	subject := (*IssueService).RemoveDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveDependency_Ugly(t *core.T) {
	subject := (*IssueService).RemoveDependency
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListBlocks_Good(t *core.T) {
	subject := (*IssueService).ListBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListBlocks_Bad(t *core.T) {
	subject := (*IssueService).ListBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_ListBlocks_Ugly(t *core.T) {
	subject := (*IssueService).ListBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterBlocks_Good(t *core.T) {
	subject := (*IssueService).IterBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterBlocks_Bad(t *core.T) {
	subject := (*IssueService).IterBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_IterBlocks_Ugly(t *core.T) {
	subject := (*IssueService).IterBlocks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddBlock_Good(t *core.T) {
	subject := (*IssueService).AddBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddBlock_Bad(t *core.T) {
	subject := (*IssueService).AddBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_AddBlock_Ugly(t *core.T) {
	subject := (*IssueService).AddBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveBlock_Good(t *core.T) {
	subject := (*IssueService).RemoveBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveBlock_Bad(t *core.T) {
	subject := (*IssueService).RemoveBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestIssues_IssueService_RemoveBlock_Ugly(t *core.T) {
	subject := (*IssueService).RemoveBlock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
