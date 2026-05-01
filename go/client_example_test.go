package forge

func ExampleAPIError_Error() {
	_ = (*APIError).Error
}

func ExampleAPIError_String() {
	_ = (*APIError).String
}

func ExampleAPIError_GoString() {
	_ = (*APIError).GoString
}

func ExampleIsNotFound() {
	_ = IsNotFound
}

func ExampleIsForbidden() {
	_ = IsForbidden
}

func ExampleIsConflict() {
	_ = IsConflict
}

func ExampleWithHTTPClient() {
	_ = WithHTTPClient
}

func ExampleWithUserAgent() {
	_ = WithUserAgent
}

func ExampleRateLimit_String() {
	_ = (*RateLimit).String
}

func ExampleRateLimit_GoString() {
	_ = (*RateLimit).GoString
}

func ExampleClient_BaseURL() {
	_ = (*Client).BaseURL
}

func ExampleClient_RateLimit() {
	_ = (*Client).RateLimit
}

func ExampleClient_UserAgent() {
	_ = (*Client).UserAgent
}

func ExampleClient_HTTPClient() {
	_ = (*Client).HTTPClient
}

func ExampleClient_String() {
	_ = (*Client).String
}

func ExampleClient_GoString() {
	_ = (*Client).GoString
}

func ExampleClient_HasToken() {
	_ = (*Client).HasToken
}

func ExampleNewClient() {
	_ = NewClient
}

func ExampleClient_Get() {
	_ = (*Client).Get
}

func ExampleClient_Post() {
	_ = (*Client).Post
}

func ExampleClient_Patch() {
	_ = (*Client).Patch
}

func ExampleClient_Put() {
	_ = (*Client).Put
}

func ExampleClient_Delete() {
	_ = (*Client).Delete
}

func ExampleClient_DeleteWithBody() {
	_ = (*Client).DeleteWithBody
}

func ExampleClient_PostRaw() {
	_ = (*Client).PostRaw
}

func ExampleClient_GetRaw() {
	_ = (*Client).GetRaw
}
