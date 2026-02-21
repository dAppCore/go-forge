package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPagination_Good_SinglePage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]map[string]int{{"id": 1}, {"id": 2}})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	result, err := ListAll[map[string]int](context.Background(), c, "/api/v1/repos", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("got %d items", len(result))
	}
}

func TestPagination_Good_MultiPage(t *testing.T) {
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page++
		w.Header().Set("X-Total-Count", "100")
		items := make([]map[string]int, 50)
		for i := range items {
			items[i] = map[string]int{"id": (page-1)*50 + i + 1}
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	result, err := ListAll[map[string]int](context.Background(), c, "/api/v1/repos", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 100 {
		t.Errorf("got %d items, want 100", len(result))
	}
}

func TestPagination_Good_EmptyResult(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "0")
		json.NewEncoder(w).Encode([]map[string]int{})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	result, err := ListAll[map[string]int](context.Background(), c, "/api/v1/repos", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("got %d items", len(result))
	}
}

func TestListPage_Good_QueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Query().Get("page")
		l := r.URL.Query().Get("limit")
		s := r.URL.Query().Get("state")
		if p != "2" || l != "25" || s != "open" {
			t.Errorf("wrong params: page=%s limit=%s state=%s", p, l, s)
		}
		w.Header().Set("X-Total-Count", "50")
		json.NewEncoder(w).Encode([]map[string]int{})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	_, err := ListPage[map[string]int](context.Background(), c, "/api/v1/repos",
		map[string]string{"state": "open"}, ListOptions{Page: 2, Limit: 25})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPagination_Bad_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]string{"message": "fail"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	_, err := ListAll[map[string]int](context.Background(), c, "/api/v1/repos", nil)
	if err == nil {
		t.Fatal("expected error")
	}
}
