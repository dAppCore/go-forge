package types

import json "github.com/goccy/go-json"

type createFileOptionsCompat CreateFileOptions

func (o CreateFileOptions) MarshalJSON() ([]byte, error) {
	aux := createFileOptionsCompat(o)
	if aux.ContentBase64 == "" && o.Content != "" {
		aux.ContentBase64 = o.Content
	}
	return json.Marshal(aux)
}

func (o *CreateFileOptions) UnmarshalJSON(data []byte) error {
	var aux createFileOptionsCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = CreateFileOptions(aux)
	if o.Content == "" {
		o.Content = o.ContentBase64
	}
	return nil
}

type updateFileOptionsCompat UpdateFileOptions

func (o UpdateFileOptions) MarshalJSON() ([]byte, error) {
	aux := updateFileOptionsCompat(o)
	if aux.ContentBase64 == "" && o.Content != "" {
		aux.ContentBase64 = o.Content
	}
	return json.Marshal(aux)
}

func (o *UpdateFileOptions) UnmarshalJSON(data []byte) error {
	var aux updateFileOptionsCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = UpdateFileOptions(aux)
	if o.Content == "" {
		o.Content = o.ContentBase64
	}
	return nil
}

type changeFileOperationCompat ChangeFileOperation

func (o ChangeFileOperation) MarshalJSON() ([]byte, error) {
	aux := changeFileOperationCompat(o)
	if aux.ContentBase64 == "" && o.Content != "" {
		aux.ContentBase64 = o.Content
	}
	return json.Marshal(aux)
}

func (o *ChangeFileOperation) UnmarshalJSON(data []byte) error {
	var aux changeFileOperationCompat
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*o = ChangeFileOperation(aux)
	if o.Content == "" {
		o.Content = o.ContentBase64
	}
	return nil
}
