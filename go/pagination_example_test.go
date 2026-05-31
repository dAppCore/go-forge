package forge

func ExampleListOptions_String() {
	_ = (*ListOptions).String
}

func ExampleListOptions_GoString() {
	_ = (*ListOptions).GoString
}

func ExamplePagedResult_String() {
	_ = (*PagedResult[int]).String
}

func ExamplePagedResult_GoString() {
	_ = (*PagedResult[int]).GoString
}

func ExampleListPage() {
	_ = ListPage[int]
}

func ExampleListAll() {
	_ = ListAll[int]
}

func ExampleListIter() {
	_ = ListIter[int]
}
