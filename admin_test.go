package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestAdminService_ListUsers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/users" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{
			{ID: 1, UserName: "alice"},
			{ID: 2, UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Admin.ListUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Errorf("got %d users, want 2", len(users))
	}
	if users[0].UserName != "alice" {
		t.Errorf("got username=%q, want %q", users[0].UserName, "alice")
	}
}

func TestAdminService_CreateUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/users" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateUserOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Username != "newuser" {
			t.Errorf("got username=%q, want %q", opts.Username, "newuser")
		}
		if opts.Email != "new@example.com" {
			t.Errorf("got email=%q, want %q", opts.Email, "new@example.com")
		}
		json.NewEncoder(w).Encode(types.User{ID: 42, UserName: "newuser", Email: "new@example.com"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Admin.CreateUser(context.Background(), &types.CreateUserOption{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != 42 {
		t.Errorf("got id=%d, want 42", user.ID)
	}
	if user.UserName != "newuser" {
		t.Errorf("got username=%q, want %q", user.UserName, "newuser")
	}
}

func TestAdminService_DeleteUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.DeleteUser(context.Background(), "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_RunCron_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/cron/repo_health_check" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.RunCron(context.Background(), "repo_health_check"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_EditUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["email"] != "alice@new.com" {
			t.Errorf("got email=%v, want %q", body["email"], "alice@new.com")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Admin.EditUser(context.Background(), "alice", map[string]any{
		"email": "alice@new.com",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_RenameUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/users/alice/rename" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.RenameUserOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.NewName != "alice2" {
			t.Errorf("got new_username=%q, want %q", opts.NewName, "alice2")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.RenameUser(context.Background(), "alice", "alice2"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_ListOrgs_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/orgs" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Organization{
			{ID: 10, Name: "myorg"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	orgs, err := f.Admin.ListOrgs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(orgs) != 1 {
		t.Errorf("got %d orgs, want 1", len(orgs))
	}
	if orgs[0].Name != "myorg" {
		t.Errorf("got name=%q, want %q", orgs[0].Name, "myorg")
	}
}

func TestAdminService_ListEmails_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/emails" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Email{
			{Email: "alice@example.com", Primary: true},
			{Email: "bob@example.com", Verified: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	emails, err := f.Admin.ListEmails(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 2 {
		t.Errorf("got %d emails, want 2", len(emails))
	}
	if emails[0].Email != "alice@example.com" || !emails[0].Primary {
		t.Errorf("got first email=%+v, want primary alice@example.com", emails[0])
	}
}

func TestAdminService_SearchEmails_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/emails/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("q"); got != "alice" {
			t.Errorf("got q=%q, want %q", got, "alice")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Email{
			{Email: "alice@example.com", Primary: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	emails, err := f.Admin.SearchEmails(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(emails) != 1 {
		t.Errorf("got %d emails, want 1", len(emails))
	}
	if emails[0].Email != "alice@example.com" {
		t.Errorf("got email=%q, want %q", emails[0].Email, "alice@example.com")
	}
}

func TestAdminService_ListCron_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/cron" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Cron{
			{Name: "repo_health_check", Schedule: "@every 24h"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	crons, err := f.Admin.ListCron(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(crons) != 1 {
		t.Errorf("got %d crons, want 1", len(crons))
	}
	if crons[0].Name != "repo_health_check" {
		t.Errorf("got name=%q, want %q", crons[0].Name, "repo_health_check")
	}
}

func TestAdminService_AdoptRepo_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/unadopted/alice/myrepo" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.AdoptRepo(context.Background(), "alice", "myrepo"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_GenerateRunnerToken_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/runners/registration-token" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"token": "abc123"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	token, err := f.Admin.GenerateRunnerToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token != "abc123" {
		t.Errorf("got token=%q, want %q", token, "abc123")
	}
}

func TestAdminService_DeleteUser_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "user not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Admin.DeleteUser(context.Background(), "nonexistent")
	if !IsNotFound(err) {
		t.Errorf("expected not found error, got %v", err)
	}
}

func TestAdminService_CreateUser_Forbidden_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "only admins can create users"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Admin.CreateUser(context.Background(), &types.CreateUserOption{
		Username: "newuser",
		Email:    "new@example.com",
	})
	if !IsForbidden(err) {
		t.Errorf("expected forbidden error, got %v", err)
	}
}
