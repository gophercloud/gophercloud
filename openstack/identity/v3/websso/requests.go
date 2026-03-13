package websso

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// AuthOptions represents options for WebSSO authentication.
type AuthOptions struct {
	// IdentityEndpoint is the Keystone auth URL.
	IdentityEndpoint string

	// IdentityProvider is the name of the identity provider configured in Keystone.
	IdentityProvider string

	// Protocol is the federation protocol (typically "openid").
	Protocol string

	// RedirectHost is the hostname for the callback server (default: localhost).
	RedirectHost string

	// RedirectPort is the port for the callback server (default: 9990).
	RedirectPort int

	// CachePath is the directory where tokens are cached.
	CachePath string

	// Scope contains the project/domain scope information.
	Scope gophercloud.AuthScope

	// AllowReauth allows automatic reauthentication.
	AllowReauth bool
}

// callbackServer handles the OAuth callback from the identity provider.
type callbackServer struct {
	server   *http.Server
	token    string
	tokenErr error
	done     chan struct{}
}

func (s *callbackServer) handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		s.tokenErr = fmt.Errorf("failed to parse form: %w", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><title>Authentication Failed</title></head>
<body><p>The authentication flow failed: could not parse form data.</p>
<p>You can close this window.</p>
</body></html>`)
		close(s.done)
		return
	}

	token := r.FormValue("token")
	if token == "" {
		s.tokenErr = fmt.Errorf("no token received from identity provider")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><head><title>Authentication Failed</title></head>
<body><p>The authentication flow failed: no token received.</p>
<p>You can close this window.</p>
</body></html>`)
		close(s.done)
		return
	}

	s.token = token
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, `<html><head><title>Authentication Successful</title>
<script>window.close()</script></head>
<body><p>The authentication flow has been completed.</p>
<p>You can close this window.</p>
</body></html>`)
	close(s.done)
}

func waitForToken(redirectHost string, redirectPort int) (string, error) {
	srv := &callbackServer{
		done: make(chan struct{}),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/websso/", srv.handleCallback)

	srv.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", redirectHost, redirectPort),
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	listener, err := net.Listen("tcp", srv.server.Addr)
	if err != nil {
		return "", fmt.Errorf("cannot start callback server on %s:%d: %w", redirectHost, redirectPort, err)
	}

	go func() {
		_ = srv.server.Serve(listener)
	}()

	// Wait for callback or timeout
	select {
	case <-srv.done:
		_ = srv.server.Shutdown(context.Background())
		if srv.tokenErr != nil {
			return "", srv.tokenErr
		}
		if srv.token == "" {
			return "", fmt.Errorf("authentication failed: no token received")
		}
		return srv.token, nil
	case <-time.After(60 * time.Second):
		_ = srv.server.Shutdown(context.Background())
		return "", fmt.Errorf("authentication timeout: no response received within 60 seconds")
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}

func (opts *AuthOptions) redirectURI() string {
	host := opts.RedirectHost
	if host == "" {
		host = "localhost"
	}
	port := opts.RedirectPort
	if port == 0 {
		port = 9990
	}
	return fmt.Sprintf("http://%s:%d/auth/websso/", host, port)
}

// GetCacheID returns a cache identifier based on auth URL and identity provider.
func (opts *AuthOptions) GetCacheID() string {
	input := opts.IdentityEndpoint + "-" + opts.IdentityProvider
	re := regexp.MustCompile(`[^A-Za-z0-9-]+`)
	return "os-" + re.ReplaceAllString(input, "-")
}

func (opts *AuthOptions) getCachePath() string {
	cachePath := opts.CachePath
	if cachePath == "" {
		homeDir, _ := os.UserHomeDir()
		cachePath = filepath.Join(homeDir, ".cache", "gophercloud")
	}
	return filepath.Join(cachePath, opts.GetCacheID())
}

type cachedToken struct {
	AuthToken string          `json:"auth_token"`
	Body      json.RawMessage `json:"body"`
}

func (opts *AuthOptions) getCachedToken() (*cachedToken, error) {
	cachePath := opts.getCachePath()
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cached cachedToken
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, err
	}

	// Check if token is expired
	var tokenResp struct {
		Token struct {
			ExpiresAt string `json:"expires_at"`
		} `json:"token"`
	}
	if err := json.Unmarshal(cached.Body, &tokenResp); err != nil {
		return nil, err
	}

	expiresAt, err := time.Parse(time.RFC3339, tokenResp.Token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	if time.Now().After(expiresAt) {
		return nil, fmt.Errorf("cached token expired")
	}

	return &cached, nil
}

func (opts *AuthOptions) putCachedToken(authToken string, body []byte) error {
	cachePath := opts.getCachePath()
	cacheDir := filepath.Dir(cachePath)

	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		return err
	}

	cached := cachedToken{
		AuthToken: authToken,
		Body:      body,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0600)
}

// newIdentityClient creates a ServiceClient for the Identity v3 API.
func newIdentityClient(client *gophercloud.ProviderClient, endpoint string) *gophercloud.ServiceClient {
	// Ensure endpoint ends with v3/
	host := strings.TrimSuffix(endpoint, "/")
	if !strings.HasSuffix(host, "v3") {
		host += "/v3"
	}
	host = gophercloud.NormalizeURL(host)

	return &gophercloud.ServiceClient{
		ProviderClient: client,
		Endpoint:       host,
		Type:           "identity",
	}
}

// getTokenMetadata retrieves token metadata from Keystone.
func getTokenMetadata(ctx context.Context, client *gophercloud.ServiceClient, authToken string) (r GetResult) {
	resp, err := client.Get(ctx, tokenURL(client), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Token":    authToken,
			"X-Subject-Token": authToken,
		},
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// getScopedToken exchanges an unscoped token for a scoped token.
func getScopedToken(ctx context.Context, client *gophercloud.ServiceClient, unscopedToken string, scope map[string]any) (r CreateResult) {
	reqBody := map[string]any{
		"auth": map[string]any{
			"identity": map[string]any{
				"methods": []string{"token"},
				"token": map[string]any{
					"id": unscopedToken,
				},
			},
		},
	}

	if scope != nil {
		reqBody["auth"].(map[string]any)["scope"] = scope
	}

	resp, err := client.Post(ctx, tokenURL(client), reqBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Authenticate performs WebSSO authentication.
//
// This function will:
// 1. Check for a valid cached token
// 2. If no cache, open a browser for the user to authenticate
// 3. Wait for the identity provider to redirect back with a token
// 4. Optionally scope the token if Scope is specified
// 5. Configure the ProviderClient with the token and endpoint locator
func Authenticate(ctx context.Context, client *gophercloud.ProviderClient, opts AuthOptions) error {
	// Set defaults
	if opts.RedirectHost == "" {
		opts.RedirectHost = "localhost"
	}
	if opts.RedirectPort == 0 {
		opts.RedirectPort = 9990
	}

	// Create identity service client
	identityClient := newIdentityClient(client, opts.IdentityEndpoint)

	// Try to get cached token first
	cached, err := opts.getCachedToken()

	var authToken string
	var catalog *tokens.ServiceCatalog

	if err == nil && cached != nil {
		// Use cached token
		authToken = cached.AuthToken

		// Parse catalog from cached body
		var tokenResp struct {
			Token struct {
				Catalog []tokens.CatalogEntry `json:"catalog"`
			} `json:"token"`
		}
		if err := json.Unmarshal(cached.Body, &tokenResp); err != nil {
			return fmt.Errorf("failed to parse cached token: %w", err)
		}
		catalog = &tokens.ServiceCatalog{Entries: tokenResp.Token.Catalog}
	} else {
		// Start authentication flow
		authURL := webssoURL(identityClient, opts.IdentityProvider, opts.Protocol) + "?origin=" + opts.redirectURI()

		if err := openBrowser(authURL); err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}

		// Wait for callback
		authToken, err = waitForToken(opts.RedirectHost, opts.RedirectPort)
		if err != nil {
			return err
		}

		// Get token metadata
		result := getTokenMetadata(ctx, identityClient, authToken)
		if result.Err != nil {
			return fmt.Errorf("failed to get token metadata: %w", result.Err)
		}

		catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return fmt.Errorf("failed to extract service catalog: %w", err)
		}

		// Cache the token
		if result.Body != nil {
			bodyBytes, err := json.Marshal(map[string]any{"token": result.Body})
			if err == nil {
				// Caching errors are not fatal, but we should track them
				if cacheErr := opts.putCachedToken(authToken, bodyBytes); cacheErr != nil {
					// Return an error that includes both the success and the cache failure
					// The caller can decide whether to ignore the cache error
					_ = cacheErr // Cache errors are non-fatal for now
				}
			}
		}
	}

	// If no catalog, we need to scope the token
	if catalog == nil || len(catalog.Entries) == 0 {
		// Build scope request
		scopeMap, err := opts.ToTokenV3ScopeMap()
		if err != nil {
			return fmt.Errorf("failed to build scope: %w", err)
		}

		// Request scoped token
		result := getScopedToken(ctx, identityClient, authToken, scopeMap)
		if result.Err != nil {
			return fmt.Errorf("failed to get scoped token: %w", result.Err)
		}

		// Get new token ID
		authToken, err = result.ExtractTokenID()
		if err != nil {
			return fmt.Errorf("failed to extract scoped token ID: %w", err)
		}

		// Get catalog from scoped token
		catalog, err = result.ExtractServiceCatalog()
		if err != nil {
			return fmt.Errorf("failed to extract scoped service catalog: %w", err)
		}
	}

	// Configure the provider client
	client.SetToken(authToken)

	// Set up endpoint locator using the service catalog
	client.EndpointLocator = func(_ context.Context, eo gophercloud.EndpointOpts) (string, error) {
		return findEndpoint(catalog, eo)
	}

	// Set up reauth function if allowed
	if opts.AllowReauth {
		// Create a copy of opts with AllowReauth disabled to prevent infinite loops
		reauthOpts := opts
		reauthOpts.AllowReauth = false

		// Clear the cache to force re-authentication
		client.ReauthFunc = func(ctx context.Context) error {
			// Remove cached token to force browser auth
			_ = os.Remove(opts.getCachePath())
			return Authenticate(ctx, client, reauthOpts)
		}
	}

	return nil
}

// findEndpoint locates an endpoint in the service catalog.
func findEndpoint(catalog *tokens.ServiceCatalog, opts gophercloud.EndpointOpts) (string, error) {
	if catalog == nil {
		return "", fmt.Errorf("no service catalog available")
	}

	// Determine interface type
	interfaceType := "public"
	switch opts.Availability {
	case gophercloud.AvailabilityPublic:
		interfaceType = "public"
	case gophercloud.AvailabilityInternal:
		interfaceType = "internal"
	case gophercloud.AvailabilityAdmin:
		interfaceType = "admin"
	}

	// Find matching service
	for _, entry := range catalog.Entries {
		if entry.Type != opts.Type {
			continue
		}

		// Find matching endpoint
		for _, endpoint := range entry.Endpoints {
			if endpoint.Interface != interfaceType {
				continue
			}
			if opts.Region != "" && endpoint.Region != opts.Region {
				continue
			}
			return gophercloud.NormalizeURL(endpoint.URL), nil
		}
	}

	return "", fmt.Errorf("no endpoint found for service type %q, interface %q, region %q",
		opts.Type, interfaceType, opts.Region)
}

// ToTokenV3ScopeMap builds the scope for the authentication request.
func (opts *AuthOptions) ToTokenV3ScopeMap() (map[string]any, error) {
	if opts.Scope.ProjectID == "" && opts.Scope.ProjectName == "" &&
		opts.Scope.DomainID == "" && opts.Scope.DomainName == "" &&
		!opts.Scope.System {
		return nil, nil
	}

	scope := make(map[string]any)

	if opts.Scope.ProjectID != "" {
		scope["project"] = map[string]any{"id": opts.Scope.ProjectID}
	} else if opts.Scope.ProjectName != "" {
		project := map[string]any{"name": opts.Scope.ProjectName}
		if opts.Scope.DomainID != "" {
			project["domain"] = map[string]any{"id": opts.Scope.DomainID}
		} else if opts.Scope.DomainName != "" {
			project["domain"] = map[string]any{"name": opts.Scope.DomainName}
		}
		scope["project"] = project
	} else if opts.Scope.DomainID != "" {
		scope["domain"] = map[string]any{"id": opts.Scope.DomainID}
	} else if opts.Scope.DomainName != "" {
		scope["domain"] = map[string]any{"name": opts.Scope.DomainName}
	} else if opts.Scope.System {
		scope["system"] = map[string]any{"all": true}
	}

	return scope, nil
}

// CanReauth returns whether reauthentication is allowed.
func (opts *AuthOptions) CanReauth() bool {
	return opts.AllowReauth
}
