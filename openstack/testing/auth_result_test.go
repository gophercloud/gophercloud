package testing

import (
	"context"
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

func (f *fakeAuthOptionsBuilder) Authenticate(ctx context.Context, httpClient *http.Client) (*auth.AuthResult, error) {
	f.calledCount++
	return &f.authResult, f.authErr
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

func TestWithReauthLeavesReauthFuncNilWhenCanReauthFalse(t *testing.T) {
	fake := &fakeAuthOptionsBuilder{}
	client, err := openstack.NewClientFromAuthResult(sampleAuthResult(), openstack.WithReauth(fake, auth.AuthResult{CanReauth: false}))
	th.AssertNoErr(t, err)

	if client.ReauthFunc != nil {
		t.Fatal("expected ReauthFunc to stay nil when CanReauth is false")
	}
}
