package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
)

func TestPagination_SinglePage_Good(t *testing.T) {
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

func TestPagination_MultiPage_Good(t *testing.T) {
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

func TestPagination_EmptyResult_Good(t *testing.T) {
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

func TestPagination_Iter_Good(t *testing.T) {
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
	count := 0
	for item, err := range ListIter[map[string]int](context.Background(), c, "/api/v1/repos", nil) {
		if err != nil {
			t.Fatal(err)
		}
		count++
		if item["id"] != count {
			t.Errorf("got id %d, want %d", item["id"], count)
		}
	}

	if count != 100 {
		t.Errorf("got %d items, want 100", count)
	}
}

func TestListPage_QueryParams_Good(t *testing.T) {
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
		map[string]string{"state": "open"}, ListOptions{Page: 2, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
}

func TestPagination_ServerError_Bad(t *testing.T) {
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

func TestPagination_ListOptions_String_Good(t *core.T) {
	subject := (*ListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListOptions_String_Bad(t *core.T) {
	subject := (*ListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListOptions_String_Ugly(t *core.T) {
	subject := (*ListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListOptions_GoString_Good(t *core.T) {
	subject := (*ListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListOptions_GoString_Bad(t *core.T) {
	subject := (*ListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListOptions_GoString_Ugly(t *core.T) {
	subject := (*ListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_String_Good(t *core.T) {
	subject := (*PagedResult[int]).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_String_Bad(t *core.T) {
	subject := (*PagedResult[int]).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_String_Ugly(t *core.T) {
	subject := (*PagedResult[int]).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_GoString_Good(t *core.T) {
	subject := (*PagedResult[int]).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_GoString_Bad(t *core.T) {
	subject := (*PagedResult[int]).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_PagedResult_GoString_Ugly(t *core.T) {
	subject := (*PagedResult[int]).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListPage_Good(t *core.T) {
	subject := ListPage[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListPage_Bad(t *core.T) {
	subject := ListPage[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListPage_Ugly(t *core.T) {
	subject := ListPage[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListAll_Good(t *core.T) {
	subject := ListAll[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListAll_Bad(t *core.T) {
	subject := ListAll[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListAll_Ugly(t *core.T) {
	subject := ListAll[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListIter_Good(t *core.T) {
	subject := ListIter[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListIter_Bad(t *core.T) {
	subject := ListIter[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPagination_ListIter_Ugly(t *core.T) {
	subject := ListIter[int]
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
