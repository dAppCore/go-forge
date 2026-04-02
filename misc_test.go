package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestMiscService_RenderMarkdown_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/markdown" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.MarkdownOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Text != "# Hello" {
			t.Errorf("got text=%q, want %q", opts.Text, "# Hello")
		}
		if opts.Mode != "gfm" {
			t.Errorf("got mode=%q, want %q", opts.Mode, "gfm")
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello</h1>\n"))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	html, err := f.Misc.RenderMarkdown(context.Background(), "# Hello", "gfm")
	if err != nil {
		t.Fatal(err)
	}
	want := "<h1>Hello</h1>\n"
	if html != want {
		t.Errorf("got %q, want %q", html, want)
	}
}

func TestMiscService_RenderMarkup_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/markup" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.MarkupOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Text != "**Hello**" {
			t.Errorf("got text=%q, want %q", opts.Text, "**Hello**")
		}
		if opts.Mode != "gfm" {
			t.Errorf("got mode=%q, want %q", opts.Mode, "gfm")
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<p><strong>Hello</strong></p>\n"))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	html, err := f.Misc.RenderMarkup(context.Background(), "**Hello**", "gfm")
	if err != nil {
		t.Fatal(err)
	}
	want := "<p><strong>Hello</strong></p>\n"
	if html != want {
		t.Errorf("got %q, want %q", html, want)
	}
}

func TestMiscService_RenderMarkdownRaw_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/markdown/raw" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); !strings.HasPrefix(got, "text/plain") {
			t.Errorf("got content-type=%q, want text/plain", got)
		}
		if got := r.Header.Get("Accept"); got != "text/html" {
			t.Errorf("got accept=%q, want text/html", got)
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "# Hello" {
			t.Errorf("got body=%q, want %q", string(data), "# Hello")
		}
		w.Header().Set("X-RateLimit-Limit", "80")
		w.Header().Set("X-RateLimit-Remaining", "79")
		w.Header().Set("X-RateLimit-Reset", "1700000003")
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<h1>Hello</h1>\n"))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	html, err := f.Misc.RenderMarkdownRaw(context.Background(), "# Hello")
	if err != nil {
		t.Fatal(err)
	}
	want := "<h1>Hello</h1>\n"
	if html != want {
		t.Errorf("got %q, want %q", html, want)
	}
	rl := f.Client().RateLimit()
	if rl.Limit != 80 || rl.Remaining != 79 || rl.Reset != 1700000003 {
		t.Fatalf("unexpected rate limit: %+v", rl)
	}
}

func TestMiscService_GetVersion_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/version" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.ServerVersion{
			Version: "1.21.0",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	ver, err := f.Misc.GetVersion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if ver.Version != "1.21.0" {
		t.Errorf("got version=%q, want %q", ver.Version, "1.21.0")
	}
}

func TestMiscService_GetAPISettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/settings/api" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GeneralAPISettings{
			DefaultGitTreesPerPage: 25,
			DefaultMaxBlobSize:     4096,
			DefaultPagingNum:       1,
			MaxResponseItems:       500,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Misc.GetAPISettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.DefaultPagingNum != 1 || settings.MaxResponseItems != 500 {
		t.Fatalf("unexpected api settings: %+v", settings)
	}
}

func TestMiscService_GetAttachmentSettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/settings/attachment" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GeneralAttachmentSettings{
			AllowedTypes: "image/*",
			Enabled:      true,
			MaxFiles:     10,
			MaxSize:      1048576,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Misc.GetAttachmentSettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !settings.Enabled || settings.MaxFiles != 10 {
		t.Fatalf("unexpected attachment settings: %+v", settings)
	}
}

func TestMiscService_GetRepositorySettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/settings/repository" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GeneralRepoSettings{
			ForksDisabled:        true,
			HTTPGitDisabled:      true,
			LFSDisabled:          true,
			MigrationsDisabled:   true,
			MirrorsDisabled:      false,
			StarsDisabled:        true,
			TimeTrackingDisabled: false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Misc.GetRepositorySettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !settings.ForksDisabled || !settings.HTTPGitDisabled {
		t.Fatalf("unexpected repository settings: %+v", settings)
	}
}

func TestMiscService_GetUISettings_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/settings/ui" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GeneralUISettings{
			AllowedReactions: []string{"+1", "-1"},
			CustomEmojis:     []string{":forgejo:"},
			DefaultTheme:     "forgejo-auto",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	settings, err := f.Misc.GetUISettings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if settings.DefaultTheme != "forgejo-auto" || len(settings.AllowedReactions) != 2 {
		t.Fatalf("unexpected ui settings: %+v", settings)
	}
}

func TestMiscService_ListLicenses_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/licenses" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.LicensesTemplateListEntry{
			{Key: "mit", Name: "MIT License"},
			{Key: "gpl-3.0", Name: "GNU General Public License v3.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	licenses, err := f.Misc.ListLicenses(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(licenses) != 2 {
		t.Fatalf("got %d licenses, want 2", len(licenses))
	}
	if licenses[0].Key != "mit" {
		t.Errorf("got key=%q, want %q", licenses[0].Key, "mit")
	}
	if licenses[1].Key != "gpl-3.0" {
		t.Errorf("got key=%q, want %q", licenses[1].Key, "gpl-3.0")
	}
}

func TestMiscService_IterLicenses_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/licenses" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.LicensesTemplateListEntry{
			{Key: "mit", Name: "MIT License"},
			{Key: "gpl-3.0", Name: "GNU General Public License v3.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var names []string
	for item, err := range f.Misc.IterLicenses(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, item.Name)
	}
	if len(names) != 2 {
		t.Fatalf("got %d licences, want 2", len(names))
	}
	if names[0] != "MIT License" || names[1] != "GNU General Public License v3.0" {
		t.Fatalf("unexpected licences: %+v", names)
	}
}

func TestMiscService_GetLicense_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/licenses/mit" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.LicenseTemplateInfo{
			Key:  "mit",
			Name: "MIT License",
			Body: "MIT License body text...",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	lic, err := f.Misc.GetLicense(context.Background(), "mit")
	if err != nil {
		t.Fatal(err)
	}
	if lic.Key != "mit" {
		t.Errorf("got key=%q, want %q", lic.Key, "mit")
	}
	if lic.Name != "MIT License" {
		t.Errorf("got name=%q, want %q", lic.Name, "MIT License")
	}
}

func TestMiscService_ListGitignoreTemplates_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/gitignore/templates" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]string{"Go", "Python", "Node"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	names, err := f.Misc.ListGitignoreTemplates(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 3 {
		t.Fatalf("got %d templates, want 3", len(names))
	}
	if names[0] != "Go" {
		t.Errorf("got [0]=%q, want %q", names[0], "Go")
	}
}

func TestMiscService_IterGitignoreTemplates_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/gitignore/templates" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]string{"Go", "Python", "Node"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var names []string
	for item, err := range f.Misc.IterGitignoreTemplates(context.Background()) {
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, item)
	}
	if len(names) != 3 {
		t.Fatalf("got %d templates, want 3", len(names))
	}
	if names[0] != "Go" {
		t.Errorf("got [0]=%q, want %q", names[0], "Go")
	}
}

func TestMiscService_GetGitignoreTemplate_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/gitignore/templates/Go" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GitignoreTemplateInfo{
			Name:   "Go",
			Source: "*.exe\n*.test\n/vendor/",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tmpl, err := f.Misc.GetGitignoreTemplate(context.Background(), "Go")
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.Name != "Go" {
		t.Errorf("got name=%q, want %q", tmpl.Name, "Go")
	}
}

func TestMiscService_GetNodeInfo_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/nodeinfo" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.NodeInfo{
			Version: "2.1",
			Software: &types.NodeInfoSoftware{
				Name:    "forgejo",
				Version: "1.21.0",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	info, err := f.Misc.GetNodeInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Version != "2.1" {
		t.Errorf("got version=%q, want %q", info.Version, "2.1")
	}
	if info.Software.Name != "forgejo" {
		t.Errorf("got software name=%q, want %q", info.Software.Name, "forgejo")
	}
}

func TestMiscService_GetSigningKey_Good(t *testing.T) {
	want := "-----BEGIN PGP PUBLIC KEY BLOCK-----\n..."
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/signing-key.gpg" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(want))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	key, err := f.Misc.GetSigningKey(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if key != want {
		t.Fatalf("got %q, want %q", key, want)
	}
}

func TestMiscService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Misc.GetLicense(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
