package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dappco.re/go/core/forge/types"
)

func TestNotificationService_List_Good(t *testing.T) {
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

func TestNotificationService_ListRepo_Good(t *testing.T) {
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

func TestNotificationService_NewAvailable_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/notifications/new" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.NotificationCount{New: 3})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	count, err := f.Notifications.NewAvailable(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if count.New != 3 {
		t.Fatalf("got new=%d, want 3", count.New)
	}
}

func TestNotificationService_GetThread_Good(t *testing.T) {
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

func TestNotificationService_MarkRead_Good(t *testing.T) {
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

func TestNotificationService_MarkThreadRead_Good(t *testing.T) {
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

func TestNotificationService_MarkRepoNotifications_Good(t *testing.T) {
	lastReadAt := time.Date(2026, time.April, 2, 15, 4, 5, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/notifications" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query()["status-types"]; len(got) != 2 || got[0] != "unread" || got[1] != "pinned" {
			t.Errorf("got status-types=%v, want [unread pinned]", got)
		}
		if got := r.URL.Query().Get("all"); got != "true" {
			t.Errorf("got all=%q, want true", got)
		}
		if got := r.URL.Query().Get("to-status"); got != "read" {
			t.Errorf("got to-status=%q, want read", got)
		}
		if got := r.URL.Query().Get("last_read_at"); got != lastReadAt.Format(time.RFC3339) {
			t.Errorf("got last_read_at=%q, want %q", got, lastReadAt.Format(time.RFC3339))
		}
		w.WriteHeader(http.StatusResetContent)
		json.NewEncoder(w).Encode([]types.NotificationThread{
			{ID: 7, Unread: false, Subject: &types.NotificationSubject{Title: "Pinned release"}},
			{ID: 8, Unread: false, Subject: &types.NotificationSubject{Title: "New docs"}},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	threads, err := f.Notifications.MarkRepoNotifications(context.Background(), "core", "go-forge", &NotificationRepoMarkOptions{
		All:         true,
		StatusTypes: []string{"unread", "pinned"},
		ToStatus:    "read",
		LastReadAt:  &lastReadAt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(threads) != 2 {
		t.Fatalf("got %d threads, want 2", len(threads))
	}
	if threads[0].ID != 7 || threads[1].ID != 8 {
		t.Fatalf("got ids=%d,%d want 7,8", threads[0].ID, threads[1].ID)
	}
	if threads[0].Subject.Title != "Pinned release" {
		t.Errorf("got title=%q, want %q", threads[0].Subject.Title, "Pinned release")
	}
}

func TestNotificationService_NotFound_Bad(t *testing.T) {
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
