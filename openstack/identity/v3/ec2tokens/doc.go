/*
Package ec2tokens provides low-level access to Keystone's EC2 and S3 token endpoints.

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
identityClient below is an already-configured Identity service client, including
the service credentials needed to authorize the EC2 token endpoint.

	func createToken(ctx context.Context, identityClient *gophercloud.ServiceClient) (*tokens.Token, error) {
	 request := tokenRequest{
	 Body: map[string]any{"credentials": map[string]any{
	 "access": "access-key",
	 "host": "ec2.example.com",
	 "path": "/",
	 "verb": "GET",
	 "params": map[string]string{"SignatureVersion": "2", "SignatureMethod": "HmacSHA256"},
	 "signature": "precomputed-signature",
	}},
	 Headers: map[string]string{"Accept": "application/json"},
	 }
	 return ec2tokens.Create(ctx, identityClient, request).ExtractToken()
	}

The signature must have been computed for the supplied AWS request. This package
does not sign credentials. ValidateS3Token similarly accepts a complete body via
ValidateS3TokenOptsBuilder and optional headers via
ValidateS3TokenOptsHeadersBuilder.
*/
package ec2tokens
