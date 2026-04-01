package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestRepoService_ListTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.TopicName{TopicNames: []string{"go", "forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	topics, err := f.Repos.ListTopics(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(topics, []string{"go", "forge"}) {
		t.Fatalf("got %#v", topics)
	}
}

func TestRepoService_UpdateTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.RepoTopicOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(body.Topics, []string{"go", "forge"}) {
			t.Fatalf("got %#v", body.Topics)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.UpdateTopics(context.Background(), "core", "go-forge", []string{"go", "forge"}); err != nil {
		t.Fatal(err)
	}
}
