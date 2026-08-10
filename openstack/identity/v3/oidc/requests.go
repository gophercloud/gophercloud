package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

type discoveryDocument struct {
	TokenEndpoint       string   `json:"token_endpoint"`
	GrantTypesSupported []string `json:"grant_types_supported"`
}

// AuthOptions contains OIDC client-credentials authentication options.
type AuthOptions struct {
	// IdentityProviderName is the Keystone federation identity provider.
	IdentityProviderName string

	// Protocol is the Keystone federation protocol.
	Protocol string

	// ClientID is the OAuth 2.0 client identifier.
	ClientID string

	// ClientSecret is the OAuth 2.0 client secret.
	ClientSecret string

	// AccessTokenEndpoint is the identity provider's token endpoint.
	AccessTokenEndpoint string

	// AccessTokenType selects the token response field. It defaults to "access_token".
	AccessTokenType string

	// DiscoveryEndpoint is used when AccessTokenEndpoint is empty.
	DiscoveryEndpoint string

	// OIDCScope is the OAuth 2.0 scope. It defaults to "openid".
	OIDCScope string

	// AllowReauth enables automatic reauthentication.
	AllowReauth bool

	// Scope controls the resulting Keystone token scope.
	Scope tokens.Scope
}

// ToTokenV3ScopeMap builds the Keystone token scope.
func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]any, error) {
	return (&tokens.AuthOptions{Scope: opts.Scope}).ToTokenV3ScopeMap()
}

// ToTokenV3HeadersMap implements tokens.AuthOptionsBuilder.
func (opts *AuthOptions) ToTokenV3HeadersMap(map[string]any) (map[string]string, error) {
	return nil, nil
}

// CanReauth reports whether automatic reauthentication is enabled.
func (opts *AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}

// ToTokenV3CreateMap implements tokens.AuthOptionsBuilder.
func (opts *AuthOptions) ToTokenV3CreateMap(map[string]any) (map[string]any, error) {
	return nil, nil
}

// Create authenticates with OIDC client credentials.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts tokens.AuthOptionsBuilder) (r tokens.CreateResult) {
	oidcOpts, ok := opts.(*AuthOptions)
	if !ok || oidcOpts == nil {
		r.Err = fmt.Errorf("expected *oidc.AuthOptions, got %T", opts)
		return
	}

	if err := validateAuthOptions(oidcOpts); err != nil {
		r.Err = err
		return
	}
	if c == nil || c.ProviderClient == nil {
		r.Err = fmt.Errorf("oidc: ServiceClient or ProviderClient is nil")
		return
	}

	accessTokenEndpoint, err := resolveAccessTokenEndpoint(ctx, c, oidcOpts)
	if err != nil {
		r.Err = fmt.Errorf("OIDC access token endpoint resolution failed: %w", err)
		return
	}

	accessToken, err := fetchAccessToken(ctx, c, oidcOpts, accessTokenEndpoint)
	if err != nil {
		r.Err = fmt.Errorf("OIDC IdP token request failed: %w", err)
		return
	}

	unscopedToken, unscopedBody, unscopedHeader, err := exchangeForKeystoneToken(ctx, c, oidcOpts, accessToken)
	if err != nil {
		r.Err = fmt.Errorf("OIDC federation auth failed: %w", err)
		return
	}

	scope, err := oidcOpts.ToTokenV3ScopeMap()
	if err != nil {
		r.Err = err
		return
	}

	if scope != nil {
		r = tokens.Create(ctx, c, &tokens.AuthOptions{TokenID: unscopedToken, Scope: oidcOpts.Scope})
	} else {
		r.Body = unscopedBody
		r.Header = unscopedHeader
	}

	return
}

func validateAuthOptions(opts *AuthOptions) error {
	if opts.IdentityProviderName == "" {
		return fmt.Errorf("missing required field: IdentityProviderName")
	}
	if opts.Protocol == "" {
		return fmt.Errorf("missing required field: Protocol")
	}
	if opts.ClientID == "" {
		return fmt.Errorf("missing required field: ClientID")
	}
	if opts.AccessTokenEndpoint == "" && opts.DiscoveryEndpoint == "" {
		return fmt.Errorf("at least one of AccessTokenEndpoint or DiscoveryEndpoint must be provided")
	}
	return nil
}

func resolveAccessTokenEndpoint(ctx context.Context, c *gophercloud.ServiceClient, opts *AuthOptions) (string, error) {
	if opts.AccessTokenEndpoint != "" {
		return opts.AccessTokenEndpoint, nil
	}

	discovery, err := fetchDiscoveryDocument(ctx, c, opts.DiscoveryEndpoint)
	if err != nil {
		return "", err
	}

	if len(discovery.GrantTypesSupported) > 0 && !slices.Contains(discovery.GrantTypesSupported, "client_credentials") {
		return "", fmt.Errorf("IdP does not support client_credentials grant type (supported: %v)", discovery.GrantTypesSupported)
	}

	if discovery.TokenEndpoint == "" {
		return "", fmt.Errorf("discovery document does not contain a valid token_endpoint")
	}

	return discovery.TokenEndpoint, nil
}

func fetchDiscoveryDocument(ctx context.Context, c *gophercloud.ServiceClient, discoveryURL string) (discoveryDocument, error) {
	var doc discoveryDocument
	if c == nil || c.ProviderClient == nil {
		return doc, fmt.Errorf("service client or provider client is nil")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", discoveryURL, nil)
	if err != nil {
		return doc, fmt.Errorf("failed to create discovery request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	httpClient := c.HTTPClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return doc, fmt.Errorf("failed to fetch discovery document from %s: %w", discoveryURL, err)
	}
	defer resp.Body.Close()

	const maxDiscoverySize = 1 << 20 // 1 MiB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxDiscoverySize+1))
	if err != nil {
		return doc, fmt.Errorf("failed to read discovery document response: %w", err)
	}
	if int64(len(body)) > maxDiscoverySize {
		return doc, fmt.Errorf("discovery document exceeds maximum allowed size of %d bytes", maxDiscoverySize)
	}

	if resp.StatusCode != http.StatusOK {
		return doc, fmt.Errorf("discovery endpoint returned HTTP %d", resp.StatusCode)
	}

	if err := json.Unmarshal(body, &doc); err != nil {
		return doc, fmt.Errorf("discovery document is not valid JSON: %w", err)
	}
	return doc, nil
}

func fetchAccessToken(ctx context.Context, c *gophercloud.ServiceClient, opts *AuthOptions, accessTokenEndpoint string) (string, error) {
	if c == nil || c.ProviderClient == nil {
		return "", fmt.Errorf("service client or provider client is nil")
	}

	scope := opts.OIDCScope
	if scope == "" {
		scope = "openid"
	}

	formData := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {scope},
	}
	if opts.ClientSecret == "" {
		formData.Set("client_id", opts.ClientID)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", accessTokenEndpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	if opts.ClientSecret != "" {
		req.SetBasicAuth(opts.ClientID, opts.ClientSecret)
	}

	httpClient := c.HTTPClient
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	const maxIDPResponseSize = 1 << 20 // 1 MiB
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxIDPResponseSize+1))
	if err != nil {
		return "", fmt.Errorf("failed to read IdP response: %w", err)
	}
	if int64(len(body)) > maxIDPResponseSize {
		return "", fmt.Errorf("IdP response body exceeded %d bytes", maxIDPResponseSize)
	}

	if resp.StatusCode != http.StatusOK {
		snippet := string(body)
		if len(snippet) > 512 {
			snippet = snippet[:512] + "..."
		}
		return "", fmt.Errorf("IdP returned HTTP %d: %s", resp.StatusCode, snippet)
	}

	var tokenResponse map[string]any
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to parse IdP token response: %w", err)
	}

	tokenField := opts.AccessTokenType
	if tokenField == "" {
		tokenField = "access_token"
	}

	accessToken, ok := tokenResponse[tokenField].(string)
	if !ok || accessToken == "" {
		return "", fmt.Errorf("IdP response missing %q field", tokenField)
	}
	if tokenField == "access_token" {
		tokenType, ok := tokenResponse["token_type"].(string)
		if !ok || !strings.EqualFold(tokenType, "Bearer") {
			return "", fmt.Errorf("IdP response has unsupported token_type %q", tokenType)
		}
	}

	return accessToken, nil
}

func exchangeForKeystoneToken(ctx context.Context, c *gophercloud.ServiceClient, opts *AuthOptions, accessToken string) (string, any, http.Header, error) {
	federationURL := authURL(c, opts.IdentityProviderName, opts.Protocol)

	var body any
	resp, err := c.Post(ctx, federationURL, nil, &body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"Authorization": "Bearer " + accessToken,
		},
		OkCodes:     []int{201},
		OmitHeaders: []string{"X-Auth-Token"},
	})
	if err != nil {
		return "", nil, nil, err
	}

	unscopedToken := resp.Header.Get("X-Subject-Token")
	if unscopedToken == "" {
		return "", nil, nil, fmt.Errorf("federation auth response missing X-Subject-Token header")
	}

	return unscopedToken, body, resp.Header, nil
}
