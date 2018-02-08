package testing

import (
	"fmt"
	"net/http"
	"testing"
	"github.com/gophercloud/gophercloud"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/pagination"
)


func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/vpnservices", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "vpnservices":[
        {
            "router_id": "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
            "status": "PENDING_CREATE",
            "name": "vpnservice1",
            "admin_state_up": true,
            "subnet_id": null,
            "project_id": "10039663455a446d8ba2cbb058b0f578",
            "tenant_id": "10039663455a446d8ba2cbb058b0f578",
            "description": "Test VPN service"
        }
   ]
}
      `)
	})

	count := 0

	services.List(fake.ServiceClient(), services.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := services.ExtractServices(page)
		if err != nil {
			t.Errorf("Failed to extract members: %v", err)
			return false, err
		}

		expected := []services.Service{
			{
				Status:       "PENDING_CREATE",
				Name:         "vpnservice1",
				AdminStateUp: gophercloud.Enabled,
				TenantID:     "10039663455a446d8ba2cbb058b0f578",
				Description:  "Test VPN service",
				SubnetID:     "",
				RouterID:     "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
				ProjectID:    "10039663455a446d8ba2cbb058b0f578",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/vpnservices/5c561d9d-eaea-45f6-ae3e-08d1a7080828", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "vpnservice": {
        "router_id": "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa",
        "status": "PENDING_CREATE",
        "name": "vpnservice1",
        "admin_state_up": true,
        "subnet_id": null,
        "project_id": "10039663455a446d8ba2cbb058b0f578",
        "tenant_id": "10039663455a446d8ba2cbb058b0f578",
        "id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "description": "VPN test service"
    }
}
        `)
	})

	serv, err := services.Get(fake.ServiceClient(), "5c561d9d-eaea-45f6-ae3e-08d1a7080828").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "PENDING_CREATE", serv.Status)
	th.AssertEquals(t, "vpnservice1", serv.Name)
	th.AssertEquals(t, "VPN test service", serv.Description)
	th.AssertEquals(t, true, *serv.AdminStateUp)
	th.AssertEquals(t, "10039663455a446d8ba2cbb058b0f578", serv.ProjectID)
	th.AssertEquals(t, "5c561d9d-eaea-45f6-ae3e-08d1a7080828", serv.ID)
	th.AssertEquals(t, "10039663455a446d8ba2cbb058b0f578", serv.TenantID)
	th.AssertEquals(t, "66e3b16c-8ce5-40fb-bb49-ab6d8dc3f2aa", serv.RouterID)
	th.AssertEquals(t, "", serv.SubnetID)

}