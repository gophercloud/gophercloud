package loadbalancers

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/rackspace/gophercloud/openstack/networking/v2/common"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestURLs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.AssertEquals(t, th.Endpoint()+"v2.0/lbaas/loadbalancers", rootURL(fake.ServiceClient()))
	th.AssertEquals(t, th.Endpoint()+"v2.0/lbaas/loadbalancers/foo", resourceURL(fake.ServiceClient(), "foo"))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "loadbalancers":[
         {
           "id": "c331058c-6a40-4144-948e-b9fb1df9db4b",
           "tenant_id": "54030507-44f7-473c-9342-b4d14a95f692",
           "name": "web_lb",
           "description": "lb config for the web tier",
           "vip_subnet_id": "8a49c438-848f-467b-9655-ea1548708154",
           "vip_address": "10.30.176.47",
           "flavor": "small",
           "provider": "provider_1",
           "admin_state_up": true,
           "provisioning_status": "ACTIVE",
           "operating_status": "ONLINE"
         },
         {
           "id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
           "tenant_id": "54030507-44f7-473c-9342-b4d14a95f692",
           "name": "db_lb",
	   "description": "lb config for the db tier",
           "vip_subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
           "vip_address": "10.30.176.48",
           "flavor": "medium",
           "provider": "provider_2",
           "admin_state_up": true,
           "provisioning_status": "PENDING_CREATE",
           "operating_status": "OFFLINE"
         }
      ]
}
			`)
	})

	count := 0

	List(fake.ServiceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractLoadbalancers(page)
		if err != nil {
			t.Errorf("Failed to extract LBs: %v", err)
			return false, err
		}

		expected := []LoadBalancer{
			{
				ID:                 "c331058c-6a40-4144-948e-b9fb1df9db4b",
				TenantID:           "54030507-44f7-473c-9342-b4d14a95f692",
				Name:               "web_lb",
				Description:        "lb config for the web tier",
				VipSubnetID:        "8a49c438-848f-467b-9655-ea1548708154",
				VipAddress:         "10.30.176.47",
				Flavor:             "small",
				Provider:           "provider_1",
				AdminStateUp:       true,
				ProvisioningStatus: "ACTIVE",
				OperatingStatus:    "ONLINE",
			},
			{
				ID:                 "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
				TenantID:           "54030507-44f7-473c-9342-b4d14a95f692",
				Name:               "db_lb",
				Description:        "lb config for the db tier",
				VipSubnetID:        "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
				VipAddress:         "10.30.176.48",
				Flavor:             "medium",
				Provider:           "provider_2",
				AdminStateUp:       true,
				ProvisioningStatus: "PENDING_CREATE",
				OperatingStatus:    "OFFLINE",
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

	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "loadbalancer": {
	"name": "NewLb",
	"vip_subnet_id": "8032909d-47a1-4715-90af-5153ffe39861",
	"vip_address": "10.0.0.11",
	"flavor": "small",
	"provider": "provider_1",
	"admin_state_up": true
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "loadbalancer": {
        "id": "04816630-0320-4f9d-a17b-402b4e145d91",
	"tenant_id": "54030507-44f7-473c-9342-b4d14a95f692",
	"name": "NewLb",
	"description": "",
	"vip_subnet_id": "8032909d-47a1-4715-90af-5153ffe39861",
	"vip_address": "10.0.0.11",
	"flavor": "small",
	"provider": "provider_1",
	"admin_state_up": true,
	"provisioning_status": "PENDING_CREATE",
	"operating_status": "OFFLINE"
    }
}
		`)
	})

	opts := CreateOpts{
		Name:         "NewLb",
		AdminStateUp: Up,
		VipSubnetID:  "8032909d-47a1-4715-90af-5153ffe39861",
		VipAddress:   "10.0.0.11",
		Flavor:       "small",
		Provider:     "provider_1",
	}

	r, err := Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "PENDING_CREATE", r.ProvisioningStatus)
	th.AssertEquals(t, "", r.Description)
	th.AssertEquals(t, true, r.AdminStateUp)
	th.AssertEquals(t, "8032909d-47a1-4715-90af-5153ffe39861", r.VipSubnetID)
	th.AssertEquals(t, "54030507-44f7-473c-9342-b4d14a95f692", r.TenantID)
	th.AssertEquals(t, "NewLb", r.Name)
	th.AssertEquals(t, "OFFLINE", r.OperatingStatus)
	th.AssertEquals(t, "small", r.Flavor)
	th.AssertEquals(t, "provider_1", r.Provider)
	th.AssertEquals(t, "10.0.0.11", r.VipAddress)
}

func TestRequiredCreateOpts(t *testing.T) {
	res := Create(fake.ServiceClient(), CreateOpts{})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo", Description: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
	res = Create(fake.ServiceClient(), CreateOpts{Name: "foo", Description: "bar", VipAddress: "bar"})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/3c073c6d-18f7-45e4-8921-84358b70542d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
   "loadbalancer": {
	"id": "3c073c6d-18f7-45e4-8921-84358b70542d",
	"tenant_id": "54030507-44f7-473c-9342-b4d14a95f692",
	"name": "web_lb",
	"description": "",
	"vip_subnet_id": "8a49c438-848f-467b-9655-ea1548708154",
	"vip_address": "10.30.176.47",
	"flavor": "small",
	"provider": "provider_1",
	"admin_state_up": true,
	"provisioning_status": "ACTIVE",
	"operating_status": "ONLINE"
	}
}
			`)
	})

	loadbalancer, err := Get(fake.ServiceClient(), "3c073c6d-18f7-45e4-8921-84358b70542d").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "ACTIVE", loadbalancer.ProvisioningStatus)
	th.AssertEquals(t, "ONLINE", loadbalancer.OperatingStatus)
	th.AssertEquals(t, "", loadbalancer.Description)
	th.AssertEquals(t, true, loadbalancer.AdminStateUp)
	th.AssertEquals(t, "web_lb", loadbalancer.Name)
	th.AssertEquals(t, "10.30.176.47", loadbalancer.VipAddress)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/038d98e2-6a7d-473a-9ce4-76f7e1e0c6a4", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "loadbalancer": {
        "name": "NewLbName"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, `
{
    "loadbalancer": {
    "id": "038d98e2-6a7d-473a-9ce4-76f7e1e0c6a4",
	"tenant_id": "54030507-44f7-473c-9342-b4d14a95f692",
	"name": "NewLbName",
	"description": "",
	"vip_subnet_id": "8032909d-47a1-4715-90af-5153ffe39861",
	"vip_address": "10.0.0.11",
	"flavor": "small",
	"provider": "provider_1",
	"admin_state_up": true,
	"provisioning_status": "PENDING_UPDATE",
	"operating_status": "OFFLINE"
    }
}
		`)
	})

	options := UpdateOpts{
		Name: "NewLbName",
	}

	loadbalancer, err := Update(fake.ServiceClient(), "038d98e2-6a7d-473a-9ce4-76f7e1e0c6a4", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "NewLbName", loadbalancer.Name)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lbaas/loadbalancers/82c54b9a-cf84-460f-bdcf-fa591e6fd6f0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := Delete(fake.ServiceClient(), "82c54b9a-cf84-460f-bdcf-fa591e6fd6f0")
	th.AssertNoErr(t, res.Err)
}
