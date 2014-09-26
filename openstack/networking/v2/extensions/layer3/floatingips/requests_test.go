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

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/floatingips/2f245a7b-796b-4f26-9cf9-9e82d248fda7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "floatingip": {
        "floating_network_id": "90f742b1-6d17-487b-ba95-71881dbc0b64",
        "fixed_ip_address": "192.0.0.2",
        "floating_ip_address": "10.0.0.3",
        "tenant_id": "017d8de156df4177889f31a9bd6edc00",
        "status": "DOWN",
        "port_id": "74a342ce-8e07-4e91-880c-9f834b68fa25",
        "id": "2f245a7b-796b-4f26-9cf9-9e82d248fda7"
    }
}
      `)
	})

	ip, err := Get(serviceClient(), "2f245a7b-796b-4f26-9cf9-9e82d248fda7").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "90f742b1-6d17-487b-ba95-71881dbc0b64", ip.FloatingNetworkID)
	th.AssertEquals(t, "10.0.0.3", ip.FloatingIP)
	th.AssertEquals(t, "74a342ce-8e07-4e91-880c-9f834b68fa25", ip.PortID)
	th.AssertEquals(t, "192.0.0.2", ip.FixedIP)
	th.AssertEquals(t, "017d8de156df4177889f31a9bd6edc00", ip.TenantID)
	th.AssertEquals(t, "DOWN", ip.Status)
	th.AssertEquals(t, "2f245a7b-796b-4f26-9cf9-9e82d248fda7", ip.ID)
}

func TestAssociate(t *testing.T) {

}

func TestDisassociate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/floatingips/2f245a7b-796b-4f26-9cf9-9e82d248fda7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "floatingip": {
      "port_id": null
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "floatingip": {
        "router_id": "d23abc8d-2991-4a55-ba98-2aaea84cc72f",
        "tenant_id": "4969c491a3c74ee4af974e6d800c62de",
        "floating_network_id": "376da547-b977-4cfe-9cba-275c80debf57",
        "fixed_ip_address": null,
        "floating_ip_address": "172.24.4.228",
        "port_id": null,
        "id": "2f245a7b-796b-4f26-9cf9-9e82d248fda7"
    }
}
    `)
	})

	ip, err := Update(serviceClient(), "2f245a7b-796b-4f26-9cf9-9e82d248fda7", UpdateOpts{}).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, "", ip.FixedIP)
	th.AssertDeepEquals(t, "", ip.PortID)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/floatingips/2f245a7b-796b-4f26-9cf9-9e82d248fda7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := Delete(serviceClient(), "2f245a7b-796b-4f26-9cf9-9e82d248fda7")
	th.AssertNoErr(t, res.Err)
}
