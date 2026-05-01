package types

import (
	json "github.com/goccy/go-json"
	"testing"
)

func TestMergePullRequestOption_MarshalJSON_CompatMergeStyle_Good(t *testing.T) {
	data, err := json.Marshal(MergePullRequestOption{
		MergeMessageField: "PR: Add feature",
		MergeStyle:        "squash",
	})
	if err != nil {
		t.Fatal(err)
	}

	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatal(err)
	}
	if got := body["Do"]; got != "squash" {
		t.Fatalf("got Do=%v, want squash", got)
	}
	if _, ok := body["MergeStyle"]; ok {
		t.Fatalf("did not expect MergeStyle in request body: %#v", body)
	}
	if got := body["MergeMessageField"]; got != "PR: Add feature" {
		t.Fatalf("got MergeMessageField=%v, want %q", got, "PR: Add feature")
	}
}

func TestMergePullRequestOption_UnmarshalJSON_CompatMergeStyle_Good(t *testing.T) {
	var opts MergePullRequestOption
	if err := json.Unmarshal([]byte(`{"Do":"rebase","MergeMessageField":"ready"}`), &opts); err != nil {
		t.Fatal(err)
	}
	if opts.Do != "rebase" {
		t.Fatalf("got Do=%q, want %q", opts.Do, "rebase")
	}
	if opts.MergeStyle != "rebase" {
		t.Fatalf("got MergeStyle=%q, want %q", opts.MergeStyle, "rebase")
	}
	if opts.MergeMessageField != "ready" {
		t.Fatalf("got MergeMessageField=%q, want %q", opts.MergeMessageField, "ready")
	}
}
