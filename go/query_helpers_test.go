package forge

import (
	"testing"
	"time"
)

// The variadic query helpers each have three branches: no filters yields nil,
// a present-but-empty filter yields nil, and a populated filter yields the map.
// These tests pin all three per helper.

func TestAdminUnadoptedQuery_AllBranches(t *testing.T) {
	if got := adminUnadoptedQuery(); got != nil {
		t.Errorf("no filters: got %v, want nil", got)
	}
	if got := adminUnadoptedQuery(AdminUnadoptedListOptions{}); got != nil {
		t.Errorf("empty filter: got %v, want nil", got)
	}
	if got := adminUnadoptedQuery(AdminUnadoptedListOptions{Pattern: "core/*"}); got["pattern"] != "core/*" {
		t.Errorf("populated: got %v", got)
	}
}

func TestRepoCommentQuery_AllBranches(t *testing.T) {
	if got := repoCommentQuery(); got != nil {
		t.Errorf("no filters: got %v, want nil", got)
	}
	if got := repoCommentQuery(RepoCommentListOptions{}); got != nil {
		t.Errorf("empty filter: got %v, want nil", got)
	}
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	if got := repoCommentQuery(RepoCommentListOptions{Since: &since}); got["since"] != since.Format(time.RFC3339) {
		t.Errorf("populated: got %v", got)
	}
}

func TestMilestoneQuery_AllBranches(t *testing.T) {
	if got := milestoneQuery(); got != nil {
		t.Errorf("no filters: got %v, want nil", got)
	}
	if got := milestoneQuery(MilestoneListOptions{}); got != nil {
		t.Errorf("empty filter: got %v, want nil", got)
	}
	if got := milestoneQuery(MilestoneListOptions{State: "open"}); got["state"] != "open" {
		t.Errorf("populated: got %v", got)
	}
}

func TestOrgActivityFeedQuery_AllBranches(t *testing.T) {
	if got := orgActivityFeedQuery(); got != nil {
		t.Errorf("no filters: got %v, want nil", got)
	}
	if got := orgActivityFeedQuery(OrgActivityFeedListOptions{}); got != nil {
		t.Errorf("empty filter: got %v, want nil", got)
	}
	day := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	if got := orgActivityFeedQuery(OrgActivityFeedListOptions{Date: &day}); got["date"] != "2026-05-31" {
		t.Errorf("populated: got %v", got)
	}
}

func TestActivityFeedQuery_AllBranches(t *testing.T) {
	if got := activityFeedQuery(); got != nil {
		t.Errorf("no filters: got %v, want nil", got)
	}
	if got := activityFeedQuery(ActivityFeedListOptions{}); got != nil {
		t.Errorf("empty filter: got %v, want nil", got)
	}
	day := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	if got := activityFeedQuery(ActivityFeedListOptions{Date: &day}); got["date"] != "2026-05-01" {
		t.Errorf("populated: got %v", got)
	}
}
