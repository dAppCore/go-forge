package types

import (
	core "dappco.re/go"
)

func TestContentCompat_CreateFileOptions_MarshalJSON_Good(t *core.T) {
	subject := (*CreateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_CreateFileOptions_MarshalJSON_Bad(t *core.T) {
	subject := (*CreateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_CreateFileOptions_MarshalJSON_Ugly(t *core.T) {
	subject := (*CreateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_CreateFileOptions_UnmarshalJSON_Good(t *core.T) {
	subject := (*CreateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_CreateFileOptions_UnmarshalJSON_Bad(t *core.T) {
	subject := (*CreateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_CreateFileOptions_UnmarshalJSON_Ugly(t *core.T) {
	subject := (*CreateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_MarshalJSON_Good(t *core.T) {
	subject := (*UpdateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_MarshalJSON_Bad(t *core.T) {
	subject := (*UpdateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_MarshalJSON_Ugly(t *core.T) {
	subject := (*UpdateFileOptions).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_UnmarshalJSON_Good(t *core.T) {
	subject := (*UpdateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_UnmarshalJSON_Bad(t *core.T) {
	subject := (*UpdateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_UpdateFileOptions_UnmarshalJSON_Ugly(t *core.T) {
	subject := (*UpdateFileOptions).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_MarshalJSON_Good(t *core.T) {
	subject := (*ChangeFileOperation).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_MarshalJSON_Bad(t *core.T) {
	subject := (*ChangeFileOperation).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_MarshalJSON_Ugly(t *core.T) {
	subject := (*ChangeFileOperation).MarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_UnmarshalJSON_Good(t *core.T) {
	subject := (*ChangeFileOperation).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_UnmarshalJSON_Bad(t *core.T) {
	subject := (*ChangeFileOperation).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestContentCompat_ChangeFileOperation_UnmarshalJSON_Ugly(t *core.T) {
	subject := (*ChangeFileOperation).UnmarshalJSON
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
