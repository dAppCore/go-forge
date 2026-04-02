package forge

import (
	"context"
	"fmt"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go/core"
)

func TestClient_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "token test-token" {
			t.Errorf("missing auth header")
		}
		if r.URL.Path != "/api/v1/user" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]string{"login": "virgil"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	var out map[string]string
	err := c.Get(context.Background(), "/api/v1/user", &out)
	if err != nil {
		t.Fatal(err)
	}
	if out["login"] != "virgil" {
		t.Errorf("got login=%q", out["login"])
	}
}

func TestClient_Post_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["name"] != "test-repo" {
			t.Errorf("wrong body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "name": "test-repo"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	body := map[string]string{"name": "test-repo"}
	var out map[string]any
	err := c.Post(context.Background(), "/api/v1/orgs/core/repos", body, &out)
	if err != nil {
		t.Fatal(err)
	}
	if out["name"] != "test-repo" {
		t.Errorf("got name=%v", out["name"])
	}
}

func TestClient_PostRaw_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if got := r.URL.Path; got != "/api/v1/markdown" {
			t.Errorf("wrong path: %s", got)
		}
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "98")
		w.Header().Set("X-RateLimit-Reset", "1700000001")
		w.Write([]byte("<p>Hello</p>"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	body := map[string]string{"text": "Hello"}
	got, err := c.PostRaw(context.Background(), "/api/v1/markdown", body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "<p>Hello</p>" {
		t.Errorf("got body=%q", string(got))
	}
	rl := c.RateLimit()
	if rl.Limit != 100 || rl.Remaining != 98 || rl.Reset != 1700000001 {
		t.Fatalf("unexpected rate limit: %+v", rl)
	}
}

func TestClient_Delete_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	err := c.Delete(context.Background(), "/api/v1/repos/core/test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetRaw_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if got := r.URL.Path; got != "/api/v1/signing-key.gpg" {
			t.Errorf("wrong path: %s", got)
		}
		w.Header().Set("X-RateLimit-Limit", "60")
		w.Header().Set("X-RateLimit-Remaining", "59")
		w.Header().Set("X-RateLimit-Reset", "1700000002")
		w.Write([]byte("key-data"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	got, err := c.GetRaw(context.Background(), "/api/v1/signing-key.gpg")
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "key-data" {
		t.Errorf("got body=%q", string(got))
	}
	rl := c.RateLimit()
	if rl.Limit != 60 || rl.Remaining != 59 || rl.Reset != 1700000002 {
		t.Fatalf("unexpected rate limit: %+v", rl)
	}
}

func TestClient_ServerError_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal error"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	err := c.Get(context.Background(), "/api/v1/user", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var apiErr *APIError
	if !core.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("got status=%d", apiErr.StatusCode)
	}
}

func TestClient_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	err := c.Get(context.Background(), "/api/v1/repos/x/y", nil)
	if !IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestClient_ContextCancellation_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "test-token")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	err := c.Get(ctx, "/api/v1/user", nil)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestClient_Options_Good(t *testing.T) {
	c := NewClient("https://forge.lthn.ai", "tok",
		WithUserAgent("go-forge/1.0"),
	)
	if c.userAgent != "go-forge/1.0" {
		t.Errorf("got user agent=%q", c.userAgent)
	}
	if got := c.UserAgent(); got != "go-forge/1.0" {
		t.Errorf("got UserAgent()=%q", got)
	}
}

func TestClient_HasToken_Good(t *testing.T) {
	c := NewClient("https://forge.lthn.ai", "tok")
	if !c.HasToken() {
		t.Fatal("expected HasToken to report configured token")
	}
}

func TestClient_HasToken_Bad(t *testing.T) {
	c := NewClient("https://forge.lthn.ai", "")
	if c.HasToken() {
		t.Fatal("expected HasToken to report missing token")
	}
}

func TestClient_WithHTTPClient_Good(t *testing.T) {
	custom := &http.Client{}
	c := NewClient("https://forge.lthn.ai", "tok", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("expected custom HTTP client to be set")
	}
	if got := c.HTTPClient(); got != custom {
		t.Error("expected HTTPClient() to return the configured HTTP client")
	}
}

func TestClient_String_Good(t *testing.T) {
	c := NewClient("https://forge.lthn.ai", "tok", WithUserAgent("go-forge/1.0"))
	got := fmt.Sprint(c)
	want := `forge.Client{baseURL="https://forge.lthn.ai", token=set, userAgent="go-forge/1.0"}`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := c.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", c); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestAPIError_Error_Good(t *testing.T) {
	e := &APIError{StatusCode: 404, Message: "not found", URL: "/api/v1/repos/x/y"}
	got := e.Error()
	want := "forge: /api/v1/repos/x/y 404: not found"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIsConflict_Match_Good(t *testing.T) {
	err := &APIError{StatusCode: http.StatusConflict, Message: "conflict", URL: "/test"}
	if !IsConflict(err) {
		t.Error("expected IsConflict to return true for 409")
	}
}

func TestIsConflict_NotConflict_Bad(t *testing.T) {
	err := &APIError{StatusCode: http.StatusNotFound, Message: "not found", URL: "/test"}
	if IsConflict(err) {
		t.Error("expected IsConflict to return false for 404")
	}
}

func TestIsForbidden_NotForbidden_Bad(t *testing.T) {
	err := &APIError{StatusCode: http.StatusNotFound, Message: "not found", URL: "/test"}
	if IsForbidden(err) {
		t.Error("expected IsForbidden to return false for 404")
	}
}

func TestClient_RateLimit_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "99")
		w.Header().Set("X-RateLimit-Reset", "1700000000")
		json.NewEncoder(w).Encode(map[string]string{"login": "test"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	var out map[string]string
	if err := c.Get(context.Background(), "/api/v1/user", &out); err != nil {
		t.Fatal(err)
	}

	rl := c.RateLimit()
	if rl.Limit != 100 {
		t.Errorf("got limit=%d, want 100", rl.Limit)
	}
	if rl.Remaining != 99 {
		t.Errorf("got remaining=%d, want 99", rl.Remaining)
	}
	if rl.Reset != 1700000000 {
		t.Errorf("got reset=%d, want 1700000000", rl.Reset)
	}
}

func TestClient_Forbidden_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"message": "forbidden"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.Get(context.Background(), "/api/v1/admin", nil)
	if !IsForbidden(err) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestClient_Conflict_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "already exists"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.Post(context.Background(), "/api/v1/repos", map[string]string{"name": "dup"}, nil)
	if !IsConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
