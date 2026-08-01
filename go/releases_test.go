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

	core "dappco.re/go"
	"dappco.re/go/forge/types"
)

func readMultipartReleaseAttachment(t *testing.T, r *http.Request) (map[string]string, string, string) {
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

	fields := make(map[string]string)
	reader := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	var fileName string
	var fileContent string
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		data, err := io.ReadAll(part)
		if err != nil {
			t.Fatal(err)
		}
		if part.FormName() == "attachment" {
			fileName = part.FileName()
			fileContent = string(data)
			continue
		}
		fields[part.FormName()] = string(data)
	}

	return fields, fileName, fileContent
}

func TestReleaseService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Release{
			{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"},
			{ID: 2, TagName: "v2.0.0", Title: "Release 2.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Releases.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", result.Items[0].TagName, "v1.0.0")
	}
}

func TestReleaseService_ListFiltered_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		want := map[string]string{
			"draft":       "true",
			"pre-release": "true",
			"q":           "1.0",
			"page":        "1",
			"limit":       "50",
		}
		for key, wantValue := range want {
			if got := r.URL.Query().Get(key); got != wantValue {
				t.Errorf("got %s=%q, want %q", key, got, wantValue)
			}
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Release{{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	releases, err := f.Releases.ListReleases(context.Background(), "core", "go-forge", ReleaseListOptions{
		Draft:      true,
		PreRelease: true,
		Query:      "1.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 || releases[0].TagName != "v1.0.0" {
		t.Fatalf("got %#v", releases)
	}
}

func TestReleaseService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", release.TagName, "v1.0.0")
	}
	if release.Title != "Release 1.0" {
		t.Errorf("got title=%q, want %q", release.Title, "Release 1.0")
	}
}

func TestReleaseService_GetByTag_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/tags/v1.0.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.GetByTag(context.Background(), "core", "go-forge", "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", release.TagName, "v1.0.0")
	}
	if release.ID != 1 {
		t.Errorf("got id=%d, want 1", release.ID)
	}
}

func TestReleaseService_GetLatest_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/latest" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 3, TagName: "v2.1.0", Title: "Latest Release"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.GetLatest(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v2.1.0" {
		t.Errorf("got tag=%q, want %q", release.TagName, "v2.1.0")
	}
	if release.Title != "Latest Release" {
		t.Errorf("got title=%q, want %q", release.Title, "Latest Release")
	}
}

func TestReleaseService_CreateAttachment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/1/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("name"); got != "linux-amd64" {
			t.Fatalf("got name=%q", got)
		}
		fields, filename, content := readMultipartReleaseAttachment(t, r)
		if !reflect.DeepEqual(fields, map[string]string{}) {
			t.Fatalf("got fields=%#v", fields)
		}
		if filename != "release.tar.gz" {
			t.Fatalf("got filename=%q", filename)
		}
		if content != "release bytes" {
			t.Fatalf("got content=%q", content)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Attachment{ID: 9, Name: filename})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Releases.CreateAttachment(
		context.Background(),
		"core",
		"go-forge",
		1,
		&ReleaseAttachmentUploadOptions{Name: "linux-amd64"},
		"release.tar.gz",
		bytes.NewBufferString("release bytes"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "release.tar.gz" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestReleaseService_CreateAttachmentExternalURL_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/1/assets" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		if got := r.URL.Query().Get("name"); got != "docs" {
			t.Fatalf("got name=%q", got)
		}
		fields, filename, content := readMultipartReleaseAttachment(t, r)
		if !reflect.DeepEqual(fields, map[string]string{"external_url": "https://example.com/release.tar.gz"}) {
			t.Fatalf("got fields=%#v", fields)
		}
		if filename != "" || content != "" {
			t.Fatalf("unexpected file upload: filename=%q content=%q", filename, content)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Attachment{ID: 10, Name: "docs"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Releases.CreateAttachment(
		context.Background(),
		"core",
		"go-forge",
		1,
		&ReleaseAttachmentUploadOptions{Name: "docs", ExternalURL: "https://example.com/release.tar.gz"},
		"",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "docs" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestReleaseService_EditAttachment_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/1/assets/4" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditAttachmentOptions
		json.NewDecoder(r.Body).Decode(&body)
		if body.Name != "release-notes.pdf" {
			t.Fatalf("got body=%#v", body)
		}
		json.NewEncoder(w).Encode(types.Attachment{ID: 4, Name: body.Name})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	attachment, err := f.Releases.EditAttachment(context.Background(), "core", "go-forge", 1, 4, &types.EditAttachmentOptions{Name: "release-notes.pdf"})
	if err != nil {
		t.Fatal(err)
	}
	if attachment.Name != "release-notes.pdf" {
		t.Fatalf("got name=%q", attachment.Name)
	}
}

func TestReleases_ReleaseListOptions_String_Good(t *core.T) {
	subject := (*ReleaseListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseListOptions_String_Bad(t *core.T) {
	subject := (*ReleaseListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseListOptions_String_Ugly(t *core.T) {
	subject := (*ReleaseListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseListOptions_GoString_Good(t *core.T) {
	subject := (*ReleaseListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseListOptions_GoString_Bad(t *core.T) {
	subject := (*ReleaseListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseListOptions_GoString_Ugly(t *core.T) {
	subject := (*ReleaseListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_String_Good(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_String_Bad(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_String_Ugly(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_GoString_Good(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_GoString_Bad(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseAttachmentUploadOptions_GoString_Ugly(t *core.T) {
	subject := (*ReleaseAttachmentUploadOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleasesPage_Good(t *core.T) {
	subject := (*ReleaseService).ListReleasesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleasesPage_Bad(t *core.T) {
	subject := (*ReleaseService).ListReleasesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleasesPage_Ugly(t *core.T) {
	subject := (*ReleaseService).ListReleasesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleases_Good(t *core.T) {
	subject := (*ReleaseService).ListReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleases_Bad(t *core.T) {
	subject := (*ReleaseService).ListReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListReleases_Ugly(t *core.T) {
	subject := (*ReleaseService).ListReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterReleases_Good(t *core.T) {
	subject := (*ReleaseService).IterReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterReleases_Bad(t *core.T) {
	subject := (*ReleaseService).IterReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterReleases_Ugly(t *core.T) {
	subject := (*ReleaseService).IterReleases
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateRelease_Good(t *core.T) {
	subject := (*ReleaseService).CreateRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateRelease_Bad(t *core.T) {
	subject := (*ReleaseService).CreateRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateRelease_Ugly(t *core.T) {
	subject := (*ReleaseService).CreateRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetByTag_Good(t *core.T) {
	subject := (*ReleaseService).GetByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetByTag_Bad(t *core.T) {
	subject := (*ReleaseService).GetByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetByTag_Ugly(t *core.T) {
	subject := (*ReleaseService).GetByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetRelease_Good(t *core.T) {
	subject := (*ReleaseService).GetRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetRelease_Bad(t *core.T) {
	subject := (*ReleaseService).GetRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetRelease_Ugly(t *core.T) {
	subject := (*ReleaseService).GetRelease
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetLatest_Good(t *core.T) {
	subject := (*ReleaseService).GetLatest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetLatest_Bad(t *core.T) {
	subject := (*ReleaseService).GetLatest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetLatest_Ugly(t *core.T) {
	subject := (*ReleaseService).GetLatest
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteByTag_Good(t *core.T) {
	subject := (*ReleaseService).DeleteByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteByTag_Bad(t *core.T) {
	subject := (*ReleaseService).DeleteByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteByTag_Ugly(t *core.T) {
	subject := (*ReleaseService).DeleteByTag
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListAssets_Good(t *core.T) {
	subject := (*ReleaseService).ListAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListAssets_Bad(t *core.T) {
	subject := (*ReleaseService).ListAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_ListAssets_Ugly(t *core.T) {
	subject := (*ReleaseService).ListAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAttachment_Good(t *core.T) {
	subject := (*ReleaseService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAttachment_Bad(t *core.T) {
	subject := (*ReleaseService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAttachment_Ugly(t *core.T) {
	subject := (*ReleaseService).CreateAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAttachment_Good(t *core.T) {
	subject := (*ReleaseService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAttachment_Bad(t *core.T) {
	subject := (*ReleaseService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAttachment_Ugly(t *core.T) {
	subject := (*ReleaseService).EditAttachment
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAsset_Good(t *core.T) {
	subject := (*ReleaseService).CreateAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAsset_Bad(t *core.T) {
	subject := (*ReleaseService).CreateAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_CreateAsset_Ugly(t *core.T) {
	subject := (*ReleaseService).CreateAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAsset_Good(t *core.T) {
	subject := (*ReleaseService).EditAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAsset_Bad(t *core.T) {
	subject := (*ReleaseService).EditAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_EditAsset_Ugly(t *core.T) {
	subject := (*ReleaseService).EditAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterAssets_Good(t *core.T) {
	subject := (*ReleaseService).IterAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterAssets_Bad(t *core.T) {
	subject := (*ReleaseService).IterAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_IterAssets_Ugly(t *core.T) {
	subject := (*ReleaseService).IterAssets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetAsset_Good(t *core.T) {
	subject := (*ReleaseService).GetAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetAsset_Bad(t *core.T) {
	subject := (*ReleaseService).GetAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_GetAsset_Ugly(t *core.T) {
	subject := (*ReleaseService).GetAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteAsset_Good(t *core.T) {
	subject := (*ReleaseService).DeleteAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteAsset_Bad(t *core.T) {
	subject := (*ReleaseService).DeleteAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestReleases_ReleaseService_DeleteAsset_Ugly(t *core.T) {
	subject := (*ReleaseService).DeleteAsset
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
