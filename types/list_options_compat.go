// Compatibility types for RFC-style list options.

package types

import "time"

// ListIssueOption is a compatibility alias for repository issue list filters.
//
// Usage:
//
//	opts := ListIssueOption{State: "open", Sort: "created"}
type ListIssueOption struct {
	AssignedBy  string     `json:"assigned_by,omitempty"`
	Before      *time.Time `json:"before,omitempty"`
	CreatedBy   string     `json:"created_by,omitempty"`
	Labels      string     `json:"labels,omitempty"`
	Limit       int        `json:"limit,omitempty"`
	MentionedBy string     `json:"mentioned_by,omitempty"`
	Milestones  string     `json:"milestones,omitempty"`
	Page        int        `json:"page,omitempty"`
	PageSize    int        `json:"page_size,omitempty"`
	Query       string     `json:"q,omitempty"`
	Sort        string     `json:"sort,omitempty"`
	State       string     `json:"state,omitempty"`
	Since       *time.Time `json:"since,omitempty"`
	Type        string     `json:"type,omitempty"`
}

// ListPullRequestsOption is a compatibility alias for repository pull request list filters.
//
// Usage:
//
//	opts := ListPullRequestsOption{State: "open"}
type ListPullRequestsOption struct {
	Labels    []int64 `json:"labels,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Milestone int64   `json:"milestone,omitempty"`
	Page      int     `json:"page,omitempty"`
	PageSize  int     `json:"page_size,omitempty"`
	Poster    string  `json:"poster,omitempty"`
	Sort      string  `json:"sort,omitempty"`
	State     string  `json:"state,omitempty"`
}

// ListCommitsOption is a compatibility alias for repository commit list filters.
//
// Usage:
//
//	opts := ListCommitsOption{Sha: "main"}
type ListCommitsOption struct {
	Files        *bool  `json:"files,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Not          string `json:"not,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	Path         string `json:"path,omitempty"`
	Sha          string `json:"sha,omitempty"`
	Stat         *bool  `json:"stat,omitempty"`
	Verification *bool  `json:"verification,omitempty"`
}
