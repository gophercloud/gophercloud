package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/vpnservices", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "vpnservice": {
        "router_id": "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
        "name": "vpn",
        "admin_state_up": true,
		"description": "OpenStack VPN service",
		"tenant_id":  "10039663455a446d8ba2cbb058b0f578"
    }
}      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "vpnservice": {
        "router_id": "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
        "status": "PENDING_CREATE",
        "name": "vpn",
        "external_v6_ip": "2001:db8::1",
        "admin_state_up": true,
        "subnet_id": null,
        "tenant_id": "10039663455a446d8ba2cbb058b0f578",
        "external_v4_ip": "172.32.1.11",
        "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "description": "OpenStack VPN service"
    }
}
    `)
	})

	options := services.CreateOpts{
		TenantID:     "10039663455a446d8ba2cbb058b0f578",
		Name:         "vpn",
		Description:  "OpenStack VPN service",
		AdminStateUp: gophercloud.Enabled,
		RouterID:     "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
	}
	actual, err := services.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	expected := services.Service{
		RouterID:     "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
		Status:       "PENDING_CREATE",
		Name:         "vpn",
		ExternalV6IP: "2001:db8::1",
		AdminStateUp: true,
		SubnetID:     "",
		TenantID:     "10039663455a446d8ba2cbb058b0f578",
		ExternalV4IP: "172.32.1.11",
		ID:           "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		Description:  "OpenStack VPN service",
	}
	th.AssertDeepEquals(t, expected, *actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/vpnservices/5c561d9d-eaea-45f6-ae3e-08d1a7080828", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := services.Delete(fake.ServiceClient(), "5c561d9d-eaea-45f6-ae3e-08d1a7080828")
	th.AssertNoErr(t, res.Err)
}
