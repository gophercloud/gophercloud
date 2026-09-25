package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthenticateDoesNotUseProviderToken(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	var requests atomic.Int32
	fakeServer.Mux.HandleFunc("/v3/auth/tokens", func(w http.ResponseWriter, r *http.Request) {
		n := requests.Add(1)
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		th.CheckEquals(t, "test-agent "+gophercloud.DefaultUserAgent, r.Header.Get("User-Agent"))

		var body struct {
			Auth struct {
				Identity struct {
					Password struct {
						User struct {
							Password string `json:"password"`
						} `json:"user"`
					} `json:"password"`
				} `json:"identity"`
			} `json:"auth"`
		}
		th.CheckNoErr(t, json.NewDecoder(r.Body).Decode(&body))
		if body.Auth.Identity.Password.User.Password != "valid" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("X-Subject-Token", fmt.Sprintf("token-%d", n))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"token":{"expires_at":"2099-01-01T00:00:00Z"}}`)
	})

	provider, err := openstack.NewClient(fakeServer.Endpoint() + "v3/")
	th.AssertNoErr(t, err)
	provider.UserAgent.Prepend("test-agent")
	options := func(password string) auth.AuthOptionsV3 {
		return auth.AuthOptionsV3{
			AuthURL: fakeServer.Endpoint(),
			Auth:    auth.V3PasswordOpts{UserID: "user", Password: password, AllowReauth: true},
		}
	}

	th.AssertNoErr(t, openstack.Authenticate(context.Background(), provider, options("valid")))
	th.AssertEquals(t, "token-1", provider.Token())

	// Rejected credentials must not reauthenticate with the previous ones.
	err = openstack.Authenticate(context.Background(), provider, options("invalid"))
	if !gophercloud.ResponseCodeIs(err, http.StatusUnauthorized) {
		t.Fatalf("expected a 401 response error, got %v", err)
	}
	th.AssertEquals(t, int32(2), requests.Load())
	th.AssertEquals(t, "token-1", provider.Token())
}
