package types_test

import (
	. "dappco.re/go"
	"dappco.re/go/forge/types"
)

func TestAX7_CreateFileOptions_MarshalJSON_Good(t *T) {
	got, err := types.CreateFileOptions{Content: "payload"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "content")
}

func TestAX7_CreateFileOptions_MarshalJSON_Bad(t *T) {
	got, err := types.CreateFileOptions{}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "content")
}

func TestAX7_CreateFileOptions_MarshalJSON_Ugly(t *T) {
	got, err := types.CreateFileOptions{Content: "payload", ContentBase64: "encoded"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "encoded")
}

func TestAX7_CreateFileOptions_UnmarshalJSON_Good(t *T) {
	var got types.CreateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content":"payload"}`))
	AssertNoError(t, err)
	AssertEqual(t, "payload", got.Content)
}

func TestAX7_CreateFileOptions_UnmarshalJSON_Bad(t *T) {
	var got types.CreateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content_base64"`))
	AssertError(t, err)
	AssertEqual(t, "", got.Content)
}

func TestAX7_CreateFileOptions_UnmarshalJSON_Ugly(t *T) {
	var got types.CreateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content":"plain"}`))
	AssertNoError(t, err)
	AssertEqual(t, "plain", got.Content)
}

func TestAX7_UpdateFileOptions_MarshalJSON_Good(t *T) {
	got, err := types.UpdateFileOptions{Content: "payload"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "content")
}

func TestAX7_UpdateFileOptions_MarshalJSON_Bad(t *T) {
	got, err := types.UpdateFileOptions{}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "content")
}

func TestAX7_UpdateFileOptions_MarshalJSON_Ugly(t *T) {
	got, err := types.UpdateFileOptions{Content: "payload", ContentBase64: "encoded"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "encoded")
}

func TestAX7_UpdateFileOptions_UnmarshalJSON_Good(t *T) {
	var got types.UpdateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content":"payload"}`))
	AssertNoError(t, err)
	AssertEqual(t, "payload", got.Content)
}

func TestAX7_UpdateFileOptions_UnmarshalJSON_Bad(t *T) {
	var got types.UpdateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content_base64"`))
	AssertError(t, err)
	AssertEqual(t, "", got.Content)
}

func TestAX7_UpdateFileOptions_UnmarshalJSON_Ugly(t *T) {
	var got types.UpdateFileOptions
	err := got.UnmarshalJSON([]byte(`{"content":"plain"}`))
	AssertNoError(t, err)
	AssertEqual(t, "plain", got.Content)
}

func TestAX7_ChangeFileOperation_MarshalJSON_Good(t *T) {
	got, err := types.ChangeFileOperation{Content: "payload"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "content")
}

func TestAX7_ChangeFileOperation_MarshalJSON_Bad(t *T) {
	got, err := types.ChangeFileOperation{}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "operation")
}

func TestAX7_ChangeFileOperation_MarshalJSON_Ugly(t *T) {
	got, err := types.ChangeFileOperation{Content: "payload", ContentBase64: "encoded"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "encoded")
}

func TestAX7_ChangeFileOperation_UnmarshalJSON_Good(t *T) {
	var got types.ChangeFileOperation
	err := got.UnmarshalJSON([]byte(`{"content":"payload"}`))
	AssertNoError(t, err)
	AssertEqual(t, "payload", got.Content)
}

func TestAX7_ChangeFileOperation_UnmarshalJSON_Bad(t *T) {
	var got types.ChangeFileOperation
	err := got.UnmarshalJSON([]byte(`{"content_base64"`))
	AssertError(t, err)
	AssertEqual(t, "", got.Content)
}

func TestAX7_ChangeFileOperation_UnmarshalJSON_Ugly(t *T) {
	var got types.ChangeFileOperation
	err := got.UnmarshalJSON([]byte(`{"content":"plain"}`))
	AssertNoError(t, err)
	AssertEqual(t, "plain", got.Content)
}

func TestAX7_MergePullRequestOption_MarshalJSON_Good(t *T) {
	got, err := types.MergePullRequestOption{MergeStyle: "merge"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "Do")
}

func TestAX7_MergePullRequestOption_MarshalJSON_Bad(t *T) {
	got, err := types.MergePullRequestOption{}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "Do")
}

func TestAX7_MergePullRequestOption_MarshalJSON_Ugly(t *T) {
	got, err := types.MergePullRequestOption{MergeStyle: "merge", Do: "rebase"}.MarshalJSON()
	AssertNoError(t, err)
	AssertContains(t, string(got), "rebase")
}

func TestAX7_MergePullRequestOption_UnmarshalJSON_Good(t *T) {
	var got types.MergePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"do":"merge"}`))
	AssertNoError(t, err)
	AssertEqual(t, "merge", got.MergeStyle)
}

func TestAX7_MergePullRequestOption_UnmarshalJSON_Bad(t *T) {
	var got types.MergePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"do"`))
	AssertError(t, err)
	AssertEqual(t, "", got.MergeStyle)
}

func TestAX7_MergePullRequestOption_UnmarshalJSON_Ugly(t *T) {
	var got types.MergePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"Do":"squash"}`))
	AssertNoError(t, err)
	AssertEqual(t, "squash", got.MergeStyle)
}

func TestAX7_CreateIssueOption_UnmarshalJSON_Good(t *T) {
	var got types.CreateIssueOption
	err := got.UnmarshalJSON([]byte(`{"labels":["bug","help"]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{"bug", "help"}, got.Labels)
}

func TestAX7_CreateIssueOption_UnmarshalJSON_Bad(t *T) {
	var got types.CreateIssueOption
	err := got.UnmarshalJSON([]byte(`{"labels"`))
	AssertError(t, err)
	AssertNil(t, got.Labels)
}

func TestAX7_CreateIssueOption_UnmarshalJSON_Ugly(t *T) {
	var got types.CreateIssueOption
	err := got.UnmarshalJSON([]byte(`{"labels":[]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{}, got.Labels)
}

func TestAX7_CreatePullRequestOption_UnmarshalJSON_Good(t *T) {
	var got types.CreatePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels":["bug","help"]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{"bug", "help"}, got.Labels)
}

func TestAX7_CreatePullRequestOption_UnmarshalJSON_Bad(t *T) {
	var got types.CreatePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels"`))
	AssertError(t, err)
	AssertNil(t, got.Labels)
}

func TestAX7_CreatePullRequestOption_UnmarshalJSON_Ugly(t *T) {
	var got types.CreatePullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels":[]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{}, got.Labels)
}

func TestAX7_EditPullRequestOption_UnmarshalJSON_Good(t *T) {
	var got types.EditPullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels":["bug","help"]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{"bug", "help"}, got.Labels)
}

func TestAX7_EditPullRequestOption_UnmarshalJSON_Bad(t *T) {
	var got types.EditPullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels"`))
	AssertError(t, err)
	AssertNil(t, got.Labels)
}

func TestAX7_EditPullRequestOption_UnmarshalJSON_Ugly(t *T) {
	var got types.EditPullRequestOption
	err := got.UnmarshalJSON([]byte(`{"labels":[]}`))
	AssertNoError(t, err)
	AssertEqual(t, []string{}, got.Labels)
}
