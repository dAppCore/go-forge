package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// The membership / collaborator / flag predicate endpoints share a three-way
// shape: a 204 means true, a 404 means false (not an error), and any other
// error propagates. The generated tests cover the 204 case; these pin the
// 404-false and error-propagation branches.

// boolEndpoint is the signature shared by all predicate endpoints under test.
type boolEndpoint func(srvURL string) (bool, error)

func runBoolEndpoint(t *testing.T, status int, body any, call boolEndpoint) (bool, error) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if status != 0 {
			w.WriteHeader(status)
		}
		if body != nil {
			json.NewEncoder(w).Encode(body)
		}
	}))
	defer srv.Close()
	return call(srv.URL)
}

func TestOrgService_IsMember_NotFound_False(t *testing.T) {
	got, err := runBoolEndpoint(t, http.StatusNotFound, map[string]string{"message": "not found"}, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsMember(context.Background(), "core", "alice")
	})
	if err != nil || got {
		t.Fatalf("404 should be (false, nil), got (%v, %v)", got, err)
	}
}

func TestOrgService_IsMember_ServerError_Bad(t *testing.T) {
	_, err := runBoolEndpoint(t, http.StatusInternalServerError, map[string]string{"message": "boom"}, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsMember(context.Background(), "core", "alice")
	})
	if err == nil {
		t.Fatal("500 should propagate as an error")
	}
}

func TestOrgService_IsBlocked_NotFound_False(t *testing.T) {
	got, err := runBoolEndpoint(t, http.StatusNotFound, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsBlocked(context.Background(), "core", "alice")
	})
	if err != nil || got {
		t.Fatalf("404 should be (false, nil), got (%v, %v)", got, err)
	}
}

func TestOrgService_IsBlocked_ServerError_Bad(t *testing.T) {
	_, err := runBoolEndpoint(t, http.StatusBadGateway, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsBlocked(context.Background(), "core", "alice")
	})
	if err == nil {
		t.Fatal("502 should propagate as an error")
	}
}

func TestOrgService_IsPublicMember_NotFound_False(t *testing.T) {
	got, err := runBoolEndpoint(t, http.StatusNotFound, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsPublicMember(context.Background(), "core", "alice")
	})
	if err != nil || got {
		t.Fatalf("404 should be (false, nil), got (%v, %v)", got, err)
	}
}

func TestOrgService_IsPublicMember_ServerError_Bad(t *testing.T) {
	_, err := runBoolEndpoint(t, http.StatusForbidden, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Orgs.IsPublicMember(context.Background(), "core", "alice")
	})
	if err == nil {
		t.Fatal("403 should propagate as an error")
	}
}

func TestRepoService_CheckCollaborator_NotFound_False(t *testing.T) {
	got, err := runBoolEndpoint(t, http.StatusNotFound, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Repos.CheckCollaborator(context.Background(), "core", "go-forge", "alice")
	})
	if err != nil || got {
		t.Fatalf("404 should be (false, nil), got (%v, %v)", got, err)
	}
}

func TestRepoService_CheckCollaborator_ServerError_Bad(t *testing.T) {
	_, err := runBoolEndpoint(t, http.StatusInternalServerError, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Repos.CheckCollaborator(context.Background(), "core", "go-forge", "alice")
	})
	if err == nil {
		t.Fatal("500 should propagate as an error")
	}
}

func TestRepoService_HasFlag_NotFound_False(t *testing.T) {
	got, err := runBoolEndpoint(t, http.StatusNotFound, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Repos.HasFlag(context.Background(), "core", "go-forge", "beta")
	})
	if err != nil || got {
		t.Fatalf("404 should be (false, nil), got (%v, %v)", got, err)
	}
}

func TestRepoService_HasFlag_ServerError_Bad(t *testing.T) {
	_, err := runBoolEndpoint(t, http.StatusInternalServerError, nil, func(u string) (bool, error) {
		return NewForge(u, "tok").Repos.HasFlag(context.Background(), "core", "go-forge", "beta")
	})
	if err == nil {
		t.Fatal("500 should propagate as an error")
	}
}
