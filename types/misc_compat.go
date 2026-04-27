package types

import json "github.com/goccy/go-json"

type mergePullRequestOptionCompat MergePullRequestOption

func (o MergePullRequestOption) MarshalJSON() ([]byte, error) {
	aux := mergePullRequestOptionCompat(o)
	if aux.Do == "" && o.MergeStyle != "" {
		aux.Do = o.MergeStyle
	}
	return json.Marshal(aux)
}

func (o *MergePullRequestOption) UnmarshalJSON(data []byte) error {
	var aux mergePullRequestOptionCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = MergePullRequestOption(aux)
	if o.MergeStyle == "" {
		o.MergeStyle = o.Do
	}
	return nil
}
