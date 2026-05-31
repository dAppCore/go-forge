package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
)

// do is the discard-response convenience wrapper around doJSON. Cover its
// happy path and its error propagation directly.

func TestClient_do_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"ok": "yes"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	var out map[string]string
	if err := c.do(context.Background(), http.MethodGet, "/api/v1/whatever", nil, &out); err != nil {
		t.Fatal(err)
	}
	if out["ok"] != "yes" {
		t.Fatalf("got %v", out)
	}
}

func TestClient_do_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	if err := c.do(context.Background(), http.MethodGet, "/api/v1/whatever", nil, nil); err == nil {
		t.Fatal("expected error from do on a 500")
	}
}

func TestClient_authorizationHeader_Good(t *testing.T) {
	c := NewClient("http://localhost", "secret")
	if got := c.authorizationHeader(); got != "Bearer secret" {
		t.Fatalf("got %q, want Bearer secret", got)
	}
}

func TestClient_authorizationHeader_Bad(t *testing.T) {
	// Empty token yields no header.
	c := NewClient("http://localhost", "")
	if got := c.authorizationHeader(); got != "" {
		t.Fatalf("empty token should yield empty header, got %q", got)
	}
}

func TestClient_authorizationHeader_Ugly(t *testing.T) {
	// Nil receiver must be handled gracefully (used during construction).
	var c *Client
	if got := c.authorizationHeader(); got != "" {
		t.Fatalf("nil client should yield empty header, got %q", got)
	}
}

func TestParseRateLimitInt_Good(t *testing.T) {
	if got := parseRateLimitInt("42"); got != 42 {
		t.Fatalf("got %d, want 42", got)
	}
}

func TestParseRateLimitInt_Bad(t *testing.T) {
	// Non-numeric input falls back to 0 rather than erroring.
	if got := parseRateLimitInt("not-a-number"); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
	if got := parseRateLimitInt(""); got != 0 {
		t.Fatalf("empty got %d, want 0", got)
	}
}

func TestParseRateLimitInt64_Good(t *testing.T) {
	if got := parseRateLimitInt64("1717200000"); got != 1717200000 {
		t.Fatalf("got %d, want 1717200000", got)
	}
}

func TestParseRateLimitInt64_Bad(t *testing.T) {
	if got := parseRateLimitInt64("garbage"); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}

func TestClient_postMultipartJSON_Fields_Good(t *testing.T) {
	// A fields map (not just a file) exercises the WriteField loop, and an
	// empty fieldName skips the file part entirely.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		if got := r.FormValue("kind"); got != "label" {
			t.Errorf("got kind=%q, want label", got)
		}
		json.NewEncoder(w).Encode(map[string]string{"ok": "yes"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	var out map[string]string
	err := c.postMultipartJSON(
		context.Background(),
		"/api/v1/x",
		nil,
		map[string]string{"kind": "label"},
		"", // no file
		"",
		nil,
		&out,
	)
	if err != nil {
		t.Fatal(err)
	}
	if out["ok"] != "yes" {
		t.Fatalf("got %v", out)
	}
}

func TestClient_postMultipartJSON_NilContent_Good(t *testing.T) {
	// A non-empty fieldName with nil content writes an empty file part.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	err := c.postMultipartJSON(
		context.Background(),
		"/api/v1/x",
		nil,
		nil,
		"attachment",
		"empty.bin",
		nil,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_postMultipartJSON_DecodeError_Bad(t *testing.T) {
	// A 2xx body that is not valid JSON surfaces a decode error.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	var out map[string]string
	err := c.postMultipartJSON(
		context.Background(),
		"/api/v1/x",
		nil,
		nil,
		"attachment",
		"f.bin",
		core.NewReader("data"),
		&out,
	)
	if err == nil {
		t.Fatal("expected a decode error from a non-JSON body")
	}
}

func TestClient_postMultipartJSON_ParseURL_Bad(t *testing.T) {
	// A control character in the base URL fails URLParse before any request.
	c := NewClient("http://exa\x7fmple.com", "tok")
	err := c.postMultipartJSON(
		context.Background(),
		"/api/v1/x",
		nil,
		nil,
		"attachment",
		"f.bin",
		core.NewReader("data"),
		nil,
	)
	if err == nil {
		t.Fatal("expected a URL parse error")
	}
}

func TestReadBody_Good(t *testing.T) {
	data, err := readBody(core.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("got %q, want hello", string(data))
	}
}

func TestReadBody_Bad(t *testing.T) {
	// A non-reader argument fails core.ReadAll, which readBody surfaces as an
	// error rather than a nil/empty body.
	if _, err := readBody("not a reader"); err == nil {
		t.Fatal("expected an error reading a non-reader")
	}
}
