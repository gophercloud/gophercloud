package pools

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

	th.AssertEquals(t, th.Endpoint()+"v2.0/lb/pools", rootURL(serviceClient()))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lb/pools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "pools":[
      {
         "status":"ACTIVE",
         "lb_method":"ROUND_ROBIN",
         "protocol":"HTTP",
         "description":"",
         "health_monitors":[
            "466c8345-28d8-4f84-a246-e04380b0461d",
            "5d4b5228-33b0-4e60-b225-9b727c1a20e7"
         ],
         "members":[
            "701b531b-111a-4f21-ad85-4795b7b12af6",
            "beb53b4d-230b-4abd-8118-575b8fa006ef"
         ],
         "status_description": null,
         "id":"72741b06-df4d-4715-b142-276b6bce75ab",
         "vip_id":"4ec89087-d057-4e2c-911f-60a3b47ee304",
         "name":"app_pool",
         "admin_state_up":true,
         "subnet_id":"8032909d-47a1-4715-90af-5153ffe39861",
         "tenant_id":"83657cfcdfe44cd5920adaf26c48ceea",
         "health_monitors_status": [],
         "provider": "haproxy"
      }
   ]
}
			`)
	})

	count := 0

	List(serviceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractPools(page)
		if err != nil {
			t.Errorf("Failed to extract pools: %v", err)
			return false, err
		}

		expected := []Pool{
			Pool{
				Status:      "ACTIVE",
				LBMethod:    "ROUND_ROBIN",
				Protocol:    "HTTP",
				Description: "",
				MonitorIDs: []string{
					"466c8345-28d8-4f84-a246-e04380b0461d",
					"5d4b5228-33b0-4e60-b225-9b727c1a20e7",
				},
				SubnetID:     "8032909d-47a1-4715-90af-5153ffe39861",
				TenantID:     "83657cfcdfe44cd5920adaf26c48ceea",
				AdminStateUp: true,
				Name:         "app_pool",
				MemberIDs: []string{
					"701b531b-111a-4f21-ad85-4795b7b12af6",
					"beb53b4d-230b-4abd-8118-575b8fa006ef",
				},
				ID:    "72741b06-df4d-4715-b142-276b6bce75ab",
				VIPID: "4ec89087-d057-4e2c-911f-60a3b47ee304",
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lb/pools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "pool": {
        "lb_method": "ROUND_ROBIN",
        "protocol": "HTTP",
        "name": "Example pool",
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "pool": {
        "status": "PENDING_CREATE",
        "lb_method": "ROUND_ROBIN",
        "protocol": "HTTP",
        "description": "",
        "health_monitors": [],
        "members": [],
        "status_description": null,
        "id": "69055154-f603-4a28-8951-7cc2d9e54a9a",
        "vip_id": null,
        "name": "Example pool",
        "admin_state_up": true,
        "subnet_id": "1981f108-3c48-48d2-b908-30f7d28532c9",
        "tenant_id": "2ffc6e22aae24e4795f87155d24c896f",
        "health_monitors_status": []
    }
}
		`)
	})

	options := CreateOpts{
		LBMethod: LBMethodRoundRobin,
		Protocol: "HTTP",
		Name:     "Example pool",
		SubnetID: "1981f108-3c48-48d2-b908-30f7d28532c9",
		TenantID: "2ffc6e22aae24e4795f87155d24c896f",
	}
	p, err := Create(serviceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "PENDING_CREATE", p.Status)
	th.AssertEquals(t, "ROUND_ROBIN", p.LBMethod)
	th.AssertEquals(t, "HTTP", p.Protocol)
	th.AssertEquals(t, "", p.Description)
	th.AssertDeepEquals(t, []string{}, p.MonitorIDs)
	th.AssertDeepEquals(t, []string{}, p.MemberIDs)
	th.AssertEquals(t, "69055154-f603-4a28-8951-7cc2d9e54a9a", p.ID)
	th.AssertEquals(t, "Example pool", p.Name)
	th.AssertEquals(t, "1981f108-3c48-48d2-b908-30f7d28532c9", p.SubnetID)
	th.AssertEquals(t, "2ffc6e22aae24e4795f87155d24c896f", p.TenantID)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lb/pools/332abe93-f488-41ba-870b-2ac66be7f853", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "pool":{
      "id":"332abe93-f488-41ba-870b-2ac66be7f853",
      "tenant_id":"19eaa775-cf5d-49bc-902e-2f85f668d995",
      "name":"Example pool",
      "description":"",
      "protocol":"tcp",
      "lb_algorithm":"ROUND_ROBIN",
      "session_persistence":{
      },
      "healthmonitor_id":null,
      "members":[
      ],
      "admin_state_up":true,
      "status":"ACTIVE"
   }
}
			`)
	})

	n, err := Get(serviceClient(), "332abe93-f488-41ba-870b-2ac66be7f853").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.ID, "332abe93-f488-41ba-870b-2ac66be7f853")
}
