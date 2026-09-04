package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/auth"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAuthResultToken(t *testing.T) {
	r := auth.AuthResult{TokenID: "abc123"}
	th.AssertEquals(t, "abc123", r.Token())
}

func TestAuthResultAuthenticatedHeaders(t *testing.T) {
	r := auth.AuthResult{TokenID: "abc123"}
	th.AssertDeepEquals(t, map[string]string{"X-Auth-Token": "abc123"}, r.AuthenticatedHeaders())
}

func TestAuthResultExpired(t *testing.T) {
	past := auth.AuthResult{ExpiresAt: time.Now().Add(-time.Hour)}
	th.AssertEquals(t, true, past.Expired())

	future := auth.AuthResult{ExpiresAt: time.Now().Add(time.Hour)}
	th.AssertEquals(t, false, future.Expired())
}

func TestAuthResultWillExpireBy(t *testing.T) {
	r := auth.AuthResult{ExpiresAt: time.Now().Add(30 * time.Minute)}
	th.AssertEquals(t, true, r.WillExpireBy(time.Hour))
	th.AssertEquals(t, false, r.WillExpireBy(time.Minute))
}

func testCatalog() auth.ServiceCatalog {
	return auth.ServiceCatalog{
		Entries: []auth.CatalogEntry{
			{
				Type: "compute",
				Name: "nova",
				Endpoints: []auth.Endpoint{
					{Interface: "public", Region: "RegionOne", URL: "http://public.example.com/compute"},
					{Interface: "internal", Region: "RegionOne", URL: "http://internal.example.com/compute"},
					{Interface: "public", Region: "RegionTwo", URL: "http://public2.example.com/compute"},
				},
			},
		},
	}
}

func TestAuthResultEndpointMatchesTypeAndAvailability(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	url, err := r.Endpoint(gophercloud.EndpointOpts{
		Type:         "compute",
		Availability: gophercloud.AvailabilityInternal,
		Region:       "RegionOne",
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://internal.example.com/compute/", url)
}

func TestAuthResultEndpointDefaultsToPublic(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	url, err := r.Endpoint(gophercloud.EndpointOpts{Type: "compute", Region: "RegionTwo"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://public2.example.com/compute/", url)
}

func TestAuthResultEndpointNotFound(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	_, err := r.Endpoint(gophercloud.EndpointOpts{Type: "does-not-exist"})
	th.AssertErr(t, err)

	_, ok := err.(*gophercloud.ErrEndpointNotFound)
	th.AssertEquals(t, true, ok)
}

func TestAuthResultEndpointLocator(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	locator := r.EndpointLocator()
	url, err := locator(context.TODO(), gophercloud.EndpointOpts{Type: "compute", Region: "RegionOne"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://public.example.com/compute/", url)
}
