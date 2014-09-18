package ports

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const TokenID = "123"

func ServiceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: TokenID,
		},
		Endpoint: th.Endpoint(),
	}
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "ports": [
        {
            "status": "ACTIVE",
            "binding:host_id": "devstack",
            "name": "",
            "allowed_address_pairs": [],
            "admin_state_up": true,
            "network_id": "70c1db1f-b701-45bd-96e0-a313ee3430b3",
            "tenant_id": "",
            "extra_dhcp_opts": [],
            "binding:vif_details": {
                "port_filter": true,
                "ovs_hybrid_plug": true
            },
            "binding:vif_type": "ovs",
            "device_owner": "network:router_gateway",
            "mac_address": "fa:16:3e:58:42:ed",
            "binding:profile": {},
            "binding:vnic_type": "normal",
            "fixed_ips": [
                {
                    "subnet_id": "008ba151-0b8c-4a67-98b5-0d2b87666062",
                    "ip_address": "172.24.4.2"
                }
            ],
            "id": "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
            "security_groups": [],
            "device_id": "9ae135f4-b6e0-4dad-9e91-3c223e385824"
        }
    ]
}
      `)
	})

	count := 0

	List(ServiceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractPorts(page)
		if err != nil {
			t.Errorf("Failed to extract subnets: %v", err)
			return false, nil
		}

		expected := []Port{
			Port{
				Status:              "ACTIVE",
				Name:                "",
				AllowedAddressPairs: []interface{}(nil),
				AdminStateUp:        true,
				NetworkID:           "70c1db1f-b701-45bd-96e0-a313ee3430b3",
				TenantID:            "",
				ExtraDHCPOpts:       []interface{}{},
				DeviceOwner:         "network:router_gateway",
				MACAddress:          "fa:16:3e:58:42:ed",
				FixedIPs: []IP{
					IP{
						SubnetID:  "008ba151-0b8c-4a67-98b5-0d2b87666062",
						IPAddress: "172.24.4.2",
					},
				},
				ID:                "d80b1a3b-4fc1-49f3-952e-1e2ab7081d8b",
				SecurityGroups:    []string{},
				DeviceID:          "9ae135f4-b6e0-4dad-9e91-3c223e385824",
				BindingHostID:     "devstack",
				BindingVIFDetails: map[string]interface{}{"port_filter": true, "ovs_hybrid_plug": true},
				BindingVIFType:    "ovs",
				BindingProfile:    map[string]interface{}{},
				BindingVNICType:   "normal",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
