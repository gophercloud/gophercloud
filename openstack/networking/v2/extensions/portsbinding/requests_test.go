package portsbinding

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/rackspace/gophercloud/openstack/networking/v2/common"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "port": {
        "status": "ACTIVE",
        "name": "",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "7e02058126cc4950b75f9970368ba177",
        "device_owner": "network:router_interface",
        "mac_address": "fa:16:3e:23:fd:d7",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.1"
            }
        ],
        "id": "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2",
        "security_groups": [],
        "device_id": "5e3898d7-11be-483e-9732-b2f5eccd2b2e",
        "binding:host_id": "HOST1",
        "binding:vnic_type": "normal"
    }
}
			`)
	})

	n, err := Get(fake.ServiceClient(), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "ACTIVE")
	th.AssertEquals(t, n.Name, "")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "7e02058126cc4950b75f9970368ba177")
	th.AssertEquals(t, n.DeviceOwner, "network:router_interface")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:23:fd:d7")
	th.AssertDeepEquals(t, n.FixedIPs, []IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.1"},
	})
	th.AssertEquals(t, n.ID, "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{})
	th.AssertEquals(t, n.DeviceID, "5e3898d7-11be-483e-9732-b2f5eccd2b2e")
	th.AssertEquals(t, n.HostID, "HOST1")
	th.AssertEquals(t, n.VNICType, "normal")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "port": {
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "name": "private-port",
        "admin_state_up": true,
		"fixed_ips": [
				{
						"subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
						"ip_address": "10.0.0.2"
				}
		],
		"security_groups": ["foo"],
		"binding:host_id": "HOST1",
        "binding:vnic_type": "normal"
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "port": {
        "status": "DOWN",
        "name": "private-port",
        "allowed_address_pairs": [],
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.2"
            }
        ],
        "binding:host_id": "HOST1",
        "binding:vnic_type": "normal",
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": ""
    }
}
		`)
	})

	asu := true
	options := CreateOpts{
		CreateOptsBuilder: ports.CreateOpts{
			Name:         "private-port",
			AdminStateUp: &asu,
			NetworkID:    "a87cc70a-3e15-4acf-8205-9b711a3531b7",
			FixedIPs: []IP{
				{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
			},
			SecurityGroups: []string{"foo"},
		},
		HostID:   "HOST1",
		VNICType: "normal",
	}
	n, err := Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.Status, "DOWN")
	th.AssertEquals(t, n.Name, "private-port")
	th.AssertEquals(t, n.AdminStateUp, true)
	th.AssertEquals(t, n.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, n.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertEquals(t, n.DeviceOwner, "")
	th.AssertEquals(t, n.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, n.FixedIPs, []IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.2"},
	})
	th.AssertEquals(t, n.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertDeepEquals(t, n.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
	th.AssertEquals(t, n.HostID, "HOST1")
	th.AssertEquals(t, n.VNICType, "normal")
}

func TestRequiredCreateOpts(t *testing.T) {
	res := Create(fake.ServiceClient(), CreateOpts{CreateOptsBuilder: ports.CreateOpts{}})
	if res.Err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/65c0ee9f-d634-4522-8954-51021b570b0d", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
		"port": {
			"name": "new_port_name",
			"fixed_ips": [
				{
					"subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
					"ip_address": "10.0.0.3"
				}
			],
			"security_groups": [
            	"f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        	],
        	"binding:host_id": "HOST1",
        	"binding:vnic_type": "normal"
		}
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "port": {
        "status": "DOWN",
        "name": "new_port_name",
        "admin_state_up": true,
        "network_id": "a87cc70a-3e15-4acf-8205-9b711a3531b7",
        "tenant_id": "d6700c0c9ffa4f1cb322cd4a1f3906fa",
        "device_owner": "",
        "mac_address": "fa:16:3e:c9:cb:f0",
        "fixed_ips": [
            {
                "subnet_id": "a0304c3a-4f08-4c43-88af-d796509c97d2",
                "ip_address": "10.0.0.3"
            }
        ],
        "id": "65c0ee9f-d634-4522-8954-51021b570b0d",
        "security_groups": [
            "f0ac4394-7e4a-4409-9701-ba8be283dbc3"
        ],
        "device_id": "",
        "binding:host_id": "HOST1",
        "binding:vnic_type": "normal"
    }
}
		`)
	})

	options := UpdateOpts{
		UpdateOptsBuilder: ports.UpdateOpts{
			Name: "new_port_name",
			FixedIPs: []IP{
				{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
			},
			SecurityGroups: []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"},
		},
		HostID:   "HOST1",
		VNICType: "normal",
	}

	s, err := Update(fake.ServiceClient(), "65c0ee9f-d634-4522-8954-51021b570b0d", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Name, "new_port_name")
	th.AssertDeepEquals(t, s.FixedIPs, []IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.3"},
	})
	th.AssertDeepEquals(t, s.SecurityGroups, []string{"f0ac4394-7e4a-4409-9701-ba8be283dbc3"})
	th.AssertEquals(t, s.HostID, "HOST1")
	th.AssertEquals(t, s.VNICType, "normal")
}
