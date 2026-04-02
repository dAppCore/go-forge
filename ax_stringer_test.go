package forge

import (
	"fmt"
	"testing"
	"time"
)

func TestParams_String_Good(t *testing.T) {
	params := Params{"repo": "go-forge", "owner": "core"}
	want := `forge.Params{owner="core", repo="go-forge"}`
	if got := params.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(params); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", params); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestParams_String_NilSafe(t *testing.T) {
	var params Params
	want := "forge.Params{<nil>}"
	if got := params.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(params); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", params); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestListOptions_String_Good(t *testing.T) {
	opts := ListOptions{Page: 2, Limit: 25}
	want := "forge.ListOptions{page=2, limit=25}"
	if got := opts.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(opts); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", opts); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestRateLimit_String_Good(t *testing.T) {
	rl := RateLimit{Limit: 80, Remaining: 79, Reset: 1700000003}
	want := "forge.RateLimit{limit=80, remaining=79, reset=1700000003}"
	if got := rl.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(rl); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", rl); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestPagedResult_String_Good(t *testing.T) {
	page := PagedResult[int]{
		Items:      []int{1, 2, 3},
		TotalCount: 10,
		Page:       2,
		HasMore:    true,
	}
	want := "forge.PagedResult{items=3, totalCount=10, page=2, hasMore=true}"
	if got := page.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(page); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", page); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestOption_Stringers_Good(t *testing.T) {
	when := time.Date(2026, time.April, 2, 8, 3, 4, 0, time.UTC)

	cases := []struct {
		name string
		got  fmt.Stringer
		want string
	}{
		{
			name: "AdminActionsRunListOptions",
			got:  AdminActionsRunListOptions{Event: "push", Status: "success"},
			want: `forge.AdminActionsRunListOptions{event="push", status="success"}`,
		},
		{
			name: "AttachmentUploadOptions",
			got:  AttachmentUploadOptions{Name: "screenshot.png", UpdatedAt: &when},
			want: `forge.AttachmentUploadOptions{name="screenshot.png", updated_at="2026-04-02T08:03:04Z"}`,
		},
		{
			name: "NotificationListOptions",
			got:  NotificationListOptions{All: true, StatusTypes: []string{"unread"}, SubjectTypes: []string{"issue"}},
			want: `forge.NotificationListOptions{all=true, status_types=[]string{"unread"}, subject_types=[]string{"issue"}}`,
		},
		{
			name: "SearchIssuesOptions",
			got:  SearchIssuesOptions{State: "open", PriorityRepoID: 99, Assigned: true, Query: "build"},
			want: `forge.SearchIssuesOptions{state="open", q="build", priority_repo_id=99, assigned=true}`,
		},
		{
			name: "ReleaseAttachmentUploadOptions",
			got:  ReleaseAttachmentUploadOptions{Name: "release.zip"},
			want: `forge.ReleaseAttachmentUploadOptions{name="release.zip"}`,
		},
		{
			name: "UserSearchOptions",
			got:  UserSearchOptions{UID: 1001},
			want: `forge.UserSearchOptions{uid=1001}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.got.String(); got != tc.want {
				t.Fatalf("got String()=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprint(tc.got); got != tc.want {
				t.Fatalf("got fmt.Sprint=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprintf("%#v", tc.got); got != tc.want {
				t.Fatalf("got GoString=%q, want %q", got, tc.want)
			}
		})
	}
}

func TestOption_Stringers_Empty(t *testing.T) {
	cases := []struct {
		name string
		got  fmt.Stringer
		want string
	}{
		{
			name: "AdminUnadoptedListOptions",
			got:  AdminUnadoptedListOptions{},
			want: `forge.AdminUnadoptedListOptions{}`,
		},
		{
			name: "MilestoneListOptions",
			got:  MilestoneListOptions{},
			want: `forge.MilestoneListOptions{}`,
		},
		{
			name: "UserKeyListOptions",
			got:  UserKeyListOptions{},
			want: `forge.UserKeyListOptions{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.got.String(); got != tc.want {
				t.Fatalf("got String()=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprint(tc.got); got != tc.want {
				t.Fatalf("got fmt.Sprint=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprintf("%#v", tc.got); got != tc.want {
				t.Fatalf("got GoString=%q, want %q", got, tc.want)
			}
		})
	}
}

func TestService_Stringers_Good(t *testing.T) {
	client := NewClient("https://forge.example", "token")

	cases := []struct {
		name string
		got  fmt.Stringer
		want string
	}{
		{
			name: "RepoService",
			got:  newRepoService(client),
			want: `forge.RepoService{resource=forge.Resource{path="/api/v1/repos/{owner}/{repo}", collection="/api/v1/repos/{owner}"}}`,
		},
		{
			name: "AdminService",
			got:  newAdminService(client),
			want: `forge.AdminService{client=forge.Client{baseURL="https://forge.example", token=set, userAgent="go-forge/0.1"}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.got.String(); got != tc.want {
				t.Fatalf("got String()=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprint(tc.got); got != tc.want {
				t.Fatalf("got fmt.Sprint=%q, want %q", got, tc.want)
			}
			if got := fmt.Sprintf("%#v", tc.got); got != tc.want {
				t.Fatalf("got GoString=%q, want %q", got, tc.want)
			}
		})
	}
}

func TestService_Stringers_NilSafe(t *testing.T) {
	var repo *RepoService
	if got, want := repo.String(), "forge.RepoService{<nil>}"; got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got, want := fmt.Sprint(repo), "forge.RepoService{<nil>}"; got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got, want := fmt.Sprintf("%#v", repo), "forge.RepoService{<nil>}"; got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}

	var admin *AdminService
	if got, want := admin.String(), "forge.AdminService{<nil>}"; got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got, want := fmt.Sprint(admin), "forge.AdminService{<nil>}"; got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got, want := fmt.Sprintf("%#v", admin), "forge.AdminService{<nil>}"; got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}
