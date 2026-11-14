package testing

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/websso"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthenticateV3WebSSOReauth(t *testing.T) {
	fakeKeystone := th.SetupHTTP()
	defer fakeKeystone.Teardown()

	fakeKeystone.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodGet)
		token := r.Header.Get("X-Subject-Token")
		w.Header().Set("X-Subject-Token", token)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"token":{"expires_at":%q,"catalog":[]}}`, time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano))
	})

	client, err := openstack.NewClient(fakeKeystone.Endpoint() + "v3/")
	th.AssertNoErr(t, err)

	var browserCalls atomic.Int32
	opts := &websso.AuthOptions{
		IdentityProviderName: "my-idp",
		Protocol:             "openid",
		AllowReauth:          true,
		RedirectHost:         "127.0.0.1",
		RedirectPort:         availablePort(t),
		Timeout:              5 * time.Second,
		BrowserOpener: func(target string) error {
			token := fmt.Sprintf("websso-token-%d", browserCalls.Add(1))
			parsed, err := url.Parse(target)
			if err != nil {
				return err
			}
			form := url.Values{"token": {token}}
			request, err := http.NewRequest(http.MethodPost, parsed.Query().Get("origin"), strings.NewReader(form.Encode()))
			if err != nil {
				return err
			}
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			request.Header.Set("Sec-Fetch-Mode", "navigate")
			request.Header.Set("Sec-Fetch-Dest", "document")
			request.Header.Set("Sec-Fetch-Site", "cross-site")
			request.Header.Set("Origin", "null")
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return err
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected callback status: %s", response.Status)
			}
			return nil
		},
	}

	err = openstack.AuthenticateV3(context.Background(), client, opts, gophercloud.EndpointOpts{})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "websso-token-1", client.TokenID)
	if client.ReauthFunc == nil {
		t.Fatal("AuthenticateV3 did not configure reauthentication")
	}

	err = client.Reauthenticate(context.Background(), client.TokenID)
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "websso-token-2", client.TokenID)
	th.CheckEquals(t, int32(2), browserCalls.Load())
}

func availablePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	th.AssertNoErr(t, err)
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
