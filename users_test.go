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

func TestUserService_GetSettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/settings" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.UserSettings{
			FullName:    "Alice",
			Language:    "en-US",
			Theme:       "forgejo-auto",
			HideEmail:   true,
			Pronouns:    "she/her",
			Website:     "https://example.com",
			Location:    "Earth",
			Description: "maintainer",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Users.GetSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.FullName != "Alice" {
		t.Errorf("got full name=%q, want %q", settings.FullName, "Alice")
	}
	if !settings.HideEmail {
		t.Errorf("got hide_email=%v, want true", settings.HideEmail)
	}
}

func TestUserService_UpdateSettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/settings" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.UserSettingsOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.FullName != "Alice" || !body.HideEmail || body.Theme != "forgejo-auto" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.UserSettings{
			FullName:  body.FullName,
			HideEmail: body.HideEmail,
			Theme:     body.Theme,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Users.UpdateSettings(context.Background(), &types.UserSettingsOptions{
		FullName:  "Alice",
		HideEmail: true,
		Theme:     "forgejo-auto",
	})
	if err != nil {
		t.Fatal(err)
	}
	if settings.FullName != "Alice" {
		t.Errorf("got full name=%q, want %q", settings.FullName, "Alice")
	}
	if !settings.HideEmail {
		t.Errorf("got hide_email=%v, want true", settings.HideEmail)
	}
}

func TestUserService_GetQuota_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.QuotaInfo{
			Groups: &types.QuotaGroupList{},
			Used: &types.QuotaUsed{
				Size: &types.QuotaUsedSize{
					Repos: &types.QuotaUsedSizeRepos{
						Public:  123,
						Private: 456,
					},
				},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	quota, err := f.Users.GetQuota(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if quota.Used == nil || quota.Used.Size == nil || quota.Used.Size.Repos == nil {
		t.Fatalf("quota usage was not decoded: %+v", quota)
	}
	if quota.Used.Size.Repos.Public != 123 || quota.Used.Size.Repos.Private != 456 {
		t.Errorf("unexpected repository quota usage: %+v", quota.Used.Size.Repos)
	}
}

func TestUserService_ListQuotaArtifacts_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/artifacts" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.QuotaUsedArtifact{
			{Name: "artifact-1", Size: 123, HTMLURL: "https://example.com/actions/runs/1"},
			{Name: "artifact-2", Size: 456, HTMLURL: "https://example.com/actions/runs/2"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	artifacts, err := f.Users.ListQuotaArtifacts(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(artifacts) != 2 {
		t.Fatalf("got %d artifacts, want 2", len(artifacts))
	}
	if artifacts[0].Name != "artifact-1" || artifacts[0].Size != 123 {
		t.Errorf("unexpected first artifact: %+v", artifacts[0])
	}
}

func TestUserService_IterQuotaArtifacts_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/artifacts" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.QuotaUsedArtifact{
			{Name: "artifact-1", Size: 123, HTMLURL: "https://example.com/actions/runs/1"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.QuotaUsedArtifact
	for artifact, err := range f.Users.IterQuotaArtifacts(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, artifact)
	}
	if len(got) != 1 {
		t.Fatalf("got %d artifacts, want 1", len(got))
	}
	if got[0].Name != "artifact-1" {
		t.Errorf("unexpected artifact: %+v", got[0])
	}
}

func TestUserService_ListQuotaAttachments_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/attachments" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.QuotaUsedAttachment{
			{Name: "issue-attachment.png", Size: 123, APIURL: "https://example.com/api/attachments/1"},
			{Name: "release-attachment.tar.gz", Size: 456, APIURL: "https://example.com/api/attachments/2"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachments, err := f.Users.ListQuotaAttachments(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(attachments) != 2 {
		t.Fatalf("got %d attachments, want 2", len(attachments))
	}
	if attachments[0].Name != "issue-attachment.png" || attachments[0].Size != 123 {
		t.Errorf("unexpected first attachment: %+v", attachments[0])
	}
}

func TestUserService_IterQuotaAttachments_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/attachments" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.QuotaUsedAttachment{
			{Name: "issue-attachment.png", Size: 123, APIURL: "https://example.com/api/attachments/1"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.QuotaUsedAttachment
	for attachment, err := range f.Users.IterQuotaAttachments(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, attachment)
	}
	if len(got) != 1 {
		t.Fatalf("got %d attachments, want 1", len(got))
	}
	if got[0].Name != "issue-attachment.png" {
		t.Errorf("unexpected attachment: %+v", got[0])
	}
}

func TestUserService_ListQuotaPackages_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/packages" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.QuotaUsedPackage{
			{Name: "pkg-one", Type: "container", Version: "1.0.0", Size: 123, HTMLURL: "https://example.com/packages/1"},
			{Name: "pkg-two", Type: "npm", Version: "2.0.0", Size: 456, HTMLURL: "https://example.com/packages/2"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	packages, err := f.Users.ListQuotaPackages(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(packages) != 2 {
		t.Fatalf("got %d packages, want 2", len(packages))
	}
	if packages[0].Name != "pkg-one" || packages[0].Type != "container" || packages[0].Size != 123 {
		t.Errorf("unexpected first package: %+v", packages[0])
	}
}

func TestUserService_IterQuotaPackages_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/quota/packages" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.QuotaUsedPackage{
			{Name: "pkg-one", Type: "container", Version: "1.0.0", Size: 123, HTMLURL: "https://example.com/packages/1"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.QuotaUsedPackage
	for pkg, err := range f.Users.IterQuotaPackages(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, pkg)
	}
	if len(got) != 1 {
		t.Fatalf("got %d packages, want 1", len(got))
	}
	if got[0].Name != "pkg-one" || got[0].Type != "container" {
		t.Errorf("unexpected package: %+v", got[0])
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

func TestUserService_ListBlockedUsers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/list_blocked" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.BlockedUser{
			{BlockID: 11},
			{BlockID: 12},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	blocked, err := f.Users.ListBlockedUsers(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(blocked) != 2 {
		t.Fatalf("got %d blocked users, want 2", len(blocked))
	}
	if blocked[0].BlockID != 11 {
		t.Errorf("unexpected first blocked user: %+v", blocked[0])
	}
}

func TestUserService_Block_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/block/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.Block(context.Background(), "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_CheckFollowing_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/following/bob" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	following, err := f.Users.CheckFollowing(context.Background(), "alice", "bob")
	if err != nil {
		t.Fatal(err)
	}
	if !following {
		t.Fatal("got following=false, want true")
	}
}

func TestUserService_CheckFollowing_Bad_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/following/bob" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	following, err := f.Users.CheckFollowing(context.Background(), "alice", "bob")
	if err != nil {
		t.Fatal(err)
	}
	if following {
		t.Fatal("got following=true, want false")
	}
}

func TestUserService_Unblock_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/unblock/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.Unblock(context.Background(), "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_IterBlockedUsers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/list_blocked" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.BlockedUser{
			{BlockID: 77},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for blocked, err := range f.Users.IterBlockedUsers(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if blocked.BlockID != 77 {
			t.Errorf("unexpected blocked user: %+v", blocked)
		}
	}
	if count != 1 {
		t.Fatalf("got %d blocked users, want 1", count)
	}
}

func TestUserService_ListMySubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Users.ListMySubscriptions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 {
		t.Fatalf("got %d repositories, want 1", len(repos))
	}
	if repos[0].FullName != "core/go-forge" {
		t.Errorf("got full name=%q, want %q", repos[0].FullName, "core/go-forge")
	}
}

func TestUserService_IterMySubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for repo, err := range f.Users.IterMySubscriptions(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if repo.FullName != "core/go-forge" {
			t.Errorf("got full name=%q, want %q", repo.FullName, "core/go-forge")
		}
	}
	if count != 1 {
		t.Fatalf("got %d repositories, want 1", count)
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

func TestUserService_UpdateAvatar_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/avatar" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.UpdateUserAvatarOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Image != "aGVsbG8=" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.UpdateAvatar(context.Background(), &types.UpdateUserAvatarOption{Image: "aGVsbG8="}); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_DeleteAvatar_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/avatar" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteAvatar(context.Background()); err != nil {
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

func TestUserService_ListSubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Users.ListSubscriptions(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 {
		t.Fatalf("got %d repositories, want 1", len(repos))
	}
	if repos[0].Name != "go-forge" {
		t.Errorf("got name=%q, want %q", repos[0].Name, "go-forge")
	}
}

func TestUserService_IterSubscriptions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/subscriptions" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for repo, err := range f.Users.IterSubscriptions(context.Background(), "alice") {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if repo.Name != "go-forge" {
			t.Errorf("got name=%q, want %q", repo.Name, "go-forge")
		}
	}
	if count != 1 {
		t.Fatalf("got %d repositories, want 1", count)
	}
}

func TestUserService_ListMyStarred_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/starred" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Users.ListMyStarred(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 {
		t.Fatalf("got %d repositories, want 1", len(repos))
	}
	if repos[0].FullName != "core/go-forge" {
		t.Errorf("got full_name=%q, want %q", repos[0].FullName, "core/go-forge")
	}
}

func TestUserService_IterMyStarred_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/starred" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{
			{Name: "go-forge", FullName: "core/go-forge"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for repo, err := range f.Users.IterMyStarred(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if repo.Name != "go-forge" {
			t.Errorf("got name=%q, want %q", repo.Name, "go-forge")
		}
	}
	if count != 1 {
		t.Fatalf("got %d repositories, want 1", count)
	}
}

func TestUserService_CheckStarring_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/starred/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	starring, err := f.Users.CheckStarring(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !starring {
		t.Fatal("got starring=false, want true")
	}
}

func TestUserService_CheckStarring_Bad_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/starred/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	starring, err := f.Users.CheckStarring(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if starring {
		t.Fatal("got starring=true, want false")
	}
}

func TestUserService_GetHeatmap_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/heatmap" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.UserHeatmapData{
			{Contributions: 3},
			{Contributions: 7},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	heatmap, err := f.Users.GetHeatmap(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(heatmap) != 2 {
		t.Fatalf("got %d heatmap points, want 2", len(heatmap))
	}
	if heatmap[0].Contributions != 3 || heatmap[1].Contributions != 7 {
		t.Errorf("unexpected heatmap data: %+v", heatmap)
	}
}
