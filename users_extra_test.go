package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestUserService_ListMyFollowing_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/following" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.User{{ID: 1, UserName: "alice"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Users.ListMyFollowing(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].UserName != "alice" {
		t.Fatalf("got %#v", users)
	}
}
