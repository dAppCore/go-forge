package forge

func ExampleResource_String() {
	_ = (*Resource[int, int, int]).String
}

func ExampleResource_GoString() {
	_ = (*Resource[int, int, int]).GoString
}

func ExampleNewResource() {
	_ = NewResource[int, int, int]
}

func ExampleResource_List() {
	_ = (*Resource[int, int, int]).List
}

func ExampleResource_ListAll() {
	_ = (*Resource[int, int, int]).ListAll
}

func ExampleResource_Iter() {
	_ = (*Resource[int, int, int]).Iter
}

func ExampleResource_Get() {
	_ = (*Resource[int, int, int]).Get
}

func ExampleResource_Create() {
	_ = (*Resource[int, int, int]).Create
}

func ExampleResource_Update() {
	_ = (*Resource[int, int, int]).Update
}

func ExampleResource_Delete() {
	_ = (*Resource[int, int, int]).Delete
}
