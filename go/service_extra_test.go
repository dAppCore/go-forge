// SPDX-License-Identifier: EUPL-1.2

package forge

import (
	"testing"

	core "dappco.re/go"
)

// NewService and Register expose the Core-registerable factory shape. These
// tests wire a Service into a Core and assert the constructed (or nil-for-
// late-wiring) Forge handle.

func TestNewService_WithURL_Good(t *testing.T) {
	c := core.New(core.WithService(NewService(ServiceOptions{
		URL:   "https://forge.example.com",
		Token: "tok",
	})))

	svc, ok := core.ServiceFor[*Service](c, "forge")
	if !ok {
		t.Fatal("forge service not registered")
	}
	if svc.Forge == nil {
		t.Fatal("expected a constructed Forge when URL is set")
	}
	if got := svc.Forge.BaseURL(); got != "https://forge.example.com" {
		t.Fatalf("got baseURL=%q", got)
	}
}

func TestNewService_EmptyURL_Ugly(t *testing.T) {
	// An empty URL registers a Service for late wiring with a nil Forge.
	c := core.New(core.WithService(NewService(ServiceOptions{})))

	svc, ok := core.ServiceFor[*Service](c, "forge")
	if !ok {
		t.Fatal("forge service not registered")
	}
	if svc.Forge != nil {
		t.Fatal("expected nil Forge for empty-URL late-wiring registration")
	}
}

func TestRegister_Good(t *testing.T) {
	c := core.New(core.WithService(Register))

	svc, ok := core.ServiceFor[*Service](c, "forge")
	if !ok {
		t.Fatal("forge service not registered via Register")
	}
	if svc.Forge != nil {
		t.Fatal("Register wires a nil Forge for later assignment")
	}

	// Consumers wire a Forge in after registration.
	svc.Forge = NewForge("https://forge.example.com", "tok")
	if svc.Forge.BaseURL() != "https://forge.example.com" {
		t.Fatalf("got baseURL=%q", svc.Forge.BaseURL())
	}
}
