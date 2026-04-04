package forge

// String returns a safe summary of the actions service.
//
// Usage:
//
//	s := &forge.ActionsService{}
//	_ = s.String()
func (s *ActionsService) String() string {
	if s == nil {
		return "forge.ActionsService{<nil>}"
	}
	return serviceString("forge.ActionsService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the actions service.
//
// Usage:
//
//	s := &forge.ActionsService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *ActionsService) GoString() string { return s.String() }

// String returns a safe summary of the ActivityPub service.
//
// Usage:
//
//	s := &forge.ActivityPubService{}
//	_ = s.String()
func (s *ActivityPubService) String() string {
	if s == nil {
		return "forge.ActivityPubService{<nil>}"
	}
	return serviceString("forge.ActivityPubService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the ActivityPub service.
//
// Usage:
//
//	s := &forge.ActivityPubService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *ActivityPubService) GoString() string { return s.String() }

// String returns a safe summary of the admin service.
//
// Usage:
//
//	s := &forge.AdminService{}
//	_ = s.String()
func (s *AdminService) String() string {
	if s == nil {
		return "forge.AdminService{<nil>}"
	}
	return serviceString("forge.AdminService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the admin service.
//
// Usage:
//
//	s := &forge.AdminService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *AdminService) GoString() string { return s.String() }

// String returns a safe summary of the branch service.
//
// Usage:
//
//	s := &forge.BranchService{}
//	_ = s.String()
func (s *BranchService) String() string {
	if s == nil {
		return "forge.BranchService{<nil>}"
	}
	return serviceString("forge.BranchService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the branch service.
//
// Usage:
//
//	s := &forge.BranchService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *BranchService) GoString() string { return s.String() }

// String returns a safe summary of the commit service.
//
// Usage:
//
//	s := &forge.CommitService{}
//	_ = s.String()
func (s *CommitService) String() string {
	if s == nil {
		return "forge.CommitService{<nil>}"
	}
	return serviceString("forge.CommitService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the commit service.
//
// Usage:
//
//	s := &forge.CommitService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *CommitService) GoString() string { return s.String() }

// String returns a safe summary of the content service.
//
// Usage:
//
//	s := &forge.ContentService{}
//	_ = s.String()
func (s *ContentService) String() string {
	if s == nil {
		return "forge.ContentService{<nil>}"
	}
	return serviceString("forge.ContentService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the content service.
//
// Usage:
//
//	s := &forge.ContentService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *ContentService) GoString() string { return s.String() }

// String returns a safe summary of the issue service.
//
// Usage:
//
//	s := &forge.IssueService{}
//	_ = s.String()
func (s *IssueService) String() string {
	if s == nil {
		return "forge.IssueService{<nil>}"
	}
	return serviceString("forge.IssueService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the issue service.
//
// Usage:
//
//	s := &forge.IssueService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *IssueService) GoString() string { return s.String() }

// String returns a safe summary of the label service.
//
// Usage:
//
//	s := &forge.LabelService{}
//	_ = s.String()
func (s *LabelService) String() string {
	if s == nil {
		return "forge.LabelService{<nil>}"
	}
	return serviceString("forge.LabelService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the label service.
//
// Usage:
//
//	s := &forge.LabelService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *LabelService) GoString() string { return s.String() }

// String returns a safe summary of the milestone service.
//
// Usage:
//
//	s := &forge.MilestoneService{}
//	_ = s.String()
func (s *MilestoneService) String() string {
	if s == nil {
		return "forge.MilestoneService{<nil>}"
	}
	return serviceString("forge.MilestoneService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the milestone service.
//
// Usage:
//
//	s := &forge.MilestoneService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *MilestoneService) GoString() string { return s.String() }

// String returns a safe summary of the misc service.
//
// Usage:
//
//	s := &forge.MiscService{}
//	_ = s.String()
func (s *MiscService) String() string {
	if s == nil {
		return "forge.MiscService{<nil>}"
	}
	return serviceString("forge.MiscService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the misc service.
//
// Usage:
//
//	s := &forge.MiscService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *MiscService) GoString() string { return s.String() }

// String returns a safe summary of the notification service.
//
// Usage:
//
//	s := &forge.NotificationService{}
//	_ = s.String()
func (s *NotificationService) String() string {
	if s == nil {
		return "forge.NotificationService{<nil>}"
	}
	return serviceString("forge.NotificationService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the notification service.
//
// Usage:
//
//	s := &forge.NotificationService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *NotificationService) GoString() string { return s.String() }

// String returns a safe summary of the organisation service.
//
// Usage:
//
//	s := &forge.OrgService{}
//	_ = s.String()
func (s *OrgService) String() string {
	if s == nil {
		return "forge.OrgService{<nil>}"
	}
	return serviceString("forge.OrgService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the organisation service.
//
// Usage:
//
//	s := &forge.OrgService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *OrgService) GoString() string { return s.String() }

// String returns a safe summary of the package service.
//
// Usage:
//
//	s := &forge.PackageService{}
//	_ = s.String()
func (s *PackageService) String() string {
	if s == nil {
		return "forge.PackageService{<nil>}"
	}
	return serviceString("forge.PackageService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the package service.
//
// Usage:
//
//	s := &forge.PackageService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *PackageService) GoString() string { return s.String() }

// String returns a safe summary of the pull request service.
//
// Usage:
//
//	s := &forge.PullService{}
//	_ = s.String()
func (s *PullService) String() string {
	if s == nil {
		return "forge.PullService{<nil>}"
	}
	return serviceString("forge.PullService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the pull request service.
//
// Usage:
//
//	s := &forge.PullService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *PullService) GoString() string { return s.String() }

// String returns a safe summary of the release service.
//
// Usage:
//
//	s := &forge.ReleaseService{}
//	_ = s.String()
func (s *ReleaseService) String() string {
	if s == nil {
		return "forge.ReleaseService{<nil>}"
	}
	return serviceString("forge.ReleaseService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the release service.
//
// Usage:
//
//	s := &forge.ReleaseService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *ReleaseService) GoString() string { return s.String() }

// String returns a safe summary of the repository service.
//
// Usage:
//
//	s := &forge.RepoService{}
//	_ = s.String()
func (s *RepoService) String() string {
	if s == nil {
		return "forge.RepoService{<nil>}"
	}
	return serviceString("forge.RepoService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the repository service.
//
// Usage:
//
//	s := &forge.RepoService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *RepoService) GoString() string { return s.String() }

// String returns a safe summary of the team service.
//
// Usage:
//
//	s := &forge.TeamService{}
//	_ = s.String()
func (s *TeamService) String() string {
	if s == nil {
		return "forge.TeamService{<nil>}"
	}
	return serviceString("forge.TeamService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the team service.
//
// Usage:
//
//	s := &forge.TeamService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *TeamService) GoString() string { return s.String() }

// String returns a safe summary of the user service.
//
// Usage:
//
//	s := &forge.UserService{}
//	_ = s.String()
func (s *UserService) String() string {
	if s == nil {
		return "forge.UserService{<nil>}"
	}
	return serviceString("forge.UserService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the user service.
//
// Usage:
//
//	s := &forge.UserService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *UserService) GoString() string { return s.String() }

// String returns a safe summary of the webhook service.
//
// Usage:
//
//	s := &forge.WebhookService{}
//	_ = s.String()
func (s *WebhookService) String() string {
	if s == nil {
		return "forge.WebhookService{<nil>}"
	}
	return serviceString("forge.WebhookService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the webhook service.
//
// Usage:
//
//	s := &forge.WebhookService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *WebhookService) GoString() string { return s.String() }

// String returns a safe summary of the wiki service.
//
// Usage:
//
//	s := &forge.WikiService{}
//	_ = s.String()
func (s *WikiService) String() string {
	if s == nil {
		return "forge.WikiService{<nil>}"
	}
	return serviceString("forge.WikiService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the wiki service.
//
// Usage:
//
//	s := &forge.WikiService{}
//	_ = fmt.Sprintf("%#v", s)
func (s *WikiService) GoString() string { return s.String() }
