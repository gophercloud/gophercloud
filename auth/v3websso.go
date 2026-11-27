package auth

import (
	"context"
	"fmt"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth/tokencache"
)

const (
	webSSODefaultRedirectPort = 9990
	webSSODefaultTimeout      = 5 * time.Minute
	webSSOCallbackPath        = "/auth/websso/"
)

// V3WebSSOOpts authenticates through Keystone's WebSSO federation flow. It
// opens the identity provider's login page in a browser and receives the
// unscoped token on a loopback address, so it needs a browser on the same
// machine.
type V3WebSSOOpts struct {
	AuthURL string

	IdentityProviderName string `required:"true"`
	Protocol             string `required:"true"`

	AllowReauth bool
	Scope       *Scope

	// RedirectHost is the loopback host that receives the token. It defaults
	// to localhost.
	RedirectHost string

	// RedirectPort defaults to 9990.
	RedirectPort int

	// Timeout limits the wait for the browser login. It defaults to five
	// minutes.
	Timeout time.Duration

	// BrowserOpener opens the WebSSO URL. It defaults to the operating
	// system's browser.
	BrowserOpener func(string) error `json:"-"`

	// TokenCache stores the unscoped token, so that later authentications,
	// for any scope, do not open the browser. It requires CacheNamespace.
	TokenCache tokencache.Cache `json:"-"`

	// CacheNamespace separates the cached tokens of different browser
	// logins, such as two accounts at the same identity provider.
	CacheNamespace string
}

// WebSSOCacheKey returns the TokenCache key for opts at identityEndpoint. It
// does not depend on the scope.
func WebSSOCacheKey(identityEndpoint string, opts *V3WebSSOOpts) string {
	return tokencache.Key(tokencache.KeyOptions{
		Flow:             "websso",
		Principal:        opts.CacheNamespace,
		IdentityEndpoint: identityEndpoint,
		IdentityProvider: opts.IdentityProviderName,
		Protocol:         opts.Protocol,
	})
}

func (opts *V3WebSSOOpts) GetAuthURL() string {
	return (AuthOptionsV3{AuthURL: opts.AuthURL}).GetAuthURL()
}

func (opts *V3WebSSOOpts) Authenticate(ctx context.Context, provider *gophercloud.ProviderClient) (*AuthResult, error) {
	if opts == nil {
		return nil, fmt.Errorf("websso: expected non-nil options")
	}
	if opts.AuthURL == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "AuthURL"}
	}
	// Check everything before opening the browser.
	if err := opts.validate(); err != nil {
		return nil, err
	}
	scope, err := opts.Scope.ToScopeMap()
	if err != nil {
		return nil, err
	}

	if provider == nil {
		provider = &gophercloud.ProviderClient{}
	}
	client := &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: opts.GetAuthURL()}

	var result *AuthResult
	cacheKey := ""
	if opts.TokenCache != nil {
		cacheKey = WebSSOCacheKey(client.Endpoint, opts)
		result, err = cachedWebSSOToken(ctx, client, opts.TokenCache, cacheKey)
		if err != nil {
			return nil, err
		}
	}
	if result == nil {
		token, err := captureWebSSOToken(ctx, client, opts)
		if err != nil {
			return nil, err
		}
		result, err = webSSOTokenDetails(ctx, client, token)
		if err != nil {
			return nil, err
		}
		if opts.TokenCache != nil {
			tokencache.Persist(opts.TokenCache, cacheKey, client.Endpoint, result.TokenID, result.ExpiresAt)
		}
	}

	if scope != nil {
		rescope := AuthOptionsV3{AuthURL: client.Endpoint, Auth: V3TokenOpts{Token: result.TokenID, Scope: opts.Scope}}
		result, err = rescope.Authenticate(ctx, provider)
		if err != nil {
			return nil, err
		}
	}
	result.CanReauth = opts.AllowReauth
	return result, nil
}

// cachedWebSSOToken returns the details of the cached unscoped token, or nil
// when there is none or Keystone no longer accepts it.
func cachedWebSSOToken(ctx context.Context, client *gophercloud.ServiceClient, cache tokencache.Cache, key string) (*AuthResult, error) {
	token, ok := tokencache.Load(cache, key, client.Endpoint)
	if !ok {
		return nil, nil
	}
	result, err := webSSOTokenDetails(ctx, client, token)
	if gophercloud.ResponseCodeIs(err, http.StatusUnauthorized) || gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
		_ = cache.Delete(key)
		return nil, nil
	}
	return result, err
}

func (opts *V3WebSSOOpts) validate() error {
	if _, err := gophercloud.BuildRequestBody(opts, ""); err != nil {
		return err
	}
	if opts.RedirectPort < 0 || opts.RedirectPort > 65535 {
		return fmt.Errorf("websso: RedirectPort must be 0 or between 1 and 65535")
	}
	if opts.Timeout < 0 {
		return fmt.Errorf("websso: Timeout must not be negative")
	}
	if opts.TokenCache != nil && opts.CacheNamespace == "" {
		return fmt.Errorf("websso: CacheNamespace is required when TokenCache is enabled")
	}
	if opts.RedirectHost != "" && opts.RedirectHost != "localhost" {
		ip := net.ParseIP(opts.RedirectHost)
		if ip == nil || !ip.IsLoopback() {
			return fmt.Errorf("websso: RedirectHost must be a loopback address")
		}
	}
	return nil
}

// webSSOTokenDetails retrieves the token's details, authenticating with the
// token itself.
func webSSOTokenDetails(ctx context.Context, client *gophercloud.ServiceClient, tokenID string) (*AuthResult, error) {
	var result gophercloud.Result
	resp, err := client.Get(ctx, client.ServiceURL("auth", "tokens"), &result.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"X-Auth-Token": tokenID, "X-Subject-Token": tokenID},
		OkCodes:     []int{http.StatusOK, http.StatusNonAuthoritativeInfo},
	})
	_, result.Header, result.Err = gophercloud.ParseResponse(resp, err)
	if result.Err == nil {
		result.Header.Set("X-Subject-Token", tokenID)
	}
	return v3AuthResult(result, false)
}

// captureWebSSOToken opens the browser and waits for Keystone to post the
// unscoped token to the loopback callback.
func captureWebSSOToken(ctx context.Context, client *gophercloud.ServiceClient, opts *V3WebSSOOpts) (string, error) {
	keystoneOrigin, err := endpointOrigin(client.Endpoint)
	if err != nil {
		return "", fmt.Errorf("websso: invalid identity endpoint: %w", err)
	}
	host := opts.RedirectHost
	if host == "" {
		host = "localhost"
	}
	port := opts.RedirectPort
	if port == 0 {
		port = webSSODefaultRedirectPort
	}
	timeout := opts.Timeout
	if timeout == 0 {
		timeout = webSSODefaultTimeout
	}
	origin := "http://" + net.JoinHostPort(host, strconv.Itoa(port)) + webSSOCallbackPath
	webSSOURL := client.ServiceURL("auth", "OS-FEDERATION", "identity_providers", opts.IdentityProviderName, "protocols", opts.Protocol, "websso") +
		"?origin=" + url.QueryEscape(origin)

	listeners, err := listenLoopback(host, port)
	if err != nil {
		return "", err
	}

	tokenCh := make(chan string, 1)
	errCh := make(chan error, len(listeners))
	var accepted atomic.Bool
	mux := http.NewServeMux()
	mux.HandleFunc(webSSOCallbackPath, func(w http.ResponseWriter, request *http.Request) {
		if request.URL.Path != webSSOCallbackPath {
			http.NotFound(w, request)
			return
		}
		if request.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
		if err != nil || mediaType != "application/x-www-form-urlencoded" {
			http.Error(w, "Unsupported Content-Type", http.StatusUnsupportedMediaType)
			return
		}
		if !validCallbackContext(request, keystoneOrigin) {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		request.Body = http.MaxBytesReader(w, request.Body, 64<<10)
		if err := request.ParseForm(); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		token := request.PostFormValue("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}
		if !accepted.CompareAndSwap(false, true) {
			http.Error(w, "Token already received", http.StatusConflict)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "<html><body><script>window.close()</script><h2>Authentication successful</h2><p>You may close this window.</p></body></html>")
		tokenCh <- token
	})

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	for _, listener := range listeners {
		go func() {
			if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
				errCh <- fmt.Errorf("websso: callback server error: %w", err)
			}
		}()
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
		}
	}()

	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("websso: authentication cancelled: %w", err)
	}
	opener := opts.BrowserOpener
	if opener == nil {
		opener = openBrowser
	}
	if err := opener(webSSOURL); err != nil {
		return "", fmt.Errorf("websso: failed to open browser: %w", err)
	}

	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case token := <-tokenCh:
		return token, nil
	case err := <-errCh:
		return "", err
	case <-timer.C:
		return "", fmt.Errorf("websso: authentication timed out after %s", timeout)
	case <-ctx.Done():
		return "", fmt.Errorf("websso: authentication cancelled: %w", ctx.Err())
	}
}

// listenLoopback listens for the callback. Browsers may resolve localhost to
// either loopback address, so it listens on both; otherwise another local
// process could receive the token on the one left unused.
func listenLoopback(host string, port int) ([]net.Listener, error) {
	listen := func(host string) (net.Listener, error) {
		address := net.JoinHostPort(host, strconv.Itoa(port))
		listener, err := net.Listen("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("websso: failed to start callback server on %s: %w", address, err)
		}
		return listener, nil
	}

	if host != "localhost" {
		listener, err := listen(host)
		if err != nil {
			return nil, err
		}
		return []net.Listener{listener}, nil
	}

	ipv4, err := listen("127.0.0.1")
	if err != nil {
		return nil, err
	}
	ipv6, err := listen("::1")
	if err == nil {
		return []net.Listener{ipv4, ipv6}, nil
	}
	probe, probeErr := net.Listen("tcp", "[::1]:0")
	if probeErr != nil {
		// Without IPv6 loopback, browsers can only use 127.0.0.1.
		return []net.Listener{ipv4}, nil
	}
	_ = probe.Close()
	_ = ipv4.Close()
	return nil, err
}

// validCallbackContext reports whether the callback is a form submission
// that a page navigated to, and, when the browser says where from, that it
// came from Keystone.
func validCallbackContext(request *http.Request, keystoneOrigin string) bool {
	if request.Header.Get("Sec-Fetch-Mode") != "navigate" ||
		request.Header.Get("Sec-Fetch-Dest") != "document" ||
		request.Header.Get("Sec-Fetch-Site") == "none" {
		return false
	}

	origin := request.Header.Get("Origin")
	if origin != "" && origin != "null" && !sameOrigin(origin, keystoneOrigin) {
		return false
	}
	referer := request.Header.Get("Referer")
	return referer == "" || sameOrigin(referer, keystoneOrigin)
}

func sameOrigin(rawURL, expected string) bool {
	origin, err := endpointOrigin(rawURL)
	return err == nil && origin == expected
}

// endpointOrigin returns the URL's scheme, host and non-default port.
func endpointOrigin(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("url must include a scheme and host")
	}
	scheme := strings.ToLower(parsed.Scheme)
	host := strings.ToLower(parsed.Hostname())
	port := parsed.Port()
	if (scheme == "https" && port == "443") || (scheme == "http" && port == "80") {
		port = ""
	}
	if port != "" {
		host = net.JoinHostPort(host, port)
	} else if strings.Contains(host, ":") {
		host = "[" + host + "]"
	}
	return scheme + "://" + host, nil
}

func openBrowser(target string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", target)
	case "linux":
		command = exec.Command("xdg-open", target)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", target)
	default:
		return fmt.Errorf("unsupported operating system %q", runtime.GOOS)
	}
	if err := command.Start(); err != nil {
		return err
	}
	go func() {
		_ = command.Wait()
	}()
	return nil
}
