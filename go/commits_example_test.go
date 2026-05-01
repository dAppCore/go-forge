package forge

func ExampleCommitListOptions_String() {
	_ = (*CommitListOptions).String
}

func ExampleCommitListOptions_GoString() {
	_ = (*CommitListOptions).GoString
}

func ExampleCommitService_List() {
	_ = (*CommitService).List
}

func ExampleCommitService_ListAll() {
	_ = (*CommitService).ListAll
}

func ExampleCommitService_Iter() {
	_ = (*CommitService).Iter
}

func ExampleCommitService_Get() {
	_ = (*CommitService).Get
}

func ExampleCommitService_ListCommitsPage() {
	_ = (*CommitService).ListCommitsPage
}

func ExampleCommitService_ListCommits() {
	_ = (*CommitService).ListCommits
}

func ExampleCommitService_IterCommits() {
	_ = (*CommitService).IterCommits
}

func ExampleCommitService_GetCommit() {
	_ = (*CommitService).GetCommit
}

func ExampleCommitService_GetDiffOrPatch() {
	_ = (*CommitService).GetDiffOrPatch
}

func ExampleCommitService_GetPullRequest() {
	_ = (*CommitService).GetPullRequest
}

func ExampleCommitService_GetCombinedStatus() {
	_ = (*CommitService).GetCombinedStatus
}

func ExampleCommitService_GetCombinedStatusByRef() {
	_ = (*CommitService).GetCombinedStatusByRef
}

func ExampleCommitService_ListStatuses() {
	_ = (*CommitService).ListStatuses
}

func ExampleCommitService_IterStatuses() {
	_ = (*CommitService).IterStatuses
}

func ExampleCommitService_CreateStatus() {
	_ = (*CommitService).CreateStatus
}

func ExampleCommitService_GetNote() {
	_ = (*CommitService).GetNote
}

func ExampleCommitService_SetNote() {
	_ = (*CommitService).SetNote
}

func ExampleCommitService_DeleteNote() {
	_ = (*CommitService).DeleteNote
}
