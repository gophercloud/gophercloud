package vips

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

	th.AssertEquals(t, th.Endpoint()+"v2.0/lb/vips", rootURL(serviceClient()))
	th.AssertEquals(t, th.Endpoint()+"v2.0/lb/vips/foo", resourceURL(serviceClient(), "foo"))
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/lb/vips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "vips":[
         {
           "id": "db902c0c-d5ff-4753-b465-668ad9656918",
           "tenant_id": "310df60f-2a10-4ee5-9554-98393092194c",
           "name": "web_vip",
           "description": "lb config for the web tier",
           "subnet_id": "96a4386a-f8c3-42ed-afce-d7954eee77b3",
           "address" : "10.30.176.47",
           "port_id" : "cd1f7a47-4fa6-449c-9ee7-632838aedfea",
           "protocol": "HTTP",
           "protocol_port": 80,
           "pool_id" : "cfc6589d-f949-4c66-99d2-c2da56ef3764",
           "admin_state_up": true,
           "status": "ACTIVE"
         },
         {
           "id": "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
           "tenant_id": "310df60f-2a10-4ee5-9554-98393092194c",
           "name": "db_vip",
					 "description": "lb config for the db tier",
           "subnet_id": "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
           "address" : "10.30.176.48",
           "port_id" : "cd1f7a47-4fa6-449c-9ee7-632838aedfea",
           "protocol": "TCP",
           "protocol_port": 3306,
           "pool_id" : "41efe233-7591-43c5-9cf7-923964759f9e",
           "session_persistence" : {"type" : "SOURCE_IP"},
           "connection_limit" : 2000,
           "admin_state_up": true,
           "status": "INACTIVE"
         }
      ]
}
			`)
	})

	count := 0

	List(serviceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractVIPs(page)
		if err != nil {
			t.Errorf("Failed to extract LBs: %v", err)
			return false, err
		}

		expected := []VirtualIP{
			VirtualIP{
				ID:           "db902c0c-d5ff-4753-b465-668ad9656918",
				TenantID:     "310df60f-2a10-4ee5-9554-98393092194c",
				Name:         "web_vip",
				Description:  "lb config for the web tier",
				SubnetID:     "96a4386a-f8c3-42ed-afce-d7954eee77b3",
				Address:      "10.30.176.47",
				PortID:       "cd1f7a47-4fa6-449c-9ee7-632838aedfea",
				Protocol:     "HTTP",
				ProtocolPort: 80,
				PoolID:       "cfc6589d-f949-4c66-99d2-c2da56ef3764",
				//Persistence:  SessionPersistence{},
				ConnLimit:    0,
				AdminStateUp: true,
				Status:       "ACTIVE",
			},
			VirtualIP{
				ID:           "36e08a3e-a78f-4b40-a229-1e7e23eee1ab",
				TenantID:     "310df60f-2a10-4ee5-9554-98393092194c",
				Name:         "db_vip",
				Description:  "lb config for the db tier",
				SubnetID:     "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
				Address:      "10.30.176.48",
				PortID:       "cd1f7a47-4fa6-449c-9ee7-632838aedfea",
				Protocol:     "TCP",
				ProtocolPort: 3306,
				PoolID:       "41efe233-7591-43c5-9cf7-923964759f9e",
				//Persistence:  SessionPersistence{Type: "SOURCE_IP"},
				ConnLimit:    2000,
				AdminStateUp: true,
				Status:       "INACTIVE",
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

	th.Mux.HandleFunc("/v2.0/lb/vips", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "vip": {
        "protocol": "HTTP",
        "name": "NewVip",
        "admin_state_up": true,
        "subnet_id": "8032909d-47a1-4715-90af-5153ffe39861",
        "pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
        "protocol_port": 80
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "vip": {
        "status": "PENDING_CREATE",
        "protocol": "HTTP",
        "description": "",
        "admin_state_up": true,
        "subnet_id": "8032909d-47a1-4715-90af-5153ffe39861",
        "tenant_id": "83657cfcdfe44cd5920adaf26c48ceea",
        "connection_limit": -1,
        "pool_id": "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
        "address": "10.0.0.11",
        "protocol_port": 80,
        "port_id": "f7e6fe6a-b8b5-43a8-8215-73456b32e0f5",
        "id": "c987d2be-9a3c-4ac9-a046-e8716b1350e2",
        "name": "NewVip"
    }
}
		`)
	})

	iTrue := true
	opts := CreateOpts{
		Protocol:     "HTTP",
		Name:         "NewVip",
		AdminStateUp: &iTrue,
		SubnetID:     "8032909d-47a1-4715-90af-5153ffe39861",
		PoolID:       "61b1f87a-7a21-4ad3-9dda-7f81d249944f",
		ProtocolPort: 80,
	}

	r, err := Create(serviceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "PENDING_CREATE", r.Status)
	th.AssertEquals(t, "HTTP", r.Protocol)
	th.AssertEquals(t, "", r.Description)
	th.AssertEquals(t, true, r.AdminStateUp)
	th.AssertEquals(t, "8032909d-47a1-4715-90af-5153ffe39861", r.SubnetID)
	th.AssertEquals(t, "83657cfcdfe44cd5920adaf26c48ceea", r.TenantID)
	th.AssertEquals(t, -1, r.ConnLimit)
	th.AssertEquals(t, "61b1f87a-7a21-4ad3-9dda-7f81d249944f", r.PoolID)
	th.AssertEquals(t, "10.0.0.11", r.Address)
	th.AssertEquals(t, 80, r.ProtocolPort)
	th.AssertEquals(t, "f7e6fe6a-b8b5-43a8-8215-73456b32e0f5", r.PortID)
	th.AssertEquals(t, "c987d2be-9a3c-4ac9-a046-e8716b1350e2", r.ID)
	th.AssertEquals(t, "NewVip", r.Name)
}
