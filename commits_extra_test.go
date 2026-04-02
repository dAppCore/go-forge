package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestCommitService_GetCombinedStatusByRef_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits/main/status" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.CombinedStatus{
			SHA:        "main",
			TotalCount: 3,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	status, err := f.Commits.GetCombinedStatusByRef(context.Background(), "core", "go-forge", "main")
	if err != nil {
		t.Fatal(err)
	}
	if status.SHA != "main" || status.TotalCount != 3 {
		t.Fatalf("got %#v", status)
	}
}
