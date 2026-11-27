/*
Package oauth1 manages OAuth1 consumers, request tokens, and access tokens.

For normal Keystone authentication, use github.com/gophercloud/gophercloud/v2/auth.
OAuth1 credential management examples below remain separate from exchanging those
credentials for a Keystone token.

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
	 "identity": map[string]any{"methods": []string{"oauth1"}, "oauth1": map[string]any{}},
	}},
	 Headers: map[string]string{"Authorization": "OAuth <precomputed signature parameters>"},
	 }
	 return oauth1.Create(ctx, identityClient, request).ExtractToken()
	}

The Authorization value must be an OAuth1 header signed for the actual POST URL.
Create does not calculate signatures or assemble authentication credentials.

Example to Create an OAuth1 Consumer

	createConsumerOpts := oauth1.CreateConsumerOpts{
		Description: "My consumer",
	}
	consumer, err := oauth1.CreateConsumer(context.TODO(), identityClient, createConsumerOpts).Extract()
	if err != nil {
		panic(err)
	}

	// NOTE: Consumer secret is available only on create response
	fmt.Printf("Consumer: %+v\n", consumer)

Example to Request an unauthorized OAuth1 token

	requestTokenOpts := oauth1.RequestTokenOpts{
		OAuthConsumerKey:     consumer.ID,
		OAuthConsumerSecret:  consumer.Secret,
		OAuthSignatureMethod: oauth1.HMACSHA1,
		RequestedProjectID:   projectID,
	}
	requestToken, err := oauth1.RequestToken(context.TODO(), identityClient, requestTokenOpts).Extract()
	if err != nil {
		panic(err)
	}

	// NOTE: Request token secret is available only on request response
	fmt.Printf("Request token: %+v\n", requestToken)

Example to Authorize an unauthorized OAuth1 token

	authorizeTokenOpts := oauth1.AuthorizeTokenOpts{
		Roles: []oauth1.Role{
			{Name: "member"},
		},
	}
	authToken, err := oauth1.AuthorizeToken(context.TODO(), identityClient, requestToken.OAuthToken, authorizeTokenOpts).Extract()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Verifier ID of the unauthorized Token: %+v\n", authToken.OAuthVerifier)

Example to Create an OAuth1 Access Token

	accessTokenOpts := oauth1.CreateAccessTokenOpts{
		OAuthConsumerKey:     consumer.ID,
		OAuthConsumerSecret:  consumer.Secret,
		OAuthToken:           requestToken.OAuthToken,
		OAuthTokenSecret:     requestToken.OAuthTokenSecret,
		OAuthVerifier:        authToken.OAuthVerifier,
		OAuthSignatureMethod: oauth1.HMACSHA1,
	}
	accessToken, err := oauth1.CreateAccessToken(context.TODO(), identityClient, accessTokenOpts).Extract()
	if err != nil {
		panic(err)
	}

	// NOTE: Access token secret is available only on create response
	fmt.Printf("OAuth1 Access Token: %+v\n", accessToken)

Example to List User's OAuth1 Access Tokens

	allPages, err := oauth1.ListAccessTokens(identityClient, userID).AllPages(context.TODO())
	if err != nil {
		panic(err)
	}
	accessTokens, err := oauth1.ExtractAccessTokens(allPages)
	if err != nil {
		panic(err)
	}

	for _, accessToken := range accessTokens {
		fmt.Printf("Access Token: %+v\n", accessToken)
	}
*/
package oauth1
