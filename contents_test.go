package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestContentService_ListContents_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("ref"); got != "main" {
			t.Errorf("got ref=%q, want %q", got, "main")
		}
		json.NewEncoder(w).Encode([]types.ContentsResponse{
			{Name: "README.md", Path: "README.md", Type: "file"},
			{Name: "docs", Path: "docs", Type: "dir"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	items, err := f.Contents.ListContents(context.Background(), "core", "go-forge", "main")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}
	if items[0].Name != "README.md" || items[1].Type != "dir" {
		t.Fatalf("unexpected results: %+v", items)
	}
}

func TestContentService_IterContents_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.ContentsResponse{
			{Name: "README.md", Path: "README.md", Type: "file"},
			{Name: "docs", Path: "docs", Type: "dir"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for item, err := range f.Contents.IterContents(context.Background(), "core", "go-forge", "") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, item.Name)
	}
	if len(got) != 2 {
		t.Fatalf("got %d items, want 2", len(got))
	}
	if got[0] != "README.md" || got[1] != "docs" {
		t.Fatalf("unexpected items: %+v", got)
	}
}

func TestContentService_GetFile_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.ContentsResponse{
			Name:     "README.md",
			Path:     "README.md",
			Type:     "file",
			Encoding: "base64",
			Content:  "IyBnby1mb3JnZQ==",
			SHA:      "abc123",
			Size:     12,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	file, err := f.Contents.GetFile(context.Background(), "core", "go-forge", "README.md")
	if err != nil {
		t.Fatal(err)
	}
	if file.Name != "README.md" {
		t.Errorf("got name=%q, want %q", file.Name, "README.md")
	}
	if file.Type != "file" {
		t.Errorf("got type=%q, want %q", file.Type, "file")
	}
	if file.SHA != "abc123" {
		t.Errorf("got sha=%q, want %q", file.SHA, "abc123")
	}
	if file.Content != "IyBnby1mb3JnZQ==" {
		t.Errorf("got content=%q, want %q", file.Content, "IyBnby1mb3JnZQ==")
	}
}

func TestContentService_CreateFile_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents/docs/new.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateFileOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.ContentBase64 != "bmV3IGZpbGU=" {
			t.Errorf("got content=%q, want %q", opts.ContentBase64, "bmV3IGZpbGU=")
		}
		json.NewEncoder(w).Encode(types.FileResponse{
			Content: &types.ContentsResponse{
				Name: "new.md",
				Path: "docs/new.md",
				Type: "file",
				SHA:  "def456",
			},
			Commit: &types.FileCommitResponse{
				SHA:     "commit789",
				Message: "create docs/new.md",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	resp, err := f.Contents.CreateFile(context.Background(), "core", "go-forge", "docs/new.md", &types.CreateFileOptions{
		ContentBase64: "bmV3IGZpbGU=",
		Message:       "create docs/new.md",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content.Name != "new.md" {
		t.Errorf("got name=%q, want %q", resp.Content.Name, "new.md")
	}
	if resp.Content.SHA != "def456" {
		t.Errorf("got sha=%q, want %q", resp.Content.SHA, "def456")
	}
	if resp.Commit.Message != "create docs/new.md" {
		t.Errorf("got commit message=%q, want %q", resp.Commit.Message, "create docs/new.md")
	}
}

func TestContentService_UpdateFile_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.UpdateFileOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.SHA != "abc123" {
			t.Errorf("got sha=%q, want %q", opts.SHA, "abc123")
		}
		json.NewEncoder(w).Encode(types.FileResponse{
			Content: &types.ContentsResponse{
				Name: "README.md",
				Path: "README.md",
				Type: "file",
				SHA:  "updated456",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	resp, err := f.Contents.UpdateFile(context.Background(), "core", "go-forge", "README.md", &types.UpdateFileOptions{
		ContentBase64: "dXBkYXRlZA==",
		SHA:           "abc123",
		Message:       "update README",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Content.SHA != "updated456" {
		t.Errorf("got sha=%q, want %q", resp.Content.SHA, "updated456")
	}
}

func TestContentService_DeleteFile_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/contents/old.txt" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.DeleteFileOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.SHA != "sha123" {
			t.Errorf("got sha=%q, want %q", opts.SHA, "sha123")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(types.FileDeleteResponse{})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Contents.DeleteFile(context.Background(), "core", "go-forge", "old.txt", &types.DeleteFileOptions{
		SHA:     "sha123",
		Message: "remove old file",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestContentService_GetRawFile_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/raw/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("# go-forge\n\nA Go client for Forgejo."))
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	data, err := f.Contents.GetRawFile(context.Background(), "core", "go-forge", "README.md")
	if err != nil {
		t.Fatal(err)
	}
	want := "# go-forge\n\nA Go client for Forgejo."
	if string(data) != want {
		t.Errorf("got %q, want %q", string(data), want)
	}
}

func TestContentService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "file not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Contents.GetFile(context.Background(), "core", "go-forge", "nonexistent.md")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestContentService_GetRawNotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "file not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Contents.GetRawFile(context.Background(), "core", "go-forge", "nonexistent.md")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
