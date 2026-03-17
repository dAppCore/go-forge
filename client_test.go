package forge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Good_Get(t *testing.T) {
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

func TestClient_Good_Post(t *testing.T) {
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

func TestClient_Good_Delete(t *testing.T) {
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

func TestClient_Bad_ServerError(t *testing.T) {
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
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("got status=%d", apiErr.StatusCode)
	}
}

func TestClient_Bad_NotFound(t *testing.T) {
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

func TestClient_Good_ContextCancellation(t *testing.T) {
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

func TestClient_Good_Options(t *testing.T) {
	c := NewClient("https://forge.lthn.ai", "tok",
		WithUserAgent("go-forge/1.0"),
	)
	if c.userAgent != "go-forge/1.0" {
		t.Errorf("got user agent=%q", c.userAgent)
	}
}

func TestAPIError_Good_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Message: "not found", URL: "/api/v1/repos/x/y"}
	got := err.Error()
	if !strings.Contains(got, "404") || !strings.Contains(got, "not found") || !strings.Contains(got, "/api/v1/repos/x/y") {
		t.Errorf("Error() = %q, want status code, message, and URL", got)
	}
}

func TestClient_Good_IsConflict(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "already exists"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.Post(context.Background(), "/api/v1/repos", nil, nil)
	if !IsConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestClient_Good_WithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c := NewClient("https://forge.lthn.ai", "tok", WithHTTPClient(custom))
	if c.httpClient != custom {
		t.Error("WithHTTPClient did not set custom client")
	}
}

func TestClient_Good_RateLimit(t *testing.T) {
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
		t.Errorf("got Limit=%d, want 100", rl.Limit)
	}
	if rl.Remaining != 99 {
		t.Errorf("got Remaining=%d, want 99", rl.Remaining)
	}
	if rl.Reset != 1700000000 {
		t.Errorf("got Reset=%d, want 1700000000", rl.Reset)
	}
}

func TestClient_Bad_IsForbidden(t *testing.T) {
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

func TestClient_Bad_ParseErrorPlainText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("bad gateway"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.Get(context.Background(), "/api/v1/user", nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 502 {
		t.Errorf("got status=%d, want 502", apiErr.StatusCode)
	}
	if apiErr.Message != "bad gateway" {
		t.Errorf("got message=%q, want %q", apiErr.Message, "bad gateway")
	}
}

func TestClient_Bad_ParseErrorEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.Get(context.Background(), "/api/v1/user", nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 503 {
		t.Errorf("got status=%d, want 503", apiErr.StatusCode)
	}
	// Empty body falls back to http.StatusText
	if apiErr.Message != "Service Unavailable" {
		t.Errorf("got message=%q, want %q", apiErr.Message, "Service Unavailable")
	}
}

func TestClient_Bad_IsNotFoundNonAPIError(t *testing.T) {
	if IsNotFound(errors.New("random error")) {
		t.Error("IsNotFound should return false for non-APIError")
	}
	if IsForbidden(errors.New("random error")) {
		t.Error("IsForbidden should return false for non-APIError")
	}
	if IsConflict(errors.New("random error")) {
		t.Error("IsConflict should return false for non-APIError")
	}
}
