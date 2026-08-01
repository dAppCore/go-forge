package forge

func ExampleActivityPubService_GetInstanceActor() {
	_ = (*ActivityPubService).GetInstanceActor
}

func ExampleActivityPubService_SendInstanceActorInbox() {
	_ = (*ActivityPubService).SendInstanceActorInbox
}

func ExampleActivityPubService_GetRepositoryActor() {
	_ = (*ActivityPubService).GetRepositoryActor
}

func ExampleActivityPubService_SendRepositoryInbox() {
	_ = (*ActivityPubService).SendRepositoryInbox
}

func ExampleActivityPubService_GetPersonActor() {
	_ = (*ActivityPubService).GetPersonActor
}

func ExampleActivityPubService_SendPersonInbox() {
	_ = (*ActivityPubService).SendPersonInbox
}
