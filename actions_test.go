package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestActionsService_ListRepoSecrets_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/secrets" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Secret{
			{Name: "DEPLOY_KEY"},
			{Name: "API_TOKEN"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	secrets, err := f.Actions.ListRepoSecrets(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(secrets) != 2 {
		t.Fatalf("got %d secrets, want 2", len(secrets))
	}
	if secrets[0].Name != "DEPLOY_KEY" {
		t.Errorf("got name=%q, want %q", secrets[0].Name, "DEPLOY_KEY")
	}
	if secrets[1].Name != "API_TOKEN" {
		t.Errorf("got name=%q, want %q", secrets[1].Name, "API_TOKEN")
	}
}

func TestActionsService_CreateRepoSecret_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/secrets/DEPLOY_KEY" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["data"] != "super-secret" {
			t.Errorf("got data=%q, want %q", body["data"], "super-secret")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.CreateRepoSecret(context.Background(), "core", "go-forge", "DEPLOY_KEY", "super-secret")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_DeleteRepoSecret_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/secrets/OLD_KEY" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.DeleteRepoSecret(context.Background(), "core", "go-forge", "OLD_KEY")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_ListRepoVariables_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/variables" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.ActionVariable{
			{Name: "CI_ENV", Data: "production"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	vars, err := f.Actions.ListRepoVariables(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(vars) != 1 {
		t.Fatalf("got %d variables, want 1", len(vars))
	}
	if vars[0].Name != "CI_ENV" {
		t.Errorf("got name=%q, want %q", vars[0].Name, "CI_ENV")
	}
}

func TestActionsService_CreateRepoVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/variables/CI_ENV" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateVariableOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Value != "staging" {
			t.Errorf("got value=%q, want %q", body.Value, "staging")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.CreateRepoVariable(context.Background(), "core", "go-forge", "CI_ENV", "staging")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_UpdateRepoVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/variables/CI_ENV" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.UpdateVariableOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "CI_ENV_NEW" {
			t.Errorf("got name=%q, want %q", body.Name, "CI_ENV_NEW")
		}
		if body.Value != "production" {
			t.Errorf("got value=%q, want %q", body.Value, "production")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.UpdateRepoVariable(context.Background(), "core", "go-forge", "CI_ENV", &types.UpdateVariableOption{
		Name:  "CI_ENV_NEW",
		Value: "production",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_DeleteRepoVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/variables/OLD_VAR" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.DeleteRepoVariable(context.Background(), "core", "go-forge", "OLD_VAR")
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_ListOrgSecrets_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/lethean/actions/secrets" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Secret{
			{Name: "ORG_SECRET"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	secrets, err := f.Actions.ListOrgSecrets(context.Background(), "lethean")
	if err != nil {
		t.Fatal(err)
	}
	if len(secrets) != 1 {
		t.Fatalf("got %d secrets, want 1", len(secrets))
	}
	if secrets[0].Name != "ORG_SECRET" {
		t.Errorf("got name=%q, want %q", secrets[0].Name, "ORG_SECRET")
	}
}

func TestActionsService_ListOrgVariables_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/lethean/actions/variables" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.ActionVariable{
			{Name: "ORG_VAR", Data: "org-value"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	vars, err := f.Actions.ListOrgVariables(context.Background(), "lethean")
	if err != nil {
		t.Fatal(err)
	}
	if len(vars) != 1 {
		t.Fatalf("got %d variables, want 1", len(vars))
	}
	if vars[0].Name != "ORG_VAR" {
		t.Errorf("got name=%q, want %q", vars[0].Name, "ORG_VAR")
	}
}

func TestActionsService_DispatchWorkflow_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/workflows/build.yml/dispatches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["ref"] != "main" {
			t.Errorf("got ref=%v, want %q", body["ref"], "main")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Actions.DispatchWorkflow(context.Background(), "core", "go-forge", "build.yml", map[string]any{
		"ref": "main",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_ListRepoTasks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/tasks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Errorf("got limit=%q, want %q", got, "25")
		}
		json.NewEncoder(w).Encode(types.ActionTaskResponse{
			Entries: []*types.ActionTask{
				{ID: 101, Name: "build"},
				{ID: 102, Name: "test"},
			},
			TotalCount: 2,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	resp, err := f.Actions.ListRepoTasks(context.Background(), "core", "go-forge", ListOptions{Page: 2, Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
	if resp.TotalCount != 2 {
		t.Fatalf("got total_count=%d, want 2", resp.TotalCount)
	}
	if len(resp.Entries) != 2 {
		t.Fatalf("got %d tasks, want 2", len(resp.Entries))
	}
	if resp.Entries[0].ID != 101 || resp.Entries[1].Name != "test" {
		t.Fatalf("unexpected tasks: %#v", resp.Entries)
	}
}

func TestActionsService_IterRepoTasks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/actions/tasks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		switch r.URL.Query().Get("page") {
		case "1":
			json.NewEncoder(w).Encode(types.ActionTaskResponse{
				Entries:    []*types.ActionTask{{ID: 1, Name: "build"}},
				TotalCount: 2,
			})
		case "2":
			json.NewEncoder(w).Encode(types.ActionTaskResponse{
				Entries:    []*types.ActionTask{{ID: 2, Name: "test"}},
				TotalCount: 2,
			})
		default:
			t.Fatalf("unexpected page %q", r.URL.Query().Get("page"))
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []types.ActionTask
	for task, err := range f.Actions.IterRepoTasks(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, task)
	}
	if len(got) != 2 {
		t.Fatalf("got %d tasks, want 2", len(got))
	}
	if got[0].ID != 1 || got[1].Name != "test" {
		t.Fatalf("unexpected tasks: %#v", got)
	}
}

func TestActionsService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Actions.ListRepoSecrets(context.Background(), "core", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
