package floatingips

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "123"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: th.Endpoint(),
	}
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/floatingips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "floatingip": {
        "floating_network_id": "376da547-b977-4cfe-9cba-275c80debf57",
        "port_id": "ce705c24-c1ef-408a-bda3-7bbd946164ab"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "floatingip": {
        "router_id": "d23abc8d-2991-4a55-ba98-2aaea84cc72f",
        "tenant_id": "4969c491a3c74ee4af974e6d800c62de",
        "floating_network_id": "376da547-b977-4cfe-9cba-275c80debf57",
        "fixed_ip_address": "10.0.0.3",
        "floating_ip_address": "",
        "port_id": "ce705c24-c1ef-408a-bda3-7bbd946164ab",
        "id": "2f245a7b-796b-4f26-9cf9-9e82d248fda7"
    }
}
		`)
	})

	options := CreateOpts{
		FloatingNetworkID: "376da547-b977-4cfe-9cba-275c80debf57",
		PortID:            "ce705c24-c1ef-408a-bda3-7bbd946164ab",
	}

	ip, err := Create(serviceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "2f245a7b-796b-4f26-9cf9-9e82d248fda7", ip.ID)
	th.AssertEquals(t, "4969c491a3c74ee4af974e6d800c62de", ip.TenantID)
	th.AssertEquals(t, "376da547-b977-4cfe-9cba-275c80debf57", ip.FloatingNetworkID)
	th.AssertEquals(t, "", ip.FloatingIP)
	th.AssertEquals(t, "ce705c24-c1ef-408a-bda3-7bbd946164ab", ip.PortID)
	th.AssertEquals(t, "10.0.0.3", ip.FixedIP)
}
