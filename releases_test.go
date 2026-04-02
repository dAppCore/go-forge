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

	"dappco.re/go/core/forge/types"
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
