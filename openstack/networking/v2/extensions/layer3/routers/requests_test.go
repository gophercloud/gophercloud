package routers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "123"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: th.Endpoint(),
	}
}

func TestURLs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.AssertEquals(t, th.Endpoint()+"v2.0/routers", rootURL(serviceClient()))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "routers": [
        {
            "status": "ACTIVE",
            "external_gateway_info": null,
            "name": "second_routers",
            "admin_state_up": true,
            "tenant_id": "6b96ff0cb17a4b859e1e575d221683d3",
            "id": "7177abc4-5ae9-4bb7-b0d4-89e94a4abf3b"
        },
        {
            "status": "ACTIVE",
            "external_gateway_info": {
                "network_id": "3c5bcddd-6af9-4e6b-9c3e-c153e521cab8"
            },
            "name": "router1",
            "admin_state_up": true,
            "tenant_id": "33a40233088643acb66ff6eb0ebea679",
            "id": "a9254bdb-2613-4a13-ac4c-adc581fba50d"
        }
    ]
}
			`)
	})

	count := 0

	List(serviceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractRouters(page)
		if err != nil {
			t.Errorf("Failed to extract routers: %v", err)
			return false, err
		}

		expected := []Router{
			Router{
				Status:         "ACTIVE",
				ExtGatewayInfo: GatewayInfo{NetworkID: ""},
				AdminStateUp:   true,
				Name:           "second_routers",
				ID:             "7177abc4-5ae9-4bb7-b0d4-89e94a4abf3b",
				TenantID:       "6b96ff0cb17a4b859e1e575d221683d3",
			},
			Router{
				Status:         "ACTIVE",
				ExtGatewayInfo: GatewayInfo{NetworkID: "3c5bcddd-6af9-4e6b-9c3e-c153e521cab8"},
				AdminStateUp:   true,
				Name:           "router1",
				ID:             "a9254bdb-2613-4a13-ac4c-adc581fba50d",
				TenantID:       "33a40233088643acb66ff6eb0ebea679",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {

}

func TestGet(t *testing.T) {

}

func TestUpdate(t *testing.T) {

}

func TestDelete(t *testing.T) {

}
