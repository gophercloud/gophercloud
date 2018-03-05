package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/vpn/ipsec-site-connections", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    
    "ipsec_site_connection": {
        "psk": "secret",
        "initiator": "bi-directional",
        "ipsecpolicy_id": "e6e23d0c-9519-4d52-8ea4-5b1f96d857b1",
        "admin_state_up": true,
        "mtu": 1500,
        "peer_ep_group_id": "9ad5a7e0-6dac-41b4-b20d-a7b8645fddf1",
        "ikepolicy_id": "9b00d6b0-6c93-4ca5-9747-b8ade7bb514f",
        "vpnservice_id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "local_ep_group_id": "3e1815dd-e212-43d0-8f13-b494fa553e68",
        "peer_address": "172.24.4.233",
        "peer_id": "172.24.4.233",
        "name": "vpnconnection1"
    
}
}      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "ipsec_site_connection": {
        "status": "PENDING_CREATE",
        "psk": "secret",
        "initiator": "bi-directional",
        "name": "vpnconnection1",
        "admin_state_up": true,
        "project_id": "10039663455a446d8ba2cbb058b0f578",
        "tenant_id": "10039663455a446d8ba2cbb058b0f578",
        "auth_mode": "psk",
        "peer_cidrs": [],
        "mtu": 1500,
        "peer_ep_group_id": "9ad5a7e0-6dac-41b4-b20d-a7b8645fddf1",
        "ikepolicy_id": "9b00d6b0-6c93-4ca5-9747-b8ade7bb514f",
        "vpnservice_id": "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
        "dpd": {
            "action": "hold",
            "interval": 30,
            "timeout": 120
        },
        "route_mode": "static",
        "ipsecpolicy_id": "e6e23d0c-9519-4d52-8ea4-5b1f96d857b1",
        "local_ep_group_id": "3e1815dd-e212-43d0-8f13-b494fa553e68",
        "peer_address": "172.24.4.233",
        "peer_id": "172.24.4.233",
        "id": "851f280f-5639-4ea3-81aa-e298525ab74b",
        "description": ""
    }
}
    `)
	})

	options := siteconnections.CreateOpts{
		Name:           "vpnconnection1",
		AdminStateUp:   gophercloud.Enabled,
		PSK:            "secret",
		Initiator:      siteconnections.InitiatorBiDirectional,
		IPSecPolicyID:  "e6e23d0c-9519-4d52-8ea4-5b1f96d857b1",
		MTU:            1500,
		PeerEPGroupID:  "9ad5a7e0-6dac-41b4-b20d-a7b8645fddf1",
		IKEPolicyID:    "9b00d6b0-6c93-4ca5-9747-b8ade7bb514f",
		VPNServiceID:   "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		LocalEPGroupID: "3e1815dd-e212-43d0-8f13-b494fa553e68",
		PeerAddress:    "172.24.4.233",
		PeerID:         "172.24.4.233",
	}
	actual, err := siteconnections.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
	expectedDPD := siteconnections.DPD{
		Action:   "hold",
		Interval: 30,
		Timeout:  120,
	}
	expected := siteconnections.Connection{
		TenantID:           "10039663455a446d8ba2cbb058b0f578",
		Name:               "vpnconnection1",
		AdminStateUp:       true,
		PSK:                "secret",
		Initiator:          "bi-directional",
		IPSecPolicyID:      "e6e23d0c-9519-4d52-8ea4-5b1f96d857b1",
		MTU:                1500,
		PeerEPGroupID:      "9ad5a7e0-6dac-41b4-b20d-a7b8645fddf1",
		IKEPolicyID:        "9b00d6b0-6c93-4ca5-9747-b8ade7bb514f",
		VPNServiceID:       "5c561d9d-eaea-45f6-ae3e-08d1a7080828",
		LocalEPGroupID:     "3e1815dd-e212-43d0-8f13-b494fa553e68",
		PeerAddress:        "172.24.4.233",
		PeerID:             "172.24.4.233",
		Status:             "PENDING_CREATE",
		ProjectID:          "10039663455a446d8ba2cbb058b0f578",
		AuthenticationMode: "psk",
		PeerCIDRs:          []string{},
		DPD:                expectedDPD,
		RouteMode:          "static",
		ID:                 "851f280f-5639-4ea3-81aa-e298525ab74b",
		Description:        "",
	}
	th.AssertDeepEquals(t, expected, *actual)
}
