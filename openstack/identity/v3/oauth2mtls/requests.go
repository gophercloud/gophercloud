package oauth2mtls

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// AuthOptions contains OAuth2 mTLS client-credentials options.
type AuthOptions struct {
	// OAuth2Endpoint specifies Keystone's OS-OAUTH2 token endpoint. When empty,
	// it is derived from the identity ServiceClient endpoint.
	OAuth2Endpoint string

	// ClientID is the Keystone user ID associated with the client certificate.
	ClientID string

	// AllowReauth enables automatic reauthentication.
	AllowReauth bool
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

// ToTokenV3ScopeMap implements tokens.AuthOptionsBuilder.
func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]any, error) {
	return nil, nil
}

// ToTokenV3HeadersMap implements tokens.AuthOptionsBuilder.
func (opts *AuthOptions) ToTokenV3HeadersMap(map[string]any) (map[string]string, error) {
	return nil, nil
}

// ToTokenV3CreateMap implements tokens.AuthOptionsBuilder.
func (opts *AuthOptions) ToTokenV3CreateMap(map[string]any) (map[string]any, error) {
	return nil, nil
}

// CanReauth reports whether automatic reauthentication is enabled.
func (opts *AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}

func (opts *AuthOptions) validate() error {
	if opts.ClientID == "" {
		return gophercloud.ErrMissingInput{Argument: "ClientID"}
	}
	return nil
}

// Create authenticates with OAuth2 mTLS client credentials.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts tokens.AuthOptionsBuilder) (r tokens.CreateResult) {
	mtlsOpts, ok := opts.(*AuthOptions)
	if !ok || mtlsOpts == nil {
		r.Err = fmt.Errorf("oauth2mtls: expected non-nil *oauth2mtls.AuthOptions, got %T", opts)
		return
	}

	if err := mtlsOpts.validate(); err != nil {
		r.Err = err
		return
	}

	if c == nil || c.ProviderClient == nil {
		r.Err = fmt.Errorf("oauth2mtls: ServiceClient or ProviderClient is nil")
		return
	}

	oauth2Endpoint := mtlsOpts.OAuth2Endpoint
	if oauth2Endpoint == "" {
		oauth2Endpoint = tokenURL(c)
	}

	formData := url.Values{
		"grant_type": {"client_credentials"},
		"client_id":  {mtlsOpts.ClientID},
	}

	var tokenResp tokenResponse
	authClient := unauthenticatedClient(c)
	resp, err := authClient.Post(ctx, oauth2Endpoint, strings.NewReader(formData.Encode()), &tokenResp, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Content-Type": "application/x-www-form-urlencoded"},
		OkCodes:     []int{http.StatusOK},
	})
	_, _, r.Err = gophercloud.ParseResponse(resp, err)
	if r.Err != nil {
		return
	}

	if tokenResp.AccessToken == "" {
		r.Err = fmt.Errorf("oauth2mtls: token response missing access_token field")
		return
	}
	if !strings.EqualFold(tokenResp.TokenType, "Bearer") {
		r.Err = fmt.Errorf("oauth2mtls: token response has unsupported token_type %q", tokenResp.TokenType)
		return
	}

	// The OS-OAUTH2 response has no catalog, so retrieve the full token.
	resp, err = authClient.Get(ctx, validateURL(authClient), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Token":    tokenResp.AccessToken,
			"X-Subject-Token": tokenResp.AccessToken,
		},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	if r.Err != nil {
		return
	}
	r.Header.Set("X-Subject-Token", tokenResp.AccessToken)

	return
}

func unauthenticatedClient(c *gophercloud.ServiceClient) *gophercloud.ServiceClient {
	client := *c
	provider := *c.ProviderClient
	provider.Throwaway = true
	provider.ReauthFunc = nil
	client.ProviderClient = &provider
	return &client
}
