package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	"github.com/gophercloud/gophercloud/v2/openstack"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func sampleAuthResult() auth.AuthResult {
	return auth.AuthResult{
		TokenID:   "cached-token-id",
		ExpiresAt: time.Now().Add(time.Hour),
		Catalog: auth.ServiceCatalog{
			Entries: []auth.CatalogEntry{
				{
					ID: "1", Name: "nova", Type: "compute",
					Endpoints: []auth.Endpoint{
						{ID: "1", Interface: "public", Region: "RegionOne", URL: "http://example.com/compute"},
					},
				},
			},
		},
	}
}

func TestNewClientFromAuthResultSetsTokenAndEndpointLocator(t *testing.T) {
	client, err := openstack.NewClientFromAuthResult(sampleAuthResult())
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "cached-token-id", client.TokenID)

	endpoint, err := client.EndpointLocator(context.TODO(), gophercloud.EndpointOpts{Type: "compute"})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "http://example.com/compute/", endpoint)

	if client.ReauthFunc != nil {
		t.Fatal("expected ReauthFunc to be nil without WithReauth")
	}
}

type fakeAuthOptionsBuilder struct {
	authResult  auth.AuthResult
	authErr     error
	calledCount int
}

func (f *fakeAuthOptionsBuilder) Authenticate(ctx context.Context, client *gophercloud.ProviderClient) (*auth.AuthResult, error) {
	f.calledCount++
	return &f.authResult, f.authErr
}

func (f *fakeAuthOptionsBuilder) GetAuthURL() string {
	return "http://example.com/identity"
}

func TestWithReauthWiresReauthFuncWhenCanReauthTrue(t *testing.T) {
	fake := &fakeAuthOptionsBuilder{authResult: auth.AuthResult{TokenID: "fresh-token"}}
	client, err := openstack.NewClientFromAuthResult(sampleAuthResult(), openstack.WithReauth(fake, auth.AuthResult{CanReauth: true}))
	th.AssertNoErr(t, err)

	if client.ReauthFunc == nil {
		t.Fatal("expected ReauthFunc to be wired when CanReauth is true")
	}
	th.AssertNoErr(t, client.ReauthFunc(context.TODO()))
	th.CheckEquals(t, "fresh-token", client.TokenID)
	th.CheckEquals(t, 1, fake.calledCount)
}

func TestWithReauthReauthenticatesWithoutWaitingOnItsOwnLock(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, http.MethodPost)
		th.TestHeaderUnset(t, r, "X-Auth-Token")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"access": {
				"token": {"id": "fresh-token", "expires": "2026-09-01T12:00:00Z"},
				"user": {"id": "user-id", "name": "testuser", "roles": []},
				"serviceCatalog": []
			}
		}`)
	})

	options := auth.AuthOptionsV2{
		AuthURL: fakeServer.Endpoint(),
		Auth: auth.V2PasswordOpts{
			Username:    "testuser",
			Password:    "testpass",
			AllowReauth: true,
		},
	}
	client, err := openstack.NewClientFromAuthResult(
		sampleAuthResult(),
		openstack.WithReauth(options, auth.AuthResult{CanReauth: true}),
	)
	th.AssertNoErr(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	th.AssertNoErr(t, client.Reauthenticate(ctx, client.Token()))
	th.CheckEquals(t, "fresh-token", client.Token())
}

func TestWithReauthLeavesReauthFuncNilWhenCanReauthFalse(t *testing.T) {
	fake := &fakeAuthOptionsBuilder{}
	client, err := openstack.NewClientFromAuthResult(sampleAuthResult(), openstack.WithReauth(fake, auth.AuthResult{CanReauth: false}))
	th.AssertNoErr(t, err)

	if client.ReauthFunc != nil {
		t.Fatal("expected ReauthFunc to stay nil when CanReauth is false")
	}
}
