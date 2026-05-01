package forge

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
)

type ax7Payload struct {
	Name string `json:"name,omitempty"`
}

type ax7Transport struct {
	status int
	body   string
	mu     sync.Mutex
	count  int
	path   string
}

func ax7NewTransport(status int) *ax7Transport {
	body := "null"
	if status >= http.StatusBadRequest {
		body = `{"message":"ax7 failure"}`
	}
	return &ax7Transport{status: status, body: body}
}

func (tr *ax7Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := req.Context().Err(); err != nil {
		return nil, err
	}
	tr.mu.Lock()
	tr.count++
	tr.path = req.URL.EscapedPath()
	tr.mu.Unlock()
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	header.Set("X-Total-Count", "0")
	header.Set("token", "ax7-token")
	header.Set("X-RateLimit-Limit", "100")
	header.Set("X-RateLimit-Remaining", "99")
	header.Set("X-RateLimit-Reset", "1700000000")
	body := tr.body
	if req.URL.EscapedPath() == "/api/v1/users/search" {
		body = `{"data":[{"id":1}]}`
	}
	return &http.Response{StatusCode: tr.status, Status: fmt.Sprintf("%d %s", tr.status, http.StatusText(tr.status)), Header: header, Body: io.NopCloser(strings.NewReader(body)), Request: req}, nil
}

func (tr *ax7Transport) Count() int       { tr.mu.Lock(); defer tr.mu.Unlock(); return tr.count }
func (tr *ax7Transport) LastPath() string { tr.mu.Lock(); defer tr.mu.Unlock(); return tr.path }
func ax7Client(status int) (*Client, *ax7Transport) {
	tr := ax7NewTransport(status)
	return NewClient("http://forge.test", "tok", WithHTTPClient(&http.Client{Transport: tr})), tr
}
func ax7Forge(status int) (*Forge, *ax7Transport) {
	tr := ax7NewTransport(status)
	return NewForge("http://forge.test", "tok", WithHTTPClient(&http.Client{Transport: tr})), tr
}
func ax7Resource(status int) (*Resource[ax7Payload, ax7Payload, ax7Payload], *ax7Transport) {
	c, tr := ax7Client(status)
	return NewResource[ax7Payload, ax7Payload, ax7Payload](c, "/api/v1/items/{id}"), tr
}
func ax7CanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}
func ax7AssertRequest(t *core.T, tr *ax7Transport) {
	core.AssertEqual(t, 1, tr.Count())
	core.AssertNotEmpty(t, tr.LastPath())
}

func TestAX7_ActionsService_ListRepoSecrets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListRepoSecrets(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoSecrets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListRepoSecrets(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoSecrets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListRepoSecrets(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterRepoSecrets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterRepoSecrets(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoSecrets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterRepoSecrets(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoSecrets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterRepoSecrets(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateRepoSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateRepoSecret(ctx, "core", "go-forge", "name", "secret")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateRepoSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateRepoSecret(ctx, "core", "go-forge", "name", "secret")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateRepoSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateRepoSecret(ctx, "core", "go-forge", "name", "secret")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteRepoSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteRepoSecret(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteRepoSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteRepoSecret(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteRepoSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteRepoSecret(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_ListRepoVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListRepoVariables(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListRepoVariables(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListRepoVariables(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterRepoVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterRepoVariables(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterRepoVariables(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterRepoVariables(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateRepoVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateRepoVariable(ctx, "core", "go-forge", "name", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateRepoVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateRepoVariable(ctx, "core", "go-forge", "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateRepoVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateRepoVariable(ctx, "core", "go-forge", "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_UpdateRepoVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.UpdateRepoVariable(ctx, "core", "go-forge", "name", &types.UpdateVariableOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateRepoVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.UpdateRepoVariable(ctx, "core", "go-forge", "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateRepoVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.UpdateRepoVariable(ctx, "core", "go-forge", "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteRepoVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteRepoVariable(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteRepoVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteRepoVariable(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteRepoVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteRepoVariable(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_ListOrgSecrets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListOrgSecrets(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListOrgSecrets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListOrgSecrets(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListOrgSecrets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListOrgSecrets(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterOrgSecrets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterOrgSecrets(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterOrgSecrets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterOrgSecrets(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterOrgSecrets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterOrgSecrets(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_ListOrgVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListOrgVariables(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListOrgVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListOrgVariables(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListOrgVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListOrgVariables(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterOrgVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterOrgVariables(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterOrgVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterOrgVariables(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterOrgVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterOrgVariables(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_GetOrgVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.GetOrgVariable(ctx, "core", "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_GetOrgVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.GetOrgVariable(ctx, "core", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_GetOrgVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.GetOrgVariable(ctx, "core", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateOrgVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateOrgVariable(ctx, "core", "name", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateOrgVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateOrgVariable(ctx, "core", "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateOrgVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateOrgVariable(ctx, "core", "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_UpdateOrgVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.UpdateOrgVariable(ctx, "core", "name", &types.UpdateVariableOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateOrgVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.UpdateOrgVariable(ctx, "core", "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateOrgVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.UpdateOrgVariable(ctx, "core", "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteOrgVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteOrgVariable(ctx, "core", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteOrgVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteOrgVariable(ctx, "core", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteOrgVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteOrgVariable(ctx, "core", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateOrgSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateOrgSecret(ctx, "core", "name", "secret")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateOrgSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateOrgSecret(ctx, "core", "name", "secret")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateOrgSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateOrgSecret(ctx, "core", "name", "secret")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteOrgSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteOrgSecret(ctx, "core", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteOrgSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteOrgSecret(ctx, "core", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteOrgSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteOrgSecret(ctx, "core", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_ListUserVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListUserVariables(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListUserVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListUserVariables(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListUserVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListUserVariables(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterUserVariables_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterUserVariables(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterUserVariables_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterUserVariables(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterUserVariables_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterUserVariables(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_GetUserVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.GetUserVariable(ctx, "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_GetUserVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.GetUserVariable(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_GetUserVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.GetUserVariable(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateUserVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateUserVariable(ctx, "name", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateUserVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateUserVariable(ctx, "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateUserVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateUserVariable(ctx, "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_UpdateUserVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.UpdateUserVariable(ctx, "name", &types.UpdateVariableOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateUserVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.UpdateUserVariable(ctx, "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_UpdateUserVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.UpdateUserVariable(ctx, "name", &types.UpdateVariableOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteUserVariable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteUserVariable(ctx, "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteUserVariable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteUserVariable(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteUserVariable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteUserVariable(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_CreateUserSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.CreateUserSecret(ctx, "name", "secret")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateUserSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.CreateUserSecret(ctx, "name", "secret")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_CreateUserSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.CreateUserSecret(ctx, "name", "secret")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DeleteUserSecret_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DeleteUserSecret(ctx, "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteUserSecret_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DeleteUserSecret(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DeleteUserSecret_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DeleteUserSecret(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_DispatchWorkflow_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Actions.DispatchWorkflow(ctx, "core", "go-forge", "ci.yml", map[string]any{"key": "value"})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DispatchWorkflow_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Actions.DispatchWorkflow(ctx, "core", "go-forge", "ci.yml", map[string]any{"key": "value"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_DispatchWorkflow_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Actions.DispatchWorkflow(ctx, "core", "go-forge", "ci.yml", map[string]any{"key": "value"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_ListRepoTasks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Actions.ListRepoTasks(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoTasks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Actions.ListRepoTasks(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_ListRepoTasks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Actions.ListRepoTasks(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_IterRepoTasks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Actions.IterRepoTasks(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoTasks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Actions.IterRepoTasks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActionsService_IterRepoTasks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Actions.IterRepoTasks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_GetInstanceActor_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.ActivityPub.GetInstanceActor(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetInstanceActor_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.ActivityPub.GetInstanceActor(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetInstanceActor_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.ActivityPub.GetInstanceActor(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_SendInstanceActorInbox_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.ActivityPub.SendInstanceActorInbox(ctx, &types.ForgeLike{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendInstanceActorInbox_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.ActivityPub.SendInstanceActorInbox(ctx, &types.ForgeLike{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendInstanceActorInbox_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.ActivityPub.SendInstanceActorInbox(ctx, &types.ForgeLike{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_GetRepositoryActor_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.ActivityPub.GetRepositoryActor(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetRepositoryActor_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.ActivityPub.GetRepositoryActor(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetRepositoryActor_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.ActivityPub.GetRepositoryActor(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_SendRepositoryInbox_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.ActivityPub.SendRepositoryInbox(ctx, 1, &types.ForgeLike{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendRepositoryInbox_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.ActivityPub.SendRepositoryInbox(ctx, 1, &types.ForgeLike{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendRepositoryInbox_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.ActivityPub.SendRepositoryInbox(ctx, 1, &types.ForgeLike{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_GetPersonActor_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.ActivityPub.GetPersonActor(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetPersonActor_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.ActivityPub.GetPersonActor(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_GetPersonActor_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.ActivityPub.GetPersonActor(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_SendPersonInbox_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.ActivityPub.SendPersonInbox(ctx, 1, &types.ForgeLike{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendPersonInbox_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.ActivityPub.SendPersonInbox(ctx, 1, &types.ForgeLike{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ActivityPubService_SendPersonInbox_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.ActivityPub.SendPersonInbox(ctx, 1, &types.ForgeLike{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminActionsRunListOptions_String_Good(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AdminActionsRunListOptions_String_Bad(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AdminActionsRunListOptions_String_Ugly(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AdminActionsRunListOptions_GoString_Good(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AdminActionsRunListOptions_GoString_Bad(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AdminActionsRunListOptions_GoString_Ugly(t *core.T) {
	value := AdminActionsRunListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AdminUnadoptedListOptions_String_Good(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AdminUnadoptedListOptions_String_Bad(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AdminUnadoptedListOptions_String_Ugly(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AdminUnadoptedListOptions_GoString_Good(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AdminUnadoptedListOptions_GoString_Bad(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AdminUnadoptedListOptions_GoString_Ugly(t *core.T) {
	value := AdminUnadoptedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AdminService_ListUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListUsers(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListUsers(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListUsers(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterUsers(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterUsers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterUsers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateUser(ctx, &types.CreateUserOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateUser(ctx, &types.CreateUserOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateUser(ctx, &types.CreateUserOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateUserKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateUserKey(ctx, "alice", &types.CreateKeyOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateUserKey(ctx, "alice", &types.CreateKeyOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateUserKey(ctx, "alice", &types.CreateKeyOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteUserKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteUserKey(ctx, "alice", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUserKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteUserKey(ctx, "alice", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUserKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteUserKey(ctx, "alice", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateUserOrg_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateUserOrg(ctx, "alice", &types.CreateOrgOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserOrg_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateUserOrg(ctx, "alice", &types.CreateOrgOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserOrg_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateUserOrg(ctx, "alice", &types.CreateOrgOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GetUserQuota_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.GetUserQuota(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetUserQuota_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.GetUserQuota(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetUserQuota_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.GetUserQuota(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_SetUserQuotaGroups_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.SetUserQuotaGroups(ctx, "alice", &types.SetUserQuotaGroupsOptions{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_SetUserQuotaGroups_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.SetUserQuotaGroups(ctx, "alice", &types.SetUserQuotaGroupsOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_SetUserQuotaGroups_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.SetUserQuotaGroups(ctx, "alice", &types.SetUserQuotaGroupsOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateUserRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateUserRepo(ctx, "alice", &types.CreateRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateUserRepo(ctx, "alice", &types.CreateRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateUserRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateUserRepo(ctx, "alice", &types.CreateRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_EditUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.EditUser(ctx, "alice", map[string]any{"key": "value"})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.EditUser(ctx, "alice", map[string]any{"key": "value"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.EditUser(ctx, "alice", map[string]any{"key": "value"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteUser(ctx, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteUser(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteUser(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_RenameUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.RenameUser(ctx, "alice", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RenameUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.RenameUser(ctx, "alice", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RenameUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.RenameUser(ctx, "alice", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListOrgs(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListOrgs(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListOrgs(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterOrgs(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListEmails(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListEmails(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListEmails(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterEmails(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterEmails(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterEmails(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListHooks(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListHooks(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListHooks(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterHooks(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterHooks(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterHooks(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GetHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.GetHook(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.GetHook(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.GetHook(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateHook(ctx, &types.CreateHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateHook(ctx, &types.CreateHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateHook(ctx, &types.CreateHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_EditHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.EditHook(ctx, 1, &types.EditHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.EditHook(ctx, 1, &types.EditHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.EditHook(ctx, 1, &types.EditHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteHook(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteHook(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteHook(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListQuotaGroups_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListQuotaGroups(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaGroups_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListQuotaGroups(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaGroups_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListQuotaGroups(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterQuotaGroups_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterQuotaGroups(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaGroups_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaGroups(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaGroups_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaGroups(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateQuotaGroup_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateQuotaGroup(ctx, &types.CreateQuotaGroupOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateQuotaGroup_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateQuotaGroup(ctx, &types.CreateQuotaGroupOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateQuotaGroup_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateQuotaGroup(ctx, &types.CreateQuotaGroupOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GetQuotaGroup_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.GetQuotaGroup(ctx, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetQuotaGroup_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.GetQuotaGroup(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetQuotaGroup_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.GetQuotaGroup(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteQuotaGroup_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteQuotaGroup(ctx, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteQuotaGroup_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteQuotaGroup(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteQuotaGroup_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteQuotaGroup(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_AddQuotaGroupRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.AddQuotaGroupRule(ctx, "value", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AddQuotaGroupRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.AddQuotaGroupRule(ctx, "value", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AddQuotaGroupRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.AddQuotaGroupRule(ctx, "value", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_RemoveQuotaGroupRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.RemoveQuotaGroupRule(ctx, "value", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RemoveQuotaGroupRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.RemoveQuotaGroupRule(ctx, "value", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RemoveQuotaGroupRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.RemoveQuotaGroupRule(ctx, "value", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListQuotaGroupUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListQuotaGroupUsers(ctx, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaGroupUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListQuotaGroupUsers(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaGroupUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListQuotaGroupUsers(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterQuotaGroupUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterQuotaGroupUsers(ctx, "value") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaGroupUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaGroupUsers(ctx, "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaGroupUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaGroupUsers(ctx, "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_AddQuotaGroupUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.AddQuotaGroupUser(ctx, "value", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AddQuotaGroupUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.AddQuotaGroupUser(ctx, "value", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AddQuotaGroupUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.AddQuotaGroupUser(ctx, "value", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_RemoveQuotaGroupUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.RemoveQuotaGroupUser(ctx, "value", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RemoveQuotaGroupUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.RemoveQuotaGroupUser(ctx, "value", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RemoveQuotaGroupUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.RemoveQuotaGroupUser(ctx, "value", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListQuotaRules_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListQuotaRules(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaRules_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListQuotaRules(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListQuotaRules_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListQuotaRules(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterQuotaRules_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterQuotaRules(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaRules_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaRules(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterQuotaRules_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterQuotaRules(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_CreateQuotaRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.CreateQuotaRule(ctx, &types.CreateQuotaRuleOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateQuotaRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.CreateQuotaRule(ctx, &types.CreateQuotaRuleOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_CreateQuotaRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.CreateQuotaRule(ctx, &types.CreateQuotaRuleOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GetQuotaRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.GetQuotaRule(ctx, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetQuotaRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.GetQuotaRule(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GetQuotaRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.GetQuotaRule(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_EditQuotaRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.EditQuotaRule(ctx, "value", &types.EditQuotaRuleOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditQuotaRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.EditQuotaRule(ctx, "value", &types.EditQuotaRuleOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_EditQuotaRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.EditQuotaRule(ctx, "value", &types.EditQuotaRuleOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteQuotaRule_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteQuotaRule(ctx, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteQuotaRule_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteQuotaRule(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteQuotaRule_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteQuotaRule(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListUnadoptedRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListUnadoptedRepos(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListUnadoptedRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListUnadoptedRepos(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListUnadoptedRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListUnadoptedRepos(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterUnadoptedRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterUnadoptedRepos(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterUnadoptedRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterUnadoptedRepos(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterUnadoptedRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterUnadoptedRepos(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_SearchEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.SearchEmails(ctx, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_SearchEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.SearchEmails(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_SearchEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.SearchEmails(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterSearchEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterSearchEmails(ctx, "value") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterSearchEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterSearchEmails(ctx, "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterSearchEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterSearchEmails(ctx, "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_RunCron_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.RunCron(ctx, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RunCron_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.RunCron(ctx, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_RunCron_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.RunCron(ctx, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListCron_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListCron(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListCron_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListCron(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListCron_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListCron(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterCron_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterCron(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterCron_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterCron(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterCron_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterCron(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_ListActionsRuns_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.ListActionsRuns(ctx, AdminActionsRunListOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListActionsRuns_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.ListActionsRuns(ctx, AdminActionsRunListOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_ListActionsRuns_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.ListActionsRuns(ctx, AdminActionsRunListOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_IterActionsRuns_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Admin.IterActionsRuns(ctx, AdminActionsRunListOptions{}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterActionsRuns_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Admin.IterActionsRuns(ctx, AdminActionsRunListOptions{}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_IterActionsRuns_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Admin.IterActionsRuns(ctx, AdminActionsRunListOptions{}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_AdoptRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.AdoptRepo(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AdoptRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.AdoptRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_AdoptRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.AdoptRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_DeleteUnadoptedRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Admin.DeleteUnadoptedRepo(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUnadoptedRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Admin.DeleteUnadoptedRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_DeleteUnadoptedRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Admin.DeleteUnadoptedRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GenerateRunnerToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Admin.GenerateRunnerToken(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GenerateRunnerToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Admin.GenerateRunnerToken(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_AdminService_GenerateRunnerToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Admin.GenerateRunnerToken(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_ListBranchesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.ListBranchesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranchesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.ListBranchesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranchesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.ListBranchesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_ListBranches_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.ListBranches(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranches_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.ListBranches(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranches_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.ListBranches(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_IterBranches_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Branches.IterBranches(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_IterBranches_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Branches.IterBranches(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_IterBranches_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Branches.IterBranches(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_CreateBranch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.CreateBranch(ctx, "core", "go-forge", &types.CreateBranchRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_CreateBranch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.CreateBranch(ctx, "core", "go-forge", &types.CreateBranchRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_CreateBranch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.CreateBranch(ctx, "core", "go-forge", &types.CreateBranchRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_GetBranch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.GetBranch(ctx, "core", "go-forge", "main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_GetBranch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.GetBranch(ctx, "core", "go-forge", "main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_GetBranch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.GetBranch(ctx, "core", "go-forge", "main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_UpdateBranch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Branches.UpdateBranch(ctx, "core", "go-forge", "main", &types.UpdateBranchRepoOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_UpdateBranch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Branches.UpdateBranch(ctx, "core", "go-forge", "main", &types.UpdateBranchRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_UpdateBranch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Branches.UpdateBranch(ctx, "core", "go-forge", "main", &types.UpdateBranchRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_DeleteBranch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Branches.DeleteBranch(ctx, "core", "go-forge", "main")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_DeleteBranch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Branches.DeleteBranch(ctx, "core", "go-forge", "main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_DeleteBranch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Branches.DeleteBranch(ctx, "core", "go-forge", "main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_ListBranchProtections_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.ListBranchProtections(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranchProtections_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.ListBranchProtections(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_ListBranchProtections_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.ListBranchProtections(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_IterBranchProtections_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Branches.IterBranchProtections(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_IterBranchProtections_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Branches.IterBranchProtections(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_IterBranchProtections_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Branches.IterBranchProtections(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_GetBranchProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.GetBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_GetBranchProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.GetBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_GetBranchProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.GetBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_CreateBranchProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.CreateBranchProtection(ctx, "core", "go-forge", &types.CreateBranchProtectionOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_CreateBranchProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.CreateBranchProtection(ctx, "core", "go-forge", &types.CreateBranchProtectionOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_CreateBranchProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.CreateBranchProtection(ctx, "core", "go-forge", &types.CreateBranchProtectionOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_EditBranchProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Branches.EditBranchProtection(ctx, "core", "go-forge", "name", &types.EditBranchProtectionOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_EditBranchProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Branches.EditBranchProtection(ctx, "core", "go-forge", "name", &types.EditBranchProtectionOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_EditBranchProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Branches.EditBranchProtection(ctx, "core", "go-forge", "name", &types.EditBranchProtectionOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_DeleteBranchProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Branches.DeleteBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_DeleteBranchProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Branches.DeleteBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_BranchService_DeleteBranchProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Branches.DeleteBranchProtection(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_APIError_Error_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusNotFound, URL: "/api/v1/items", Message: "missing"}
	got := err.Error()
	core.AssertContains(t, got, "missing")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_Error_Bad(t *core.T) {
	var err *APIError
	got := err.Error()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_Error_Ugly(t *core.T) {
	err := &APIError{}
	got := fmt.Sprintf("%#v", err)
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_APIError_String_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusNotFound, URL: "/api/v1/items", Message: "missing"}
	got := err.String()
	core.AssertContains(t, got, "missing")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_String_Bad(t *core.T) {
	var err *APIError
	got := err.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_String_Ugly(t *core.T) {
	err := &APIError{}
	got := fmt.Sprintf("%#v", err)
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_APIError_GoString_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusNotFound, URL: "/api/v1/items", Message: "missing"}
	got := err.GoString()
	core.AssertContains(t, got, "missing")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_GoString_Bad(t *core.T) {
	var err *APIError
	got := err.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_APIError_GoString_Ugly(t *core.T) {
	err := &APIError{}
	got := fmt.Sprintf("%#v", err)
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_IsNotFound_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusNotFound, URL: "/api/v1/items", Message: "status"}
	got := IsNotFound(err)
	core.AssertTrue(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsNotFound_Bad(t *core.T) {
	err := &APIError{StatusCode: http.StatusInternalServerError, URL: "/api/v1/items", Message: "status"}
	got := IsNotFound(err)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsNotFound_Ugly(t *core.T) {
	got := IsNotFound(nil)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_IsForbidden_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusForbidden, URL: "/api/v1/items", Message: "status"}
	got := IsForbidden(err)
	core.AssertTrue(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsForbidden_Bad(t *core.T) {
	err := &APIError{StatusCode: http.StatusInternalServerError, URL: "/api/v1/items", Message: "status"}
	got := IsForbidden(err)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsForbidden_Ugly(t *core.T) {
	got := IsForbidden(nil)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_IsConflict_Good(t *core.T) {
	err := &APIError{StatusCode: http.StatusConflict, URL: "/api/v1/items", Message: "status"}
	got := IsConflict(err)
	core.AssertTrue(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsConflict_Bad(t *core.T) {
	err := &APIError{StatusCode: http.StatusInternalServerError, URL: "/api/v1/items", Message: "status"}
	got := IsConflict(err)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IsConflict_Ugly(t *core.T) {
	got := IsConflict(nil)
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_WithHTTPClient_Good(t *core.T) {
	custom := &http.Client{Transport: ax7NewTransport(http.StatusOK)}
	c := NewClient("http://forge.test", "tok", WithHTTPClient(custom))
	core.AssertEqual(t, custom, c.HTTPClient())
}
func TestAX7_WithHTTPClient_Bad(t *core.T) {
	c := NewClient("http://forge.test", "tok", WithHTTPClient(nil))
	got := c.HTTPClient()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", c))
}
func TestAX7_WithHTTPClient_Ugly(t *core.T) {
	custom := &http.Client{Transport: ax7NewTransport(http.StatusInternalServerError)}
	c := NewClient("http://forge.test", "tok", WithHTTPClient(custom))
	core.AssertEqual(t, custom, c.HTTPClient())
}

func TestAX7_WithUserAgent_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok", WithUserAgent("ax7-agent"))
	got := c.UserAgent()
	core.AssertEqual(t, "ax7-agent", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_WithUserAgent_Bad(t *core.T) {
	c := NewClient("http://forge.test", "tok", WithUserAgent(""))
	got := c.UserAgent()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_WithUserAgent_Ugly(t *core.T) {
	c := NewClient("http://forge.test", "tok", WithUserAgent("agent/with spaces"))
	got := c.UserAgent()
	core.AssertContains(t, got, "spaces")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_RateLimit_String_Good(t *core.T) {
	value := RateLimit{Limit: 100, Remaining: 99, Reset: 1700000000}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RateLimit_String_Bad(t *core.T) {
	value := RateLimit{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RateLimit_String_Ugly(t *core.T) {
	value := RateLimit{Limit: 100, Remaining: 99, Reset: 1700000000}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RateLimit_GoString_Good(t *core.T) {
	value := RateLimit{Limit: 100, Remaining: 99, Reset: 1700000000}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RateLimit_GoString_Bad(t *core.T) {
	value := RateLimit{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RateLimit_GoString_Ugly(t *core.T) {
	value := RateLimit{Limit: 100, Remaining: 99, Reset: 1700000000}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_Client_BaseURL_Good(t *core.T) {
	c := NewClient("http://forge.test///", "tok")
	core.AssertEqual(t, "http://forge.test", c.BaseURL())
	core.AssertNotEmpty(t, core.Sprintf("%T", c))
}
func TestAX7_Client_BaseURL_Bad(t *core.T) {
	var c *Client
	got := c.BaseURL()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", c))
}
func TestAX7_Client_BaseURL_Ugly(t *core.T) {
	c := NewClient("", "tok")
	got := c.BaseURL()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", c))
}

func TestAX7_Client_RateLimit_Good(t *core.T) {
	c, _ := ax7Client(http.StatusOK)
	got := c.RateLimit()
	core.AssertEqual(t, 0, got.Limit)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_RateLimit_Bad(t *core.T) {
	var c *Client
	got := c.RateLimit()
	core.AssertEqual(t, RateLimit{}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_RateLimit_Ugly(t *core.T) {
	c, _ := ax7Client(http.StatusOK)
	var out map[string]any
	core.AssertNoError(t, c.Get(context.Background(), "/api/v1/user", &out))
	core.AssertEqual(t, 100, c.RateLimit().Limit)
}

func TestAX7_Client_UserAgent_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok")
	got := c.UserAgent()
	core.AssertContains(t, got, "go-forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_UserAgent_Bad(t *core.T) {
	var c *Client
	got := c.UserAgent()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_UserAgent_Ugly(t *core.T) {
	c := NewClient("http://forge.test", "tok", WithUserAgent(""))
	got := c.UserAgent()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Client_HTTPClient_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok")
	got := c.HTTPClient()
	core.AssertNotNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_HTTPClient_Bad(t *core.T) {
	var c *Client
	got := c.HTTPClient()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", c))
}
func TestAX7_Client_HTTPClient_Ugly(t *core.T) {
	custom := &http.Client{Transport: ax7NewTransport(http.StatusOK)}
	c := NewClient("http://forge.test", "tok", WithHTTPClient(custom))
	core.AssertEqual(t, custom, c.HTTPClient())
}

func TestAX7_Client_String_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok")
	got := c.String()
	core.AssertContains(t, got, "forge.Client")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_String_Bad(t *core.T) {
	var c *Client
	got := c.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_String_Ugly(t *core.T) {
	c := NewClient("http://forge.test", "")
	got := fmt.Sprintf("%#v", c)
	core.AssertContains(t, got, "token=unset")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Client_GoString_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok")
	got := c.GoString()
	core.AssertContains(t, got, "forge.Client")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_GoString_Bad(t *core.T) {
	var c *Client
	got := c.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_GoString_Ugly(t *core.T) {
	c := NewClient("http://forge.test", "")
	got := fmt.Sprintf("%#v", c)
	core.AssertContains(t, got, "token=unset")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Client_HasToken_Good(t *core.T) {
	c := NewClient("http://forge.test", "tok")
	got := c.HasToken()
	core.AssertTrue(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_HasToken_Bad(t *core.T) {
	c := NewClient("http://forge.test", "")
	got := c.HasToken()
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Client_HasToken_Ugly(t *core.T) {
	var c *Client
	got := c.HasToken()
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_NewClient_Good(t *core.T) {
	c := NewClient("http://forge.test///", "tok")
	core.AssertEqual(t, "http://forge.test", c.BaseURL())
	core.AssertTrue(t, c.HasToken())
}
func TestAX7_NewClient_Bad(t *core.T) {
	c := NewClient("", "")
	core.AssertEqual(t, "", c.BaseURL())
	core.AssertFalse(t, c.HasToken())
}
func TestAX7_NewClient_Ugly(t *core.T) {
	custom := &http.Client{Transport: ax7NewTransport(http.StatusOK)}
	c := NewClient("http://forge.test", "tok", WithHTTPClient(custom))
	core.AssertEqual(t, custom, c.HTTPClient())
}

func TestAX7_Client_Get_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.Get(ctx, "/api/v1/items", &map[string]string{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Get_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.Get(ctx, "/api/v1/items", &map[string]string{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Get_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.Get(ctx, "/api/v1/items", &map[string]string{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_Post_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.Post(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Post_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.Post(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Post_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.Post(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_Patch_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.Patch(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Patch_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.Patch(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Patch_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.Patch(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_Put_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.Put(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Put_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.Put(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Put_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.Put(ctx, "/api/v1/items", map[string]string{"name": "value"}, &map[string]string{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_Delete_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.Delete(ctx, "/api/v1/items")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Delete_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.Delete(ctx, "/api/v1/items")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_Delete_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.Delete(ctx, "/api/v1/items")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_DeleteWithBody_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	err := c.DeleteWithBody(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_DeleteWithBody_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	err := c.DeleteWithBody(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_DeleteWithBody_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	err := c.DeleteWithBody(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_PostRaw_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	got, err := c.PostRaw(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_PostRaw_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := c.PostRaw(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_PostRaw_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := c.PostRaw(ctx, "/api/v1/items", map[string]string{"name": "value"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Client_GetRaw_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	got, err := c.GetRaw(ctx, "/api/v1/items")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_GetRaw_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := c.GetRaw(ctx, "/api/v1/items")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Client_GetRaw_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := c.GetRaw(ctx, "/api/v1/items")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitListOptions_String_Good(t *core.T) {
	value := CommitListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_CommitListOptions_String_Bad(t *core.T) {
	value := CommitListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_CommitListOptions_String_Ugly(t *core.T) {
	value := CommitListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_CommitListOptions_GoString_Good(t *core.T) {
	value := CommitListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_CommitListOptions_GoString_Bad(t *core.T) {
	value := CommitListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_CommitListOptions_GoString_Ugly(t *core.T) {
	value := CommitListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_CommitService_List_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_List_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_List_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_ListAll_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListAll_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListAll_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_Iter_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Commits.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_Iter_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Commits.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_Iter_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Commits.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_Get_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_Get_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_Get_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_ListCommitsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.ListCommitsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListCommitsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.ListCommitsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListCommitsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.ListCommitsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_ListCommits_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.ListCommits(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListCommits_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.ListCommits(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListCommits_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.ListCommits(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_IterCommits_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Commits.IterCommits(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_IterCommits_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Commits.IterCommits(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_IterCommits_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Commits.IterCommits(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetCommit_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetCommit(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCommit_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetCommit(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCommit_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetCommit(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetDiffOrPatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetDiffOrPatch(ctx, "core", "go-forge", "abc123", "diff")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetDiffOrPatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetDiffOrPatch(ctx, "core", "go-forge", "abc123", "diff")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetDiffOrPatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetDiffOrPatch(ctx, "core", "go-forge", "abc123", "diff")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetPullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetPullRequest(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetPullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetPullRequest(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetPullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetPullRequest(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetCombinedStatus_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetCombinedStatus(ctx, "core", "go-forge", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCombinedStatus_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetCombinedStatus(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCombinedStatus_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetCombinedStatus(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetCombinedStatusByRef_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetCombinedStatusByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCombinedStatusByRef_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetCombinedStatusByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetCombinedStatusByRef_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetCombinedStatusByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_ListStatuses_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.ListStatuses(ctx, "core", "go-forge", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListStatuses_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.ListStatuses(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_ListStatuses_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.ListStatuses(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_IterStatuses_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Commits.IterStatuses(ctx, "core", "go-forge", "heads/main") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_IterStatuses_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Commits.IterStatuses(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_IterStatuses_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Commits.IterStatuses(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_CreateStatus_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.CreateStatus(ctx, "core", "go-forge", "abc123", &types.CreateStatusOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_CreateStatus_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.CreateStatus(ctx, "core", "go-forge", "abc123", &types.CreateStatusOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_CreateStatus_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.CreateStatus(ctx, "core", "go-forge", "abc123", &types.CreateStatusOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GetNote_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.GetNote(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetNote_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.GetNote(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_GetNote_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.GetNote(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_SetNote_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Commits.SetNote(ctx, "core", "go-forge", "abc123", "message")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_SetNote_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Commits.SetNote(ctx, "core", "go-forge", "abc123", "message")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_SetNote_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Commits.SetNote(ctx, "core", "go-forge", "abc123", "message")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_DeleteNote_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Commits.DeleteNote(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_DeleteNote_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Commits.DeleteNote(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_CommitService_DeleteNote_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Commits.DeleteNote(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ConfigPath_Good(t *core.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	got, err := ConfigPath()
	core.AssertNoError(t, err)
	core.AssertEqual(t, filepath.Join(home, ".config", "forge", "config.json"), got)
}
func TestAX7_ConfigPath_Bad(t *core.T) {
	t.Setenv("HOME", "")
	got, err := ConfigPath()
	core.AssertError(t, err)
	core.AssertEqual(t, "", got)
}
func TestAX7_ConfigPath_Ugly(t *core.T) {
	home := filepath.Join(t.TempDir(), "space dir")
	t.Setenv("HOME", home)
	got, err := ConfigPath()
	core.AssertNoError(t, err)
	core.AssertContains(t, got, "space dir")
}

func TestAX7_SaveConfig_Good(t *core.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	err := SaveConfig("http://forge.test", "tok")
	core.AssertNoError(t, err)
	_, statErr := os.Stat(filepath.Join(home, ".config", "forge", "config.json"))
	core.AssertNoError(t, statErr)
}
func TestAX7_SaveConfig_Bad(t *core.T) {
	t.Setenv("HOME", "")
	err := SaveConfig("http://forge.test", "tok")
	core.AssertError(t, err)
	core.AssertNotNil(t, err)
}
func TestAX7_SaveConfig_Ugly(t *core.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	err := SaveConfig("http://forge.test/path?q=1", "tok with spaces")
	core.AssertNoError(t, err)
	_, statErr := os.Stat(filepath.Join(home, ".config", "forge", "config.json"))
	core.AssertNoError(t, statErr)
}

func TestAX7_ResolveConfig_Good(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	url, token, err := ResolveConfig("http://forge.test", "tok")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "http://forge.test", url)
	core.AssertEqual(t, "tok", token)
}
func TestAX7_ResolveConfig_Bad(t *core.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	core.AssertNoError(t, os.MkdirAll(filepath.Join(home, ".config", "forge"), 0700))
	core.AssertNoError(t, os.WriteFile(filepath.Join(home, ".config", "forge", "config.json"), []byte("{bad"), 0600))
	_, _, err := ResolveConfig("", "")
	core.AssertError(t, err)
}
func TestAX7_ResolveConfig_Ugly(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "http://env.test")
	t.Setenv("FORGE_TOKEN", "env-token")
	url, token, err := ResolveConfig("", "")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "http://env.test", url)
	core.AssertEqual(t, "env-token", token)
}

func TestAX7_NewFromConfig_Good(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	fg, err := NewFromConfig("http://forge.test", "tok")
	core.AssertNoError(t, err)
	core.AssertTrue(t, fg.HasToken())
}
func TestAX7_NewFromConfig_Bad(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")
	fg, err := NewFromConfig("http://forge.test", "")
	core.AssertError(t, err)
	core.AssertNil(t, fg)
}
func TestAX7_NewFromConfig_Ugly(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "http://env.test")
	t.Setenv("FORGE_TOKEN", "env-token")
	fg, err := NewFromConfig("", "")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "http://env.test", fg.BaseURL())
}

func TestAX7_NewForgeFromConfig_Good(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	fg, err := NewForgeFromConfig("http://forge.test", "tok")
	core.AssertNoError(t, err)
	core.AssertTrue(t, fg.HasToken())
}
func TestAX7_NewForgeFromConfig_Bad(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")
	fg, err := NewForgeFromConfig("http://forge.test", "")
	core.AssertError(t, err)
	core.AssertNil(t, fg)
}
func TestAX7_NewForgeFromConfig_Ugly(t *core.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "http://env.test")
	t.Setenv("FORGE_TOKEN", "env-token")
	fg, err := NewForgeFromConfig("", "")
	core.AssertNoError(t, err)
	core.AssertEqual(t, "http://env.test", fg.BaseURL())
}

func TestAX7_ContentService_ListContents_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.ListContents(ctx, "core", "go-forge", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_ListContents_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.ListContents(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_ListContents_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.ListContents(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_IterContents_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Contents.IterContents(ctx, "core", "go-forge", "heads/main") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_IterContents_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Contents.IterContents(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_IterContents_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Contents.IterContents(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_GetFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.GetFile(ctx, "core", "go-forge", "README.md")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.GetFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.GetFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_GetContents_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.GetContents(ctx, "core", "go-forge", "README.md")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetContents_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.GetContents(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetContents_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.GetContents(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_CreateFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.CreateFile(ctx, "core", "go-forge", "README.md", &types.CreateFileOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_CreateFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.CreateFile(ctx, "core", "go-forge", "README.md", &types.CreateFileOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_CreateFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.CreateFile(ctx, "core", "go-forge", "README.md", &types.CreateFileOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_UpdateFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.UpdateFile(ctx, "core", "go-forge", "README.md", &types.UpdateFileOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_UpdateFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.UpdateFile(ctx, "core", "go-forge", "README.md", &types.UpdateFileOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_UpdateFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.UpdateFile(ctx, "core", "go-forge", "README.md", &types.UpdateFileOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_DeleteFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Contents.DeleteFile(ctx, "core", "go-forge", "README.md", &types.DeleteFileOptions{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_DeleteFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Contents.DeleteFile(ctx, "core", "go-forge", "README.md", &types.DeleteFileOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_DeleteFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Contents.DeleteFile(ctx, "core", "go-forge", "README.md", &types.DeleteFileOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_GetRawFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Contents.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetRawFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Contents.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ContentService_GetRawFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Contents.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NewForge_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	core.AssertNotNil(t, fg.Repos)
	core.AssertTrue(t, fg.HasToken())
}
func TestAX7_NewForge_Bad(t *core.T) {
	fg := NewForge("", "")
	core.AssertEqual(t, "", fg.BaseURL())
	core.AssertFalse(t, fg.HasToken())
}
func TestAX7_NewForge_Ugly(t *core.T) {
	fg := NewForge("http://forge.test///", "tok", WithUserAgent("ax7"))
	core.AssertEqual(t, "http://forge.test", fg.BaseURL())
	core.AssertEqual(t, "ax7", fg.UserAgent())
}

func TestAX7_Forge_Client_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.Client()
	core.AssertNotNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_Client_Bad(t *core.T) {
	var fg *Forge
	got := fg.Client()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", fg))
}
func TestAX7_Forge_Client_Ugly(t *core.T) {
	fg := &Forge{}
	got := fg.Client()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", fg))
}

func TestAX7_Forge_BaseURL_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.BaseURL()
	core.AssertNotEmpty(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_BaseURL_Bad(t *core.T) {
	var fg *Forge
	got := fg.BaseURL()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_BaseURL_Ugly(t *core.T) {
	fg := &Forge{}
	got := fg.BaseURL()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Forge_RateLimit_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.RateLimit()
	core.AssertEqual(t, RateLimit{}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_RateLimit_Bad(t *core.T) {
	var fg *Forge
	got := fg.RateLimit()
	core.AssertEqual(t, RateLimit{}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_RateLimit_Ugly(t *core.T) {
	fg := &Forge{}
	got := fg.RateLimit()
	core.AssertEqual(t, RateLimit{}, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Forge_UserAgent_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.UserAgent()
	core.AssertNotEmpty(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_UserAgent_Bad(t *core.T) {
	var fg *Forge
	got := fg.UserAgent()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_UserAgent_Ugly(t *core.T) {
	fg := &Forge{}
	got := fg.UserAgent()
	core.AssertEqual(t, "", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Forge_HTTPClient_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.HTTPClient()
	core.AssertNotNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_HTTPClient_Bad(t *core.T) {
	var fg *Forge
	got := fg.HTTPClient()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", fg))
}
func TestAX7_Forge_HTTPClient_Ugly(t *core.T) {
	fg := &Forge{}
	got := fg.HTTPClient()
	core.AssertNil(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", fg))
}

func TestAX7_Forge_HasToken_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.HasToken()
	core.AssertTrue(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_HasToken_Bad(t *core.T) {
	fg := NewForge("http://forge.test", "")
	got := fg.HasToken()
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_HasToken_Ugly(t *core.T) {
	var fg *Forge
	got := fg.HasToken()
	core.AssertFalse(t, got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Forge_String_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.String()
	core.AssertContains(t, got, "forge.Forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_String_Bad(t *core.T) {
	var fg *Forge
	got := fg.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_String_Ugly(t *core.T) {
	fg := &Forge{}
	got := fmt.Sprintf("%#v", fg)
	core.AssertContains(t, got, "client=<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Forge_GoString_Good(t *core.T) {
	fg := NewForge("http://forge.test", "tok")
	got := fg.GoString()
	core.AssertContains(t, got, "forge.Forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_GoString_Bad(t *core.T) {
	var fg *Forge
	got := fg.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Forge_GoString_Ugly(t *core.T) {
	fg := &Forge{}
	got := fmt.Sprintf("%#v", fg)
	core.AssertContains(t, got, "client=<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Builder_Set_Good(t *core.T) {
	b := newQueryBuilder()
	b.Set("owner", "core")
	core.AssertEqual(t, "owner=core", b.Encode())
}

func TestAX7_Builder_Set_Bad(t *core.T) {
	b := newQueryBuilder()
	b.Set("owner", "core")
	b.Set("owner", "forge")
	core.AssertEqual(t, "owner=forge", b.Encode())
}

func TestAX7_Builder_Set_Ugly(t *core.T) {
	b := newQueryBuilder()
	b.Set("owner name", "core/repo")
	core.AssertContains(t, b.Encode(), "owner+name=core%2Frepo")
}

func TestAX7_Builder_Add_Good(t *core.T) {
	b := newQueryBuilder()
	b.Add("label", "bug")
	core.AssertEqual(t, "label=bug", b.Encode())
}

func TestAX7_Builder_Add_Bad(t *core.T) {
	b := newQueryBuilder()
	b.Add("label", "bug")
	b.Add("label", "help")
	core.AssertEqual(t, "label=bug&label=help", b.Encode())
}

func TestAX7_Builder_Add_Ugly(t *core.T) {
	b := newQueryBuilder()
	b.Add("label name", "bug/help")
	core.AssertContains(t, b.Encode(), "label+name=bug%2Fhelp")
}

func TestAX7_Builder_Encode_Good(t *core.T) {
	b := newQueryBuilder()
	b.Set("repo", "go-forge")
	core.AssertEqual(t, "repo=go-forge", b.Encode())
}

func TestAX7_Builder_Encode_Bad(t *core.T) {
	b := newQueryBuilder()
	got := b.Encode()
	core.AssertEqual(t, "", got)
}

func TestAX7_Builder_Encode_Ugly(t *core.T) {
	var b *queryBuilder
	got := b.Encode()
	core.AssertEqual(t, "", got)
}

func TestAX7_IssueListOptions_String_Good(t *core.T) {
	value := IssueListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IssueListOptions_String_Bad(t *core.T) {
	value := IssueListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_IssueListOptions_String_Ugly(t *core.T) {
	value := IssueListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_IssueListOptions_GoString_Good(t *core.T) {
	value := IssueListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_IssueListOptions_GoString_Bad(t *core.T) {
	value := IssueListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_IssueListOptions_GoString_Ugly(t *core.T) {
	value := IssueListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AttachmentUploadOptions_String_Good(t *core.T) {
	value := AttachmentUploadOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AttachmentUploadOptions_String_Bad(t *core.T) {
	value := AttachmentUploadOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AttachmentUploadOptions_String_Ugly(t *core.T) {
	value := AttachmentUploadOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_AttachmentUploadOptions_GoString_Good(t *core.T) {
	value := AttachmentUploadOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_AttachmentUploadOptions_GoString_Bad(t *core.T) {
	value := AttachmentUploadOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_AttachmentUploadOptions_GoString_Ugly(t *core.T) {
	value := AttachmentUploadOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoCommentListOptions_String_Good(t *core.T) {
	value := RepoCommentListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoCommentListOptions_String_Bad(t *core.T) {
	value := RepoCommentListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoCommentListOptions_String_Ugly(t *core.T) {
	value := RepoCommentListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoCommentListOptions_GoString_Good(t *core.T) {
	value := RepoCommentListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoCommentListOptions_GoString_Bad(t *core.T) {
	value := RepoCommentListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoCommentListOptions_GoString_Ugly(t *core.T) {
	value := RepoCommentListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_IssueService_GetIssue_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.GetIssue(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetIssue_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.GetIssue(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetIssue_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.GetIssue(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditIssue_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditIssue(ctx, "core", "go-forge", 1, &types.EditIssueOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditIssue_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditIssue(ctx, "core", "go-forge", 1, &types.EditIssueOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditIssue_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditIssue(ctx, "core", "go-forge", 1, &types.EditIssueOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteIssue_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteIssue(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteIssue_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteIssue(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteIssue_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteIssue(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_SearchIssuesOptions_String_Good(t *core.T) {
	value := SearchIssuesOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_SearchIssuesOptions_String_Bad(t *core.T) {
	value := SearchIssuesOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_SearchIssuesOptions_String_Ugly(t *core.T) {
	value := SearchIssuesOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_SearchIssuesOptions_GoString_Good(t *core.T) {
	value := SearchIssuesOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_SearchIssuesOptions_GoString_Bad(t *core.T) {
	value := SearchIssuesOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_SearchIssuesOptions_GoString_Ugly(t *core.T) {
	value := SearchIssuesOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_IssueService_SearchIssuesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.SearchIssuesPage(ctx, SearchIssuesOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SearchIssuesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.SearchIssuesPage(ctx, SearchIssuesOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SearchIssuesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.SearchIssuesPage(ctx, SearchIssuesOptions{}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_SearchIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.SearchIssues(ctx, SearchIssuesOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SearchIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.SearchIssues(ctx, SearchIssuesOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SearchIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.SearchIssues(ctx, SearchIssuesOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterSearchIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterSearchIssues(ctx, SearchIssuesOptions{}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterSearchIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterSearchIssues(ctx, SearchIssuesOptions{}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterSearchIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterSearchIssues(ctx, SearchIssuesOptions{}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListIssuesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssuesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssuesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListIssues(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterIssues(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListRepoIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListRepoIssues(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListRepoIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListRepoIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListRepoIssuesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListRepoIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoIssuesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListRepoIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoIssuesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListRepoIssuesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterRepoIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterRepoIssues(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterRepoIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterRepoIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterRepoIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterRepoIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_CreateIssue_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.CreateIssue(ctx, "core", "go-forge", &types.CreateIssueOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateIssue_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.CreateIssue(ctx, "core", "go-forge", &types.CreateIssueOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateIssue_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.CreateIssue(ctx, "core", "go-forge", &types.CreateIssueOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_Pin_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.Pin(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_Pin_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.Pin(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_Pin_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.Pin(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_MovePin_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.MovePin(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_MovePin_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.MovePin(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_MovePin_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.MovePin(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListPinnedIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListPinnedIssues(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListPinnedIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListPinnedIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListPinnedIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListPinnedIssues(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterPinnedIssues_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterPinnedIssues(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterPinnedIssues_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterPinnedIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterPinnedIssues_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterPinnedIssues(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_Unpin_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.Unpin(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_Unpin_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.Unpin(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_Unpin_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.Unpin(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_SetDeadline_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.SetDeadline(ctx, "core", "go-forge", 1, "2026-04-28T00:00:00Z")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SetDeadline_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.SetDeadline(ctx, "core", "go-forge", 1, "2026-04-28T00:00:00Z")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SetDeadline_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.SetDeadline(ctx, "core", "go-forge", 1, "2026-04-28T00:00:00Z")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddReaction_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.AddReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddReaction_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.AddReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddReaction_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.AddReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListReactions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListReactions(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListReactions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListReactions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListReactions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListReactions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterReactions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterReactions(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterReactions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterReactions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterReactions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterReactions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteReaction_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteReaction_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteReaction_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_StartStopwatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.StartStopwatch(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_StartStopwatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.StartStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_StartStopwatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.StartStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_StopStopwatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.StopStopwatch(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_StopStopwatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.StopStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_StopStopwatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.StopStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteStopwatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteStopwatch(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteStopwatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteStopwatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteStopwatch(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListTimes(ctx, "core", "go-forge", 1, "value", nil, nil)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListTimes(ctx, "core", "go-forge", 1, "value", nil, nil)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListTimes(ctx, "core", "go-forge", 1, "value", nil, nil)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterTimes(ctx, "core", "go-forge", 1, "value", nil, nil) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterTimes(ctx, "core", "go-forge", 1, "value", nil, nil) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterTimes(ctx, "core", "go-forge", 1, "value", nil, nil) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddTime_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.AddTime(ctx, "core", "go-forge", 1, &types.AddTimeOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddTime_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.AddTime(ctx, "core", "go-forge", 1, &types.AddTimeOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddTime_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.AddTime(ctx, "core", "go-forge", 1, &types.AddTimeOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ResetTime_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.ResetTime(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ResetTime_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.ResetTime(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ResetTime_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.ResetTime(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteTime_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteTime(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteTime_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteTime(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteTime_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteTime(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.AddLabels(ctx, "core", "go-forge", 1, []int64{1})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.AddLabels(ctx, "core", "go-forge", 1, []int64{1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.AddLabels(ctx, "core", "go-forge", 1, []int64{1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_RemoveLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.RemoveLabel(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.RemoveLabel(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.RemoveLabel(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListComments(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListComments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListComments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListIssueComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListIssueComments(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssueComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListIssueComments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListIssueComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListIssueComments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterComments(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterComments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterComments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterIssueComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterIssueComments(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterIssueComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterIssueComments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterIssueComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterIssueComments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GetIssueComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.GetIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetIssueComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.GetIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetIssueComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.GetIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditIssueComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditIssueComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditIssueComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditIssueComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditIssueComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditIssueComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteIssueComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteIssueComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteIssueComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteIssueComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_CreateComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.CreateComment(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.CreateComment(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.CreateComment(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditComment(ctx, "core", "go-forge", 1, 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteComment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteComment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListRepoComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListRepoComments(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListRepoComments(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListRepoComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListRepoComments(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterRepoComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterRepoComments(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterRepoComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterRepoComments(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterRepoComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterRepoComments(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GetRepoComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.GetRepoComment(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetRepoComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.GetRepoComment(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetRepoComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.GetRepoComment(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditRepoComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditRepoComment(ctx, "core", "go-forge", 1, &types.EditIssueCommentOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditRepoComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditRepoComment(ctx, "core", "go-forge", 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditRepoComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditRepoComment(ctx, "core", "go-forge", 1, &types.EditIssueCommentOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteRepoComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteRepoComment(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteRepoComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteRepoComment(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteRepoComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteRepoComment(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListCommentReactions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListCommentReactions(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListCommentReactions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListCommentReactions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListCommentReactions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListCommentReactions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterCommentReactions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterCommentReactions(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterCommentReactions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterCommentReactions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterCommentReactions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterCommentReactions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddCommentReaction_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.AddCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddCommentReaction_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.AddCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddCommentReaction_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.AddCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteCommentReaction_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteCommentReaction_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteCommentReaction_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteCommentReaction(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_CreateAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.CreateAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.CreateAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.CreateAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListAttachments(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListAttachments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListAttachments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterAttachments(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterAttachments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterAttachments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GetAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.GetAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.GetAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.GetAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListCommentAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListCommentAttachments(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListCommentAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListCommentAttachments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListCommentAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListCommentAttachments(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterCommentAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterCommentAttachments(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterCommentAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterCommentAttachments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterCommentAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterCommentAttachments(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GetCommentAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.GetCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetCommentAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.GetCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_GetCommentAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.GetCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_CreateCommentAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.CreateCommentAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateCommentAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.CreateCommentAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CreateCommentAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.CreateCommentAttachment(ctx, "core", "go-forge", 1, &AttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_EditCommentAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.EditCommentAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditCommentAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.EditCommentAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_EditCommentAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.EditCommentAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_DeleteCommentAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.DeleteCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteCommentAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.DeleteCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_DeleteCommentAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.DeleteCommentAttachment(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListTimeline_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListTimeline(ctx, "core", "go-forge", 1, nil, nil)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListTimeline_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListTimeline(ctx, "core", "go-forge", 1, nil, nil)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListTimeline_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListTimeline(ctx, "core", "go-forge", 1, nil, nil)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterTimeline_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterTimeline(ctx, "core", "go-forge", 1, nil, nil) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterTimeline_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterTimeline(ctx, "core", "go-forge", 1, nil, nil) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterTimeline_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterTimeline(ctx, "core", "go-forge", 1, nil, nil) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListSubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListSubscriptions(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListSubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListSubscriptions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListSubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListSubscriptions(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterSubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterSubscriptions(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterSubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterSubscriptions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterSubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterSubscriptions(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_CheckSubscription_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.CheckSubscription(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CheckSubscription_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.CheckSubscription(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_CheckSubscription_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.CheckSubscription(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_SubscribeUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.SubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SubscribeUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.SubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_SubscribeUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.SubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_UnsubscribeUser_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.UnsubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_UnsubscribeUser_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.UnsubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_UnsubscribeUser_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.UnsubscribeUser(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListDependencies_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListDependencies(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListDependencies_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListDependencies(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListDependencies_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListDependencies(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterDependencies_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterDependencies(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterDependencies_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterDependencies(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterDependencies_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterDependencies(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddDependency_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.AddDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddDependency_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.AddDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddDependency_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.AddDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_RemoveDependency_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.RemoveDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveDependency_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.RemoveDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveDependency_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.RemoveDependency(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_ListBlocks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Issues.ListBlocks(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListBlocks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Issues.ListBlocks(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_ListBlocks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Issues.ListBlocks(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_IterBlocks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Issues.IterBlocks(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterBlocks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Issues.IterBlocks(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_IterBlocks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Issues.IterBlocks(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_AddBlock_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.AddBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddBlock_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.AddBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_AddBlock_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.AddBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_RemoveBlock_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Issues.RemoveBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveBlock_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Issues.RemoveBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_IssueService_RemoveBlock_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Issues.RemoveBlock(ctx, "core", "go-forge", 1, types.IssueMeta{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListRepoLabelsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListRepoLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListRepoLabelsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListRepoLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListRepoLabelsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListRepoLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListRepoLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListRepoLabels(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListRepoLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListRepoLabels(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListRepoLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListRepoLabels(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListLabelsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabelsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabelsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListLabelsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListLabels(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListLabels(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListLabels(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_IterRepoLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Labels.IterRepoLabels(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterRepoLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Labels.IterRepoLabels(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterRepoLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Labels.IterRepoLabels(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_IterLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Labels.IterLabels(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Labels.IterLabels(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Labels.IterLabels(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_GetRepoLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.GetRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetRepoLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.GetRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetRepoLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.GetRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_CreateRepoLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.CreateRepoLabel(ctx, "core", "go-forge", &types.CreateLabelOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_CreateRepoLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.CreateRepoLabel(ctx, "core", "go-forge", &types.CreateLabelOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_CreateRepoLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.CreateRepoLabel(ctx, "core", "go-forge", &types.CreateLabelOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_EditRepoLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.EditRepoLabel(ctx, "core", "go-forge", 1, &types.EditLabelOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_EditRepoLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.EditRepoLabel(ctx, "core", "go-forge", 1, &types.EditLabelOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_EditRepoLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.EditRepoLabel(ctx, "core", "go-forge", 1, &types.EditLabelOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_DeleteRepoLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Labels.DeleteRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_DeleteRepoLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Labels.DeleteRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_DeleteRepoLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Labels.DeleteRepoLabel(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListOrgLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListOrgLabels(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListOrgLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListOrgLabels(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListOrgLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListOrgLabels(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_IterOrgLabels_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Labels.IterOrgLabels(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterOrgLabels_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Labels.IterOrgLabels(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterOrgLabels_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Labels.IterOrgLabels(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_CreateOrgLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.CreateOrgLabel(ctx, "core", &types.CreateLabelOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_CreateOrgLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.CreateOrgLabel(ctx, "core", &types.CreateLabelOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_CreateOrgLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.CreateOrgLabel(ctx, "core", &types.CreateLabelOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_GetOrgLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.GetOrgLabel(ctx, "core", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetOrgLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.GetOrgLabel(ctx, "core", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetOrgLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.GetOrgLabel(ctx, "core", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_EditOrgLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.EditOrgLabel(ctx, "core", 1, &types.EditLabelOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_EditOrgLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.EditOrgLabel(ctx, "core", 1, &types.EditLabelOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_EditOrgLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.EditOrgLabel(ctx, "core", 1, &types.EditLabelOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_DeleteOrgLabel_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Labels.DeleteOrgLabel(ctx, "core", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_DeleteOrgLabel_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Labels.DeleteOrgLabel(ctx, "core", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_DeleteOrgLabel_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Labels.DeleteOrgLabel(ctx, "core", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_ListLabelTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.ListLabelTemplates(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabelTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.ListLabelTemplates(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_ListLabelTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.ListLabelTemplates(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_IterLabelTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Labels.IterLabelTemplates(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterLabelTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Labels.IterLabelTemplates(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_IterLabelTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Labels.IterLabelTemplates(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_GetLabelTemplate_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Labels.GetLabelTemplate(ctx, "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetLabelTemplate_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Labels.GetLabelTemplate(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_LabelService_GetLabelTemplate_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Labels.GetLabelTemplate(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneListOptions_String_Good(t *core.T) {
	value := MilestoneListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_MilestoneListOptions_String_Bad(t *core.T) {
	value := MilestoneListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_MilestoneListOptions_String_Ugly(t *core.T) {
	value := MilestoneListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_MilestoneListOptions_GoString_Good(t *core.T) {
	value := MilestoneListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_MilestoneListOptions_GoString_Bad(t *core.T) {
	value := MilestoneListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_MilestoneListOptions_GoString_Ugly(t *core.T) {
	value := MilestoneListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_MilestoneService_ListMilestonesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.ListMilestonesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListMilestonesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.ListMilestonesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListMilestonesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.ListMilestonesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_ListMilestones_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.ListMilestones(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListMilestones_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.ListMilestones(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListMilestones_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.ListMilestones(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_IterMilestones_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Milestones.IterMilestones(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_IterMilestones_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Milestones.IterMilestones(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_IterMilestones_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Milestones.IterMilestones(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_List_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_List_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_List_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_Iter_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Milestones.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Iter_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Milestones.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Iter_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Milestones.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_ListAll_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListAll_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_ListAll_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_Get_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.Get(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Get_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.Get(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Get_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.Get(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_GetMilestone_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.GetMilestone(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_GetMilestone_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.GetMilestone(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_GetMilestone_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.GetMilestone(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_Create_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.Create(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Create_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.Create(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Create_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.Create(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_CreateMilestone_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.CreateMilestone(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_CreateMilestone_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.CreateMilestone(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_CreateMilestone_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.CreateMilestone(ctx, "core", "go-forge", &types.CreateMilestoneOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_Edit_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.Edit(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Edit_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.Edit(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Edit_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.Edit(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_EditMilestone_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Milestones.EditMilestone(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_EditMilestone_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Milestones.EditMilestone(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_EditMilestone_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Milestones.EditMilestone(ctx, "core", "go-forge", 1, &types.EditMilestoneOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_Delete_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Milestones.Delete(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Delete_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Milestones.Delete(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_Delete_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Milestones.Delete(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_DeleteMilestone_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Milestones.DeleteMilestone(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_DeleteMilestone_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Milestones.DeleteMilestone(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MilestoneService_DeleteMilestone_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Milestones.DeleteMilestone(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_RenderMarkdown_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.RenderMarkdown(ctx, "hello", "gfm")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkdown_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.RenderMarkdown(ctx, "hello", "gfm")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkdown_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.RenderMarkdown(ctx, "hello", "gfm")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_RenderMarkup_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.RenderMarkup(ctx, "hello", "gfm")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkup_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.RenderMarkup(ctx, "hello", "gfm")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkup_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.RenderMarkup(ctx, "hello", "gfm")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_RenderMarkdownRaw_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.RenderMarkdownRaw(ctx, "hello")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkdownRaw_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.RenderMarkdownRaw(ctx, "hello")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_RenderMarkdownRaw_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.RenderMarkdownRaw(ctx, "hello")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_ListLicenses_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.ListLicenses(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_ListLicenses_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.ListLicenses(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_ListLicenses_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.ListLicenses(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_IterLicenses_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Misc.IterLicenses(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_IterLicenses_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Misc.IterLicenses(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_IterLicenses_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Misc.IterLicenses(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetLicense_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetLicense(ctx, "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetLicense_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetLicense(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetLicense_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetLicense(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_ListGitignoreTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.ListGitignoreTemplates(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_ListGitignoreTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.ListGitignoreTemplates(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_ListGitignoreTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.ListGitignoreTemplates(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_IterGitignoreTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Misc.IterGitignoreTemplates(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_IterGitignoreTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Misc.IterGitignoreTemplates(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_IterGitignoreTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Misc.IterGitignoreTemplates(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetGitignoreTemplate_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetGitignoreTemplate(ctx, "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetGitignoreTemplate_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetGitignoreTemplate(ctx, "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetGitignoreTemplate_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetGitignoreTemplate(ctx, "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetNodeInfo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetNodeInfo(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetNodeInfo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetNodeInfo(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetNodeInfo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetNodeInfo(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetSigningKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetSigningKey(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetSigningKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetSigningKey(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetSigningKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetSigningKey(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetAPISettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetAPISettings(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetAPISettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetAPISettings(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetAPISettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetAPISettings(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetAttachmentSettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetAttachmentSettings(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetAttachmentSettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetAttachmentSettings(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetAttachmentSettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetAttachmentSettings(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetRepositorySettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetRepositorySettings(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetRepositorySettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetRepositorySettings(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetRepositorySettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetRepositorySettings(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetUISettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetUISettings(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetUISettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetUISettings(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetUISettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetUISettings(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GetVersion_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Misc.GetVersion(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetVersion_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Misc.GetVersion(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_MiscService_GetVersion_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Misc.GetVersion(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationListOptions_String_Good(t *core.T) {
	value := NotificationListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationListOptions_String_Bad(t *core.T) {
	value := NotificationListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationListOptions_String_Ugly(t *core.T) {
	value := NotificationListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationListOptions_GoString_Good(t *core.T) {
	value := NotificationListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationListOptions_GoString_Bad(t *core.T) {
	value := NotificationListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationListOptions_GoString_Ugly(t *core.T) {
	value := NotificationListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationRepoMarkOptions_String_Good(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationRepoMarkOptions_String_Bad(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationRepoMarkOptions_String_Ugly(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationRepoMarkOptions_GoString_Good(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationRepoMarkOptions_GoString_Bad(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationRepoMarkOptions_GoString_Ugly(t *core.T) {
	value := NotificationRepoMarkOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationMarkOptions_String_Good(t *core.T) {
	value := NotificationMarkOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationMarkOptions_String_Bad(t *core.T) {
	value := NotificationMarkOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationMarkOptions_String_Ugly(t *core.T) {
	value := NotificationMarkOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationMarkOptions_GoString_Good(t *core.T) {
	value := NotificationMarkOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NotificationMarkOptions_GoString_Bad(t *core.T) {
	value := NotificationMarkOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_NotificationMarkOptions_GoString_Ugly(t *core.T) {
	value := NotificationMarkOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_NotificationService_List_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.List(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_List_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.List(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_List_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.List(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_Iter_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Notifications.Iter(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_Iter_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Notifications.Iter(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_Iter_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Notifications.Iter(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_NewAvailable_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.NewAvailable(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_NewAvailable_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.NewAvailable(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_NewAvailable_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.NewAvailable(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_ListRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.ListRepo(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_ListRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.ListRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_ListRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.ListRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_IterRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Notifications.IterRepo(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_IterRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Notifications.IterRepo(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_IterRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Notifications.IterRepo(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_MarkNotifications_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.MarkNotifications(ctx, &NotificationMarkOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkNotifications_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.MarkNotifications(ctx, &NotificationMarkOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkNotifications_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.MarkNotifications(ctx, &NotificationMarkOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_MarkRepoNotifications_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.MarkRepoNotifications(ctx, "core", "go-forge", &NotificationRepoMarkOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkRepoNotifications_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.MarkRepoNotifications(ctx, "core", "go-forge", &NotificationRepoMarkOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkRepoNotifications_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.MarkRepoNotifications(ctx, "core", "go-forge", &NotificationRepoMarkOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_MarkRead_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Notifications.MarkRead(ctx)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkRead_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Notifications.MarkRead(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkRead_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Notifications.MarkRead(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_GetThread_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Notifications.GetThread(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_GetThread_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Notifications.GetThread(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_GetThread_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Notifications.GetThread(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_MarkThreadRead_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Notifications.MarkThreadRead(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkThreadRead_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Notifications.MarkThreadRead(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_NotificationService_MarkThreadRead_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Notifications.MarkThreadRead(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgActivityFeedListOptions_String_Good(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_OrgActivityFeedListOptions_String_Bad(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_OrgActivityFeedListOptions_String_Ugly(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_OrgActivityFeedListOptions_GoString_Good(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_OrgActivityFeedListOptions_GoString_Bad(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_OrgActivityFeedListOptions_GoString_Ugly(t *core.T) {
	value := OrgActivityFeedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_OrgService_GetOrg_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.GetOrg(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetOrg_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.GetOrg(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetOrg_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.GetOrg(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_UpdateOrg_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.UpdateOrg(ctx, "core", &types.EditOrgOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_UpdateOrg_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.UpdateOrg(ctx, "core", &types.EditOrgOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_UpdateOrg_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.UpdateOrg(ctx, "core", &types.EditOrgOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_DeleteOrg_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.DeleteOrg(ctx, "core")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_DeleteOrg_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.DeleteOrg(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_DeleteOrg_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.DeleteOrg(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgsPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgsPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgsPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgs(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgs(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgs(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgTeamsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgTeamsPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgTeamsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgTeamsPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgTeamsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgTeamsPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgTeams(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgTeams(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgTeams(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterOrgTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterOrgTeams(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgTeams(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgTeams(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterOrgs(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_CreateOrg_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.CreateOrg(ctx, &types.CreateOrgOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_CreateOrg_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.CreateOrg(ctx, &types.CreateOrgOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_CreateOrg_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.CreateOrg(ctx, &types.CreateOrgOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListMembersPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMembersPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMembersPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListMembers(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListMembers(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListMembers(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgMembersPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgMembersPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgMembersPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgMembersPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListOrgMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListOrgMembers(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListOrgMembers(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListOrgMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListOrgMembers(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterMembers(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterOrgMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterOrgMembers(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterOrgMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterOrgMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_AddMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.AddMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_AddMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.AddMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_AddMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.AddMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_RemoveMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.RemoveMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_RemoveMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.RemoveMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_RemoveMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.RemoveMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IsMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.IsMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.IsMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.IsMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListBlockedUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListBlockedUsers(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListBlockedUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListBlockedUsers(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListBlockedUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListBlockedUsers(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterBlockedUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterBlockedUsers(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterBlockedUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterBlockedUsers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterBlockedUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterBlockedUsers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IsBlocked_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.IsBlocked(ctx, "core", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsBlocked_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.IsBlocked(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsBlocked_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.IsBlocked(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListPublicMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListPublicMembers(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListPublicMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListPublicMembers(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListPublicMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListPublicMembers(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterPublicMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterPublicMembers(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterPublicMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterPublicMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterPublicMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterPublicMembers(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IsPublicMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.IsPublicMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsPublicMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.IsPublicMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IsPublicMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.IsPublicMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_PublicizeMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.PublicizeMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_PublicizeMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.PublicizeMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_PublicizeMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.PublicizeMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ConcealMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.ConcealMember(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ConcealMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.ConcealMember(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ConcealMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.ConcealMember(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_Block_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.Block(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_Block_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.Block(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_Block_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.Block(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_Unblock_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.Unblock(ctx, "core", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_Unblock_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.Unblock(ctx, "core", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_Unblock_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.Unblock(ctx, "core", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_GetQuota_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.GetQuota(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetQuota_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.GetQuota(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetQuota_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.GetQuota(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_CheckQuota_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.CheckQuota(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_CheckQuota_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.CheckQuota(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_CheckQuota_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.CheckQuota(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListQuotaArtifacts_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListQuotaArtifacts(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaArtifacts_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListQuotaArtifacts(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaArtifacts_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListQuotaArtifacts(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterQuotaArtifacts_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterQuotaArtifacts(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaArtifacts_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaArtifacts(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaArtifacts_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaArtifacts(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListQuotaAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListQuotaAttachments(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListQuotaAttachments(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListQuotaAttachments(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterQuotaAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterQuotaAttachments(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaAttachments(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaAttachments(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListQuotaPackages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListQuotaPackages(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaPackages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListQuotaPackages(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListQuotaPackages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListQuotaPackages(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterQuotaPackages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterQuotaPackages(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaPackages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaPackages(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterQuotaPackages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterQuotaPackages(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_GetRunnerRegistrationToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.GetRunnerRegistrationToken(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetRunnerRegistrationToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.GetRunnerRegistrationToken(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetRunnerRegistrationToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.GetRunnerRegistrationToken(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_UpdateAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.UpdateAvatar(ctx, "core", &types.UpdateUserAvatarOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_UpdateAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.UpdateAvatar(ctx, "core", &types.UpdateUserAvatarOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_UpdateAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.UpdateAvatar(ctx, "core", &types.UpdateUserAvatarOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_DeleteAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Orgs.DeleteAvatar(ctx, "core")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_DeleteAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Orgs.DeleteAvatar(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_DeleteAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Orgs.DeleteAvatar(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_SearchTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.SearchTeams(ctx, "core", "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_SearchTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.SearchTeams(ctx, "core", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_SearchTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.SearchTeams(ctx, "core", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterSearchTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterSearchTeams(ctx, "core", "value") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterSearchTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterSearchTeams(ctx, "core", "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterSearchTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterSearchTeams(ctx, "core", "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_GetUserPermissions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.GetUserPermissions(ctx, "alice", "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetUserPermissions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.GetUserPermissions(ctx, "alice", "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_GetUserPermissions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.GetUserPermissions(ctx, "alice", "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListActivityFeeds(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListActivityFeeds(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListActivityFeeds(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterActivityFeeds(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterActivityFeeds(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterActivityFeeds(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListUserOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListUserOrgs(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListUserOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListUserOrgs(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListUserOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListUserOrgs(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterUserOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterUserOrgs(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterUserOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterUserOrgs(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterUserOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterUserOrgs(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_ListMyOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Orgs.ListMyOrgs(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMyOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Orgs.ListMyOrgs(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_ListMyOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Orgs.ListMyOrgs(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_IterMyOrgs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Orgs.IterMyOrgs(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterMyOrgs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Orgs.IterMyOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_OrgService_IterMyOrgs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Orgs.IterMyOrgs(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_List_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Packages.List(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_List_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Packages.List(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_List_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Packages.List(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_Iter_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Packages.Iter(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Iter_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Packages.Iter(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Iter_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Packages.Iter(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_Get_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Packages.Get(ctx, "core", "value", "name", "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Get_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Packages.Get(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Get_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Packages.Get(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_Delete_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Packages.Delete(ctx, "core", "value", "name", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Delete_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Packages.Delete(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_Delete_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Packages.Delete(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_ListFiles_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Packages.ListFiles(ctx, "core", "value", "name", "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_ListFiles_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Packages.ListFiles(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_ListFiles_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Packages.ListFiles(ctx, "core", "value", "name", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_IterFiles_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Packages.IterFiles(ctx, "core", "value", "name", "value") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_IterFiles_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Packages.IterFiles(ctx, "core", "value", "name", "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PackageService_IterFiles_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Packages.IterFiles(ctx, "core", "value", "name", "value") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ListOptions_String_Good(t *core.T) {
	value := ListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ListOptions_String_Bad(t *core.T) {
	value := ListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ListOptions_String_Ugly(t *core.T) {
	value := ListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ListOptions_GoString_Good(t *core.T) {
	value := ListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ListOptions_GoString_Bad(t *core.T) {
	value := ListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ListOptions_GoString_Ugly(t *core.T) {
	value := ListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_PagedResult_String_Good(t *core.T) {
	value := PagedResult[ax7Payload]{Items: []ax7Payload{{Name: "one"}}, TotalCount: 1, Page: 1}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_PagedResult_String_Bad(t *core.T) {
	value := PagedResult[ax7Payload]{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_PagedResult_String_Ugly(t *core.T) {
	value := PagedResult[ax7Payload]{Items: []ax7Payload{{Name: "one"}}, TotalCount: 1, Page: 1}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_PagedResult_GoString_Good(t *core.T) {
	value := PagedResult[ax7Payload]{Items: []ax7Payload{{Name: "one"}}, TotalCount: 1, Page: 1}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_PagedResult_GoString_Bad(t *core.T) {
	value := PagedResult[ax7Payload]{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_PagedResult_GoString_Ugly(t *core.T) {
	value := PagedResult[ax7Payload]{Items: []ax7Payload{{Name: "one"}}, TotalCount: 1, Page: 1}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ListPage_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	got, err := ListPage[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ListPage_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := ListPage[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ListPage_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := ListPage[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ListAll_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	got, err := ListAll[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ListAll_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := ListAll[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ListAll_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := ListAll[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ListIter_Good(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := context.Background()
	for _, err := range ListIter[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ListIter_Bad(t *core.T) {
	c, tr := ax7Client(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range ListIter[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ListIter_Ugly(t *core.T) {
	c, tr := ax7Client(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range ListIter[ax7Payload](ctx, c, "/api/v1/items", map[string]string{"q": "go"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Params_String_Good(t *core.T) {
	value := Params{"owner": "core"}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Params_String_Bad(t *core.T) {
	value := Params{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_Params_String_Ugly(t *core.T) {
	value := Params{"owner": "core"}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_Params_GoString_Good(t *core.T) {
	value := Params{"owner": "core"}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Params_GoString_Bad(t *core.T) {
	value := Params{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_Params_GoString_Ugly(t *core.T) {
	value := Params{"owner": "core"}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ResolvePath_Good(t *core.T) {
	got := ResolvePath("/api/v1/repos/{owner}/{repo}", Params{"owner": "core", "repo": "go forge"})
	core.AssertEqual(t, "/api/v1/repos/core/go%20forge", got)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ResolvePath_Bad(t *core.T) {
	got := ResolvePath("/api/v1/repos/{owner}/{repo}", Params{"owner": "core"})
	core.AssertContains(t, got, "{repo}")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ResolvePath_Ugly(t *core.T) {
	got := ResolvePath("/api/v1/items/{id}", Params{"id": "a/b?c"})
	core.AssertContains(t, got, "a%2Fb%3Fc")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_PullListOptions_String_Good(t *core.T) {
	value := PullListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_PullListOptions_String_Bad(t *core.T) {
	value := PullListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_PullListOptions_String_Ugly(t *core.T) {
	value := PullListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_PullListOptions_GoString_Good(t *core.T) {
	value := PullListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_PullListOptions_GoString_Bad(t *core.T) {
	value := PullListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_PullListOptions_GoString_Ugly(t *core.T) {
	value := PullListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_PullService_ListPullRequestsPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListPullRequestsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullRequestsPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListPullRequestsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullRequestsPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListPullRequestsPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListPullRequests_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListPullRequests(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullRequests_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListPullRequests(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullRequests_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListPullRequests(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterPullRequests_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterPullRequests(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterPullRequests_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterPullRequests(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterPullRequests_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterPullRequests(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_CreatePullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.CreatePullRequest(ctx, "core", "go-forge", &types.CreatePullRequestOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CreatePullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.CreatePullRequest(ctx, "core", "go-forge", &types.CreatePullRequestOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CreatePullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.CreatePullRequest(ctx, "core", "go-forge", &types.CreatePullRequestOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetPullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetPullRequest(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetPullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetPullRequest(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetPullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetPullRequest(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_EditPullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.EditPullRequest(ctx, "core", "go-forge", 1, &types.EditPullRequestOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_EditPullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.EditPullRequest(ctx, "core", "go-forge", 1, &types.EditPullRequestOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_EditPullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.EditPullRequest(ctx, "core", "go-forge", 1, &types.EditPullRequestOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_DeletePullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.DeletePullRequest(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeletePullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.DeletePullRequest(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeletePullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.DeletePullRequest(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_Merge_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.Merge(ctx, "core", "go-forge", 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_Merge_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.Merge(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_Merge_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.Merge(ctx, "core", "go-forge", 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_MergePullRequest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.MergePullRequest(ctx, "core", "go-forge", 1, &types.MergePullRequestOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_MergePullRequest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.MergePullRequest(ctx, "core", "go-forge", 1, &types.MergePullRequestOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_MergePullRequest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.MergePullRequest(ctx, "core", "go-forge", 1, &types.MergePullRequestOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_CancelScheduledAutoMerge_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.CancelScheduledAutoMerge(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CancelScheduledAutoMerge_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.CancelScheduledAutoMerge(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CancelScheduledAutoMerge_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.CancelScheduledAutoMerge(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_Update_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.Update(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_Update_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.Update(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_Update_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.Update(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetDiffOrPatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetDiffOrPatch(ctx, "core", "go-forge", 1, "diff")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetDiffOrPatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetDiffOrPatch(ctx, "core", "go-forge", 1, "diff")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetDiffOrPatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetDiffOrPatch(ctx, "core", "go-forge", 1, "diff")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListCommits_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListCommits(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListCommits_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListCommits(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListCommits_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListCommits(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterCommits_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterCommits(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterCommits_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterCommits(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterCommits_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterCommits(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListReviews_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListReviews(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviews_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListReviews(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviews_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListReviews(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListPullReviews_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListPullReviews(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullReviews_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListPullReviews(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListPullReviews_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListPullReviews(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterReviews_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterReviews(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviews_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterReviews(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviews_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterReviews(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterPullReviews_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterPullReviews(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterPullReviews_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterPullReviews(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterPullReviews_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterPullReviews(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListFiles_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListFiles(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListFiles_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListFiles(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListFiles_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListFiles(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterFiles_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterFiles(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterFiles_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterFiles(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterFiles_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterFiles(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetByBaseHead_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetByBaseHead(ctx, "core", "go-forge", "value", "value")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetByBaseHead_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetByBaseHead(ctx, "core", "go-forge", "value", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetByBaseHead_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetByBaseHead(ctx, "core", "go-forge", "value", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListReviewers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListReviewers(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviewers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListReviewers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviewers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListReviewers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterReviewers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterReviewers(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviewers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterReviewers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviewers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterReviewers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_RequestReviewers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.RequestReviewers(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_RequestReviewers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.RequestReviewers(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_RequestReviewers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.RequestReviewers(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_CancelReviewRequests_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.CancelReviewRequests(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CancelReviewRequests_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.CancelReviewRequests(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CancelReviewRequests_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.CancelReviewRequests(ctx, "core", "go-forge", 1, &types.PullReviewRequestOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_SubmitReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.SubmitReview(ctx, "core", "go-forge", 1, &types.SubmitPullReviewOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_SubmitReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.SubmitReview(ctx, "core", "go-forge", 1, &types.SubmitPullReviewOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_SubmitReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.SubmitReview(ctx, "core", "go-forge", 1, &types.SubmitPullReviewOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetReview(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetPullReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetPullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetPullReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetPullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetPullReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetPullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_DeleteReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.DeleteReview(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeleteReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.DeleteReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeleteReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.DeleteReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_DeletePullReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.DeletePullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeletePullReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.DeletePullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeletePullReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.DeletePullReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_ListReviewComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.ListReviewComments(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviewComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.ListReviewComments(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_ListReviewComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.ListReviewComments(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_IterReviewComments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Pulls.IterReviewComments(ctx, "core", "go-forge", 1, 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviewComments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Pulls.IterReviewComments(ctx, "core", "go-forge", 1, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_IterReviewComments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Pulls.IterReviewComments(ctx, "core", "go-forge", 1, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GetReviewComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.GetReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetReviewComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.GetReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_GetReviewComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.GetReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_CreateReviewComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Pulls.CreateReviewComment(ctx, "core", "go-forge", 1, 1, &types.CreatePullReviewCommentOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CreateReviewComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Pulls.CreateReviewComment(ctx, "core", "go-forge", 1, 1, &types.CreatePullReviewCommentOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_CreateReviewComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Pulls.CreateReviewComment(ctx, "core", "go-forge", 1, 1, &types.CreatePullReviewCommentOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_DeleteReviewComment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.DeleteReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeleteReviewComment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.DeleteReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DeleteReviewComment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.DeleteReviewComment(ctx, "core", "go-forge", 1, 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_DismissReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.DismissReview(ctx, "core", "go-forge", 1, 1, "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DismissReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.DismissReview(ctx, "core", "go-forge", 1, 1, "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_DismissReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.DismissReview(ctx, "core", "go-forge", 1, 1, "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_UndismissReview_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Pulls.UndismissReview(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_UndismissReview_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Pulls.UndismissReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_PullService_UndismissReview_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Pulls.UndismissReview(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseListOptions_String_Good(t *core.T) {
	value := ReleaseListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ReleaseListOptions_String_Bad(t *core.T) {
	value := ReleaseListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ReleaseListOptions_String_Ugly(t *core.T) {
	value := ReleaseListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ReleaseListOptions_GoString_Good(t *core.T) {
	value := ReleaseListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ReleaseListOptions_GoString_Bad(t *core.T) {
	value := ReleaseListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ReleaseListOptions_GoString_Ugly(t *core.T) {
	value := ReleaseListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ReleaseAttachmentUploadOptions_String_Good(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ReleaseAttachmentUploadOptions_String_Bad(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ReleaseAttachmentUploadOptions_String_Ugly(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ReleaseAttachmentUploadOptions_GoString_Good(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ReleaseAttachmentUploadOptions_GoString_Bad(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ReleaseAttachmentUploadOptions_GoString_Ugly(t *core.T) {
	value := ReleaseAttachmentUploadOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ReleaseService_ListReleasesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.ListReleasesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListReleasesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.ListReleasesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListReleasesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.ListReleasesPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_ListReleases_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.ListReleases(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListReleases_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.ListReleases(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListReleases_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.ListReleases(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_IterReleases_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Releases.IterReleases(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_IterReleases_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Releases.IterReleases(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_IterReleases_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Releases.IterReleases(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_CreateRelease_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.CreateRelease(ctx, "core", "go-forge", &types.CreateReleaseOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateRelease_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.CreateRelease(ctx, "core", "go-forge", &types.CreateReleaseOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateRelease_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.CreateRelease(ctx, "core", "go-forge", &types.CreateReleaseOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GetByTag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.GetByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetByTag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.GetByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetByTag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.GetByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GetRelease_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.GetRelease(ctx, "core", "go-forge", "v1.0.0")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetRelease_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.GetRelease(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetRelease_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.GetRelease(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GetLatest_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.GetLatest(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetLatest_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.GetLatest(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetLatest_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.GetLatest(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_DeleteByTag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Releases.DeleteByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_DeleteByTag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Releases.DeleteByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_DeleteByTag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Releases.DeleteByTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_ListAssets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.ListAssets(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListAssets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.ListAssets(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_ListAssets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.ListAssets(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_CreateAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.CreateAttachment(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.CreateAttachment(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.CreateAttachment(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_EditAttachment_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_EditAttachment_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_EditAttachment_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.EditAttachment(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_CreateAsset_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.CreateAsset(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateAsset_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.CreateAsset(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_CreateAsset_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.CreateAsset(ctx, "core", "go-forge", 1, &ReleaseAttachmentUploadOptions{}, "asset.txt", strings.NewReader("payload"))
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_EditAsset_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.EditAsset(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_EditAsset_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.EditAsset(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_EditAsset_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.EditAsset(ctx, "core", "go-forge", 1, 1, &types.EditAttachmentOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_IterAssets_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Releases.IterAssets(ctx, "core", "go-forge", 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_IterAssets_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Releases.IterAssets(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_IterAssets_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Releases.IterAssets(ctx, "core", "go-forge", 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GetAsset_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Releases.GetAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetAsset_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Releases.GetAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_GetAsset_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Releases.GetAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_DeleteAsset_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Releases.DeleteAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_DeleteAsset_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Releases.DeleteAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_ReleaseService_DeleteAsset_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Releases.DeleteAsset(ctx, "core", "go-forge", 1, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoKeyListOptions_String_Good(t *core.T) {
	value := RepoKeyListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoKeyListOptions_String_Bad(t *core.T) {
	value := RepoKeyListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoKeyListOptions_String_Ugly(t *core.T) {
	value := RepoKeyListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoKeyListOptions_GoString_Good(t *core.T) {
	value := RepoKeyListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoKeyListOptions_GoString_Bad(t *core.T) {
	value := RepoKeyListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoKeyListOptions_GoString_Ugly(t *core.T) {
	value := RepoKeyListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ActivityFeedListOptions_String_Good(t *core.T) {
	value := ActivityFeedListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ActivityFeedListOptions_String_Bad(t *core.T) {
	value := ActivityFeedListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ActivityFeedListOptions_String_Ugly(t *core.T) {
	value := ActivityFeedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_ActivityFeedListOptions_GoString_Good(t *core.T) {
	value := ActivityFeedListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_ActivityFeedListOptions_GoString_Bad(t *core.T) {
	value := ActivityFeedListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_ActivityFeedListOptions_GoString_Ugly(t *core.T) {
	value := ActivityFeedListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoTimeListOptions_String_Good(t *core.T) {
	value := RepoTimeListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoTimeListOptions_String_Bad(t *core.T) {
	value := RepoTimeListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoTimeListOptions_String_Ugly(t *core.T) {
	value := RepoTimeListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoTimeListOptions_GoString_Good(t *core.T) {
	value := RepoTimeListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_RepoTimeListOptions_GoString_Bad(t *core.T) {
	value := RepoTimeListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_RepoTimeListOptions_GoString_Ugly(t *core.T) {
	value := RepoTimeListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_RepoService_GetRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRepo(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_UpdateRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.UpdateRepo(ctx, "core", "go-forge", &types.EditRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.UpdateRepo(ctx, "core", "go-forge", &types.EditRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.UpdateRepo(ctx, "core", "go-forge", &types.EditRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteRepo(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteRepo(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Migrate_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Migrate(ctx, &types.MigrateRepoOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Migrate_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Migrate(ctx, &types.MigrateRepoOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Migrate_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Migrate(ctx, &types.MigrateRepoOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreateCurrentUserRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreateCurrentUserRepo(ctx, &types.CreateRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateCurrentUserRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreateCurrentUserRepo(ctx, &types.CreateRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateCurrentUserRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreateCurrentUserRepo(ctx, &types.CreateRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreateOrgRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreateOrgRepo(ctx, "core", &types.CreateRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateOrgRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreateOrgRepo(ctx, "core", &types.CreateRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateOrgRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreateOrgRepo(ctx, "core", &types.CreateRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreateOrgRepoDeprecated_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreateOrgRepoDeprecated(ctx, "core", &types.CreateRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateOrgRepoDeprecated_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreateOrgRepoDeprecated(ctx, "core", &types.CreateRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateOrgRepoDeprecated_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreateOrgRepoDeprecated(ctx, "core", &types.CreateRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListOrgReposPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListOrgReposPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListOrgReposPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListOrgReposPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListOrgReposPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListOrgReposPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListOrgRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListOrgRepos(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListOrgRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListOrgRepos(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListOrgRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListOrgRepos(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterOrgRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterOrgRepos(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterOrgRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterOrgRepos(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterOrgRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterOrgRepos(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListCurrentUserReposPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListCurrentUserReposPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListCurrentUserReposPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListCurrentUserReposPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListCurrentUserReposPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListCurrentUserReposPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListUserReposPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListUserReposPage(ctx, "alice", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserReposPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListUserReposPage(ctx, "alice", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserReposPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListUserReposPage(ctx, "alice", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListUserRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListUserRepos(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListUserRepos(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListUserRepos(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterUserRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterUserRepos(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterUserRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterUserRepos(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterUserRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterUserRepos(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetByID_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetByID(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetByID_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetByID(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetByID_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetByID(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListTags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListTags(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListTags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListTags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterTags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterTags(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterTags(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterTags(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetTag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteTag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteTag(ctx, "core", "go-forge", "v1.0.0")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListTagProtections_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListTagProtections(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTagProtections_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListTagProtections(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTagProtections_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListTagProtections(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterTagProtections_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterTagProtections(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTagProtections_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterTagProtections(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTagProtections_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterTagProtections(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetTagProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetTagProtection(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTagProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetTagProtection(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTagProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetTagProtection(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreateTagProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreateTagProtection(ctx, "core", "go-forge", &types.CreateTagProtectionOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateTagProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreateTagProtection(ctx, "core", "go-forge", &types.CreateTagProtectionOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateTagProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreateTagProtection(ctx, "core", "go-forge", &types.CreateTagProtectionOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_EditTagProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.EditTagProtection(ctx, "core", "go-forge", 1, &types.EditTagProtectionOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_EditTagProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.EditTagProtection(ctx, "core", "go-forge", 1, &types.EditTagProtectionOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_EditTagProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.EditTagProtection(ctx, "core", "go-forge", 1, &types.EditTagProtectionOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteTagProtection_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteTagProtection(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTagProtection_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteTagProtection(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTagProtection_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteTagProtection(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListKeys(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListKeys(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListKeys(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterKeys(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterKeys(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterKeys(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetKey(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetKey(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetKey(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreateKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreateKey(ctx, "core", "go-forge", &types.CreateKeyOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreateKey(ctx, "core", "go-forge", &types.CreateKeyOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreateKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreateKey(ctx, "core", "go-forge", &types.CreateKeyOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteKey(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteKey(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteKey(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListStargazers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListStargazers(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListStargazers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListStargazers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListStargazers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListStargazers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterStargazers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterStargazers(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterStargazers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterStargazers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterStargazers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterStargazers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListSubscribers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListSubscribers(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListSubscribers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListSubscribers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListSubscribers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListSubscribers(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterSubscribers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterSubscribers(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSubscribers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterSubscribers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSubscribers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterSubscribers(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListAssignees_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListAssignees(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListAssignees_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListAssignees(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListAssignees_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListAssignees(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterAssignees_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterAssignees(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterAssignees_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterAssignees(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterAssignees_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterAssignees(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListCollaborators_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListCollaborators(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListCollaborators_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListCollaborators(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListCollaborators_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListCollaborators(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterCollaborators_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterCollaborators(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterCollaborators_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterCollaborators(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterCollaborators_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterCollaborators(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListRepoTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListRepoTeams(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListRepoTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListRepoTeams(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListRepoTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListRepoTeams(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterRepoTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterRepoTeams(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterRepoTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterRepoTeams(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterRepoTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterRepoTeams(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetRepoTeam_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepoTeam_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepoTeam_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_AddRepoTeam_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.AddRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddRepoTeam_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.AddRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddRepoTeam_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.AddRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteRepoTeam_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteRepoTeam_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteRepoTeam_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteRepoTeam(ctx, "core", "go-forge", "maintainers")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CheckCollaborator_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CheckCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CheckCollaborator_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CheckCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CheckCollaborator_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CheckCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_AddCollaborator_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.AddCollaborator(ctx, "core", "go-forge", "alice", &types.AddCollaboratorOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddCollaborator_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.AddCollaborator(ctx, "core", "go-forge", "alice", &types.AddCollaboratorOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddCollaborator_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.AddCollaborator(ctx, "core", "go-forge", "alice", &types.AddCollaboratorOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteCollaborator_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteCollaborator_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteCollaborator_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteCollaborator(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetCollaboratorPermission_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetCollaboratorPermission(ctx, "core", "go-forge", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetCollaboratorPermission_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetCollaboratorPermission(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetCollaboratorPermission_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetCollaboratorPermission(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetRepoPermissions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRepoPermissions(ctx, "core", "go-forge", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepoPermissions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRepoPermissions(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRepoPermissions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRepoPermissions(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetArchive_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetArchive(ctx, "core", "go-forge", "main.tar.gz")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetArchive_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetArchive(ctx, "core", "go-forge", "main.tar.gz")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetArchive_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetArchive(ctx, "core", "go-forge", "main.tar.gz")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Compare_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Compare(ctx, "core", "go-forge", "main...feature")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Compare_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Compare(ctx, "core", "go-forge", "main...feature")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Compare_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Compare(ctx, "core", "go-forge", "main...feature")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetRawFile_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRawFile_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRawFile_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRawFile(ctx, "core", "go-forge", "README.md")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetRawFileOrLFS_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRawFileOrLFS(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRawFileOrLFS_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRawFileOrLFS(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRawFileOrLFS_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRawFileOrLFS(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetEditorConfig_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.GetEditorConfig(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetEditorConfig_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.GetEditorConfig(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetEditorConfig_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.GetEditorConfig(ctx, "core", "go-forge", "README.md", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ApplyDiffPatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ApplyDiffPatch(ctx, "core", "go-forge", &types.UpdateFileOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ApplyDiffPatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ApplyDiffPatch(ctx, "core", "go-forge", &types.UpdateFileOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ApplyDiffPatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ApplyDiffPatch(ctx, "core", "go-forge", &types.UpdateFileOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetLanguages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetLanguages(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetLanguages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetLanguages(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetLanguages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetLanguages(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListFlags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListFlags(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListFlags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListFlags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListFlags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListFlags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterFlags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterFlags(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterFlags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterFlags(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterFlags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterFlags(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ReplaceFlags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.ReplaceFlags(ctx, "core", "go-forge", &types.ReplaceFlagsOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ReplaceFlags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.ReplaceFlags(ctx, "core", "go-forge", &types.ReplaceFlagsOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ReplaceFlags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.ReplaceFlags(ctx, "core", "go-forge", &types.ReplaceFlagsOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteFlags_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteFlags(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteFlags_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteFlags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteFlags_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteFlags(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetSigningKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetSigningKey(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetSigningKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetSigningKey(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetSigningKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetSigningKey(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListIssueTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListIssueTemplates(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListIssueTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListIssueTemplates(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListIssueTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListIssueTemplates(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterIssueTemplates_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterIssueTemplates(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterIssueTemplates_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterIssueTemplates(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterIssueTemplates_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterIssueTemplates(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetIssueConfig_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetIssueConfig(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetIssueConfig_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetIssueConfig(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetIssueConfig_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetIssueConfig(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ValidateIssueConfig_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ValidateIssueConfig(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ValidateIssueConfig_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ValidateIssueConfig(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ValidateIssueConfig_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ValidateIssueConfig(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListActivityFeeds(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListActivityFeeds(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListActivityFeeds(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterActivityFeeds(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterActivityFeeds(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterActivityFeeds(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListTopics_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListTopics(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTopics_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListTopics(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTopics_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListTopics(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterTopics_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterTopics(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTopics_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterTopics(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTopics_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterTopics(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_SearchTopics_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.SearchTopics(ctx, "go")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchTopics_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.SearchTopics(ctx, "go")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchTopics_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.SearchTopics(ctx, "go")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterSearchTopics_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterSearchTopics(ctx, "go") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSearchTopics_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterSearchTopics(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSearchTopics_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterSearchTopics(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_SearchRepositoriesPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.SearchRepositoriesPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchRepositoriesPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.SearchRepositoriesPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchRepositoriesPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.SearchRepositoriesPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_SearchRepositories_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.SearchRepositories(ctx, "go")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchRepositories_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.SearchRepositories(ctx, "go")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SearchRepositories_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.SearchRepositories(ctx, "go")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterSearchRepositories_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterSearchRepositories(ctx, "go") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSearchRepositories_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterSearchRepositories(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterSearchRepositories_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterSearchRepositories(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_UpdateTopics_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.UpdateTopics(ctx, "core", "go-forge", []string{"value"})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateTopics_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.UpdateTopics(ctx, "core", "go-forge", []string{"value"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateTopics_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.UpdateTopics(ctx, "core", "go-forge", []string{"value"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_AddTopic_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.AddTopic(ctx, "core", "go-forge", "go")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddTopic_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.AddTopic(ctx, "core", "go-forge", "go")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddTopic_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.AddTopic(ctx, "core", "go-forge", "go")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteTopic_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteTopic(ctx, "core", "go-forge", "go")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTopic_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteTopic(ctx, "core", "go-forge", "go")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteTopic_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteTopic(ctx, "core", "go-forge", "go")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_AddFlag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.AddFlag(ctx, "core", "go-forge", "triaged")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddFlag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.AddFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AddFlag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.AddFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_HasFlag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.HasFlag(ctx, "core", "go-forge", "triaged")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_HasFlag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.HasFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_HasFlag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.HasFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_RemoveFlag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.RemoveFlag(ctx, "core", "go-forge", "triaged")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_RemoveFlag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.RemoveFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_RemoveFlag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.RemoveFlag(ctx, "core", "go-forge", "triaged")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetNewPinAllowed_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetNewPinAllowed(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetNewPinAllowed_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetNewPinAllowed(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetNewPinAllowed_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetNewPinAllowed(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListPinnedPullRequests_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListPinnedPullRequests(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListPinnedPullRequests_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListPinnedPullRequests(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListPinnedPullRequests_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListPinnedPullRequests(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterPinnedPullRequests_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterPinnedPullRequests(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterPinnedPullRequests_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterPinnedPullRequests(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterPinnedPullRequests_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterPinnedPullRequests(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_UpdateAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.UpdateAvatar(ctx, "core", "go-forge", &types.UpdateRepoAvatarOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.UpdateAvatar(ctx, "core", "go-forge", &types.UpdateRepoAvatarOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_UpdateAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.UpdateAvatar(ctx, "core", "go-forge", &types.UpdateRepoAvatarOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeleteAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeleteAvatar(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeleteAvatar(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeleteAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeleteAvatar(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListPushMirrors_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListPushMirrors(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListPushMirrors_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListPushMirrors(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListPushMirrors_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListPushMirrors(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterPushMirrors_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterPushMirrors(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterPushMirrors_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterPushMirrors(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterPushMirrors_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterPushMirrors(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetPushMirror_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetPushMirror(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetPushMirror_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetPushMirror(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetPushMirror_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetPushMirror(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_CreatePushMirror_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.CreatePushMirror(ctx, "core", "go-forge", &types.CreatePushMirrorOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreatePushMirror_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.CreatePushMirror(ctx, "core", "go-forge", &types.CreatePushMirrorOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_CreatePushMirror_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.CreatePushMirror(ctx, "core", "go-forge", &types.CreatePushMirrorOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_DeletePushMirror_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.DeletePushMirror(ctx, "core", "go-forge", "name")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeletePushMirror_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.DeletePushMirror(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_DeletePushMirror_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.DeletePushMirror(ctx, "core", "go-forge", "name")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetSubscription_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetSubscription(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetSubscription_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetSubscription(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetSubscription_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetSubscription(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Watch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Watch(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Watch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Watch(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Watch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Watch(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Unwatch_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.Unwatch(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Unwatch_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.Unwatch(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Unwatch_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.Unwatch(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Fork_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Fork(ctx, "core", "go-forge", "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Fork_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Fork(ctx, "core", "go-forge", "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Fork_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Fork(ctx, "core", "go-forge", "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ForkWithOptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ForkWithOptions(ctx, "core", "go-forge", &types.CreateForkOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ForkWithOptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ForkWithOptions(ctx, "core", "go-forge", &types.CreateForkOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ForkWithOptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ForkWithOptions(ctx, "core", "go-forge", &types.CreateForkOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Generate_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Generate(ctx, "value", "value", &types.GenerateRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Generate_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Generate(ctx, "value", "value", &types.GenerateRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Generate_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Generate(ctx, "value", "value", &types.GenerateRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListForks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListForks(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListForks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListForks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListForks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListForks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterForks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterForks(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterForks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterForks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterForks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterForks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_Transfer_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.Transfer(ctx, "core", "go-forge", &types.TransferRepoOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Transfer_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.Transfer(ctx, "core", "go-forge", &types.TransferRepoOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_Transfer_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.Transfer(ctx, "core", "go-forge", &types.TransferRepoOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_AcceptTransfer_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.AcceptTransfer(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AcceptTransfer_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.AcceptTransfer(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_AcceptTransfer_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.AcceptTransfer(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_RejectTransfer_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.RejectTransfer(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_RejectTransfer_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.RejectTransfer(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_RejectTransfer_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.RejectTransfer(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_MirrorSync_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.MirrorSync(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_MirrorSync_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.MirrorSync(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_MirrorSync_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.MirrorSync(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetRunnerRegistrationToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetRunnerRegistrationToken(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRunnerRegistrationToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetRunnerRegistrationToken(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetRunnerRegistrationToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetRunnerRegistrationToken(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_SyncPushMirrors_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Repos.SyncPushMirrors(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SyncPushMirrors_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Repos.SyncPushMirrors(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_SyncPushMirrors_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Repos.SyncPushMirrors(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetBlob_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetBlob(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetBlob_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetBlob(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetBlob_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetBlob(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListGitRefs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListGitRefs(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListGitRefs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListGitRefs(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListGitRefs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListGitRefs(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterGitRefs_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterGitRefs(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterGitRefs_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterGitRefs(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterGitRefs_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterGitRefs(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListGitRefsByRef_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListGitRefsByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListGitRefsByRef_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListGitRefsByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListGitRefsByRef_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListGitRefsByRef(ctx, "core", "go-forge", "heads/main")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterGitRefsByRef_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterGitRefsByRef(ctx, "core", "go-forge", "heads/main") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterGitRefsByRef_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterGitRefsByRef(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterGitRefsByRef_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterGitRefsByRef(ctx, "core", "go-forge", "heads/main") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetAnnotatedTag_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetAnnotatedTag(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetAnnotatedTag_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetAnnotatedTag(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetAnnotatedTag_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetAnnotatedTag(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GetTree_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.GetTree(ctx, "core", "go-forge", "abc123")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTree_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.GetTree(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_GetTree_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.GetTree(ctx, "core", "go-forge", "abc123")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListTimes(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListTimes(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListTimes(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterTimes(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterTimes(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterTimes(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_ListUserTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Repos.ListUserTimes(ctx, "core", "go-forge", "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Repos.ListUserTimes(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_ListUserTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Repos.ListUserTimes(ctx, "core", "go-forge", "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_IterUserTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Repos.IterUserTimes(ctx, "core", "go-forge", "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterUserTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Repos.IterUserTimes(ctx, "core", "go-forge", "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_RepoService_IterUserTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Repos.IterUserTimes(ctx, "core", "go-forge", "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_String_Good(t *core.T) {
	res, _ := ax7Resource(http.StatusOK)
	got := res.String()
	core.AssertContains(t, got, "forge.Resource")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Resource_String_Bad(t *core.T) {
	var res *Resource[ax7Payload, ax7Payload, ax7Payload]
	got := res.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Resource_String_Ugly(t *core.T) {
	res := NewResource[ax7Payload, ax7Payload, ax7Payload](nil, "/api/v1/static")
	got := fmt.Sprintf("%#v", res)
	core.AssertContains(t, got, "forge.Resource")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_Resource_GoString_Good(t *core.T) {
	res, _ := ax7Resource(http.StatusOK)
	got := res.GoString()
	core.AssertContains(t, got, "forge.Resource")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Resource_GoString_Bad(t *core.T) {
	var res *Resource[ax7Payload, ax7Payload, ax7Payload]
	got := res.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_Resource_GoString_Ugly(t *core.T) {
	res := NewResource[ax7Payload, ax7Payload, ax7Payload](nil, "/api/v1/static")
	got := fmt.Sprintf("%#v", res)
	core.AssertContains(t, got, "forge.Resource")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}

func TestAX7_NewResource_Good(t *core.T) {
	client, _ := ax7Client(http.StatusOK)
	res := NewResource[ax7Payload, ax7Payload, ax7Payload](client, "/api/v1/items/{id}")
	core.AssertContains(t, res.String(), "collection")
}
func TestAX7_NewResource_Bad(t *core.T) {
	res := NewResource[ax7Payload, ax7Payload, ax7Payload](nil, "")
	got := res.String()
	core.AssertContains(t, got, "path")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_NewResource_Ugly(t *core.T) {
	client, _ := ax7Client(http.StatusOK)
	res := NewResource[ax7Payload, ax7Payload, ax7Payload](client, "/api/v1/static")
	core.AssertContains(t, res.String(), "static")
}

func TestAX7_Resource_List_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	got, err := res.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_List_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := res.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_List_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := res.List(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_ListAll_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	got, err := res.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_ListAll_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := res.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_ListAll_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := res.ListAll(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_Iter_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	for _, err := range res.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Iter_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range res.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Iter_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range res.Iter(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_Get_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	got, err := res.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Get_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := res.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Get_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := res.Get(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_Create_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	got, err := res.Create(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "created"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Create_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := res.Create(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "created"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Create_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := res.Create(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "created"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_Update_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	got, err := res.Update(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "updated"})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Update_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := res.Update(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "updated"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Update_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := res.Update(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"}, &ax7Payload{Name: "updated"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_Resource_Delete_Good(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := context.Background()
	err := res.Delete(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Delete_Bad(t *core.T) {
	res, tr := ax7Resource(http.StatusInternalServerError)
	ctx := context.Background()
	err := res.Delete(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_Resource_Delete_Ugly(t *core.T) {
	res, tr := ax7Resource(http.StatusOK)
	ctx := ax7CanceledContext()
	err := res.Delete(ctx, Params{"owner": "core", "repo": "go-forge", "id": "1", "index": "1"})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Actions.String()
	core.AssertContains(t, got, "forge.ActionsService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_String_Bad(t *core.T) {
	var svc *ActionsService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ActionsService")
}

func TestAX7_ActionsService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Actions)
	core.AssertContains(t, got, "forge.ActionsService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Actions.GoString()
	core.AssertContains(t, got, "forge.ActionsService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActionsService_GoString_Bad(t *core.T) {
	var svc *ActionsService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ActionsService")
}

func TestAX7_ActionsService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Actions)
	core.AssertContains(t, got, "forge.ActionsService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.ActivityPub.String()
	core.AssertContains(t, got, "forge.ActivityPubService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_String_Bad(t *core.T) {
	var svc *ActivityPubService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ActivityPubService")
}

func TestAX7_ActivityPubService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.ActivityPub)
	core.AssertContains(t, got, "forge.ActivityPubService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.ActivityPub.GoString()
	core.AssertContains(t, got, "forge.ActivityPubService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ActivityPubService_GoString_Bad(t *core.T) {
	var svc *ActivityPubService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ActivityPubService")
}

func TestAX7_ActivityPubService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.ActivityPub)
	core.AssertContains(t, got, "forge.ActivityPubService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Admin.String()
	core.AssertContains(t, got, "forge.AdminService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_String_Bad(t *core.T) {
	var svc *AdminService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.AdminService")
}

func TestAX7_AdminService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Admin)
	core.AssertContains(t, got, "forge.AdminService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Admin.GoString()
	core.AssertContains(t, got, "forge.AdminService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_AdminService_GoString_Bad(t *core.T) {
	var svc *AdminService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.AdminService")
}

func TestAX7_AdminService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Admin)
	core.AssertContains(t, got, "forge.AdminService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Branches.String()
	core.AssertContains(t, got, "forge.BranchService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_String_Bad(t *core.T) {
	var svc *BranchService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.BranchService")
}

func TestAX7_BranchService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Branches)
	core.AssertContains(t, got, "forge.BranchService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Branches.GoString()
	core.AssertContains(t, got, "forge.BranchService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_BranchService_GoString_Bad(t *core.T) {
	var svc *BranchService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.BranchService")
}

func TestAX7_BranchService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Branches)
	core.AssertContains(t, got, "forge.BranchService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Commits.String()
	core.AssertContains(t, got, "forge.CommitService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_String_Bad(t *core.T) {
	var svc *CommitService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.CommitService")
}

func TestAX7_CommitService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Commits)
	core.AssertContains(t, got, "forge.CommitService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Commits.GoString()
	core.AssertContains(t, got, "forge.CommitService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_CommitService_GoString_Bad(t *core.T) {
	var svc *CommitService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.CommitService")
}

func TestAX7_CommitService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Commits)
	core.AssertContains(t, got, "forge.CommitService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Contents.String()
	core.AssertContains(t, got, "forge.ContentService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_String_Bad(t *core.T) {
	var svc *ContentService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ContentService")
}

func TestAX7_ContentService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Contents)
	core.AssertContains(t, got, "forge.ContentService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Contents.GoString()
	core.AssertContains(t, got, "forge.ContentService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ContentService_GoString_Bad(t *core.T) {
	var svc *ContentService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ContentService")
}

func TestAX7_ContentService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Contents)
	core.AssertContains(t, got, "forge.ContentService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Issues.String()
	core.AssertContains(t, got, "forge.IssueService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_String_Bad(t *core.T) {
	var svc *IssueService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.IssueService")
}

func TestAX7_IssueService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Issues)
	core.AssertContains(t, got, "forge.IssueService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Issues.GoString()
	core.AssertContains(t, got, "forge.IssueService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_IssueService_GoString_Bad(t *core.T) {
	var svc *IssueService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.IssueService")
}

func TestAX7_IssueService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Issues)
	core.AssertContains(t, got, "forge.IssueService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Labels.String()
	core.AssertContains(t, got, "forge.LabelService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_String_Bad(t *core.T) {
	var svc *LabelService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.LabelService")
}

func TestAX7_LabelService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Labels)
	core.AssertContains(t, got, "forge.LabelService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Labels.GoString()
	core.AssertContains(t, got, "forge.LabelService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_LabelService_GoString_Bad(t *core.T) {
	var svc *LabelService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.LabelService")
}

func TestAX7_LabelService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Labels)
	core.AssertContains(t, got, "forge.LabelService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Milestones.String()
	core.AssertContains(t, got, "forge.MilestoneService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_String_Bad(t *core.T) {
	var svc *MilestoneService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.MilestoneService")
}

func TestAX7_MilestoneService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Milestones)
	core.AssertContains(t, got, "forge.MilestoneService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Milestones.GoString()
	core.AssertContains(t, got, "forge.MilestoneService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MilestoneService_GoString_Bad(t *core.T) {
	var svc *MilestoneService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.MilestoneService")
}

func TestAX7_MilestoneService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Milestones)
	core.AssertContains(t, got, "forge.MilestoneService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Misc.String()
	core.AssertContains(t, got, "forge.MiscService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_String_Bad(t *core.T) {
	var svc *MiscService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.MiscService")
}

func TestAX7_MiscService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Misc)
	core.AssertContains(t, got, "forge.MiscService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Misc.GoString()
	core.AssertContains(t, got, "forge.MiscService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_MiscService_GoString_Bad(t *core.T) {
	var svc *MiscService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.MiscService")
}

func TestAX7_MiscService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Misc)
	core.AssertContains(t, got, "forge.MiscService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Notifications.String()
	core.AssertContains(t, got, "forge.NotificationService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_String_Bad(t *core.T) {
	var svc *NotificationService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.NotificationService")
}

func TestAX7_NotificationService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Notifications)
	core.AssertContains(t, got, "forge.NotificationService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Notifications.GoString()
	core.AssertContains(t, got, "forge.NotificationService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_NotificationService_GoString_Bad(t *core.T) {
	var svc *NotificationService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.NotificationService")
}

func TestAX7_NotificationService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Notifications)
	core.AssertContains(t, got, "forge.NotificationService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Orgs.String()
	core.AssertContains(t, got, "forge.OrgService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_String_Bad(t *core.T) {
	var svc *OrgService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.OrgService")
}

func TestAX7_OrgService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Orgs)
	core.AssertContains(t, got, "forge.OrgService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Orgs.GoString()
	core.AssertContains(t, got, "forge.OrgService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_OrgService_GoString_Bad(t *core.T) {
	var svc *OrgService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.OrgService")
}

func TestAX7_OrgService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Orgs)
	core.AssertContains(t, got, "forge.OrgService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Packages.String()
	core.AssertContains(t, got, "forge.PackageService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_String_Bad(t *core.T) {
	var svc *PackageService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.PackageService")
}

func TestAX7_PackageService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Packages)
	core.AssertContains(t, got, "forge.PackageService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Packages.GoString()
	core.AssertContains(t, got, "forge.PackageService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PackageService_GoString_Bad(t *core.T) {
	var svc *PackageService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.PackageService")
}

func TestAX7_PackageService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Packages)
	core.AssertContains(t, got, "forge.PackageService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Pulls.String()
	core.AssertContains(t, got, "forge.PullService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_String_Bad(t *core.T) {
	var svc *PullService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.PullService")
}

func TestAX7_PullService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Pulls)
	core.AssertContains(t, got, "forge.PullService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Pulls.GoString()
	core.AssertContains(t, got, "forge.PullService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_PullService_GoString_Bad(t *core.T) {
	var svc *PullService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.PullService")
}

func TestAX7_PullService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Pulls)
	core.AssertContains(t, got, "forge.PullService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Releases.String()
	core.AssertContains(t, got, "forge.ReleaseService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_String_Bad(t *core.T) {
	var svc *ReleaseService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ReleaseService")
}

func TestAX7_ReleaseService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Releases)
	core.AssertContains(t, got, "forge.ReleaseService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Releases.GoString()
	core.AssertContains(t, got, "forge.ReleaseService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_ReleaseService_GoString_Bad(t *core.T) {
	var svc *ReleaseService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.ReleaseService")
}

func TestAX7_ReleaseService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Releases)
	core.AssertContains(t, got, "forge.ReleaseService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Repos.String()
	core.AssertContains(t, got, "forge.RepoService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_String_Bad(t *core.T) {
	var svc *RepoService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.RepoService")
}

func TestAX7_RepoService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Repos)
	core.AssertContains(t, got, "forge.RepoService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Repos.GoString()
	core.AssertContains(t, got, "forge.RepoService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_RepoService_GoString_Bad(t *core.T) {
	var svc *RepoService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.RepoService")
}

func TestAX7_RepoService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Repos)
	core.AssertContains(t, got, "forge.RepoService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Teams.String()
	core.AssertContains(t, got, "forge.TeamService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_String_Bad(t *core.T) {
	var svc *TeamService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.TeamService")
}

func TestAX7_TeamService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Teams)
	core.AssertContains(t, got, "forge.TeamService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Teams.GoString()
	core.AssertContains(t, got, "forge.TeamService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_GoString_Bad(t *core.T) {
	var svc *TeamService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.TeamService")
}

func TestAX7_TeamService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Teams)
	core.AssertContains(t, got, "forge.TeamService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Users.String()
	core.AssertContains(t, got, "forge.UserService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_String_Bad(t *core.T) {
	var svc *UserService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.UserService")
}

func TestAX7_UserService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Users)
	core.AssertContains(t, got, "forge.UserService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Users.GoString()
	core.AssertContains(t, got, "forge.UserService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GoString_Bad(t *core.T) {
	var svc *UserService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.UserService")
}

func TestAX7_UserService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Users)
	core.AssertContains(t, got, "forge.UserService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Webhooks.String()
	core.AssertContains(t, got, "forge.WebhookService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_String_Bad(t *core.T) {
	var svc *WebhookService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.WebhookService")
}

func TestAX7_WebhookService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Webhooks)
	core.AssertContains(t, got, "forge.WebhookService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Webhooks.GoString()
	core.AssertContains(t, got, "forge.WebhookService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GoString_Bad(t *core.T) {
	var svc *WebhookService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.WebhookService")
}

func TestAX7_WebhookService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Webhooks)
	core.AssertContains(t, got, "forge.WebhookService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_String_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Wiki.String()
	core.AssertContains(t, got, "forge.WikiService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_String_Bad(t *core.T) {
	var svc *WikiService
	got := svc.String()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.WikiService")
}

func TestAX7_WikiService_String_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Wiki)
	core.AssertContains(t, got, "forge.WikiService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_GoString_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fg.Wiki.GoString()
	core.AssertContains(t, got, "forge.WikiService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_GoString_Bad(t *core.T) {
	var svc *WikiService
	got := svc.GoString()
	core.AssertContains(t, got, "<nil>")
	core.AssertContains(t, got, "forge.WikiService")
}

func TestAX7_WikiService_GoString_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	got := fmt.Sprintf("%#v", fg.Wiki)
	core.AssertContains(t, got, "forge.WikiService")
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_CreateOrgTeam_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.CreateOrgTeam(ctx, "core", &types.CreateTeamOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_CreateOrgTeam_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.CreateOrgTeam(ctx, "core", &types.CreateTeamOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_CreateOrgTeam_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.CreateOrgTeam(ctx, "core", &types.CreateTeamOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_ListMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.ListMembers(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.ListMembers(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.ListMembers(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_IterMembers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Teams.IterMembers(ctx, 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterMembers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Teams.IterMembers(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterMembers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Teams.IterMembers(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_AddMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Teams.AddMember(ctx, 1, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_AddMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Teams.AddMember(ctx, 1, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_AddMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Teams.AddMember(ctx, 1, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_GetMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.GetMember(ctx, 1, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_GetMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.GetMember(ctx, 1, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_GetMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.GetMember(ctx, 1, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_RemoveMember_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Teams.RemoveMember(ctx, 1, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_RemoveMember_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Teams.RemoveMember(ctx, 1, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_RemoveMember_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Teams.RemoveMember(ctx, 1, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_ListRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.ListRepos(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.ListRepos(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.ListRepos(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_IterRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Teams.IterRepos(ctx, 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Teams.IterRepos(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Teams.IterRepos(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_AddRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Teams.AddRepo(ctx, 1, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_AddRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Teams.AddRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_AddRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Teams.AddRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_RemoveRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Teams.RemoveRepo(ctx, 1, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_RemoveRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Teams.RemoveRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_RemoveRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Teams.RemoveRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_GetRepo_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.GetRepo(ctx, 1, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_GetRepo_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.GetRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_GetRepo_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.GetRepo(ctx, 1, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_ListOrgTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.ListOrgTeams(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListOrgTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.ListOrgTeams(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListOrgTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.ListOrgTeams(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_IterOrgTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Teams.IterOrgTeams(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterOrgTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Teams.IterOrgTeams(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterOrgTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Teams.IterOrgTeams(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_ListActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Teams.ListActivityFeeds(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Teams.ListActivityFeeds(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_ListActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Teams.ListActivityFeeds(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_TeamService_IterActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Teams.IterActivityFeeds(ctx, 1) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Teams.IterActivityFeeds(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_TeamService_IterActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Teams.IterActivityFeeds(ctx, 1) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserSearchOptions_String_Good(t *core.T) {
	value := UserSearchOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_UserSearchOptions_String_Bad(t *core.T) {
	value := UserSearchOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_UserSearchOptions_String_Ugly(t *core.T) {
	value := UserSearchOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_UserSearchOptions_GoString_Good(t *core.T) {
	value := UserSearchOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_UserSearchOptions_GoString_Bad(t *core.T) {
	value := UserSearchOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_UserSearchOptions_GoString_Ugly(t *core.T) {
	value := UserSearchOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_UserKeyListOptions_String_Good(t *core.T) {
	value := UserKeyListOptions{}
	got := value.String()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_UserKeyListOptions_String_Bad(t *core.T) {
	value := UserKeyListOptions{}
	got := value.String()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_UserKeyListOptions_String_Ugly(t *core.T) {
	value := UserKeyListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_UserKeyListOptions_GoString_Good(t *core.T) {
	value := UserKeyListOptions{}
	got := value.GoString()
	core.AssertContains(t, got, "forge")
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
}
func TestAX7_UserKeyListOptions_GoString_Bad(t *core.T) {
	value := UserKeyListOptions{}
	got := value.GoString()
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}
func TestAX7_UserKeyListOptions_GoString_Ugly(t *core.T) {
	value := UserKeyListOptions{}
	got := fmt.Sprintf("%#v", value)
	core.AssertNotEmpty(t, got)
	core.AssertContains(t, got, "forge")
}

func TestAX7_UserService_GetUserByID_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetUserByID(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetUserByID_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetUserByID(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetUserByID_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetUserByID(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetUserByUsername_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetUserByUsername(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetUserByUsername_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetUserByUsername(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetUserByUsername_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetUserByUsername(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetCurrent_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetCurrent(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetCurrent_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetCurrent(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetCurrent_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetCurrent(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetSettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetSettings(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetSettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetSettings(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetSettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetSettings(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_UpdateSettings_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.UpdateSettings(ctx, &types.UserSettingsOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateSettings_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.UpdateSettings(ctx, &types.UserSettingsOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateSettings_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.UpdateSettings(ctx, &types.UserSettingsOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetQuota_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetQuota(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetQuota_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetQuota(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetQuota_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetQuota(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_SearchUsersPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.SearchUsersPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_SearchUsersPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.SearchUsersPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_SearchUsersPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.SearchUsersPage(ctx, "go", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_SearchUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.SearchUsers(ctx, "go")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_SearchUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.SearchUsers(ctx, "go")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_SearchUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.SearchUsers(ctx, "go")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterSearchUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterSearchUsers(ctx, "go") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterSearchUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterSearchUsers(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterSearchUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterSearchUsers(ctx, "go") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListQuotaArtifacts_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListQuotaArtifacts(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaArtifacts_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListQuotaArtifacts(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaArtifacts_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListQuotaArtifacts(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterQuotaArtifacts_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterQuotaArtifacts(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaArtifacts_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterQuotaArtifacts(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaArtifacts_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterQuotaArtifacts(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListQuotaAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListQuotaAttachments(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListQuotaAttachments(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListQuotaAttachments(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterQuotaAttachments_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterQuotaAttachments(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaAttachments_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterQuotaAttachments(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaAttachments_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterQuotaAttachments(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListQuotaPackages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListQuotaPackages(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaPackages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListQuotaPackages(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListQuotaPackages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListQuotaPackages(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterQuotaPackages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterQuotaPackages(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaPackages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterQuotaPackages(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterQuotaPackages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterQuotaPackages(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListEmails(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListEmails(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListEmails(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterEmails(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterEmails(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterEmails(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_AddEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.AddEmails(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_AddEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.AddEmails(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_AddEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.AddEmails(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteEmails_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteEmails(ctx)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteEmails_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteEmails(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteEmails_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteEmails(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_UpdateAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.UpdateAvatar(ctx, &types.UpdateUserAvatarOption{})
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.UpdateAvatar(ctx, &types.UpdateUserAvatarOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.UpdateAvatar(ctx, &types.UpdateUserAvatarOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteAvatar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteAvatar(ctx)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteAvatar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteAvatar(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteAvatar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteAvatar(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListKeys(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListKeys(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListKeys(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterKeys(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterKeys(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterKeys(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CreateKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CreateKey(ctx, &types.CreateKeyOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CreateKey(ctx, &types.CreateKeyOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CreateKey(ctx, &types.CreateKeyOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetKey(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetKey(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetKey(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteKey(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteKey(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteKey(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListUserKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListUserKeys(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListUserKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListUserKeys(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListUserKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListUserKeys(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterUserKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterUserKeys(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterUserKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterUserKeys(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterUserKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterUserKeys(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListGPGKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListGPGKeys(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListGPGKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListGPGKeys(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListGPGKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListGPGKeys(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterGPGKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterGPGKeys(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterGPGKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterGPGKeys(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterGPGKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterGPGKeys(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CreateGPGKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CreateGPGKey(ctx, &types.CreateGPGKeyOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateGPGKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CreateGPGKey(ctx, &types.CreateGPGKeyOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateGPGKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CreateGPGKey(ctx, &types.CreateGPGKeyOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetGPGKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetGPGKey(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetGPGKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetGPGKey(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetGPGKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetGPGKey(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteGPGKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteGPGKey(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteGPGKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteGPGKey(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteGPGKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteGPGKey(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListUserGPGKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListUserGPGKeys(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListUserGPGKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListUserGPGKeys(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListUserGPGKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListUserGPGKeys(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterUserGPGKeys_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterUserGPGKeys(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterUserGPGKeys_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterUserGPGKeys(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterUserGPGKeys_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterUserGPGKeys(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetGPGKeyVerificationToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetGPGKeyVerificationToken(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetGPGKeyVerificationToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetGPGKeyVerificationToken(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetGPGKeyVerificationToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetGPGKeyVerificationToken(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_VerifyGPGKey_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.VerifyGPGKey(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_VerifyGPGKey_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.VerifyGPGKey(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_VerifyGPGKey_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.VerifyGPGKey(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListTokens_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListTokens(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListTokens_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListTokens(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListTokens_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListTokens(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterTokens_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterTokens(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterTokens_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterTokens(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterTokens_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterTokens(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CreateToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CreateToken(ctx, "alice", &types.CreateAccessTokenOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CreateToken(ctx, "alice", &types.CreateAccessTokenOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CreateToken(ctx, "alice", &types.CreateAccessTokenOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteToken(ctx, "alice", "value")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteToken(ctx, "alice", "value")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteToken(ctx, "alice", "value")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListOAuth2Applications_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListOAuth2Applications(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListOAuth2Applications_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListOAuth2Applications(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListOAuth2Applications_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListOAuth2Applications(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterOAuth2Applications_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterOAuth2Applications(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterOAuth2Applications_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterOAuth2Applications(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterOAuth2Applications_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterOAuth2Applications(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CreateOAuth2Application_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CreateOAuth2Application(ctx, &types.CreateOAuth2ApplicationOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateOAuth2Application_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CreateOAuth2Application(ctx, &types.CreateOAuth2ApplicationOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CreateOAuth2Application_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CreateOAuth2Application(ctx, &types.CreateOAuth2ApplicationOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetOAuth2Application_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetOAuth2Application(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetOAuth2Application_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetOAuth2Application(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetOAuth2Application_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetOAuth2Application(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_UpdateOAuth2Application_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.UpdateOAuth2Application(ctx, 1, &types.CreateOAuth2ApplicationOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateOAuth2Application_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.UpdateOAuth2Application(ctx, 1, &types.CreateOAuth2ApplicationOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_UpdateOAuth2Application_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.UpdateOAuth2Application(ctx, 1, &types.CreateOAuth2ApplicationOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_DeleteOAuth2Application_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.DeleteOAuth2Application(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteOAuth2Application_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.DeleteOAuth2Application(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_DeleteOAuth2Application_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.DeleteOAuth2Application(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListStopwatches_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListStopwatches(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListStopwatches_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListStopwatches(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListStopwatches_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListStopwatches(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterStopwatches_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterStopwatches(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterStopwatches_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterStopwatches(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterStopwatches_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterStopwatches(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListBlockedUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListBlockedUsers(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListBlockedUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListBlockedUsers(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListBlockedUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListBlockedUsers(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterBlockedUsers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterBlockedUsers(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterBlockedUsers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterBlockedUsers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterBlockedUsers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterBlockedUsers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Block_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Block(ctx, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Block_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Block(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Block_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Block(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Unblock_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Unblock(ctx, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unblock_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Unblock(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unblock_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Unblock(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMySubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMySubscriptions(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMySubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMySubscriptions(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMySubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMySubscriptions(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMySubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMySubscriptions(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMySubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMySubscriptions(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMySubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMySubscriptions(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMyStarred_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMyStarred(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyStarred_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMyStarred(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyStarred_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMyStarred(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMyStarred_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMyStarred(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyStarred_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMyStarred(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyStarred_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMyStarred(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMyFollowers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMyFollowers(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyFollowers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMyFollowers(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyFollowers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMyFollowers(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMyFollowers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMyFollowers(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyFollowers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMyFollowers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyFollowers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMyFollowers(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMyFollowing_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMyFollowing(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyFollowing_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMyFollowing(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyFollowing_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMyFollowing(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMyFollowing_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMyFollowing(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyFollowing_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMyFollowing(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyFollowing_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMyFollowing(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMyTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMyTeams(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMyTeams(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMyTeams(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMyTeams_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMyTeams(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyTeams_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMyTeams(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyTeams_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMyTeams(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListMyTrackedTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListMyTrackedTimes(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyTrackedTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListMyTrackedTimes(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListMyTrackedTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListMyTrackedTimes(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterMyTrackedTimes_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterMyTrackedTimes(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyTrackedTimes_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterMyTrackedTimes(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterMyTrackedTimes_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterMyTrackedTimes(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CheckQuota_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CheckQuota(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckQuota_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CheckQuota(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckQuota_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CheckQuota(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetRunnerRegistrationToken_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetRunnerRegistrationToken(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetRunnerRegistrationToken_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetRunnerRegistrationToken(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetRunnerRegistrationToken_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetRunnerRegistrationToken(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListFollowers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListFollowers(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListFollowers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListFollowers(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListFollowers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListFollowers(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterFollowers_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterFollowers(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterFollowers_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterFollowers(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterFollowers_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterFollowers(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListSubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListSubscriptions(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListSubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListSubscriptions(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListSubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListSubscriptions(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterSubscriptions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterSubscriptions(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterSubscriptions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterSubscriptions(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterSubscriptions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterSubscriptions(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListFollowing_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListFollowing(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListFollowing_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListFollowing(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListFollowing_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListFollowing(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterFollowing_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterFollowing(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterFollowing_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterFollowing(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterFollowing_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterFollowing(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListActivityFeeds(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListActivityFeeds(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListActivityFeeds(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterActivityFeeds_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterActivityFeeds(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterActivityFeeds_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterActivityFeeds(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterActivityFeeds_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterActivityFeeds(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListRepos(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListRepos(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListRepos(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterRepos_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterRepos(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterRepos_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterRepos(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterRepos_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterRepos(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Follow_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Follow(ctx, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Follow_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Follow(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Follow_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Follow(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CheckFollowing_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CheckFollowing(ctx, "alice", "bob")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckFollowing_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CheckFollowing(ctx, "alice", "bob")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckFollowing_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CheckFollowing(ctx, "alice", "bob")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Unfollow_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Unfollow(ctx, "alice")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unfollow_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Unfollow(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unfollow_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Unfollow(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_ListStarred_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.ListStarred(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListStarred_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.ListStarred(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_ListStarred_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.ListStarred(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_IterStarred_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Users.IterStarred(ctx, "alice") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterStarred_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Users.IterStarred(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_IterStarred_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Users.IterStarred(ctx, "alice") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_GetHeatmap_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.GetHeatmap(ctx, "alice")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetHeatmap_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.GetHeatmap(ctx, "alice")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_GetHeatmap_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.GetHeatmap(ctx, "alice")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Star_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Star(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Star_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Star(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Star_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Star(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_Unstar_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Users.Unstar(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unstar_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Users.Unstar(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_Unstar_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Users.Unstar(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_UserService_CheckStarring_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Users.CheckStarring(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckStarring_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Users.CheckStarring(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_UserService_CheckStarring_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Users.CheckStarring(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListHooksPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListHooksPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListHooksPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListHooks(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListRepoHooksPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListRepoHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListRepoHooksPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListRepoHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListRepoHooksPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListRepoHooksPage(ctx, "core", "go-forge", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListRepoHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListRepoHooks(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListRepoHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListRepoHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListRepoHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListRepoHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_IterHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Webhooks.IterHooks(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Webhooks.IterHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Webhooks.IterHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_IterRepoHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Webhooks.IterRepoHooks(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterRepoHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Webhooks.IterRepoHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterRepoHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Webhooks.IterRepoHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_CreateHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.CreateHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.CreateHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.CreateHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_CreateRepoHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.CreateRepoHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateRepoHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.CreateRepoHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateRepoHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.CreateRepoHook(ctx, "core", "go-forge", &types.CreateHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GetRepoHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.GetRepoHook(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetRepoHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.GetRepoHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetRepoHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.GetRepoHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_EditRepoHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.EditRepoHook(ctx, "core", "go-forge", 1, &types.EditHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditRepoHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.EditRepoHook(ctx, "core", "go-forge", 1, &types.EditHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditRepoHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.EditRepoHook(ctx, "core", "go-forge", 1, &types.EditHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_DeleteRepoHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Webhooks.DeleteRepoHook(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteRepoHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Webhooks.DeleteRepoHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteRepoHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Webhooks.DeleteRepoHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_TestHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Webhooks.TestHook(ctx, "core", "go-forge", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_TestHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Webhooks.TestHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_TestHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Webhooks.TestHook(ctx, "core", "go-forge", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListGitHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListGitHooks(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListGitHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListGitHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListGitHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListGitHooks(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_IterGitHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Webhooks.IterGitHooks(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterGitHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Webhooks.IterGitHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterGitHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Webhooks.IterGitHooks(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GetGitHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.GetGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetGitHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.GetGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetGitHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.GetGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_EditGitHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.EditGitHook(ctx, "core", "go-forge", "hook-id", &types.EditGitHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditGitHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.EditGitHook(ctx, "core", "go-forge", "hook-id", &types.EditGitHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditGitHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.EditGitHook(ctx, "core", "go-forge", "hook-id", &types.EditGitHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_DeleteGitHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Webhooks.DeleteGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteGitHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Webhooks.DeleteGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteGitHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Webhooks.DeleteGitHook(ctx, "core", "go-forge", "hook-id")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListUserHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListUserHooks(ctx)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListUserHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListUserHooks(ctx)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListUserHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListUserHooks(ctx)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListUserHooksPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListUserHooksPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListUserHooksPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListUserHooksPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListUserHooksPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListUserHooksPage(ctx, ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_IterUserHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Webhooks.IterUserHooks(ctx) {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterUserHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Webhooks.IterUserHooks(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterUserHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Webhooks.IterUserHooks(ctx) {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GetUserHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.GetUserHook(ctx, 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetUserHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.GetUserHook(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetUserHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.GetUserHook(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_CreateUserHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.CreateUserHook(ctx, &types.CreateHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateUserHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.CreateUserHook(ctx, &types.CreateHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateUserHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.CreateUserHook(ctx, &types.CreateHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_EditUserHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.EditUserHook(ctx, 1, &types.EditHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditUserHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.EditUserHook(ctx, 1, &types.EditHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditUserHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.EditUserHook(ctx, 1, &types.EditHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_DeleteUserHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Webhooks.DeleteUserHook(ctx, 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteUserHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Webhooks.DeleteUserHook(ctx, 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteUserHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Webhooks.DeleteUserHook(ctx, 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListOrgHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListOrgHooks(ctx, "core")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListOrgHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListOrgHooks(ctx, "core")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListOrgHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListOrgHooks(ctx, "core")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_ListOrgHooksPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.ListOrgHooksPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListOrgHooksPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.ListOrgHooksPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_ListOrgHooksPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.ListOrgHooksPage(ctx, "core", ListOptions{Page: 1, PageSize: 1})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_IterOrgHooks_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Webhooks.IterOrgHooks(ctx, "core") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterOrgHooks_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Webhooks.IterOrgHooks(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_IterOrgHooks_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Webhooks.IterOrgHooks(ctx, "core") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_GetOrgHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.GetOrgHook(ctx, "core", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetOrgHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.GetOrgHook(ctx, "core", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_GetOrgHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.GetOrgHook(ctx, "core", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_CreateOrgHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.CreateOrgHook(ctx, "core", &types.CreateHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateOrgHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.CreateOrgHook(ctx, "core", &types.CreateHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_CreateOrgHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.CreateOrgHook(ctx, "core", &types.CreateHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_EditOrgHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Webhooks.EditOrgHook(ctx, "core", 1, &types.EditHookOption{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditOrgHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Webhooks.EditOrgHook(ctx, "core", 1, &types.EditHookOption{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_EditOrgHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Webhooks.EditOrgHook(ctx, "core", 1, &types.EditHookOption{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WebhookService_DeleteOrgHook_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Webhooks.DeleteOrgHook(ctx, "core", 1)
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteOrgHook_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Webhooks.DeleteOrgHook(ctx, "core", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WebhookService_DeleteOrgHook_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Webhooks.DeleteOrgHook(ctx, "core", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_ListPages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Wiki.ListPages(ctx, "core", "go-forge")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_ListPages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Wiki.ListPages(ctx, "core", "go-forge")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_ListPages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Wiki.ListPages(ctx, "core", "go-forge")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_IterPages_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	for _, err := range fg.Wiki.IterPages(ctx, "core", "go-forge") {
		core.AssertNoError(t, err)
	}
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_IterPages_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	var seenErr error
	for _, err := range fg.Wiki.IterPages(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_IterPages_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	var seenErr error
	for _, err := range fg.Wiki.IterPages(ctx, "core", "go-forge") {
		seenErr = err
		break
	}
	core.AssertError(t, seenErr)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_GetPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Wiki.GetPage(ctx, "core", "go-forge", "Home")
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_GetPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Wiki.GetPage(ctx, "core", "go-forge", "Home")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_GetPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Wiki.GetPage(ctx, "core", "go-forge", "Home")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_GetPageRevisions_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Wiki.GetPageRevisions(ctx, "core", "go-forge", "Home", 1)
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_GetPageRevisions_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Wiki.GetPageRevisions(ctx, "core", "go-forge", "Home", 1)
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_GetPageRevisions_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Wiki.GetPageRevisions(ctx, "core", "go-forge", "Home", 1)
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_CreatePage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Wiki.CreatePage(ctx, "core", "go-forge", &types.CreateWikiPageOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_CreatePage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Wiki.CreatePage(ctx, "core", "go-forge", &types.CreateWikiPageOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_CreatePage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Wiki.CreatePage(ctx, "core", "go-forge", &types.CreateWikiPageOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_EditPage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	got, err := fg.Wiki.EditPage(ctx, "core", "go-forge", "Home", &types.CreateWikiPageOptions{})
	core.AssertNoError(t, err)
	core.AssertNotEmpty(t, core.Sprintf("%T", got))
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_EditPage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	_, err := fg.Wiki.EditPage(ctx, "core", "go-forge", "Home", &types.CreateWikiPageOptions{})
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_EditPage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	_, err := fg.Wiki.EditPage(ctx, "core", "go-forge", "Home", &types.CreateWikiPageOptions{})
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}

func TestAX7_WikiService_DeletePage_Good(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := context.Background()
	err := fg.Wiki.DeletePage(ctx, "core", "go-forge", "Home")
	core.AssertNoError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_DeletePage_Bad(t *core.T) {
	fg, tr := ax7Forge(http.StatusInternalServerError)
	ctx := context.Background()
	err := fg.Wiki.DeletePage(ctx, "core", "go-forge", "Home")
	core.AssertError(t, err)
	ax7AssertRequest(t, tr)
}

func TestAX7_WikiService_DeletePage_Ugly(t *core.T) {
	fg, tr := ax7Forge(http.StatusOK)
	ctx := ax7CanceledContext()
	err := fg.Wiki.DeletePage(ctx, "core", "go-forge", "Home")
	core.AssertError(t, err)
	core.AssertEqual(t, 0, tr.Count())
}
