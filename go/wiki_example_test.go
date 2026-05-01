package forge

func ExampleWikiService_ListPages() {
	_ = (*WikiService).ListPages
}

func ExampleWikiService_IterPages() {
	_ = (*WikiService).IterPages
}

func ExampleWikiService_GetPage() {
	_ = (*WikiService).GetPage
}

func ExampleWikiService_GetPageRevisions() {
	_ = (*WikiService).GetPageRevisions
}

func ExampleWikiService_CreatePage() {
	_ = (*WikiService).CreatePage
}

func ExampleWikiService_EditPage() {
	_ = (*WikiService).EditPage
}

func ExampleWikiService_DeletePage() {
	_ = (*WikiService).DeletePage
}
