package forge

// String returns a safe summary of the actions service.
func (s *ActionsService) String() string {
	if s == nil {
		return "forge.ActionsService{<nil>}"
	}
	return serviceString("forge.ActionsService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the actions service.
func (s *ActionsService) GoString() string { return s.String() }

// String returns a safe summary of the ActivityPub service.
func (s *ActivityPubService) String() string {
	if s == nil {
		return "forge.ActivityPubService{<nil>}"
	}
	return serviceString("forge.ActivityPubService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the ActivityPub service.
func (s *ActivityPubService) GoString() string { return s.String() }

// String returns a safe summary of the admin service.
func (s *AdminService) String() string {
	if s == nil {
		return "forge.AdminService{<nil>}"
	}
	return serviceString("forge.AdminService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the admin service.
func (s *AdminService) GoString() string { return s.String() }

// String returns a safe summary of the branch service.
func (s *BranchService) String() string {
	if s == nil {
		return "forge.BranchService{<nil>}"
	}
	return serviceString("forge.BranchService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the branch service.
func (s *BranchService) GoString() string { return s.String() }

// String returns a safe summary of the commit service.
func (s *CommitService) String() string {
	if s == nil {
		return "forge.CommitService{<nil>}"
	}
	return serviceString("forge.CommitService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the commit service.
func (s *CommitService) GoString() string { return s.String() }

// String returns a safe summary of the content service.
func (s *ContentService) String() string {
	if s == nil {
		return "forge.ContentService{<nil>}"
	}
	return serviceString("forge.ContentService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the content service.
func (s *ContentService) GoString() string { return s.String() }

// String returns a safe summary of the issue service.
func (s *IssueService) String() string {
	if s == nil {
		return "forge.IssueService{<nil>}"
	}
	return serviceString("forge.IssueService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the issue service.
func (s *IssueService) GoString() string { return s.String() }

// String returns a safe summary of the label service.
func (s *LabelService) String() string {
	if s == nil {
		return "forge.LabelService{<nil>}"
	}
	return serviceString("forge.LabelService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the label service.
func (s *LabelService) GoString() string { return s.String() }

// String returns a safe summary of the milestone service.
func (s *MilestoneService) String() string {
	if s == nil {
		return "forge.MilestoneService{<nil>}"
	}
	return serviceString("forge.MilestoneService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the milestone service.
func (s *MilestoneService) GoString() string { return s.String() }

// String returns a safe summary of the misc service.
func (s *MiscService) String() string {
	if s == nil {
		return "forge.MiscService{<nil>}"
	}
	return serviceString("forge.MiscService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the misc service.
func (s *MiscService) GoString() string { return s.String() }

// String returns a safe summary of the notification service.
func (s *NotificationService) String() string {
	if s == nil {
		return "forge.NotificationService{<nil>}"
	}
	return serviceString("forge.NotificationService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the notification service.
func (s *NotificationService) GoString() string { return s.String() }

// String returns a safe summary of the organisation service.
func (s *OrgService) String() string {
	if s == nil {
		return "forge.OrgService{<nil>}"
	}
	return serviceString("forge.OrgService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the organisation service.
func (s *OrgService) GoString() string { return s.String() }

// String returns a safe summary of the package service.
func (s *PackageService) String() string {
	if s == nil {
		return "forge.PackageService{<nil>}"
	}
	return serviceString("forge.PackageService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the package service.
func (s *PackageService) GoString() string { return s.String() }

// String returns a safe summary of the pull request service.
func (s *PullService) String() string {
	if s == nil {
		return "forge.PullService{<nil>}"
	}
	return serviceString("forge.PullService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the pull request service.
func (s *PullService) GoString() string { return s.String() }

// String returns a safe summary of the release service.
func (s *ReleaseService) String() string {
	if s == nil {
		return "forge.ReleaseService{<nil>}"
	}
	return serviceString("forge.ReleaseService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the release service.
func (s *ReleaseService) GoString() string { return s.String() }

// String returns a safe summary of the repository service.
func (s *RepoService) String() string {
	if s == nil {
		return "forge.RepoService{<nil>}"
	}
	return serviceString("forge.RepoService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the repository service.
func (s *RepoService) GoString() string { return s.String() }

// String returns a safe summary of the team service.
func (s *TeamService) String() string {
	if s == nil {
		return "forge.TeamService{<nil>}"
	}
	return serviceString("forge.TeamService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the team service.
func (s *TeamService) GoString() string { return s.String() }

// String returns a safe summary of the user service.
func (s *UserService) String() string {
	if s == nil {
		return "forge.UserService{<nil>}"
	}
	return serviceString("forge.UserService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the user service.
func (s *UserService) GoString() string { return s.String() }

// String returns a safe summary of the webhook service.
func (s *WebhookService) String() string {
	if s == nil {
		return "forge.WebhookService{<nil>}"
	}
	return serviceString("forge.WebhookService", "resource", &s.Resource)
}

// GoString returns a safe Go-syntax summary of the webhook service.
func (s *WebhookService) GoString() string { return s.String() }

// String returns a safe summary of the wiki service.
func (s *WikiService) String() string {
	if s == nil {
		return "forge.WikiService{<nil>}"
	}
	return serviceString("forge.WikiService", "client", s.client)
}

// GoString returns a safe Go-syntax summary of the wiki service.
func (s *WikiService) GoString() string { return s.String() }
