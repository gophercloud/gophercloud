/*
Package tokens provides low-level access to Identity v2 token creation and retrieval.

For normal Keystone authentication, use github.com/gophercloud/gophercloud/v2/auth.

# Manually creating a Keystone token

Create accepts a complete request body through CreateOptsBuilder. It optionally
accepts headers through CreateOptsHeadersBuilder. Define the request type and its
methods at package scope; auth credential options do not implement these interfaces.

	type tokenRequest struct {
	 // Body is the complete JSON request, including its outer envelope.
	 Body map[string]any
	 // Headers contains additional request headers.
	 Headers map[string]string
	}

	func (opts tokenRequest) ToTokenCreateMap() (map[string]any, error) {
	 return opts.Body, nil
	}

	func (opts tokenRequest) ToTokenCreateHeaders(method, url string) (map[string]string, error) {
	 return opts.Headers, nil
	}

The caller supplies the complete JSON body and any additional headers. The
identityClient below is an already-configured Identity service client.

	func createToken(ctx context.Context, identityClient *gophercloud.ServiceClient) (*tokens.Token, error) {
	 request := tokenRequest{
	 Body: map[string]any{"auth": map[string]any{
	 "passwordCredentials": map[string]any{"username": "username", "password": "password"},
	 "tenantId": "tenant-id",
	}},
	 Headers: map[string]string{"Accept": "application/json"},
	 }
	 return tokens.Create(ctx, identityClient, request).ExtractToken()
	}
*/
package tokens
