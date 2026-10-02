package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	tests := []struct {
		name     string
		result   auth.AuthResult
		expected map[string]string
	}{
		{
			name:     "Keystone token",
			result:   auth.AuthResult{TokenID: "abc123"},
			expected: map[string]string{"X-Auth-Token": "abc123"},
		},
		{
			name:     "bearer token",
			result:   auth.AuthResult{TokenID: "abc123", TokenType: "Bearer"},
			expected: map[string]string{"X-Auth-Token": "abc123", "Authorization": "Bearer abc123"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			th.AssertDeepEquals(t, test.expected, test.result.AuthenticatedHeaders())
		})
	}
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
	url, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{
		Type:         "compute",
		Availability: gophercloud.AvailabilityInternal,
		Region:       "RegionOne",
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://internal.example.com/compute/", url)
}

func TestAuthResultEndpointDefaultsToPublic(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	url, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "compute", Region: "RegionTwo"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://public2.example.com/compute/", url)
}

func TestAuthResultEndpointNotFound(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	_, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "does-not-exist"})
	th.AssertErr(t, err)

	_, ok := err.(*gophercloud.ErrEndpointNotFound)
	th.AssertEquals(t, true, ok)
}

func TestAuthResultEndpointLocator(t *testing.T) {
	r := auth.AuthResult{Catalog: testCatalog()}
	locator := r.EndpointLocator(nil)
	url, err := locator(context.TODO(), gophercloud.EndpointOpts{Type: "compute", Region: "RegionOne"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "http://public.example.com/compute/", url)
}

func TestAuthResultEndpointMatchesServiceName(t *testing.T) {
	r := auth.AuthResult{Catalog: auth.ServiceCatalog{Entries: []auth.CatalogEntry{
		{Type: "compute", Name: "nova", Endpoints: []auth.Endpoint{{Interface: "public", URL: "https://nova.example.com"}}},
		{Type: "compute", Name: "other", Endpoints: []auth.Endpoint{{Interface: "public", URL: "https://other.example.com"}}},
	}}}

	url, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "compute", Name: "other"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "https://other.example.com/", url)
}

func TestAuthResultEndpointMatchesAliasAndRegionID(t *testing.T) {
	r := auth.AuthResult{Catalog: auth.ServiceCatalog{Entries: []auth.CatalogEntry{
		{Type: "volumev3", Endpoints: []auth.Endpoint{{Interface: "public", Region: "Display Region", RegionID: "region-id", URL: "https://volume.example.com"}}},
	}}}

	url, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "block-storage", Region: "region-id"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "https://volume.example.com/", url)
}

func TestAuthResultEndpointRejectsAliasWithDifferentVersion(t *testing.T) {
	r := auth.AuthResult{Catalog: auth.ServiceCatalog{Entries: []auth.CatalogEntry{
		{Type: "volumev3", Endpoints: []auth.Endpoint{{Interface: "public", URL: "https://volume.example.com"}}},
	}}}

	_, err := r.Endpoint(context.Background(), nil, gophercloud.EndpointOpts{Type: "block-storage", Version: 2})
	th.AssertErr(t, err)
	_, ok := err.(*gophercloud.ErrEndpointNotFound)
	th.AssertEquals(t, true, ok)
}

func TestAuthResultEndpointLocatorDiscoversServiceVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/first/":
			fmt.Fprint(w, `{"versions":[{"id":"v3.0","status":"CURRENT"}]}`)
		case "/second/":
			fmt.Fprint(w, `{"versions":[{"id":"v2.0","status":"SUPPORTED"}]}`)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	r := auth.AuthResult{Catalog: auth.ServiceCatalog{Entries: []auth.CatalogEntry{
		{Type: "block-storage", Endpoints: []auth.Endpoint{{Interface: "public", URL: server.URL + "/first/"}}},
		{Type: "block-storage", Endpoints: []auth.Endpoint{{Interface: "public", URL: server.URL + "/second/"}}},
	}}}
	provider := &gophercloud.ProviderClient{}

	url, err := r.EndpointLocator(provider)(context.Background(), gophercloud.EndpointOpts{Type: "block-storage", Version: 2})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, server.URL+"/second/", url)
}
