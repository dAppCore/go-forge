package forge

func ExampleReleaseListOptions_String() {
	_ = (*ReleaseListOptions).String
}

func ExampleReleaseListOptions_GoString() {
	_ = (*ReleaseListOptions).GoString
}

func ExampleReleaseAttachmentUploadOptions_String() {
	_ = (*ReleaseAttachmentUploadOptions).String
}

func ExampleReleaseAttachmentUploadOptions_GoString() {
	_ = (*ReleaseAttachmentUploadOptions).GoString
}

func ExampleReleaseService_ListReleasesPage() {
	_ = (*ReleaseService).ListReleasesPage
}

func ExampleReleaseService_ListReleases() {
	_ = (*ReleaseService).ListReleases
}

func ExampleReleaseService_IterReleases() {
	_ = (*ReleaseService).IterReleases
}

func ExampleReleaseService_CreateRelease() {
	_ = (*ReleaseService).CreateRelease
}

func ExampleReleaseService_GetByTag() {
	_ = (*ReleaseService).GetByTag
}

func ExampleReleaseService_GetRelease() {
	_ = (*ReleaseService).GetRelease
}

func ExampleReleaseService_GetLatest() {
	_ = (*ReleaseService).GetLatest
}

func ExampleReleaseService_DeleteByTag() {
	_ = (*ReleaseService).DeleteByTag
}

func ExampleReleaseService_ListAssets() {
	_ = (*ReleaseService).ListAssets
}

func ExampleReleaseService_CreateAttachment() {
	_ = (*ReleaseService).CreateAttachment
}

func ExampleReleaseService_EditAttachment() {
	_ = (*ReleaseService).EditAttachment
}

func ExampleReleaseService_CreateAsset() {
	_ = (*ReleaseService).CreateAsset
}

func ExampleReleaseService_EditAsset() {
	_ = (*ReleaseService).EditAsset
}

func ExampleReleaseService_IterAssets() {
	_ = (*ReleaseService).IterAssets
}

func ExampleReleaseService_GetAsset() {
	_ = (*ReleaseService).GetAsset
}

func ExampleReleaseService_DeleteAsset() {
	_ = (*ReleaseService).DeleteAsset
}
