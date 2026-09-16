package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// V3OAuth2MTLSOpts authenticates with the OAuth 2.0 client credentials grant
// over mutual TLS. The client certificate comes from the provider's HTTP
// client.
type V3OAuth2MTLSOpts struct {
	AuthURL string

	// OAuth2Endpoint defaults to OS-OAUTH2/token under AuthURL.
	OAuth2Endpoint string

	// ClientID is the ID of the Keystone user the certificate is mapped to.
	ClientID string `required:"true"`

	AllowReauth bool
}

func (opts *V3OAuth2MTLSOpts) GetAuthURL() string {
	return (AuthOptionsV3{AuthURL: opts.AuthURL}).GetAuthURL()
}

func (opts *V3OAuth2MTLSOpts) Authenticate(ctx context.Context, provider *gophercloud.ProviderClient) (*AuthResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("oauth2mtls: expected non-nil options")
	}
	if opts.AuthURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}
	if _, err := gophercloud.BuildRequestBody(opts, ""); err != nil {
		return nil, err
	}

	if provider == nil {
		provider = &gophercloud.ProviderClient{}
	}
	client := &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: opts.GetAuthURL()}

	endpoint := opts.OAuth2Endpoint
	if endpoint == "" {
		endpoint = client.ServiceURL("OS-OAUTH2", "token")
	}
	form := url.Values{
		"grant_type": {"client_credentials"},
		"client_id":  {opts.ClientID},
	}

	var token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	resp, err := client.Post(ctx, endpoint, strings.NewReader(form.Encode()), &token, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		OkCodes:     []int{http.StatusOK},
	})
	if _, _, err := gophercloud.ParseResponse(resp, err); err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("oauth2mtls: token response missing access_token field")
	}
	if !strings.EqualFold(token.TokenType, "Bearer") {
		return nil, fmt.Errorf("oauth2mtls: token response has unsupported token_type %q", token.TokenType)
	}

	// The OS-OAUTH2 response has no catalog, so retrieve the full token.
	var result gophercloud.Result
	resp, err = client.Get(ctx, client.ServiceURL("auth", "tokens"), &result.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Token":    token.AccessToken,
			"X-Subject-Token": token.AccessToken,
		},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(resp, err)
	if result.Err == nil {
		result.Header.Set("X-Subject-Token", token.AccessToken)
	}
	authResult, err := v3AuthResult(result, opts.AllowReauth)
	if err != nil {
		return nil, err
	}
	// Services expect the certificate-bound token in the Authorization header.
	authResult.TokenType = "Bearer"
	return authResult, nil
}
