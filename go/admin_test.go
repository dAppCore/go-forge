package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
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

func TestAdminService_ListHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Hook{
			{ID: 7, Type: "forgejo", URL: "https://example.com/admin-hook", Active: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hooks, err := f.Admin.ListHooks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 {
		t.Fatalf("got %d hooks, want 1", len(hooks))
	}
	if hooks[0].ID != 7 || hooks[0].URL != "https://example.com/admin-hook" {
		t.Errorf("unexpected hook: %+v", hooks[0])
	}
}

func TestAdminService_CreateHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Type != "forgejo" {
			t.Errorf("got type=%q, want %q", opts.Type, "forgejo")
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     12,
			Type:   opts.Type,
			Active: opts.Active,
			Events: opts.Events,
			URL:    "https://example.com/admin-hook",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Admin.CreateHook(context.Background(), &types.CreateHookOption{
		Type:   "forgejo",
		Active: true,
		Events: []string{"push"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 12 {
		t.Errorf("got id=%d, want 12", hook.ID)
	}
	if hook.Type != "forgejo" {
		t.Errorf("got type=%q, want %q", hook.Type, "forgejo")
	}
}

func TestAdminService_GetHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/hooks/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     7,
			Type:   "forgejo",
			Active: true,
			URL:    "https://example.com/admin-hook",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Admin.GetHook(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 7 {
		t.Errorf("got id=%d, want 7", hook.ID)
	}
}

func TestAdminService_EditHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/hooks/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if !opts.Active {
			t.Error("expected active=true")
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     7,
			Type:   "forgejo",
			Active: opts.Active,
			URL:    "https://example.com/admin-hook",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Admin.EditHook(context.Background(), 7, &types.EditHookOption{Active: true})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 7 || !hook.Active {
		t.Errorf("unexpected hook: %+v", hook)
	}
}

func TestAdminService_DeleteHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/hooks/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.DeleteHook(context.Background(), 7); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_ListQuotaGroups_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.QuotaGroup{
			{
				Name: "default",
				Rules: []*types.QuotaRuleInfo{
					{Name: "git", Limit: 200000000, Subjects: []string{"size:repos:all"}},
				},
			},
			{
				Name: "premium",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	groups, err := f.Admin.ListQuotaGroups(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2", len(groups))
	}
	if groups[0].Name != "default" {
		t.Errorf("got name=%q, want %q", groups[0].Name, "default")
	}
	if len(groups[0].Rules) != 1 || groups[0].Rules[0].Name != "git" {
		t.Errorf("unexpected rules: %+v", groups[0].Rules)
	}
}

func TestAdminService_IterQuotaGroups_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.QuotaGroup{
			{Name: "default"},
			{Name: "premium"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for group, err := range f.Admin.IterQuotaGroups(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, group.Name)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(got) != 2 || got[0] != "default" || got[1] != "premium" {
		t.Fatalf("got %#v", got)
	}
}

func TestAdminService_CreateQuotaGroup_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateQuotaGroupOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Name != "newgroup" {
			t.Errorf("got name=%q, want %q", opts.Name, "newgroup")
		}
		if len(opts.Rules) != 1 || opts.Rules[0].Name != "git" {
			t.Fatalf("unexpected rules: %+v", opts.Rules)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.QuotaGroup{
			Name: opts.Name,
			Rules: []*types.QuotaRuleInfo{
				{
					Name:     opts.Rules[0].Name,
					Limit:    opts.Rules[0].Limit,
					Subjects: opts.Rules[0].Subjects,
				},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	group, err := f.Admin.CreateQuotaGroup(context.Background(), &types.CreateQuotaGroupOptions{
		Name: "newgroup",
		Rules: []*types.CreateQuotaRuleOptions{
			{
				Name:     "git",
				Limit:    200000000,
				Subjects: []string{"size:repos:all"},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if group.Name != "newgroup" {
		t.Errorf("got name=%q, want %q", group.Name, "newgroup")
	}
	if len(group.Rules) != 1 || group.Rules[0].Limit != 200000000 {
		t.Errorf("unexpected rules: %+v", group.Rules)
	}
}

func TestAdminService_GetQuotaGroup_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups/default" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.QuotaGroup{
			Name: "default",
			Rules: []*types.QuotaRuleInfo{
				{Name: "git", Limit: 200000000, Subjects: []string{"size:repos:all"}},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	group, err := f.Admin.GetQuotaGroup(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	if group.Name != "default" {
		t.Errorf("got name=%q, want %q", group.Name, "default")
	}
	if len(group.Rules) != 1 || group.Rules[0].Name != "git" {
		t.Fatalf("unexpected rules: %+v", group.Rules)
	}
}

func TestAdminService_DeleteQuotaGroup_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups/default" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.DeleteQuotaGroup(context.Background(), "default"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_ListQuotaGroupUsers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups/default/users" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.User{
			{ID: 1, UserName: "alice"},
			{ID: 2, UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Admin.ListQuotaGroupUsers(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatalf("got %d users, want 2", len(users))
	}
	if users[0].UserName != "alice" {
		t.Errorf("got username=%q, want %q", users[0].UserName, "alice")
	}
}

func TestAdminService_AddQuotaGroupUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups/default/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.AddQuotaGroupUser(context.Background(), "default", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_RemoveQuotaGroupUser_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/groups/default/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.RemoveQuotaGroupUser(context.Background(), "default", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_ListQuotaRules_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.QuotaRuleInfo{
			{Name: "git", Limit: 200000000, Subjects: []string{"size:repos:all"}},
			{Name: "artifacts", Limit: 50000000, Subjects: []string{"size:assets:artifacts"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	rules, err := f.Admin.ListQuotaRules(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 2 {
		t.Fatalf("got %d rules, want 2", len(rules))
	}
	if rules[0].Name != "git" {
		t.Errorf("got name=%q, want %q", rules[0].Name, "git")
	}
}

func TestAdminService_IterQuotaRules_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.QuotaRuleInfo{
			{Name: "git"},
			{Name: "artifacts"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for rule, err := range f.Admin.IterQuotaRules(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, rule.Name)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(got) != 2 || got[0] != "git" || got[1] != "artifacts" {
		t.Fatalf("got %#v", got)
	}
}

func TestAdminService_CreateQuotaRule_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateQuotaRuleOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Name != "git" || opts.Limit != 200000000 {
			t.Fatalf("unexpected options: %+v", opts)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.QuotaRuleInfo{
			Name:     opts.Name,
			Limit:    opts.Limit,
			Subjects: opts.Subjects,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	rule, err := f.Admin.CreateQuotaRule(context.Background(), &types.CreateQuotaRuleOptions{
		Name:     "git",
		Limit:    200000000,
		Subjects: []string{"size:repos:all"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rule.Name != "git" || rule.Limit != 200000000 {
		t.Errorf("unexpected rule: %+v", rule)
	}
}

func TestAdminService_GetQuotaRule_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules/git" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.QuotaRuleInfo{
			Name:     "git",
			Limit:    200000000,
			Subjects: []string{"size:repos:all"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	rule, err := f.Admin.GetQuotaRule(context.Background(), "git")
	if err != nil {
		t.Fatal(err)
	}
	if rule.Name != "git" {
		t.Errorf("got name=%q, want %q", rule.Name, "git")
	}
}

func TestAdminService_EditQuotaRule_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules/git" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditQuotaRuleOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Limit != 500000000 {
			t.Fatalf("unexpected options: %+v", opts)
		}
		json.NewEncoder(w).Encode(types.QuotaRuleInfo{
			Name:     "git",
			Limit:    opts.Limit,
			Subjects: opts.Subjects,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	rule, err := f.Admin.EditQuotaRule(context.Background(), "git", &types.EditQuotaRuleOptions{
		Limit:    500000000,
		Subjects: []string{"size:repos:all", "size:assets:packages"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if rule.Limit != 500000000 {
		t.Errorf("got limit=%d, want 500000000", rule.Limit)
	}
}

func TestAdminService_DeleteQuotaRule_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/quota/rules/git" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.DeleteQuotaRule(context.Background(), "git"); err != nil {
		t.Fatal(err)
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

func TestAdminService_ListUnadoptedRepos_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/unadopted" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("pattern"); got != "core/*" {
			t.Errorf("got pattern=%q, want %q", got, "core/*")
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("got limit=%q, want %q", got, "50")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]string{"core/myrepo"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Admin.ListUnadoptedRepos(context.Background(), AdminUnadoptedListOptions{Pattern: "core/*"})
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 || repos[0] != "core/myrepo" {
		t.Fatalf("unexpected result: %#v", repos)
	}
}

func TestAdminService_DeleteUnadoptedRepo_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/unadopted/alice/myrepo" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Admin.DeleteUnadoptedRepo(context.Background(), "alice", "myrepo"); err != nil {
		t.Fatal(err)
	}
}

func TestAdminService_ListActionsRuns_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/admin/actions/runs" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("status"); got != "in_progress" {
			t.Errorf("got status=%q, want %q", got, "in_progress")
		}
		if got := r.URL.Query().Get("branch"); got != "main" {
			t.Errorf("got branch=%q, want %q", got, "main")
		}
		if got := r.URL.Query().Get("actor"); got != "alice" {
			t.Errorf("got actor=%q, want %q", got, "alice")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Errorf("got limit=%q, want %q", got, "25")
		}
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode(types.ActionTaskResponse{
			Entries: []*types.ActionTask{
				{ID: 101, Name: "build", Status: "in_progress", Event: "push"},
				{ID: 102, Name: "test", Status: "queued", Event: "push"},
			},
			TotalCount: 3,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Admin.ListActionsRuns(context.Background(), AdminActionsRunListOptions{
		Status: "in_progress",
		Branch: "main",
		Actor:  "alice",
	}, ListOptions{Page: 2, Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalCount != 3 {
		t.Fatalf("got total count=%d, want 3", result.TotalCount)
	}
	if len(result.Items) != 2 {
		t.Fatalf("got %d runs, want 2", len(result.Items))
	}
	if result.Items[0].ID != 101 || result.Items[0].Name != "build" {
		t.Errorf("unexpected first run: %+v", result.Items[0])
	}
}

func TestAdminService_IterActionsRuns_Good(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/v1/admin/actions/runs" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		switch calls {
		case 1:
			json.NewEncoder(w).Encode(types.ActionTaskResponse{
				Entries: []*types.ActionTask{
					{ID: 201, Name: "build"},
				},
				TotalCount: 2,
			})
		default:
			json.NewEncoder(w).Encode(types.ActionTaskResponse{
				Entries: []*types.ActionTask{
					{ID: 202, Name: "test"},
				},
				TotalCount: 2,
			})
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var ids []int64
	for run, err := range f.Admin.IterActionsRuns(context.Background(), AdminActionsRunListOptions{}) {
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, run.ID)
	}
	if len(ids) != 2 || ids[0] != 201 || ids[1] != 202 {
		t.Fatalf("unexpected run ids: %v", ids)
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

func TestAdmin_AdminActionsRunListOptions_String_Good(t *core.T) {
	subject := (*AdminActionsRunListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminActionsRunListOptions_String_Bad(t *core.T) {
	subject := (*AdminActionsRunListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminActionsRunListOptions_String_Ugly(t *core.T) {
	subject := (*AdminActionsRunListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminActionsRunListOptions_GoString_Good(t *core.T) {
	subject := (*AdminActionsRunListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminActionsRunListOptions_GoString_Bad(t *core.T) {
	subject := (*AdminActionsRunListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminActionsRunListOptions_GoString_Ugly(t *core.T) {
	subject := (*AdminActionsRunListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_String_Good(t *core.T) {
	subject := (*AdminUnadoptedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_String_Bad(t *core.T) {
	subject := (*AdminUnadoptedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_String_Ugly(t *core.T) {
	subject := (*AdminUnadoptedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_GoString_Good(t *core.T) {
	subject := (*AdminUnadoptedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_GoString_Bad(t *core.T) {
	subject := (*AdminUnadoptedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminUnadoptedListOptions_GoString_Ugly(t *core.T) {
	subject := (*AdminUnadoptedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUsers_Good(t *core.T) {
	subject := (*AdminService).ListUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUsers_Bad(t *core.T) {
	subject := (*AdminService).ListUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUsers_Ugly(t *core.T) {
	subject := (*AdminService).ListUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUsers_Good(t *core.T) {
	subject := (*AdminService).IterUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUsers_Bad(t *core.T) {
	subject := (*AdminService).IterUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUsers_Ugly(t *core.T) {
	subject := (*AdminService).IterUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUser_Good(t *core.T) {
	subject := (*AdminService).CreateUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUser_Bad(t *core.T) {
	subject := (*AdminService).CreateUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUser_Ugly(t *core.T) {
	subject := (*AdminService).CreateUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserKey_Good(t *core.T) {
	subject := (*AdminService).CreateUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserKey_Bad(t *core.T) {
	subject := (*AdminService).CreateUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserKey_Ugly(t *core.T) {
	subject := (*AdminService).CreateUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUserKey_Good(t *core.T) {
	subject := (*AdminService).DeleteUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUserKey_Bad(t *core.T) {
	subject := (*AdminService).DeleteUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUserKey_Ugly(t *core.T) {
	subject := (*AdminService).DeleteUserKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserOrg_Good(t *core.T) {
	subject := (*AdminService).CreateUserOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserOrg_Bad(t *core.T) {
	subject := (*AdminService).CreateUserOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserOrg_Ugly(t *core.T) {
	subject := (*AdminService).CreateUserOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetUserQuota_Good(t *core.T) {
	subject := (*AdminService).GetUserQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetUserQuota_Bad(t *core.T) {
	subject := (*AdminService).GetUserQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetUserQuota_Ugly(t *core.T) {
	subject := (*AdminService).GetUserQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SetUserQuotaGroups_Good(t *core.T) {
	subject := (*AdminService).SetUserQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SetUserQuotaGroups_Bad(t *core.T) {
	subject := (*AdminService).SetUserQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SetUserQuotaGroups_Ugly(t *core.T) {
	subject := (*AdminService).SetUserQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserRepo_Good(t *core.T) {
	subject := (*AdminService).CreateUserRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserRepo_Bad(t *core.T) {
	subject := (*AdminService).CreateUserRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateUserRepo_Ugly(t *core.T) {
	subject := (*AdminService).CreateUserRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditUser_Good(t *core.T) {
	subject := (*AdminService).EditUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditUser_Bad(t *core.T) {
	subject := (*AdminService).EditUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditUser_Ugly(t *core.T) {
	subject := (*AdminService).EditUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUser_Good(t *core.T) {
	subject := (*AdminService).DeleteUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUser_Bad(t *core.T) {
	subject := (*AdminService).DeleteUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUser_Ugly(t *core.T) {
	subject := (*AdminService).DeleteUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RenameUser_Good(t *core.T) {
	subject := (*AdminService).RenameUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RenameUser_Bad(t *core.T) {
	subject := (*AdminService).RenameUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RenameUser_Ugly(t *core.T) {
	subject := (*AdminService).RenameUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListOrgs_Good(t *core.T) {
	subject := (*AdminService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListOrgs_Bad(t *core.T) {
	subject := (*AdminService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListOrgs_Ugly(t *core.T) {
	subject := (*AdminService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterOrgs_Good(t *core.T) {
	subject := (*AdminService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterOrgs_Bad(t *core.T) {
	subject := (*AdminService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterOrgs_Ugly(t *core.T) {
	subject := (*AdminService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListEmails_Good(t *core.T) {
	subject := (*AdminService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListEmails_Bad(t *core.T) {
	subject := (*AdminService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListEmails_Ugly(t *core.T) {
	subject := (*AdminService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterEmails_Good(t *core.T) {
	subject := (*AdminService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterEmails_Bad(t *core.T) {
	subject := (*AdminService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterEmails_Ugly(t *core.T) {
	subject := (*AdminService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListHooks_Good(t *core.T) {
	subject := (*AdminService).ListHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListHooks_Bad(t *core.T) {
	subject := (*AdminService).ListHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListHooks_Ugly(t *core.T) {
	subject := (*AdminService).ListHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterHooks_Good(t *core.T) {
	subject := (*AdminService).IterHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterHooks_Bad(t *core.T) {
	subject := (*AdminService).IterHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterHooks_Ugly(t *core.T) {
	subject := (*AdminService).IterHooks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetHook_Good(t *core.T) {
	subject := (*AdminService).GetHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetHook_Bad(t *core.T) {
	subject := (*AdminService).GetHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetHook_Ugly(t *core.T) {
	subject := (*AdminService).GetHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateHook_Good(t *core.T) {
	subject := (*AdminService).CreateHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateHook_Bad(t *core.T) {
	subject := (*AdminService).CreateHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateHook_Ugly(t *core.T) {
	subject := (*AdminService).CreateHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditHook_Good(t *core.T) {
	subject := (*AdminService).EditHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditHook_Bad(t *core.T) {
	subject := (*AdminService).EditHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditHook_Ugly(t *core.T) {
	subject := (*AdminService).EditHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteHook_Good(t *core.T) {
	subject := (*AdminService).DeleteHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteHook_Bad(t *core.T) {
	subject := (*AdminService).DeleteHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteHook_Ugly(t *core.T) {
	subject := (*AdminService).DeleteHook
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroups_Good(t *core.T) {
	subject := (*AdminService).ListQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroups_Bad(t *core.T) {
	subject := (*AdminService).ListQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroups_Ugly(t *core.T) {
	subject := (*AdminService).ListQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroups_Good(t *core.T) {
	subject := (*AdminService).IterQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroups_Bad(t *core.T) {
	subject := (*AdminService).IterQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroups_Ugly(t *core.T) {
	subject := (*AdminService).IterQuotaGroups
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaGroup_Good(t *core.T) {
	subject := (*AdminService).CreateQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaGroup_Bad(t *core.T) {
	subject := (*AdminService).CreateQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaGroup_Ugly(t *core.T) {
	subject := (*AdminService).CreateQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaGroup_Good(t *core.T) {
	subject := (*AdminService).GetQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaGroup_Bad(t *core.T) {
	subject := (*AdminService).GetQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaGroup_Ugly(t *core.T) {
	subject := (*AdminService).GetQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaGroup_Good(t *core.T) {
	subject := (*AdminService).DeleteQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaGroup_Bad(t *core.T) {
	subject := (*AdminService).DeleteQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaGroup_Ugly(t *core.T) {
	subject := (*AdminService).DeleteQuotaGroup
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupRule_Good(t *core.T) {
	subject := (*AdminService).AddQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupRule_Bad(t *core.T) {
	subject := (*AdminService).AddQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupRule_Ugly(t *core.T) {
	subject := (*AdminService).AddQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupRule_Good(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupRule_Bad(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupRule_Ugly(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroupUsers_Good(t *core.T) {
	subject := (*AdminService).ListQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroupUsers_Bad(t *core.T) {
	subject := (*AdminService).ListQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaGroupUsers_Ugly(t *core.T) {
	subject := (*AdminService).ListQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroupUsers_Good(t *core.T) {
	subject := (*AdminService).IterQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroupUsers_Bad(t *core.T) {
	subject := (*AdminService).IterQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaGroupUsers_Ugly(t *core.T) {
	subject := (*AdminService).IterQuotaGroupUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupUser_Good(t *core.T) {
	subject := (*AdminService).AddQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupUser_Bad(t *core.T) {
	subject := (*AdminService).AddQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AddQuotaGroupUser_Ugly(t *core.T) {
	subject := (*AdminService).AddQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupUser_Good(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupUser_Bad(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RemoveQuotaGroupUser_Ugly(t *core.T) {
	subject := (*AdminService).RemoveQuotaGroupUser
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaRules_Good(t *core.T) {
	subject := (*AdminService).ListQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaRules_Bad(t *core.T) {
	subject := (*AdminService).ListQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListQuotaRules_Ugly(t *core.T) {
	subject := (*AdminService).ListQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaRules_Good(t *core.T) {
	subject := (*AdminService).IterQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaRules_Bad(t *core.T) {
	subject := (*AdminService).IterQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterQuotaRules_Ugly(t *core.T) {
	subject := (*AdminService).IterQuotaRules
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaRule_Good(t *core.T) {
	subject := (*AdminService).CreateQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaRule_Bad(t *core.T) {
	subject := (*AdminService).CreateQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_CreateQuotaRule_Ugly(t *core.T) {
	subject := (*AdminService).CreateQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaRule_Good(t *core.T) {
	subject := (*AdminService).GetQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaRule_Bad(t *core.T) {
	subject := (*AdminService).GetQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GetQuotaRule_Ugly(t *core.T) {
	subject := (*AdminService).GetQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditQuotaRule_Good(t *core.T) {
	subject := (*AdminService).EditQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditQuotaRule_Bad(t *core.T) {
	subject := (*AdminService).EditQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_EditQuotaRule_Ugly(t *core.T) {
	subject := (*AdminService).EditQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaRule_Good(t *core.T) {
	subject := (*AdminService).DeleteQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaRule_Bad(t *core.T) {
	subject := (*AdminService).DeleteQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteQuotaRule_Ugly(t *core.T) {
	subject := (*AdminService).DeleteQuotaRule
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUnadoptedRepos_Good(t *core.T) {
	subject := (*AdminService).ListUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUnadoptedRepos_Bad(t *core.T) {
	subject := (*AdminService).ListUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListUnadoptedRepos_Ugly(t *core.T) {
	subject := (*AdminService).ListUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUnadoptedRepos_Good(t *core.T) {
	subject := (*AdminService).IterUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUnadoptedRepos_Bad(t *core.T) {
	subject := (*AdminService).IterUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterUnadoptedRepos_Ugly(t *core.T) {
	subject := (*AdminService).IterUnadoptedRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SearchEmails_Good(t *core.T) {
	subject := (*AdminService).SearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SearchEmails_Bad(t *core.T) {
	subject := (*AdminService).SearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_SearchEmails_Ugly(t *core.T) {
	subject := (*AdminService).SearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterSearchEmails_Good(t *core.T) {
	subject := (*AdminService).IterSearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterSearchEmails_Bad(t *core.T) {
	subject := (*AdminService).IterSearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterSearchEmails_Ugly(t *core.T) {
	subject := (*AdminService).IterSearchEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RunCron_Good(t *core.T) {
	subject := (*AdminService).RunCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RunCron_Bad(t *core.T) {
	subject := (*AdminService).RunCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_RunCron_Ugly(t *core.T) {
	subject := (*AdminService).RunCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListCron_Good(t *core.T) {
	subject := (*AdminService).ListCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListCron_Bad(t *core.T) {
	subject := (*AdminService).ListCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListCron_Ugly(t *core.T) {
	subject := (*AdminService).ListCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterCron_Good(t *core.T) {
	subject := (*AdminService).IterCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterCron_Bad(t *core.T) {
	subject := (*AdminService).IterCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterCron_Ugly(t *core.T) {
	subject := (*AdminService).IterCron
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListActionsRuns_Good(t *core.T) {
	subject := (*AdminService).ListActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListActionsRuns_Bad(t *core.T) {
	subject := (*AdminService).ListActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_ListActionsRuns_Ugly(t *core.T) {
	subject := (*AdminService).ListActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterActionsRuns_Good(t *core.T) {
	subject := (*AdminService).IterActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterActionsRuns_Bad(t *core.T) {
	subject := (*AdminService).IterActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_IterActionsRuns_Ugly(t *core.T) {
	subject := (*AdminService).IterActionsRuns
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AdoptRepo_Good(t *core.T) {
	subject := (*AdminService).AdoptRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AdoptRepo_Bad(t *core.T) {
	subject := (*AdminService).AdoptRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_AdoptRepo_Ugly(t *core.T) {
	subject := (*AdminService).AdoptRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUnadoptedRepo_Good(t *core.T) {
	subject := (*AdminService).DeleteUnadoptedRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUnadoptedRepo_Bad(t *core.T) {
	subject := (*AdminService).DeleteUnadoptedRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_DeleteUnadoptedRepo_Ugly(t *core.T) {
	subject := (*AdminService).DeleteUnadoptedRepo
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GenerateRunnerToken_Good(t *core.T) {
	subject := (*AdminService).GenerateRunnerToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GenerateRunnerToken_Bad(t *core.T) {
	subject := (*AdminService).GenerateRunnerToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestAdmin_AdminService_GenerateRunnerToken_Ugly(t *core.T) {
	subject := (*AdminService).GenerateRunnerToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
