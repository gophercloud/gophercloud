package testing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// redirectTarget records the headers of every request it receives.
type redirectTarget struct {
	mu   sync.Mutex
	seen map[string]http.Header
}

func (t *redirectTarget) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.mu.Lock()
	t.seen[r.URL.Path] = r.Header.Clone()
	t.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{}`))
}

func (t *redirectTarget) header(path, name string) string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.seen[path].Get(name)
}

// localhostURL re-addresses an httptest server, which listens on 127.0.0.1,
// as localhost. net/http compares host names when it decides whether a
// redirect crosses hosts, so localhost -> 127.0.0.1 is a cross-origin hop
// while two servers addressed as 127.0.0.1 would not be.
func localhostURL(t *testing.T, server *httptest.Server) string {
	t.Helper()
	u, err := url.Parse(server.URL)
	th.AssertNoErr(t, err)
	return "http://localhost:" + u.Port() + "/"
}

func newRedirectingClient(t *testing.T, redirectTo func(*http.Request) string) (*gophercloud.ServiceClient, *httptest.Server) {
	t.Helper()
	origin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTo(r), http.StatusFound)
	}))
	endpoint := localhostURL(t, origin)
	provider := &gophercloud.ProviderClient{IdentityBase: endpoint, IdentityEndpoint: endpoint}
	provider.SetToken("TOKEN")
	return &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: endpoint}, origin
}

func TestTokenDoesNotFollowCrossOriginRedirect(t *testing.T) {
	target := &redirectTarget{seen: map[string]http.Header{}}
	targetServer := httptest.NewServer(target)
	defer targetServer.Close()

	client, origin := newRedirectingClient(t, func(r *http.Request) string {
		return targetServer.URL + "/landed" + r.URL.Path
	})
	defer origin.Close()

	var out map[string]any
	_, err := client.Get(context.Background(), client.Endpoint+"servers", &out, nil)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "", target.header("/landed/servers", "X-Auth-Token"))
	// Other headers still follow the redirect.
	th.AssertEquals(t, true, target.header("/landed/servers", "User-Agent") != "")
}

func TestTokenFollowsSameOriginRedirect(t *testing.T) {
	target := &redirectTarget{seen: map[string]http.Header{}}
	mux := http.NewServeMux()
	mux.Handle("/landed/", target)
	var self *httptest.Server
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, localhostURL(t, self)+"landed"+r.URL.Path, http.StatusFound)
	})
	self = httptest.NewServer(mux)
	defer self.Close()

	endpoint := localhostURL(t, self)
	provider := &gophercloud.ProviderClient{IdentityBase: endpoint, IdentityEndpoint: endpoint}
	provider.SetToken("TOKEN")
	client := &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: endpoint}

	var out map[string]any
	_, err := client.Get(context.Background(), endpoint+"servers", &out, nil)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "TOKEN", target.header("/landed/servers", "X-Auth-Token"))
}

// net/http rebuilds every redirected request from the original one, so a
// chain that leaves the origin and comes back must not regain the token.
func TestTokenStaysOffWhenChainReturnsToOrigin(t *testing.T) {
	target := &redirectTarget{seen: map[string]http.Header{}}
	var origin, other *httptest.Server
	mux := http.NewServeMux()
	mux.Handle("/landed/", target)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, other.URL+"/bounce"+r.URL.Path, http.StatusFound)
	})
	origin = httptest.NewServer(mux)
	defer origin.Close()
	other = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, localhostURL(t, origin)+"landed"+r.URL.Path[len("/bounce"):], http.StatusFound)
	}))
	defer other.Close()

	endpoint := localhostURL(t, origin)
	provider := &gophercloud.ProviderClient{IdentityBase: endpoint, IdentityEndpoint: endpoint}
	provider.SetToken("TOKEN")
	client := &gophercloud.ServiceClient{ProviderClient: provider, Endpoint: endpoint}

	var out map[string]any
	_, err := client.Get(context.Background(), endpoint+"servers", &out, nil)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", target.header("/landed/servers", "X-Auth-Token"))
}

func TestCallerRedirectPolicyIsKept(t *testing.T) {
	target := &redirectTarget{seen: map[string]http.Header{}}
	targetServer := httptest.NewServer(target)
	defer targetServer.Close()

	client, origin := newRedirectingClient(t, func(r *http.Request) string {
		return targetServer.URL + "/landed" + r.URL.Path
	})
	defer origin.Close()
	client.HTTPClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// The caller's policy stops at the 302, which the request then reports as its result.
	resp, err := client.Get(context.Background(), client.Endpoint+"servers", nil, &gophercloud.RequestOpts{OkCodes: []int{302}})
	th.AssertNoErr(t, err)
	resp.Body.Close()
	th.AssertEquals(t, http.StatusFound, resp.StatusCode)
	target.mu.Lock()
	_, reached := target.seen["/landed/servers"]
	target.mu.Unlock()
	th.AssertEquals(t, false, reached)
}
