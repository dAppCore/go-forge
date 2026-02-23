package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forge.lthn.ai/core/go-forge/types"
)

func TestNotificationService_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/notifications" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.NotificationThread{
			{ID: 1, Unread: true, Subject: &types.NotificationSubject{Title: "Issue opened"}},
			{ID: 2, Unread: false, Subject: &types.NotificationSubject{Title: "PR merged"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	threads, err := f.Notifications.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(threads) != 2 {
		t.Fatalf("got %d threads, want 2", len(threads))
	}
	if threads[0].ID != 1 {
		t.Errorf("got id=%d, want 1", threads[0].ID)
	}
	if threads[0].Subject.Title != "Issue opened" {
		t.Errorf("got title=%q, want %q", threads[0].Subject.Title, "Issue opened")
	}
	if !threads[0].Unread {
		t.Error("expected thread 1 to be unread")
	}
}

func TestNotificationService_Good_ListRepo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/notifications" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.NotificationThread{
			{ID: 10, Unread: true, Subject: &types.NotificationSubject{Title: "New commit"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	threads, err := f.Notifications.ListRepo(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(threads) != 1 {
		t.Fatalf("got %d threads, want 1", len(threads))
	}
	if threads[0].ID != 10 {
		t.Errorf("got id=%d, want 10", threads[0].ID)
	}
}

func TestNotificationService_Good_GetThread(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/notifications/threads/42" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.NotificationThread{
			ID:     42,
			Unread: true,
			Subject: &types.NotificationSubject{
				Title: "Build failed",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	thread, err := f.Notifications.GetThread(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
	if thread.ID != 42 {
		t.Errorf("got id=%d, want 42", thread.ID)
	}
	if thread.Subject.Title != "Build failed" {
		t.Errorf("got title=%q, want %q", thread.Subject.Title, "Build failed")
	}
	if !thread.Unread {
		t.Error("expected thread to be unread")
	}
}

func TestNotificationService_Good_MarkRead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/notifications" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusResetContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Notifications.MarkRead(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationService_Good_MarkThreadRead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/notifications/threads/42" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusResetContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Notifications.MarkThreadRead(context.Background(), 42)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationService_Bad_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "thread not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Notifications.GetThread(context.Background(), 9999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
