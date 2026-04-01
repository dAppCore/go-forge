package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestUserService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.User{ID: 1, UserName: "alice"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.Get(context.Background(), Params{"username": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	if user.UserName != "alice" {
		t.Errorf("got username=%q, want %q", user.UserName, "alice")
	}
}

func TestUserService_GetCurrent_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.User{ID: 1, UserName: "me"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.GetCurrent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if user.UserName != "me" {
		t.Errorf("got username=%q, want %q", user.UserName, "me")
	}
}

func TestUserService_ListEmails_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/emails" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Email{
			{Email: "alice@example.com", Primary: true},
			{Email: "alice+alt@example.com", Verified: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	emails, err := f.Users.ListEmails(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 2 {
		t.Fatalf("got %d emails, want 2", len(emails))
	}
	if emails[0].Email != "alice@example.com" || !emails[0].Primary {
		t.Errorf("unexpected first email: %+v", emails[0])
	}
}

func TestUserService_ListStopwatches_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/stopwatches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.StopWatch{
			{IssueIndex: 12, IssueTitle: "First issue", RepoOwnerName: "core", RepoName: "go-forge", Seconds: 30},
			{IssueIndex: 13, IssueTitle: "Second issue", RepoOwnerName: "core", RepoName: "go-forge", Seconds: 90},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	stopwatches, err := f.Users.ListStopwatches(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(stopwatches) != 2 {
		t.Fatalf("got %d stopwatches, want 2", len(stopwatches))
	}
	if stopwatches[0].IssueIndex != 12 || stopwatches[0].Seconds != 30 {
		t.Errorf("unexpected first stopwatch: %+v", stopwatches[0])
	}
}

func TestUserService_IterStopwatches_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/stopwatches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.StopWatch{
			{IssueIndex: 99, IssueTitle: "Running task", RepoOwnerName: "core", RepoName: "go-forge", Seconds: 300},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for sw, err := range f.Users.IterStopwatches(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if sw.IssueIndex != 99 || sw.Seconds != 300 {
			t.Errorf("unexpected stopwatch: %+v", sw)
		}
	}
	if count != 1 {
		t.Fatalf("got %d stopwatches, want 1", count)
	}
}

func TestUserService_AddEmails_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/emails" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateEmailOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Emails) != 2 || body.Emails[0] != "alice@example.com" || body.Emails[1] != "alice+alt@example.com" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode([]types.Email{
			{Email: "alice@example.com", Primary: true},
			{Email: "alice+alt@example.com", Verified: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	emails, err := f.Users.AddEmails(context.Background(), "alice@example.com", "alice+alt@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 2 {
		t.Fatalf("got %d emails, want 2", len(emails))
	}
	if emails[1].Email != "alice+alt@example.com" || !emails[1].Verified {
		t.Errorf("unexpected second email: %+v", emails[1])
	}
}

func TestUserService_DeleteEmails_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/emails" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.DeleteEmailOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Emails) != 1 || body.Emails[0] != "alice+alt@example.com" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteEmails(context.Background(), "alice+alt@example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_ListFollowers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/followers" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{
			{ID: 2, UserName: "bob"},
			{ID: 3, UserName: "charlie"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	followers, err := f.Users.ListFollowers(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(followers) != 2 {
		t.Errorf("got %d followers, want 2", len(followers))
	}
	if followers[0].UserName != "bob" {
		t.Errorf("got username=%q, want %q", followers[0].UserName, "bob")
	}
}
