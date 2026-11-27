/*
Package tokens provides low-level access to Identity v3 token creation, retrieval, validation, and revocation.

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
	 "identity": map[string]any{
	  "methods": []string{"password"},
	  "password": map[string]any{"user": map[string]any{"id": "user-id", "password": "password"}},
	 },
	 "scope": map[string]any{"project": map[string]any{"id": "project-id"}},
	}},
	 Headers: map[string]string{"Accept": "application/json"},
	 }
	 return tokens.Create(ctx, identityClient, request).ExtractToken()
	}

Example to Get a Token

	token, err := tokens.Get(context.TODO(), identityClient, "token_id", nil).ExtractToken()
	if err != nil {
		panic(err)
	}

# Example to Get a Token Created with Application Credentials Access Rules

When validating or retrieving tokens that were created using application
credentials with access rules, the OpenStack-Identity-Access-Rules header
must be sent. Without this header, Keystone will return a 404 Not Found error.

See https://docs.openstack.org/keystone/latest/user/application_credentials.html

	getOpts := tokens.GetOpts{
		AccessRulesVersion: "1.0",
	}
	token, err := tokens.Get(context.TODO(), identityClient, "token_id", getOpts).ExtractToken()
	if err != nil {
		panic(err)
	}

Example to Validate a Token

	ok, err := tokens.Validate(context.TODO(), identityClient, "token_id", nil)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("Token is valid")
	}

Example to Validate a Token Created with Application Credentials Access Rules

	validateOpts := tokens.ValidateOpts{
		AccessRulesVersion: "1.0",
	}
	ok, err := tokens.Validate(context.TODO(), identityClient, "token_id", validateOpts)
	if err != nil {
		panic(err)
	}

	if ok {
		fmt.Println("Token is valid")
	}

# Revoking a token

	err := tokens.Revoke(context.Background(), identityClient, "token-id").Err
	if err != nil {
		panic(err)
	}
*/
package tokens
