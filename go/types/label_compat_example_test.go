package types

func ExampleCreateIssueOption_UnmarshalJSON() {
	_ = (*CreateIssueOption).UnmarshalJSON
}

func ExampleCreatePullRequestOption_UnmarshalJSON() {
	_ = (*CreatePullRequestOption).UnmarshalJSON
}

func ExampleEditPullRequestOption_UnmarshalJSON() {
	_ = (*EditPullRequestOption).UnmarshalJSON
}
