package forge

import (
	"testing"
	"time"
)

// The queryParams() builders are pure option-to-map translators. Each one
// follows the same contract: populate the map with only the non-zero fields,
// and return nil when nothing was set. These tests pin that contract per
// builder with Good (all fields populated), Bad (zero value yields nil), and
// Ugly (partial / edge-value population) shapes.

func wantQuery(t *testing.T, name string, got map[string]string, want map[string]string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: got %d keys %v, want %d keys %v", name, len(got), got, len(want), want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s: key %q = %q, want %q", name, k, got[k], v)
		}
	}
}

func TestQueryParams_AdminActionsRun_Good(t *testing.T) {
	o := AdminActionsRunListOptions{
		Event:   "push",
		Branch:  "dev",
		Status:  "success",
		Actor:   "cladius",
		HeadSHA: "abc123",
	}
	wantQuery(t, "AdminActionsRunListOptions", o.queryParams(), map[string]string{
		"event":    "push",
		"branch":   "dev",
		"status":   "success",
		"actor":    "cladius",
		"head_sha": "abc123",
	})
}

func TestQueryParams_AdminActionsRun_Bad(t *testing.T) {
	if got := (AdminActionsRunListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_AdminActionsRun_Ugly(t *testing.T) {
	// Only a subset of fields set — the rest must be omitted, not blank.
	o := AdminActionsRunListOptions{Event: "pull_request", Actor: "virgil"}
	wantQuery(t, "AdminActionsRunListOptions partial", o.queryParams(), map[string]string{
		"event": "pull_request",
		"actor": "virgil",
	})
}

func TestQueryParams_AdminUnadopted_Good(t *testing.T) {
	wantQuery(t, "AdminUnadoptedListOptions", AdminUnadoptedListOptions{Pattern: "core/*"}.queryParams(), map[string]string{
		"pattern": "core/*",
	})
}

func TestQueryParams_AdminUnadopted_Bad(t *testing.T) {
	if got := (AdminUnadoptedListOptions{}).queryParams(); got != nil {
		t.Errorf("empty pattern should yield nil, got %v", got)
	}
}

func TestQueryParams_Commit_Good(t *testing.T) {
	yes := true
	no := false
	o := CommitListOptions{Sha: "deadbeef", Path: "go/forge.go", Stat: &yes, Verification: &no, Files: &yes, Not: "main"}
	wantQuery(t, "CommitListOptions", o.queryParams(), map[string]string{
		"sha":          "deadbeef",
		"path":         "go/forge.go",
		"stat":         "true",
		"verification": "false",
		"files":        "true",
		"not":          "main",
	})
}

func TestQueryParams_Commit_Bad(t *testing.T) {
	if got := (CommitListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_Commit_Ugly(t *testing.T) {
	// A pointer-bool explicitly set to false must still appear in the query —
	// the distinction between "unset" (nil) and "false" is load-bearing.
	no := false
	o := CommitListOptions{Stat: &no}
	wantQuery(t, "CommitListOptions false-pointer", o.queryParams(), map[string]string{
		"stat": "false",
	})
}

func TestQueryParams_Milestone_Good(t *testing.T) {
	wantQuery(t, "MilestoneListOptions", MilestoneListOptions{State: "open", Name: "beta"}.queryParams(), map[string]string{
		"state": "open",
		"name":  "beta",
	})
}

func TestQueryParams_Milestone_Bad(t *testing.T) {
	if got := (MilestoneListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_IssueList_Good(t *testing.T) {
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	before := time.Date(2026, 6, 7, 8, 9, 10, 0, time.UTC)
	o := IssueListOptions{
		State:       "open",
		Sort:        "created",
		Labels:      "bug,help",
		Query:       "panic",
		Type:        "issues",
		Milestones:  "beta",
		Since:       &since,
		Before:      &before,
		CreatedBy:   "snider",
		AssignedBy:  "cladius",
		MentionedBy: "virgil",
	}
	wantQuery(t, "IssueListOptions", o.queryParams(), map[string]string{
		"state":        "open",
		"sort":         "created",
		"labels":       "bug,help",
		"q":            "panic",
		"type":         "issues",
		"milestones":   "beta",
		"since":        since.Format(time.RFC3339),
		"before":       before.Format(time.RFC3339),
		"created_by":   "snider",
		"assigned_by":  "cladius",
		"mentioned_by": "virgil",
	})
}

func TestQueryParams_IssueList_Bad(t *testing.T) {
	if got := (IssueListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_IssueList_Ugly(t *testing.T) {
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	o := IssueListOptions{State: "closed", Since: &since}
	wantQuery(t, "IssueListOptions partial", o.queryParams(), map[string]string{
		"state": "closed",
		"since": since.Format(time.RFC3339),
	})
}

func TestQueryParams_RepoComment_Good(t *testing.T) {
	since := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	before := time.Date(2026, 7, 8, 9, 10, 11, 0, time.UTC)
	o := RepoCommentListOptions{Since: &since, Before: &before}
	wantQuery(t, "RepoCommentListOptions", o.queryParams(), map[string]string{
		"since":  since.Format(time.RFC3339),
		"before": before.Format(time.RFC3339),
	})
}

func TestQueryParams_RepoComment_Bad(t *testing.T) {
	if got := (RepoCommentListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_SearchIssues_Good(t *testing.T) {
	since := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	before := time.Date(2026, 8, 9, 10, 11, 12, 0, time.UTC)
	o := SearchIssuesOptions{
		State:           "open",
		Labels:          "bug",
		Milestones:      "beta",
		Query:           "leak",
		PriorityRepoID:  42,
		Type:            "pulls",
		Since:           &since,
		Before:          &before,
		Assigned:        true,
		Created:         true,
		Mentioned:       true,
		ReviewRequested: true,
		Reviewed:        true,
		Owner:           "core",
		Team:            "agents",
	}
	wantQuery(t, "SearchIssuesOptions", o.queryParams(), map[string]string{
		"state":            "open",
		"labels":           "bug",
		"milestones":       "beta",
		"q":                "leak",
		"priority_repo_id": int64String(42),
		"type":             "pulls",
		"since":            since.Format(time.RFC3339),
		"before":           before.Format(time.RFC3339),
		"assigned":         "true",
		"created":          "true",
		"mentioned":        "true",
		"review_requested": "true",
		"reviewed":         "true",
		"owner":            "core",
		"team":             "agents",
	})
}

func TestQueryParams_SearchIssues_Bad(t *testing.T) {
	if got := (SearchIssuesOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_SearchIssues_Ugly(t *testing.T) {
	// Boolean flags left false must be omitted entirely, not sent as "false".
	o := SearchIssuesOptions{Query: "panic"}
	wantQuery(t, "SearchIssuesOptions bools-false", o.queryParams(), map[string]string{
		"q": "panic",
	})
}

func TestQueryParams_OrgActivityFeed_Good(t *testing.T) {
	day := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	o := OrgActivityFeedListOptions{Date: &day}
	wantQuery(t, "OrgActivityFeedListOptions", o.queryParams(), map[string]string{
		"date": "2026-05-31",
	})
}

func TestQueryParams_OrgActivityFeed_Bad(t *testing.T) {
	if got := (OrgActivityFeedListOptions{}).queryParams(); got != nil {
		t.Errorf("nil date should yield nil, got %v", got)
	}
}

func TestQueryParams_RepoKey_Good(t *testing.T) {
	o := RepoKeyListOptions{KeyID: 7, Fingerprint: "AB:CD"}
	wantQuery(t, "RepoKeyListOptions", o.queryParams(), map[string]string{
		"key_id":      int64String(7),
		"fingerprint": "AB:CD",
	})
}

func TestQueryParams_RepoKey_Bad(t *testing.T) {
	if got := (RepoKeyListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_RepoKey_Ugly(t *testing.T) {
	// KeyID == 0 is treated as unset; only fingerprint should survive.
	o := RepoKeyListOptions{Fingerprint: "EF:01"}
	wantQuery(t, "RepoKeyListOptions zero-id", o.queryParams(), map[string]string{
		"fingerprint": "EF:01",
	})
}

func TestQueryParams_ActivityFeed_Good(t *testing.T) {
	day := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	o := ActivityFeedListOptions{Date: &day}
	wantQuery(t, "ActivityFeedListOptions", o.queryParams(), map[string]string{
		"date": "2026-05-01",
	})
}

func TestQueryParams_ActivityFeed_Bad(t *testing.T) {
	if got := (ActivityFeedListOptions{}).queryParams(); got != nil {
		t.Errorf("nil date should yield nil, got %v", got)
	}
}

func TestQueryParams_RepoTime_Good(t *testing.T) {
	since := time.Date(2026, 4, 5, 6, 7, 8, 0, time.UTC)
	before := time.Date(2026, 9, 10, 11, 12, 13, 0, time.UTC)
	o := RepoTimeListOptions{User: "alice", Since: &since, Before: &before}
	wantQuery(t, "RepoTimeListOptions", o.queryParams(), map[string]string{
		"user":   "alice",
		"since":  since.Format(time.RFC3339),
		"before": before.Format(time.RFC3339),
	})
}

func TestQueryParams_RepoTime_Bad(t *testing.T) {
	if got := (RepoTimeListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_Release_Good(t *testing.T) {
	o := ReleaseListOptions{Draft: true, PreRelease: true, Query: "v1"}
	wantQuery(t, "ReleaseListOptions", o.queryParams(), map[string]string{
		"draft":       "true",
		"pre-release": "true",
		"q":           "v1",
	})
}

func TestQueryParams_Release_Bad(t *testing.T) {
	if got := (ReleaseListOptions{}).queryParams(); got != nil {
		t.Errorf("empty options should yield nil, got %v", got)
	}
}

func TestQueryParams_Release_Ugly(t *testing.T) {
	// Draft false + PreRelease false must drop those keys; only q survives.
	o := ReleaseListOptions{Query: "rc"}
	wantQuery(t, "ReleaseListOptions bools-false", o.queryParams(), map[string]string{
		"q": "rc",
	})
}

func TestQueryParams_UserSearch_Good(t *testing.T) {
	o := UserSearchOptions{UID: 99}
	wantQuery(t, "UserSearchOptions", o.queryParams(), map[string]string{
		"uid": int64String(99),
	})
}

func TestQueryParams_UserSearch_Bad(t *testing.T) {
	if got := (UserSearchOptions{}).queryParams(); got != nil {
		t.Errorf("zero UID should yield nil, got %v", got)
	}
}

func TestQueryParams_UserKey_Good(t *testing.T) {
	o := UserKeyListOptions{Fingerprint: "12:34"}
	wantQuery(t, "UserKeyListOptions", o.queryParams(), map[string]string{
		"fingerprint": "12:34",
	})
}

func TestQueryParams_UserKey_Bad(t *testing.T) {
	if got := (UserKeyListOptions{}).queryParams(); got != nil {
		t.Errorf("empty fingerprint should yield nil, got %v", got)
	}
}
