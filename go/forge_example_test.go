package forge

func ExampleNewForge() {
	_ = NewForge
}

func ExampleForge_Client() {
	_ = (*Forge).Client
}

func ExampleForge_BaseURL() {
	_ = (*Forge).BaseURL
}

func ExampleForge_RateLimit() {
	_ = (*Forge).RateLimit
}

func ExampleForge_UserAgent() {
	_ = (*Forge).UserAgent
}

func ExampleForge_HTTPClient() {
	_ = (*Forge).HTTPClient
}

func ExampleForge_HasToken() {
	_ = (*Forge).HasToken
}

func ExampleForge_String() {
	_ = (*Forge).String
}

func ExampleForge_GoString() {
	_ = (*Forge).GoString
}
