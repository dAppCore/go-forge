package types

import json "github.com/goccy/go-json"

type createIssueOptionCompat CreateIssueOption

func (o *CreateIssueOption) UnmarshalJSON(data []byte) error {
	var aux createIssueOptionCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = CreateIssueOption(aux)
	o.Labels = normaliseLabelRefs(o.Labels)
	return nil
}

type createPullRequestOptionCompat CreatePullRequestOption

func (o *CreatePullRequestOption) UnmarshalJSON(data []byte) error {
	var aux createPullRequestOptionCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = CreatePullRequestOption(aux)
	o.Labels = normaliseLabelRefs(o.Labels)
	return nil
}

type editPullRequestOptionCompat EditPullRequestOption

func (o *EditPullRequestOption) UnmarshalJSON(data []byte) error {
	var aux editPullRequestOptionCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = EditPullRequestOption(aux)
	o.Labels = normaliseLabelRefs(o.Labels)
	return nil
}

func normaliseLabelRefs(v any) any {
	items, ok := v.([]any)
	if !ok {
		return v
	}
	if len(items) == 0 {
		return []string{}
	}

	strs := make([]string, 0, len(items))
	for _, item := range items {
		s, ok := item.(string)
		if !ok {
			strs = nil
			break
		}
		strs = append(strs, s)
	}
	if strs != nil {
		return strs
	}

	ints := make([]int64, 0, len(items))
	for _, item := range items {
		switch x := item.(type) {
		case float64:
			if float64(int64(x)) != x {
				return v
			}
			ints = append(ints, int64(x))
		case int64:
			ints = append(ints, x)
		case int:
			ints = append(ints, int64(x))
		default:
			return v
		}
	}
	return ints
}
