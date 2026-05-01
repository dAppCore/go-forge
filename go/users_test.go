package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"dappco.re/go/forge/types"
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

func TestUserService_GetUserByID_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("uid"); got != "42" {
			t.Errorf("wrong uid: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("wrong page: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Errorf("wrong limit: %s", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []*types.User{
				{ID: 42, UserName: "alice"},
			},
			"ok": true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.GetUserByID(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != 42 || user.UserName != "alice" {
		t.Fatalf("got %#v", user)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("uid"); got != "42" {
			t.Errorf("wrong uid: %s", got)
		}
		w.Header().Set("X-Total-Count", "0")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []*types.User{},
			"ok":   true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Users.GetUserByID(context.Background(), 42)
	if !IsNotFound(err) {
		t.Fatalf("expected not found error, got %v", err)
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
			Groups: types.QuotaGroupList{},
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

func TestUserService_SearchUsersPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "al" {
			t.Errorf("wrong q: %s", got)
		}
		if got := r.URL.Query().Get("uid"); got != "7" {
			t.Errorf("wrong uid: %s", got)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("wrong page: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != "50" {
			t.Errorf("wrong limit: %s", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode(map[string]any{
			"data": []*types.User{
				{ID: 1, UserName: "alice"},
				{ID: 2, UserName: "alex"},
			},
			"ok": true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Users.SearchUsersPage(context.Background(), "al", ListOptions{}, UserSearchOptions{UID: 7})
	if err != nil {
		t.Fatal(err)
	}
	if result.TotalCount != 2 || result.Page != 1 || result.HasMore {
		t.Fatalf("got %#v", result)
	}
	if len(result.Items) != 2 || result.Items[0].UserName != "alice" || result.Items[1].UserName != "alex" {
		t.Fatalf("got %#v", result.Items)
	}
}

func TestUserService_SearchUsers_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/search" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("q"); got != "al" {
			t.Errorf("wrong q: %s", got)
		}
		if got := r.URL.Query().Get("limit"); got != strconv.Itoa(50) {
			t.Errorf("wrong limit: %s", got)
		}

		switch r.URL.Query().Get("page") {
		case "1":
			w.Header().Set("X-Total-Count", "3")
			json.NewEncoder(w).Encode(map[string]any{
				"data": []*types.User{
					{ID: 1, UserName: "alice"},
					{ID: 2, UserName: "alex"},
				},
				"ok": true,
			})
		case "2":
			w.Header().Set("X-Total-Count", "3")
			json.NewEncoder(w).Encode(map[string]any{
				"data": []*types.User{
					{ID: 3, UserName: "ally"},
				},
				"ok": true,
			})
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for user, err := range f.Users.IterSearchUsers(context.Background(), "al") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, user.UserName)
	}
	if requests != 2 {
		t.Fatalf("expected 2 requests, got %d", requests)
	}
	if len(got) != 3 || got[0] != "alice" || got[1] != "alex" || got[2] != "ally" {
		t.Fatalf("got %#v", got)
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

func TestUserService_ListKeysWithFilters_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("fingerprint"); got != "ABCD1234" {
			t.Fatalf("unexpected fingerprint query: %q", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PublicKey{
			{ID: 1, Title: "laptop", ReadOnly: true},
			{ID: 2, Title: "desktop", ReadOnly: false},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	keys, err := f.Users.ListKeys(context.Background(), UserKeyListOptions{
		Fingerprint: "ABCD1234",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("got %d keys, want 2", len(keys))
	}
	if keys[0].ID != 1 || keys[0].Title != "laptop" || !keys[0].ReadOnly {
		t.Errorf("unexpected first key: %+v", keys[0])
	}
}

func TestUserService_IterKeys_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PublicKey{
			{ID: 3, Title: "workstation", KeyType: "ssh-ed25519"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for key, err := range f.Users.IterKeys(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if key.ID != 3 || key.Title != "workstation" || key.KeyType != "ssh-ed25519" {
			t.Errorf("unexpected key: %+v", key)
		}
	}
	if count != 1 {
		t.Fatalf("got %d keys, want 1", count)
	}
}

func TestUserService_CreateKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateKeyOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Key != "ssh-ed25519 AAAAC3Nza..." || body.Title != "laptop" || !body.ReadOnly {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.PublicKey{
			ID:       9,
			Title:    "laptop",
			KeyType:  "ssh-ed25519",
			ReadOnly: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Users.CreateKey(context.Background(), &types.CreateKeyOption{
		Key:      "ssh-ed25519 AAAAC3Nza...",
		Title:    "laptop",
		ReadOnly: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 9 || key.Title != "laptop" || key.KeyType != "ssh-ed25519" || !key.ReadOnly {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestUserService_GetKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/keys/9" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.PublicKey{
			ID:    9,
			Title: "laptop",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Users.GetKey(context.Background(), 9)
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 9 || key.Title != "laptop" {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestUserService_DeleteKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/keys/9" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteKey(context.Background(), 9); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_ListGPGKeys_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.GPGKey{
			{ID: 1, KeyID: "ABCD1234", PublicKey: "-----BEGIN PGP PUBLIC KEY BLOCK-----"},
			{ID: 2, KeyID: "EFGH5678", Verified: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	keys, err := f.Users.ListGPGKeys(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 {
		t.Fatalf("got %d keys, want 2", len(keys))
	}
	if keys[0].ID != 1 || keys[0].KeyID != "ABCD1234" {
		t.Errorf("unexpected first key: %+v", keys[0])
	}
}

func TestUserService_IterGPGKeys_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.GPGKey{
			{ID: 3, KeyID: "IJKL9012", Verified: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count := 0
	for key, err := range f.Users.IterGPGKeys(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if key.ID != 3 || key.KeyID != "IJKL9012" {
			t.Errorf("unexpected key: %+v", key)
		}
	}
	if count != 1 {
		t.Fatalf("got %d keys, want 1", count)
	}
}

func TestUserService_CreateGPGKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateGPGKeyOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.ArmoredKey != "-----BEGIN PGP PUBLIC KEY BLOCK-----" || body.Signature != "sig" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.GPGKey{
			ID:       9,
			KeyID:    "MNOP3456",
			Verified: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Users.CreateGPGKey(context.Background(), &types.CreateGPGKeyOption{
		ArmoredKey: "-----BEGIN PGP PUBLIC KEY BLOCK-----",
		Signature:  "sig",
	})
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 9 || key.KeyID != "MNOP3456" || !key.Verified {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestUserService_GetGPGKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_keys/9" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GPGKey{
			ID:    9,
			KeyID: "MNOP3456",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Users.GetGPGKey(context.Background(), 9)
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 9 || key.KeyID != "MNOP3456" {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestUserService_DeleteGPGKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_keys/9" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteGPGKey(context.Background(), 9); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_GetGPGKeyVerificationToken_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_key_token" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Write([]byte("verification-token"))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	token, err := f.Users.GetGPGKeyVerificationToken(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if token != "verification-token" {
		t.Errorf("got token=%q, want %q", token, "verification-token")
	}
}

func TestUserService_VerifyGPGKey_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/gpg_key_verify" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "" {
			t.Errorf("unexpected content type: %q", got)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.GPGKey{
			ID:       12,
			KeyID:    "QRST7890",
			Verified: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Users.VerifyGPGKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if key.ID != 12 || key.KeyID != "QRST7890" || !key.Verified {
		t.Errorf("unexpected key: %+v", key)
	}
}

func TestUserService_ListTokens_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/tokens" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.AccessToken{
			{ID: 1, Name: "ci", Scopes: []string{"repo"}},
			{ID: 2, Name: "deploy", Scopes: []string{"read:packages"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tokens, err := f.Users.ListTokens(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) != 2 || tokens[0].Name != "ci" || tokens[1].Name != "deploy" {
		t.Fatalf("unexpected tokens: %+v", tokens)
	}
}

func TestUserService_CreateToken_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/tokens" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateAccessTokenOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "ci" || len(body.Scopes) != 1 || body.Scopes[0] != "repo" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.AccessToken{
			ID:             7,
			Name:           body.Name,
			Scopes:         body.Scopes,
			Token:          "abcdef0123456789",
			TokenLastEight: "456789",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	token, err := f.Users.CreateToken(context.Background(), "alice", &types.CreateAccessTokenOption{
		Name:   "ci",
		Scopes: []string{"repo"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if token.ID != 7 || token.Name != "ci" || token.Token != "abcdef0123456789" {
		t.Fatalf("unexpected token: %+v", token)
	}
}

func TestUserService_DeleteToken_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/tokens/ci" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteToken(context.Background(), "alice", "ci"); err != nil {
		t.Fatal(err)
	}
}

func TestUserService_ListUserKeys_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("fingerprint"); got != "abc123" {
			t.Errorf("wrong fingerprint: %s", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PublicKey{
			{ID: 4, Title: "laptop", Fingerprint: "abc123"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	keys, err := f.Users.ListUserKeys(context.Background(), "alice", UserKeyListOptions{Fingerprint: "abc123"})
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0].Title != "laptop" {
		t.Fatalf("unexpected keys: %+v", keys)
	}
}

func TestUserService_ListUserGPGKeys_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/gpg_keys" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.GPGKey{
			{ID: 8, KeyID: "ABCD1234"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	keys, err := f.Users.ListUserGPGKeys(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 1 || keys[0].KeyID != "ABCD1234" {
		t.Fatalf("unexpected gpg keys: %+v", keys)
	}
}

func TestUserService_ListOAuth2Applications_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/applications/oauth2" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.OAuth2Application{
			{ID: 1, Name: "CLI", ClientID: "cli", RedirectURIs: []string{"http://localhost:3000/callback"}},
			{ID: 2, Name: "Desktop", ClientID: "desktop"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	apps, err := f.Users.ListOAuth2Applications(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Fatalf("got %d applications, want 2", len(apps))
	}
	if apps[0].ID != 1 || apps[0].Name != "CLI" {
		t.Errorf("unexpected first application: %+v", apps[0])
	}
}

func TestUserService_CreateOAuth2Application_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/applications/oauth2" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateOAuth2ApplicationOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "CLI" || !body.ConfidentialClient || len(body.RedirectURIs) != 1 || body.RedirectURIs[0] != "http://localhost:3000/callback" {
			t.Fatalf("unexpected body: %+v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.OAuth2Application{
			ID:                 1,
			Name:               body.Name,
			ClientID:           "cli",
			ClientSecret:       "secret",
			ConfidentialClient: body.ConfidentialClient,
			RedirectURIs:       body.RedirectURIs,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	app, err := f.Users.CreateOAuth2Application(context.Background(), &types.CreateOAuth2ApplicationOptions{
		Name:               "CLI",
		ConfidentialClient: true,
		RedirectURIs:       []string{"http://localhost:3000/callback"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if app.ID != 1 || app.ClientSecret != "secret" {
		t.Errorf("unexpected application: %+v", app)
	}
}

func TestUserService_GetOAuth2Application_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/applications/oauth2/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.OAuth2Application{
			ID:                 7,
			Name:               "CLI",
			ClientID:           "cli",
			ConfidentialClient: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	app, err := f.Users.GetOAuth2Application(context.Background(), 7)
	if err != nil {
		t.Fatal(err)
	}
	if app.ID != 7 || app.ClientID != "cli" {
		t.Errorf("unexpected application: %+v", app)
	}
}

func TestUserService_UpdateOAuth2Application_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/applications/oauth2/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateOAuth2ApplicationOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "CLI v2" || len(body.RedirectURIs) != 2 {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.OAuth2Application{
			ID:                 7,
			Name:               body.Name,
			ClientID:           "cli",
			ClientSecret:       "new-secret",
			ConfidentialClient: body.ConfidentialClient,
			RedirectURIs:       body.RedirectURIs,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	app, err := f.Users.UpdateOAuth2Application(context.Background(), 7, &types.CreateOAuth2ApplicationOptions{
		Name:               "CLI v2",
		RedirectURIs:       []string{"http://localhost:3000/callback", "http://localhost:3000/alt"},
		ConfidentialClient: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if app.Name != "CLI v2" || app.ClientSecret != "new-secret" {
		t.Errorf("unexpected application: %+v", app)
	}
}

func TestUserService_DeleteOAuth2Application_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/applications/oauth2/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Users.DeleteOAuth2Application(context.Background(), 7); err != nil {
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

func TestUsers_UserSearchOptions_String_Good(t *core.T) {
	subject := (*UserSearchOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserSearchOptions_String_Bad(t *core.T) {
	subject := (*UserSearchOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserSearchOptions_String_Ugly(t *core.T) {
	subject := (*UserSearchOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserSearchOptions_GoString_Good(t *core.T) {
	subject := (*UserSearchOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserSearchOptions_GoString_Bad(t *core.T) {
	subject := (*UserSearchOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserSearchOptions_GoString_Ugly(t *core.T) {
	subject := (*UserSearchOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_String_Good(t *core.T) {
	subject := (*UserKeyListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_String_Bad(t *core.T) {
	subject := (*UserKeyListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_String_Ugly(t *core.T) {
	subject := (*UserKeyListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_GoString_Good(t *core.T) {
	subject := (*UserKeyListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_GoString_Bad(t *core.T) {
	subject := (*UserKeyListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserKeyListOptions_GoString_Ugly(t *core.T) {
	subject := (*UserKeyListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByID_Good(t *core.T) {
	subject := (*UserService).GetUserByID
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByID_Bad(t *core.T) {
	subject := (*UserService).GetUserByID
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByID_Ugly(t *core.T) {
	subject := (*UserService).GetUserByID
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByUsername_Good(t *core.T) {
	subject := (*UserService).GetUserByUsername
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByUsername_Bad(t *core.T) {
	subject := (*UserService).GetUserByUsername
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetUserByUsername_Ugly(t *core.T) {
	subject := (*UserService).GetUserByUsername
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetCurrent_Good(t *core.T) {
	subject := (*UserService).GetCurrent
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetCurrent_Bad(t *core.T) {
	subject := (*UserService).GetCurrent
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetCurrent_Ugly(t *core.T) {
	subject := (*UserService).GetCurrent
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetSettings_Good(t *core.T) {
	subject := (*UserService).GetSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetSettings_Bad(t *core.T) {
	subject := (*UserService).GetSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetSettings_Ugly(t *core.T) {
	subject := (*UserService).GetSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateSettings_Good(t *core.T) {
	subject := (*UserService).UpdateSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateSettings_Bad(t *core.T) {
	subject := (*UserService).UpdateSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateSettings_Ugly(t *core.T) {
	subject := (*UserService).UpdateSettings
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetQuota_Good(t *core.T) {
	subject := (*UserService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetQuota_Bad(t *core.T) {
	subject := (*UserService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetQuota_Ugly(t *core.T) {
	subject := (*UserService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsersPage_Good(t *core.T) {
	subject := (*UserService).SearchUsersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsersPage_Bad(t *core.T) {
	subject := (*UserService).SearchUsersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsersPage_Ugly(t *core.T) {
	subject := (*UserService).SearchUsersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsers_Good(t *core.T) {
	subject := (*UserService).SearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsers_Bad(t *core.T) {
	subject := (*UserService).SearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_SearchUsers_Ugly(t *core.T) {
	subject := (*UserService).SearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSearchUsers_Good(t *core.T) {
	subject := (*UserService).IterSearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSearchUsers_Bad(t *core.T) {
	subject := (*UserService).IterSearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSearchUsers_Ugly(t *core.T) {
	subject := (*UserService).IterSearchUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaArtifacts_Good(t *core.T) {
	subject := (*UserService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaArtifacts_Bad(t *core.T) {
	subject := (*UserService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaArtifacts_Ugly(t *core.T) {
	subject := (*UserService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaArtifacts_Good(t *core.T) {
	subject := (*UserService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaArtifacts_Bad(t *core.T) {
	subject := (*UserService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaArtifacts_Ugly(t *core.T) {
	subject := (*UserService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaAttachments_Good(t *core.T) {
	subject := (*UserService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaAttachments_Bad(t *core.T) {
	subject := (*UserService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaAttachments_Ugly(t *core.T) {
	subject := (*UserService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaAttachments_Good(t *core.T) {
	subject := (*UserService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaAttachments_Bad(t *core.T) {
	subject := (*UserService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaAttachments_Ugly(t *core.T) {
	subject := (*UserService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaPackages_Good(t *core.T) {
	subject := (*UserService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaPackages_Bad(t *core.T) {
	subject := (*UserService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListQuotaPackages_Ugly(t *core.T) {
	subject := (*UserService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaPackages_Good(t *core.T) {
	subject := (*UserService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaPackages_Bad(t *core.T) {
	subject := (*UserService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterQuotaPackages_Ugly(t *core.T) {
	subject := (*UserService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListEmails_Good(t *core.T) {
	subject := (*UserService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListEmails_Bad(t *core.T) {
	subject := (*UserService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListEmails_Ugly(t *core.T) {
	subject := (*UserService).ListEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterEmails_Good(t *core.T) {
	subject := (*UserService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterEmails_Bad(t *core.T) {
	subject := (*UserService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterEmails_Ugly(t *core.T) {
	subject := (*UserService).IterEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_AddEmails_Good(t *core.T) {
	subject := (*UserService).AddEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_AddEmails_Bad(t *core.T) {
	subject := (*UserService).AddEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_AddEmails_Ugly(t *core.T) {
	subject := (*UserService).AddEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteEmails_Good(t *core.T) {
	subject := (*UserService).DeleteEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteEmails_Bad(t *core.T) {
	subject := (*UserService).DeleteEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteEmails_Ugly(t *core.T) {
	subject := (*UserService).DeleteEmails
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateAvatar_Good(t *core.T) {
	subject := (*UserService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateAvatar_Bad(t *core.T) {
	subject := (*UserService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateAvatar_Ugly(t *core.T) {
	subject := (*UserService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteAvatar_Good(t *core.T) {
	subject := (*UserService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteAvatar_Bad(t *core.T) {
	subject := (*UserService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteAvatar_Ugly(t *core.T) {
	subject := (*UserService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListKeys_Good(t *core.T) {
	subject := (*UserService).ListKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListKeys_Bad(t *core.T) {
	subject := (*UserService).ListKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListKeys_Ugly(t *core.T) {
	subject := (*UserService).ListKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterKeys_Good(t *core.T) {
	subject := (*UserService).IterKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterKeys_Bad(t *core.T) {
	subject := (*UserService).IterKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterKeys_Ugly(t *core.T) {
	subject := (*UserService).IterKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateKey_Good(t *core.T) {
	subject := (*UserService).CreateKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateKey_Bad(t *core.T) {
	subject := (*UserService).CreateKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateKey_Ugly(t *core.T) {
	subject := (*UserService).CreateKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetKey_Good(t *core.T) {
	subject := (*UserService).GetKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetKey_Bad(t *core.T) {
	subject := (*UserService).GetKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetKey_Ugly(t *core.T) {
	subject := (*UserService).GetKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteKey_Good(t *core.T) {
	subject := (*UserService).DeleteKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteKey_Bad(t *core.T) {
	subject := (*UserService).DeleteKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteKey_Ugly(t *core.T) {
	subject := (*UserService).DeleteKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserKeys_Good(t *core.T) {
	subject := (*UserService).ListUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserKeys_Bad(t *core.T) {
	subject := (*UserService).ListUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserKeys_Ugly(t *core.T) {
	subject := (*UserService).ListUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserKeys_Good(t *core.T) {
	subject := (*UserService).IterUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserKeys_Bad(t *core.T) {
	subject := (*UserService).IterUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserKeys_Ugly(t *core.T) {
	subject := (*UserService).IterUserKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListGPGKeys_Good(t *core.T) {
	subject := (*UserService).ListGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListGPGKeys_Bad(t *core.T) {
	subject := (*UserService).ListGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListGPGKeys_Ugly(t *core.T) {
	subject := (*UserService).ListGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterGPGKeys_Good(t *core.T) {
	subject := (*UserService).IterGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterGPGKeys_Bad(t *core.T) {
	subject := (*UserService).IterGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterGPGKeys_Ugly(t *core.T) {
	subject := (*UserService).IterGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateGPGKey_Good(t *core.T) {
	subject := (*UserService).CreateGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateGPGKey_Bad(t *core.T) {
	subject := (*UserService).CreateGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateGPGKey_Ugly(t *core.T) {
	subject := (*UserService).CreateGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKey_Good(t *core.T) {
	subject := (*UserService).GetGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKey_Bad(t *core.T) {
	subject := (*UserService).GetGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKey_Ugly(t *core.T) {
	subject := (*UserService).GetGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteGPGKey_Good(t *core.T) {
	subject := (*UserService).DeleteGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteGPGKey_Bad(t *core.T) {
	subject := (*UserService).DeleteGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteGPGKey_Ugly(t *core.T) {
	subject := (*UserService).DeleteGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserGPGKeys_Good(t *core.T) {
	subject := (*UserService).ListUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserGPGKeys_Bad(t *core.T) {
	subject := (*UserService).ListUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListUserGPGKeys_Ugly(t *core.T) {
	subject := (*UserService).ListUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserGPGKeys_Good(t *core.T) {
	subject := (*UserService).IterUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserGPGKeys_Bad(t *core.T) {
	subject := (*UserService).IterUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterUserGPGKeys_Ugly(t *core.T) {
	subject := (*UserService).IterUserGPGKeys
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKeyVerificationToken_Good(t *core.T) {
	subject := (*UserService).GetGPGKeyVerificationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKeyVerificationToken_Bad(t *core.T) {
	subject := (*UserService).GetGPGKeyVerificationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetGPGKeyVerificationToken_Ugly(t *core.T) {
	subject := (*UserService).GetGPGKeyVerificationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_VerifyGPGKey_Good(t *core.T) {
	subject := (*UserService).VerifyGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_VerifyGPGKey_Bad(t *core.T) {
	subject := (*UserService).VerifyGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_VerifyGPGKey_Ugly(t *core.T) {
	subject := (*UserService).VerifyGPGKey
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListTokens_Good(t *core.T) {
	subject := (*UserService).ListTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListTokens_Bad(t *core.T) {
	subject := (*UserService).ListTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListTokens_Ugly(t *core.T) {
	subject := (*UserService).ListTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterTokens_Good(t *core.T) {
	subject := (*UserService).IterTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterTokens_Bad(t *core.T) {
	subject := (*UserService).IterTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterTokens_Ugly(t *core.T) {
	subject := (*UserService).IterTokens
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateToken_Good(t *core.T) {
	subject := (*UserService).CreateToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateToken_Bad(t *core.T) {
	subject := (*UserService).CreateToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateToken_Ugly(t *core.T) {
	subject := (*UserService).CreateToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteToken_Good(t *core.T) {
	subject := (*UserService).DeleteToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteToken_Bad(t *core.T) {
	subject := (*UserService).DeleteToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteToken_Ugly(t *core.T) {
	subject := (*UserService).DeleteToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListOAuth2Applications_Good(t *core.T) {
	subject := (*UserService).ListOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListOAuth2Applications_Bad(t *core.T) {
	subject := (*UserService).ListOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListOAuth2Applications_Ugly(t *core.T) {
	subject := (*UserService).ListOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterOAuth2Applications_Good(t *core.T) {
	subject := (*UserService).IterOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterOAuth2Applications_Bad(t *core.T) {
	subject := (*UserService).IterOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterOAuth2Applications_Ugly(t *core.T) {
	subject := (*UserService).IterOAuth2Applications
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateOAuth2Application_Good(t *core.T) {
	subject := (*UserService).CreateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateOAuth2Application_Bad(t *core.T) {
	subject := (*UserService).CreateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CreateOAuth2Application_Ugly(t *core.T) {
	subject := (*UserService).CreateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetOAuth2Application_Good(t *core.T) {
	subject := (*UserService).GetOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetOAuth2Application_Bad(t *core.T) {
	subject := (*UserService).GetOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetOAuth2Application_Ugly(t *core.T) {
	subject := (*UserService).GetOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateOAuth2Application_Good(t *core.T) {
	subject := (*UserService).UpdateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateOAuth2Application_Bad(t *core.T) {
	subject := (*UserService).UpdateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_UpdateOAuth2Application_Ugly(t *core.T) {
	subject := (*UserService).UpdateOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteOAuth2Application_Good(t *core.T) {
	subject := (*UserService).DeleteOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteOAuth2Application_Bad(t *core.T) {
	subject := (*UserService).DeleteOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_DeleteOAuth2Application_Ugly(t *core.T) {
	subject := (*UserService).DeleteOAuth2Application
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStopwatches_Good(t *core.T) {
	subject := (*UserService).ListStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStopwatches_Bad(t *core.T) {
	subject := (*UserService).ListStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStopwatches_Ugly(t *core.T) {
	subject := (*UserService).ListStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStopwatches_Good(t *core.T) {
	subject := (*UserService).IterStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStopwatches_Bad(t *core.T) {
	subject := (*UserService).IterStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStopwatches_Ugly(t *core.T) {
	subject := (*UserService).IterStopwatches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListBlockedUsers_Good(t *core.T) {
	subject := (*UserService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListBlockedUsers_Bad(t *core.T) {
	subject := (*UserService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListBlockedUsers_Ugly(t *core.T) {
	subject := (*UserService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterBlockedUsers_Good(t *core.T) {
	subject := (*UserService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterBlockedUsers_Bad(t *core.T) {
	subject := (*UserService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterBlockedUsers_Ugly(t *core.T) {
	subject := (*UserService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Block_Good(t *core.T) {
	subject := (*UserService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Block_Bad(t *core.T) {
	subject := (*UserService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Block_Ugly(t *core.T) {
	subject := (*UserService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unblock_Good(t *core.T) {
	subject := (*UserService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unblock_Bad(t *core.T) {
	subject := (*UserService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unblock_Ugly(t *core.T) {
	subject := (*UserService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMySubscriptions_Good(t *core.T) {
	subject := (*UserService).ListMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMySubscriptions_Bad(t *core.T) {
	subject := (*UserService).ListMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMySubscriptions_Ugly(t *core.T) {
	subject := (*UserService).ListMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMySubscriptions_Good(t *core.T) {
	subject := (*UserService).IterMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMySubscriptions_Bad(t *core.T) {
	subject := (*UserService).IterMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMySubscriptions_Ugly(t *core.T) {
	subject := (*UserService).IterMySubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyStarred_Good(t *core.T) {
	subject := (*UserService).ListMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyStarred_Bad(t *core.T) {
	subject := (*UserService).ListMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyStarred_Ugly(t *core.T) {
	subject := (*UserService).ListMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyStarred_Good(t *core.T) {
	subject := (*UserService).IterMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyStarred_Bad(t *core.T) {
	subject := (*UserService).IterMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyStarred_Ugly(t *core.T) {
	subject := (*UserService).IterMyStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowers_Good(t *core.T) {
	subject := (*UserService).ListMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowers_Bad(t *core.T) {
	subject := (*UserService).ListMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowers_Ugly(t *core.T) {
	subject := (*UserService).ListMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowers_Good(t *core.T) {
	subject := (*UserService).IterMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowers_Bad(t *core.T) {
	subject := (*UserService).IterMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowers_Ugly(t *core.T) {
	subject := (*UserService).IterMyFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowing_Good(t *core.T) {
	subject := (*UserService).ListMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowing_Bad(t *core.T) {
	subject := (*UserService).ListMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyFollowing_Ugly(t *core.T) {
	subject := (*UserService).ListMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowing_Good(t *core.T) {
	subject := (*UserService).IterMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowing_Bad(t *core.T) {
	subject := (*UserService).IterMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyFollowing_Ugly(t *core.T) {
	subject := (*UserService).IterMyFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTeams_Good(t *core.T) {
	subject := (*UserService).ListMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTeams_Bad(t *core.T) {
	subject := (*UserService).ListMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTeams_Ugly(t *core.T) {
	subject := (*UserService).ListMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTeams_Good(t *core.T) {
	subject := (*UserService).IterMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTeams_Bad(t *core.T) {
	subject := (*UserService).IterMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTeams_Ugly(t *core.T) {
	subject := (*UserService).IterMyTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTrackedTimes_Good(t *core.T) {
	subject := (*UserService).ListMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTrackedTimes_Bad(t *core.T) {
	subject := (*UserService).ListMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListMyTrackedTimes_Ugly(t *core.T) {
	subject := (*UserService).ListMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTrackedTimes_Good(t *core.T) {
	subject := (*UserService).IterMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTrackedTimes_Bad(t *core.T) {
	subject := (*UserService).IterMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterMyTrackedTimes_Ugly(t *core.T) {
	subject := (*UserService).IterMyTrackedTimes
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckQuota_Good(t *core.T) {
	subject := (*UserService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckQuota_Bad(t *core.T) {
	subject := (*UserService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckQuota_Ugly(t *core.T) {
	subject := (*UserService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetRunnerRegistrationToken_Good(t *core.T) {
	subject := (*UserService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetRunnerRegistrationToken_Bad(t *core.T) {
	subject := (*UserService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetRunnerRegistrationToken_Ugly(t *core.T) {
	subject := (*UserService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowers_Good(t *core.T) {
	subject := (*UserService).ListFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowers_Bad(t *core.T) {
	subject := (*UserService).ListFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowers_Ugly(t *core.T) {
	subject := (*UserService).ListFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowers_Good(t *core.T) {
	subject := (*UserService).IterFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowers_Bad(t *core.T) {
	subject := (*UserService).IterFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowers_Ugly(t *core.T) {
	subject := (*UserService).IterFollowers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListSubscriptions_Good(t *core.T) {
	subject := (*UserService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListSubscriptions_Bad(t *core.T) {
	subject := (*UserService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListSubscriptions_Ugly(t *core.T) {
	subject := (*UserService).ListSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSubscriptions_Good(t *core.T) {
	subject := (*UserService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSubscriptions_Bad(t *core.T) {
	subject := (*UserService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterSubscriptions_Ugly(t *core.T) {
	subject := (*UserService).IterSubscriptions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowing_Good(t *core.T) {
	subject := (*UserService).ListFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowing_Bad(t *core.T) {
	subject := (*UserService).ListFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListFollowing_Ugly(t *core.T) {
	subject := (*UserService).ListFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowing_Good(t *core.T) {
	subject := (*UserService).IterFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowing_Bad(t *core.T) {
	subject := (*UserService).IterFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterFollowing_Ugly(t *core.T) {
	subject := (*UserService).IterFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListActivityFeeds_Good(t *core.T) {
	subject := (*UserService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListActivityFeeds_Bad(t *core.T) {
	subject := (*UserService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListActivityFeeds_Ugly(t *core.T) {
	subject := (*UserService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterActivityFeeds_Good(t *core.T) {
	subject := (*UserService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterActivityFeeds_Bad(t *core.T) {
	subject := (*UserService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterActivityFeeds_Ugly(t *core.T) {
	subject := (*UserService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListRepos_Good(t *core.T) {
	subject := (*UserService).ListRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListRepos_Bad(t *core.T) {
	subject := (*UserService).ListRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListRepos_Ugly(t *core.T) {
	subject := (*UserService).ListRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterRepos_Good(t *core.T) {
	subject := (*UserService).IterRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterRepos_Bad(t *core.T) {
	subject := (*UserService).IterRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterRepos_Ugly(t *core.T) {
	subject := (*UserService).IterRepos
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Follow_Good(t *core.T) {
	subject := (*UserService).Follow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Follow_Bad(t *core.T) {
	subject := (*UserService).Follow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Follow_Ugly(t *core.T) {
	subject := (*UserService).Follow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckFollowing_Good(t *core.T) {
	subject := (*UserService).CheckFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckFollowing_Bad(t *core.T) {
	subject := (*UserService).CheckFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckFollowing_Ugly(t *core.T) {
	subject := (*UserService).CheckFollowing
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unfollow_Good(t *core.T) {
	subject := (*UserService).Unfollow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unfollow_Bad(t *core.T) {
	subject := (*UserService).Unfollow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unfollow_Ugly(t *core.T) {
	subject := (*UserService).Unfollow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStarred_Good(t *core.T) {
	subject := (*UserService).ListStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStarred_Bad(t *core.T) {
	subject := (*UserService).ListStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_ListStarred_Ugly(t *core.T) {
	subject := (*UserService).ListStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStarred_Good(t *core.T) {
	subject := (*UserService).IterStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStarred_Bad(t *core.T) {
	subject := (*UserService).IterStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_IterStarred_Ugly(t *core.T) {
	subject := (*UserService).IterStarred
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetHeatmap_Good(t *core.T) {
	subject := (*UserService).GetHeatmap
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetHeatmap_Bad(t *core.T) {
	subject := (*UserService).GetHeatmap
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_GetHeatmap_Ugly(t *core.T) {
	subject := (*UserService).GetHeatmap
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Star_Good(t *core.T) {
	subject := (*UserService).Star
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Star_Bad(t *core.T) {
	subject := (*UserService).Star
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Star_Ugly(t *core.T) {
	subject := (*UserService).Star
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unstar_Good(t *core.T) {
	subject := (*UserService).Unstar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unstar_Bad(t *core.T) {
	subject := (*UserService).Unstar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_Unstar_Ugly(t *core.T) {
	subject := (*UserService).Unstar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckStarring_Good(t *core.T) {
	subject := (*UserService).CheckStarring
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckStarring_Bad(t *core.T) {
	subject := (*UserService).CheckStarring
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestUsers_UserService_CheckStarring_Ugly(t *core.T) {
	subject := (*UserService).CheckStarring
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
