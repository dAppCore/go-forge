package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestOrgService_ListOrgs_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Organization{{ID: 1, Name: "core"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	orgs, err := f.Orgs.ListOrgs(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(orgs) != 1 || orgs[0].Name != "core" {
		t.Fatalf("got %#v", orgs)
	}
}

func TestOrgService_CreateOrg_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateOrgOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.UserName != "core" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.Organization{ID: 1, Name: body.UserName})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	org, err := f.Orgs.CreateOrg(context.Background(), &types.CreateOrgOption{UserName: "core"})
	if err != nil {
		t.Fatal(err)
	}
	if org.Name != "core" {
		t.Fatalf("got name=%q", org.Name)
	}
}
