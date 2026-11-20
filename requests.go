package gophercloud

// RawRequest contains the complete JSON body and optional headers for a
// token request. The caller is responsible for constructing the auth body.
type RawRequest struct {
	// Body is the complete JSON body, including the auth key.
	Body map[string]any
	// Headers contains additional HTTP headers to send.
	Headers map[string]string
}
