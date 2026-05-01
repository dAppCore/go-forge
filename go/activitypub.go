package forge

import (
	"context"

	"dappco.re/go/forge/types"
)

// ActivityPubService handles ActivityPub actor and inbox endpoints.
//
// Usage:
//
//	f := forge.NewForge("https://forge.lthn.ai", "token")
//	_, err := f.ActivityPub.GetInstanceActor(ctx)
type ActivityPubService struct {
	client *Client
}

func newActivityPubService(c *Client) *ActivityPubService {
	return &ActivityPubService{client: c}
}

// GetInstanceActor returns the instance's ActivityPub actor.
func (s *ActivityPubService) GetInstanceActor(ctx context.Context) (*types.ActivityPub, error) {
	var out types.ActivityPub
	if err := s.client.Get(ctx, "/activitypub/actor", &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendInstanceActorInbox sends an ActivityPub object to the instance inbox.
func (s *ActivityPubService) SendInstanceActorInbox(ctx context.Context, body *types.ForgeLike) error {
	return s.client.Post(ctx, "/activitypub/actor/inbox", body, nil)
}

// GetRepositoryActor returns the ActivityPub actor for a repository.
func (s *ActivityPubService) GetRepositoryActor(ctx context.Context, repositoryID int64) (*types.ActivityPub, error) {
	path := ResolvePath("/activitypub/repository-id/{repository-id}", Params{"repository-id": int64String(repositoryID)})
	var out types.ActivityPub
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendRepositoryInbox sends an ActivityPub object to a repository inbox.
func (s *ActivityPubService) SendRepositoryInbox(ctx context.Context, repositoryID int64, body *types.ForgeLike) error {
	path := ResolvePath("/activitypub/repository-id/{repository-id}/inbox", Params{"repository-id": int64String(repositoryID)})
	return s.client.Post(ctx, path, body, nil)
}

// GetPersonActor returns the Person actor for a user.
func (s *ActivityPubService) GetPersonActor(ctx context.Context, userID int64) (*types.ActivityPub, error) {
	path := ResolvePath("/activitypub/user-id/{user-id}", Params{"user-id": int64String(userID)})
	var out types.ActivityPub
	if err := s.client.Get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// SendPersonInbox sends an ActivityPub object to a user's inbox.
func (s *ActivityPubService) SendPersonInbox(ctx context.Context, userID int64, body *types.ForgeLike) error {
	path := ResolvePath("/activitypub/user-id/{user-id}/inbox", Params{"user-id": int64String(userID)})
	return s.client.Post(ctx, path, body, nil)
}
