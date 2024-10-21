package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

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
            "distributed": false,
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
            "distributed": false,
            "id": "a9254bdb-2613-4a13-ac4c-adc581fba50d"
        },
        {
            "status": "ACTIVE",
            "external_gateway_info": {
                "network_id": "2b37576e-b050-4891-8b20-e1e37a93942a",
                "external_fixed_ips": [
                    {"ip_address": "192.0.2.17", "subnet_id": "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"},
                    {"ip_address": "198.51.100.33", "subnet_id": "1d699529-bdfd-43f8-bcaa-bff00c547af2"}
                ],
                "qos_policy_id": "6601bae5-f15a-4687-8be9-ddec9a2f8a8b"
            },
            "name": "gateway",
            "admin_state_up": true,
            "tenant_id": "a3e881e0a6534880c5473d95b9442099",
            "distributed": false,
            "id": "308a035c-005d-4452-a9fe-6f8f2f0c28d8"
        }
    ]
}
			`)
	})

	count := 0

	err := routers.List(fake.ServiceClient(), routers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := routers.ExtractRouters(page)
		if err != nil {
			t.Errorf("Failed to extract routers: %v", err)
			return false, err
		}

		expected := []routers.Router{
			{
				Status:       "ACTIVE",
				GatewayInfo:  routers.GatewayInfo{NetworkID: ""},
				AdminStateUp: true,
				Distributed:  false,
				Name:         "second_routers",
				ID:           "7177abc4-5ae9-4bb7-b0d4-89e94a4abf3b",
				TenantID:     "6b96ff0cb17a4b859e1e575d221683d3",
			},
			{
				Status:       "ACTIVE",
				GatewayInfo:  routers.GatewayInfo{NetworkID: "3c5bcddd-6af9-4e6b-9c3e-c153e521cab8"},
				AdminStateUp: true,
				Distributed:  false,
				Name:         "router1",
				ID:           "a9254bdb-2613-4a13-ac4c-adc581fba50d",
				TenantID:     "33a40233088643acb66ff6eb0ebea679",
			},
			{
				Status: "ACTIVE",
				GatewayInfo: routers.GatewayInfo{
					NetworkID: "2b37576e-b050-4891-8b20-e1e37a93942a",
					ExternalFixedIPs: []routers.ExternalFixedIP{
						{IPAddress: "192.0.2.17", SubnetID: "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"},
						{IPAddress: "198.51.100.33", SubnetID: "1d699529-bdfd-43f8-bcaa-bff00c547af2"},
					},
					QoSPolicyID: "6601bae5-f15a-4687-8be9-ddec9a2f8a8b",
				},
				AdminStateUp: true,
				Distributed:  false,
				Name:         "gateway",
				ID:           "308a035c-005d-4452-a9fe-6f8f2f0c28d8",
				TenantID:     "a3e881e0a6534880c5473d95b9442099",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
   "router":{
      "name": "foo_router",
      "admin_state_up": false,
      "external_gateway_info":{
         "enable_snat": false,
         "network_id":"8ca37218-28ff-41cb-9b10-039601ea7e6b",
         "external_fixed_ips": [
             {"subnet_id": "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"}
         ],
         "qos_policy_id": "6601bae5-f15a-4687-8be9-ddec9a2f8a8b"
	  },
	  "availability_zone_hints": ["zone1", "zone2"]
   }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "router": {
        "status": "ACTIVE",
        "external_gateway_info": {
            "network_id": "8ca37218-28ff-41cb-9b10-039601ea7e6b",
            "enable_snat": false,
            "external_fixed_ips": [
                {"ip_address": "192.0.2.17", "subnet_id": "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"}
            ],
            "qos_policy_id": "6601bae5-f15a-4687-8be9-ddec9a2f8a8b"
        },
        "name": "foo_router",
        "admin_state_up": false,
        "tenant_id": "6b96ff0cb17a4b859e1e575d221683d3",
		"distributed": false,
		"availability_zone_hints": ["zone1", "zone2"],
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e"
    }
}
		`)
	})

	asu := false
	enableSNAT := false
	qosID := "6601bae5-f15a-4687-8be9-ddec9a2f8a8b"
	efi := []routers.ExternalFixedIP{
		{
			SubnetID: "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def",
		},
	}
	gwi := routers.GatewayInfo{
		NetworkID:        "8ca37218-28ff-41cb-9b10-039601ea7e6b",
		EnableSNAT:       &enableSNAT,
		ExternalFixedIPs: efi,
		QoSPolicyID:      qosID,
	}
	options := routers.CreateOpts{
		Name:                  "foo_router",
		AdminStateUp:          &asu,
		GatewayInfo:           &gwi,
		AvailabilityZoneHints: []string{"zone1", "zone2"},
	}
	r, err := routers.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	gwi.ExternalFixedIPs = []routers.ExternalFixedIP{{
		IPAddress: "192.0.2.17",
		SubnetID:  "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def",
	}}

	th.AssertEquals(t, "foo_router", r.Name)
	th.AssertEquals(t, false, r.AdminStateUp)
	th.AssertDeepEquals(t, gwi, r.GatewayInfo)
	th.AssertDeepEquals(t, []string{"zone1", "zone2"}, r.AvailabilityZoneHints)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/a07eea83-7710-4860-931b-5fe220fae533", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "status": "ACTIVE",
        "external_gateway_info": {
            "network_id": "85d76829-6415-48ff-9c63-5c5ca8c61ac6",
            "external_fixed_ips": [
                {"ip_address": "198.51.100.33", "subnet_id": "1d699529-bdfd-43f8-bcaa-bff00c547af2"}
            ],
            "qos_policy_id": "6601bae5-f15a-4687-8be9-ddec9a2f8a8b"
        },
        "routes": [
            {
                "nexthop": "10.1.0.10",
                "destination": "40.0.1.0/24"
            }
        ],
        "name": "router1",
        "admin_state_up": true,
        "tenant_id": "d6554fe62e2f41efbb6e026fad5c1542",
		"distributed": false,
		"availability_zone_hints": ["zone1", "zone2"],
        "id": "a07eea83-7710-4860-931b-5fe220fae533"
    }
}
			`)
	})

	n, err := routers.Get(context.TODO(), fake.ServiceClient(), "a07eea83-7710-4860-931b-5fe220fae533").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertDeepEquals(t, n.GatewayInfo, routers.GatewayInfo{
		NetworkID: "85d76829-6415-48ff-9c63-5c5ca8c61ac6",
		ExternalFixedIPs: []routers.ExternalFixedIP{
			{IPAddress: "198.51.100.33", SubnetID: "1d699529-bdfd-43f8-bcaa-bff00c547af2"},
		},
		QoSPolicyID: "6601bae5-f15a-4687-8be9-ddec9a2f8a8b",
	})
	th.AssertEquals(t, n.Name, "router1")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.TenantID, "d6554fe62e2f41efbb6e026fad5c1542")
	th.AssertEquals(t, n.ID, "a07eea83-7710-4860-931b-5fe220fae533")
	th.AssertDeepEquals(t, n.Routes, []routers.Route{{DestinationCIDR: "40.0.1.0/24", NextHop: "10.1.0.10"}})
	th.AssertDeepEquals(t, n.AvailabilityZoneHints, []string{"zone1", "zone2"})
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "router": {
			"name": "new_name",
        "external_gateway_info": {
            "network_id": "8ca37218-28ff-41cb-9b10-039601ea7e6b",
            "qos_policy_id": "01ba32e5-f15a-4687-8be9-ddec92a2f8a8"
		},
        "routes": [
            {
                "nexthop": "10.1.0.10",
                "destination": "40.0.1.0/24"
            }
        ]
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "status": "ACTIVE",
        "external_gateway_info": {
            "network_id": "8ca37218-28ff-41cb-9b10-039601ea7e6b",
            "external_fixed_ips": [
                {"ip_address": "192.0.2.17", "subnet_id": "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"}
            ],
            "qos_policy_id": "01ba32e5-f15a-4687-8be9-ddec92a2f8a8"
        },
        "name": "new_name",
        "admin_state_up": true,
        "tenant_id": "6b96ff0cb17a4b859e1e575d221683d3",
        "distributed": false,
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e",
        "routes": [
            {
                "nexthop": "10.1.0.10",
                "destination": "40.0.1.0/24"
            }
        ]
    }
}
		`)
	})

	gwi := routers.GatewayInfo{
		NetworkID:   "8ca37218-28ff-41cb-9b10-039601ea7e6b",
		QoSPolicyID: "01ba32e5-f15a-4687-8be9-ddec92a2f8a8",
	}
	r := []routers.Route{{DestinationCIDR: "40.0.1.0/24", NextHop: "10.1.0.10"}}
	options := routers.UpdateOpts{Name: "new_name", GatewayInfo: &gwi, Routes: &r}

	n, err := routers.Update(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	gwi.ExternalFixedIPs = []routers.ExternalFixedIP{
		{IPAddress: "192.0.2.17", SubnetID: "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"},
	}

	th.AssertEquals(t, n.Name, "new_name")
	th.AssertDeepEquals(t, n.GatewayInfo, gwi)
	th.AssertDeepEquals(t, n.Routes, []routers.Route{{DestinationCIDR: "40.0.1.0/24", NextHop: "10.1.0.10"}})
}

func TestUpdateWithoutRoutes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "router": {
        "name": "new_name"
    }
}
		`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "status": "ACTIVE",
        "external_gateway_info": {
            "network_id": "8ca37218-28ff-41cb-9b10-039601ea7e6b",
            "external_fixed_ips": [
                {"ip_address": "192.0.2.17", "subnet_id": "ab561bc4-1a8e-48f2-9fbd-376fcb1a1def"}
            ]
        },
        "name": "new_name",
        "admin_state_up": true,
        "tenant_id": "6b96ff0cb17a4b859e1e575d221683d3",
        "distributed": false,
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e",
        "routes": [
            {
                "nexthop": "10.1.0.10",
                "destination": "40.0.1.0/24"
            }
        ]
    }
}
		`)
	})

	options := routers.UpdateOpts{Name: "new_name"}

	n, err := routers.Update(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Name, "new_name")
	th.AssertDeepEquals(t, n.Routes, []routers.Route{{DestinationCIDR: "40.0.1.0/24", NextHop: "10.1.0.10"}})
}

func TestAllRoutesRemoved(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "router": {
        "routes": []
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "router": {
        "status": "ACTIVE",
        "external_gateway_info": {
            "network_id": "8ca37218-28ff-41cb-9b10-039601ea7e6b"
        },
        "name": "name",
        "admin_state_up": true,
        "tenant_id": "6b96ff0cb17a4b859e1e575d221683d3",
        "distributed": false,
        "id": "8604a0de-7f6b-409a-a47c-a1cc7bc77b2e",
        "routes": []
    }
}
		`)
	})

	r := []routers.Route{}
	options := routers.UpdateOpts{Routes: &r}

	n, err := routers.Update(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, n.Routes, []routers.Route{})
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := routers.Delete(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c")
	th.AssertNoErr(t, res.Err)
}

func TestAddInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c/add_router_interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "subnet_id": "a2f1f29d-571b-4533-907f-5803ab96ead1"
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "subnet_id": "0d32a837-8069-4ec3-84c4-3eef3e10b188",
    "tenant_id": "017d8de156df4177889f31a9bd6edc00",
    "port_id": "3f990102-4485-4df1-97a0-2c35bdb85b31",
    "id": "9a83fa11-8da5-436e-9afe-3d3ac5ce7770"
}
`)
	})

	opts := routers.AddInterfaceOpts{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1"}
	res, err := routers.AddInterface(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "0d32a837-8069-4ec3-84c4-3eef3e10b188", res.SubnetID)
	th.AssertEquals(t, "017d8de156df4177889f31a9bd6edc00", res.TenantID)
	th.AssertEquals(t, "3f990102-4485-4df1-97a0-2c35bdb85b31", res.PortID)
	th.AssertEquals(t, "9a83fa11-8da5-436e-9afe-3d3ac5ce7770", res.ID)
}

func TestAddInterfaceRequiredOpts(t *testing.T) {
	_, err := routers.AddInterface(context.TODO(), fake.ServiceClient(), "foo", routers.AddInterfaceOpts{}).Extract()
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
	_, err = routers.AddInterface(context.TODO(), fake.ServiceClient(), "foo", routers.AddInterfaceOpts{SubnetID: "bar", PortID: "baz"}).Extract()
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestRemoveInterface(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/4e8e5957-649f-477b-9e5b-f1f75b21c03c/remove_router_interface", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
		"subnet_id": "a2f1f29d-571b-4533-907f-5803ab96ead1"
}
	`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
		"subnet_id": "0d32a837-8069-4ec3-84c4-3eef3e10b188",
		"tenant_id": "017d8de156df4177889f31a9bd6edc00",
		"port_id": "3f990102-4485-4df1-97a0-2c35bdb85b31",
		"id": "9a83fa11-8da5-436e-9afe-3d3ac5ce7770"
}
`)
	})

	opts := routers.RemoveInterfaceOpts{SubnetID: "a2f1f29d-571b-4533-907f-5803ab96ead1"}
	res, err := routers.RemoveInterface(context.TODO(), fake.ServiceClient(), "4e8e5957-649f-477b-9e5b-f1f75b21c03c", opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "0d32a837-8069-4ec3-84c4-3eef3e10b188", res.SubnetID)
	th.AssertEquals(t, "017d8de156df4177889f31a9bd6edc00", res.TenantID)
	th.AssertEquals(t, "3f990102-4485-4df1-97a0-2c35bdb85b31", res.PortID)
	th.AssertEquals(t, "9a83fa11-8da5-436e-9afe-3d3ac5ce7770", res.ID)
}

func TestListL3Agents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/routers/fa3a4aaa-c73f-48aa-a603-8c8bf642b7c0/l3-agents", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "agents": [
        {
            "id": "ddbf087c-e38f-4a73-bcb3-c38f2a719a03",
            "agent_type": "L3 agent",
            "binary": "neutron-l3-agent",
            "topic": "l3_agent",
            "host": "os-ctrl-02",
            "admin_state_up": true,
            "created_at": "2017-07-26 23:15:44",
            "started_at": "2018-06-26 21:46:19",
            "heartbeat_timestamp": "2019-01-09 10:28:53",
            "description": "My L3 agent for OpenStack",
            "resources_synced": true,
            "availability_zone": "nova",
            "alive": true,
            "configurations": {
                "agent_mode": "legacy",
                "ex_gw_ports": 2,
                "floating_ips": 2,
                "handle_internal_only_routers": true,
                "interface_driver": "linuxbridge",
                "interfaces": 1,
                "log_agent_heartbeats": false,
                "routers": 2
            },
            "resource_versions": {},
            "ha_state": "standby"
        },
        {
            "id": "4541cc6c-87bc-4cee-bad2-36ca78836c91",
            "agent_type": "L3 agent",
            "binary": "neutron-l3-agent",
            "topic": "l3_agent",
            "host": "os-ctrl-03",
            "admin_state_up": true,
            "created_at": "2017-01-22 14:00:50",
            "started_at": "2018-11-06 12:09:17",
            "heartbeat_timestamp": "2019-01-09 10:28:50",
            "description": "My L3 agent for OpenStack",
            "resources_synced": true,
            "availability_zone": "nova",
            "alive": true,
            "configurations": {
                "agent_mode": "legacy",
                "ex_gw_ports": 2,
                "floating_ips": 2,
                "handle_internal_only_routers": true,
                "interface_driver": "linuxbridge",
                "interfaces": 1,
                "log_agent_heartbeats": false,
                "routers": 2
            },
            "resource_versions": {},
            "ha_state": "active"
        }
    ]
}
			`)
	})

	l3AgentsPages, err := routers.ListL3Agents(fake.ServiceClient(), "fa3a4aaa-c73f-48aa-a603-8c8bf642b7c0").AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := routers.ExtractL3Agents(l3AgentsPages)
	th.AssertNoErr(t, err)

	expected := []routers.L3Agent{
		{
			ID:               "ddbf087c-e38f-4a73-bcb3-c38f2a719a03",
			AdminStateUp:     true,
			AgentType:        "L3 agent",
			Description:      "My L3 agent for OpenStack",
			Alive:            true,
			ResourcesSynced:  true,
			Binary:           "neutron-l3-agent",
			AvailabilityZone: "nova",
			Configurations: map[string]any{
				"agent_mode":                   "legacy",
				"ex_gw_ports":                  float64(2),
				"floating_ips":                 float64(2),
				"handle_internal_only_routers": true,
				"interface_driver":             "linuxbridge",
				"interfaces":                   float64(1),
				"log_agent_heartbeats":         false,
				"routers":                      float64(2),
			},
			CreatedAt:          time.Date(2017, 7, 26, 23, 15, 44, 0, time.UTC),
			StartedAt:          time.Date(2018, 6, 26, 21, 46, 19, 0, time.UTC),
			HeartbeatTimestamp: time.Date(2019, 1, 9, 10, 28, 53, 0, time.UTC),
			Host:               "os-ctrl-02",
			Topic:              "l3_agent",
			HAState:            "standby",
			ResourceVersions:   map[string]any{},
		},
		{
			ID:               "4541cc6c-87bc-4cee-bad2-36ca78836c91",
			AdminStateUp:     true,
			AgentType:        "L3 agent",
			Description:      "My L3 agent for OpenStack",
			Alive:            true,
			ResourcesSynced:  true,
			Binary:           "neutron-l3-agent",
			AvailabilityZone: "nova",
			Configurations: map[string]any{
				"agent_mode":                   "legacy",
				"ex_gw_ports":                  float64(2),
				"floating_ips":                 float64(2),
				"handle_internal_only_routers": true,
				"interface_driver":             "linuxbridge",
				"interfaces":                   float64(1),
				"log_agent_heartbeats":         false,
				"routers":                      float64(2),
			},
			CreatedAt:          time.Date(2017, 1, 22, 14, 00, 50, 0, time.UTC),
			StartedAt:          time.Date(2018, 11, 6, 12, 9, 17, 0, time.UTC),
			HeartbeatTimestamp: time.Date(2019, 1, 9, 10, 28, 50, 0, time.UTC),
			Host:               "os-ctrl-03",
			Topic:              "l3_agent",
			HAState:            "active",
			ResourceVersions:   map[string]any{},
		},
	}
	th.CheckDeepEquals(t, expected, actual)
}
