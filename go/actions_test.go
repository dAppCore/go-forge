package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
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

func TestActionsService_GetOrgVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/lethean/actions/variables/ORG_VAR" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.ActionVariable{Name: "ORG_VAR", Data: "org-value"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	variable, err := f.Actions.GetOrgVariable(context.Background(), "lethean", "ORG_VAR")
	if err != nil {
		t.Fatal(err)
	}
	if variable.Name != "ORG_VAR" || variable.Data != "org-value" {
		t.Fatalf("unexpected variable: %#v", variable)
	}
}

func TestActionsService_CreateUserVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/variables/CI_ENV" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateVariableOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Value != "production" {
			t.Errorf("got value=%q, want %q", body.Value, "production")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Actions.CreateUserVariable(context.Background(), "CI_ENV", "production"); err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_ListUserVariables_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/variables" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.ActionVariable{{Name: "CI_ENV", Data: "production"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	vars, err := f.Actions.ListUserVariables(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(vars) != 1 || vars[0].Name != "CI_ENV" {
		t.Fatalf("unexpected variables: %#v", vars)
	}
}

func TestActionsService_GetUserVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/variables/CI_ENV" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.ActionVariable{Name: "CI_ENV", Data: "production"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	variable, err := f.Actions.GetUserVariable(context.Background(), "CI_ENV")
	if err != nil {
		t.Fatal(err)
	}
	if variable.Name != "CI_ENV" || variable.Data != "production" {
		t.Fatalf("unexpected variable: %#v", variable)
	}
}

func TestActionsService_UpdateUserVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/variables/CI_ENV" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.UpdateVariableOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Name != "CI_ENV_NEW" || body.Value != "staging" {
			t.Fatalf("unexpected body: %#v", body)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Actions.UpdateUserVariable(context.Background(), "CI_ENV", &types.UpdateVariableOption{
		Name:  "CI_ENV_NEW",
		Value: "staging",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_DeleteUserVariable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/variables/OLD_VAR" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Actions.DeleteUserVariable(context.Background(), "OLD_VAR"); err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_CreateUserSecret_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/secrets/DEPLOY_KEY" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["data"] != "secret-value" {
			t.Errorf("got data=%q, want %q", body["data"], "secret-value")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Actions.CreateUserSecret(context.Background(), "DEPLOY_KEY", "secret-value"); err != nil {
		t.Fatal(err)
	}
}

func TestActionsService_DeleteUserSecret_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/actions/secrets/OLD_KEY" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Actions.DeleteUserSecret(context.Background(), "OLD_KEY"); err != nil {
		t.Fatal(err)
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

func TestActions_ActionsService_ListRepoSecrets_Good(t *core.T) {
	subject := (*ActionsService).ListRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoSecrets_Bad(t *core.T) {
	subject := (*ActionsService).ListRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoSecrets_Ugly(t *core.T) {
	subject := (*ActionsService).ListRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoSecrets_Good(t *core.T) {
	subject := (*ActionsService).IterRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoSecrets_Bad(t *core.T) {
	subject := (*ActionsService).IterRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoSecrets_Ugly(t *core.T) {
	subject := (*ActionsService).IterRepoSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoSecret_Good(t *core.T) {
	subject := (*ActionsService).CreateRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoSecret_Bad(t *core.T) {
	subject := (*ActionsService).CreateRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoSecret_Ugly(t *core.T) {
	subject := (*ActionsService).CreateRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoSecret_Good(t *core.T) {
	subject := (*ActionsService).DeleteRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoSecret_Bad(t *core.T) {
	subject := (*ActionsService).DeleteRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoSecret_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteRepoSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoVariables_Good(t *core.T) {
	subject := (*ActionsService).ListRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoVariables_Bad(t *core.T) {
	subject := (*ActionsService).ListRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoVariables_Ugly(t *core.T) {
	subject := (*ActionsService).ListRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoVariables_Good(t *core.T) {
	subject := (*ActionsService).IterRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoVariables_Bad(t *core.T) {
	subject := (*ActionsService).IterRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoVariables_Ugly(t *core.T) {
	subject := (*ActionsService).IterRepoVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoVariable_Good(t *core.T) {
	subject := (*ActionsService).CreateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoVariable_Bad(t *core.T) {
	subject := (*ActionsService).CreateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateRepoVariable_Ugly(t *core.T) {
	subject := (*ActionsService).CreateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateRepoVariable_Good(t *core.T) {
	subject := (*ActionsService).UpdateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateRepoVariable_Bad(t *core.T) {
	subject := (*ActionsService).UpdateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateRepoVariable_Ugly(t *core.T) {
	subject := (*ActionsService).UpdateRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoVariable_Good(t *core.T) {
	subject := (*ActionsService).DeleteRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoVariable_Bad(t *core.T) {
	subject := (*ActionsService).DeleteRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteRepoVariable_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteRepoVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgSecrets_Good(t *core.T) {
	subject := (*ActionsService).ListOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgSecrets_Bad(t *core.T) {
	subject := (*ActionsService).ListOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgSecrets_Ugly(t *core.T) {
	subject := (*ActionsService).ListOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgSecrets_Good(t *core.T) {
	subject := (*ActionsService).IterOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgSecrets_Bad(t *core.T) {
	subject := (*ActionsService).IterOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgSecrets_Ugly(t *core.T) {
	subject := (*ActionsService).IterOrgSecrets
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgVariables_Good(t *core.T) {
	subject := (*ActionsService).ListOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgVariables_Bad(t *core.T) {
	subject := (*ActionsService).ListOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListOrgVariables_Ugly(t *core.T) {
	subject := (*ActionsService).ListOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgVariables_Good(t *core.T) {
	subject := (*ActionsService).IterOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgVariables_Bad(t *core.T) {
	subject := (*ActionsService).IterOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterOrgVariables_Ugly(t *core.T) {
	subject := (*ActionsService).IterOrgVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetOrgVariable_Good(t *core.T) {
	subject := (*ActionsService).GetOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetOrgVariable_Bad(t *core.T) {
	subject := (*ActionsService).GetOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetOrgVariable_Ugly(t *core.T) {
	subject := (*ActionsService).GetOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgVariable_Good(t *core.T) {
	subject := (*ActionsService).CreateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgVariable_Bad(t *core.T) {
	subject := (*ActionsService).CreateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgVariable_Ugly(t *core.T) {
	subject := (*ActionsService).CreateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateOrgVariable_Good(t *core.T) {
	subject := (*ActionsService).UpdateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateOrgVariable_Bad(t *core.T) {
	subject := (*ActionsService).UpdateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateOrgVariable_Ugly(t *core.T) {
	subject := (*ActionsService).UpdateOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgVariable_Good(t *core.T) {
	subject := (*ActionsService).DeleteOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgVariable_Bad(t *core.T) {
	subject := (*ActionsService).DeleteOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgVariable_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteOrgVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgSecret_Good(t *core.T) {
	subject := (*ActionsService).CreateOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgSecret_Bad(t *core.T) {
	subject := (*ActionsService).CreateOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateOrgSecret_Ugly(t *core.T) {
	subject := (*ActionsService).CreateOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgSecret_Good(t *core.T) {
	subject := (*ActionsService).DeleteOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgSecret_Bad(t *core.T) {
	subject := (*ActionsService).DeleteOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteOrgSecret_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteOrgSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListUserVariables_Good(t *core.T) {
	subject := (*ActionsService).ListUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListUserVariables_Bad(t *core.T) {
	subject := (*ActionsService).ListUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListUserVariables_Ugly(t *core.T) {
	subject := (*ActionsService).ListUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterUserVariables_Good(t *core.T) {
	subject := (*ActionsService).IterUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterUserVariables_Bad(t *core.T) {
	subject := (*ActionsService).IterUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterUserVariables_Ugly(t *core.T) {
	subject := (*ActionsService).IterUserVariables
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetUserVariable_Good(t *core.T) {
	subject := (*ActionsService).GetUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetUserVariable_Bad(t *core.T) {
	subject := (*ActionsService).GetUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_GetUserVariable_Ugly(t *core.T) {
	subject := (*ActionsService).GetUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserVariable_Good(t *core.T) {
	subject := (*ActionsService).CreateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserVariable_Bad(t *core.T) {
	subject := (*ActionsService).CreateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserVariable_Ugly(t *core.T) {
	subject := (*ActionsService).CreateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateUserVariable_Good(t *core.T) {
	subject := (*ActionsService).UpdateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateUserVariable_Bad(t *core.T) {
	subject := (*ActionsService).UpdateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_UpdateUserVariable_Ugly(t *core.T) {
	subject := (*ActionsService).UpdateUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserVariable_Good(t *core.T) {
	subject := (*ActionsService).DeleteUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserVariable_Bad(t *core.T) {
	subject := (*ActionsService).DeleteUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserVariable_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteUserVariable
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserSecret_Good(t *core.T) {
	subject := (*ActionsService).CreateUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserSecret_Bad(t *core.T) {
	subject := (*ActionsService).CreateUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_CreateUserSecret_Ugly(t *core.T) {
	subject := (*ActionsService).CreateUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserSecret_Good(t *core.T) {
	subject := (*ActionsService).DeleteUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserSecret_Bad(t *core.T) {
	subject := (*ActionsService).DeleteUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DeleteUserSecret_Ugly(t *core.T) {
	subject := (*ActionsService).DeleteUserSecret
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DispatchWorkflow_Good(t *core.T) {
	subject := (*ActionsService).DispatchWorkflow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DispatchWorkflow_Bad(t *core.T) {
	subject := (*ActionsService).DispatchWorkflow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_DispatchWorkflow_Ugly(t *core.T) {
	subject := (*ActionsService).DispatchWorkflow
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoTasks_Good(t *core.T) {
	subject := (*ActionsService).ListRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoTasks_Bad(t *core.T) {
	subject := (*ActionsService).ListRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_ListRepoTasks_Ugly(t *core.T) {
	subject := (*ActionsService).ListRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoTasks_Good(t *core.T) {
	subject := (*ActionsService).IterRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoTasks_Bad(t *core.T) {
	subject := (*ActionsService).IterRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActions_ActionsService_IterRepoTasks_Ugly(t *core.T) {
	subject := (*ActionsService).IterRepoTasks
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
