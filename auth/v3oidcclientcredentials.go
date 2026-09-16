package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
)

// V3OIDCClientCredentialsOpts authenticates through Keystone federation with
// an access token obtained from an OpenID Connect identity provider using
// the client credentials grant.
type V3OIDCClientCredentialsOpts struct {
	AuthURL string

	// IdentityProviderName and Protocol name the identity provider and the
	// federation protocol registered in Keystone.
	IdentityProviderName string `required:"true"`
	Protocol             string `required:"true"`

	ClientID string `required:"true"`

	// ClientSecret is sent with ClientID in HTTP Basic authentication. When
	// empty, ClientID is sent in the request body instead.
	ClientSecret string

	// AccessTokenEndpoint is the identity provider's token endpoint. When
	// empty, it is read from the document at DiscoveryEndpoint. Both must
	// use HTTPS, or HTTP on a loopback address.
	AccessTokenEndpoint string `or:"DiscoveryEndpoint"`
	DiscoveryEndpoint   string

	// AccessTokenType is the token response field sent to Keystone. It
	// defaults to "access_token".
	AccessTokenType string

	// OIDCScope defaults to "openid".
	OIDCScope string

	AllowReauth bool

	// Scope is the scope of the Keystone token. A nil scope requests an
	// unscoped token.
	Scope *Scope
}

func (opts *V3OIDCClientCredentialsOpts) GetAuthURL() string {
	return (AuthOptionsV3{AuthURL: opts.AuthURL}).GetAuthURL()
}

func (opts *V3OIDCClientCredentialsOpts) Authenticate(ctx context.Context, provider *gophercloud.ProviderClient) (*AuthResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("oidc: expected non-nil options")
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

	endpoint, err := resolveAccessTokenEndpoint(ctx, client, opts)
	if err != nil {
		return nil, fmt.Errorf("oidc: resolving the access token endpoint: %w", err)
	}
	accessToken, err := fetchAccessToken(ctx, client, opts, endpoint)
	if err != nil {
		return nil, fmt.Errorf("oidc: requesting an access token: %w", err)
	}

	var result gophercloud.Result
	resp, err := client.Post(ctx, client.ServiceURL("OS-FEDERATION", "identity_providers", opts.IdentityProviderName, "protocols", opts.Protocol, "auth"), nil, &result.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"Authorization": "Bearer " + accessToken},
		OmitHeaders: []string{"X-Auth-Token"},
		OkCodes:     []int{http.StatusCreated},
	})
	if err != nil {
		return nil, fmt.Errorf("oidc: exchanging the access token: %w", err)
	}
	result.Header = resp.Header
	unscopedToken := resp.Header.Get("X-Subject-Token")
	if unscopedToken == "" {
		return nil, fmt.Errorf("oidc: exchanging the access token: response missing X-Subject-Token header")
	}

	scope, err := opts.Scope.ToScopeMap()
	if err != nil {
		return nil, err
	}
	if scope == nil {
		return v3AuthResult(result, opts.AllowReauth)
	}
	rescope := AuthOptionsV3{AuthURL: opts.AuthURL, Auth: V3TokenOpts{Token: unscopedToken, Scope: opts.Scope, AllowReauth: opts.AllowReauth}}
	return rescope.Authenticate(ctx, provider)
}

func resolveAccessTokenEndpoint(ctx context.Context, client *gophercloud.ServiceClient, opts *V3OIDCClientCredentialsOpts) (string, error) {
	if opts.AccessTokenEndpoint != "" {
		if err := checkOIDCEndpoint("access token endpoint", opts.AccessTokenEndpoint, true); err != nil {
			return "", err
		}
		return opts.AccessTokenEndpoint, nil
	}

	// The discovery document says where the client secret is sent, so it
	// must not be readable or replaceable in transit.
	if err := checkOIDCEndpoint("discovery endpoint", opts.DiscoveryEndpoint, true); err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, opts.DiscoveryEndpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create discovery request: %w", err)
	}
	status, body, err := doOIDCRequest(client.HTTPClient, req, "discovery document")
	if err != nil {
		return "", fmt.Errorf("failed to fetch discovery document from %s: %w", opts.DiscoveryEndpoint, err)
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("discovery endpoint returned HTTP %d", status)
	}
	var discovery struct {
		TokenEndpoint       string   `json:"token_endpoint"`
		GrantTypesSupported []string `json:"grant_types_supported"`
	}
	if err := json.Unmarshal(body, &discovery); err != nil {
		return "", fmt.Errorf("discovery document is not valid JSON: %w", err)
	}

	if len(discovery.GrantTypesSupported) > 0 && !slices.Contains(discovery.GrantTypesSupported, "client_credentials") {
		return "", fmt.Errorf("IdP does not support client_credentials grant type (supported: %v)", discovery.GrantTypesSupported)
	}
	if discovery.TokenEndpoint == "" {
		return "", fmt.Errorf("discovery document does not contain a valid token_endpoint")
	}
	// A remote identity provider must not point the client secret at a
	// plaintext port on this machine.
	discoveryURL, _ := url.Parse(opts.DiscoveryEndpoint)
	if err := checkOIDCEndpoint("discovery document token_endpoint", discovery.TokenEndpoint, isLoopbackHost(discoveryURL.Hostname())); err != nil {
		return "", err
	}
	return discovery.TokenEndpoint, nil
}

// checkOIDCEndpoint refuses an endpoint that would carry the client secret,
// or say where to send it, without TLS. Plain HTTP is allowed only on
// loopback, and only when allowLoopbackHTTP is set.
func checkOIDCEndpoint(name, endpoint string, allowLoopbackHTTP bool) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid %s: %w", name, err)
	}
	if u.Scheme == "https" || (allowLoopbackHTTP && u.Scheme == "http" && isLoopbackHost(u.Hostname())) {
		return nil
	}
	return fmt.Errorf("%s %q does not use https", name, endpoint)
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func fetchAccessToken(ctx context.Context, client *gophercloud.ServiceClient, opts *V3OIDCClientCredentialsOpts, endpoint string) (string, error) {
	scope := opts.OIDCScope
	if scope == "" {
		scope = "openid"
	}
	form := url.Values{
		"grant_type": {"client_credentials"},
		"scope":      {scope},
	}
	if opts.ClientSecret == "" {
		form.Set("client_id", opts.ClientID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Send the credentials as they are, like keystoneauth. Identity providers
	// that work with the OpenStack CLI then work here too.
	if opts.ClientSecret != "" {
		req.SetBasicAuth(opts.ClientID, opts.ClientSecret)
	}

	status, body, err := doOIDCRequest(client.HTTPClient, req, "IdP response")
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		snippet := string(body)
		if len(snippet) > 512 {
			snippet = snippet[:512] + "..."
		}
		return "", fmt.Errorf("IdP returned HTTP %d: %s", status, snippet)
	}

	var response map[string]any
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to parse IdP token response: %w", err)
	}
	field := opts.AccessTokenType
	if field == "" {
		field = "access_token"
	}
	token, ok := response[field].(string)
	if !ok || token == "" {
		return "", fmt.Errorf("IdP response missing %q field", field)
	}
	if field == "access_token" {
		tokenType, ok := response["token_type"].(string)
		if !ok || !strings.EqualFold(tokenType, "Bearer") {
			return "", fmt.Errorf("IdP response has unsupported token_type %q", tokenType)
		}
	}
	return token, nil
}

const oidcMaxResponseSize = 1 << 20 // 1 MiB

// doOIDCRequest sends a request to the identity provider and returns the
// response status and body. A token request must stay on its original
// origin. Discovery may redirect across hosts, but neither request may be
// redirected from HTTPS to plaintext.
func doOIDCRequest(client http.Client, req *http.Request, what string) (int, []byte, error) {
	req.Header.Set("Accept", "application/json")

	checkRedirect := client.CheckRedirect
	client.CheckRedirect = func(next *http.Request, via []*http.Request) error {
		if checkRedirect != nil {
			if err := checkRedirect(next, via); err != nil {
				return err
			}
		} else if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 redirects")
		}
		if via[len(via)-1].URL.Scheme == "https" && next.URL.Scheme != "https" {
			return fmt.Errorf("request cannot redirect from HTTPS to plaintext")
		}
		if req.Method == http.MethodPost && oidcOrigin(req.URL) != oidcOrigin(next.URL) {
			return fmt.Errorf("token request cannot redirect to another origin")
		}
		return nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, oidcMaxResponseSize+1))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read %s: %w", what, err)
	}
	if len(body) > oidcMaxResponseSize {
		return 0, nil, fmt.Errorf("%s exceeds %d bytes", what, oidcMaxResponseSize)
	}
	return resp.StatusCode, body, nil
}

func oidcOrigin(u *url.URL) string {
	port := u.Port()
	if port == "" {
		port = "80"
		if u.Scheme == "https" {
			port = "443"
		}
	}
	return strings.ToLower(u.Scheme) + "://" + net.JoinHostPort(strings.ToLower(u.Hostname()), port)
}
