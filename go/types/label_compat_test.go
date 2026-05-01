package types

import (
	json "github.com/goccy/go-json"
	"testing"
)

func TestCreateIssueOption_Unmarshal_LabelNamesCompat_Good(t *testing.T) {
	var opts CreateIssueOption
	if err := json.Unmarshal([]byte(`{"title":"issue","labels":["enhancement","bug"]}`), &opts); err != nil {
		t.Fatal(err)
	}

	labels, ok := opts.Labels.([]string)
	if !ok {
		t.Fatalf("expected []string labels, got %T", opts.Labels)
	}
	if len(labels) != 2 || labels[0] != "enhancement" || labels[1] != "bug" {
		t.Fatalf("unexpected labels: %#v", labels)
	}
}

func TestCreatePullRequestOption_Unmarshal_LabelIDsCompat_Good(t *testing.T) {
	var opts CreatePullRequestOption
	if err := json.Unmarshal([]byte(`{"title":"pr","labels":[1,2]}`), &opts); err != nil {
		t.Fatal(err)
	}

	labels, ok := opts.Labels.([]int64)
	if !ok {
		t.Fatalf("expected []int64 labels, got %T", opts.Labels)
	}
	if len(labels) != 2 || labels[0] != 1 || labels[1] != 2 {
		t.Fatalf("unexpected labels: %#v", labels)
	}
}
