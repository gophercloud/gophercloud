package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/routerinsertion"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "name": "fw",
        "description": "OpenStack firewall",
        "admin_state_up": true,
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ports": [
          "8a3a0d6a-34b5-4a92-b65d-6375a4c1e9e8"
        ]
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "firewall_group":{
        "status": "PENDING_CREATE",
        "name": "fw",
        "description": "OpenStack firewall",
        "admin_state_up": true,
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c"
    }
}
    `)
	})

	firewallCreateOpts := firewall_groups.CreateOpts{
		TenantID:     		"b4eedccc6fb74fa8a7ad6b08382b852b",
		Name:         		"fw",
		Description:  		"OpenStack firewall",
		AdminStateUp: 		gophercloud.Enabled,
		IngressPolicyID:	"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
		EgressPolicyID:		"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
	}
	createOpts := routerinsertion.CreateOptsExt{
		CreateOptsBuilder:	firewallCreateOpts,
		PortIDs:			[]string{"8a3a0d6a-34b5-4a92-b65d-6375a4c1e9e8"},
	}

	_, err := firewall_groups.Create(fake.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
}

func TestCreateWithNoRouters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "name": "fw",
        "description": "OpenStack firewall",
        "admin_state_up": true,
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ports": []
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "firewall_group":{
        "status": "PENDING_CREATE",
        "name": "fw",
        "description": "OpenStack firewall",
        "admin_state_up": true,
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c"
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c"
    }
}
    `)
	})

	firewallCreateOpts := firewall_groups.CreateOpts{
		TenantID:     		"b4eedccc6fb74fa8a7ad6b08382b852b",
		Name:         		"fw",
		Description:  		"OpenStack firewall group",
		AdminStateUp: 		gophercloud.Enabled,
		IngressPolicyID:	"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
		EgressPolicyID:		"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
	}
	createOpts := routerinsertion.CreateOptsExt{
		CreateOptsBuilder: firewallCreateOpts,
		PortIDs:         []string{},
	}

	_, err := firewall_groups.Create(fake.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/ea5b5315-64f6-4ea3-8e58-981cc37c6576", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "name": "fw",
        "description": "updated fw",
        "admin_state_up":false,
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "ports": [
          "8a3a0d6a-34b5-4a92-b65d-6375a4c1e9e8"
        ]
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "firewall_group": {
        "status": "ACTIVE",
        "name": "fw",
        "admin_state_up": false,
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "id": "ea5b5315-64f6-4ea3-8e58-981cc37c6576",
        "description": "OpenStack firewall group"
    }
}
    `)
	})

	firewallUpdateOpts := firewall_groups.UpdateOpts{
		Name:         		"fw",
		Description:  		"updated fw",
		AdminStateUp: 		gophercloud.Disabled,
		IngressPolicyID:	"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
		EgressPolicyID:     "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
	}
	updateOpts := routerinsertion.UpdateOptsExt{
		UpdateOptsBuilder: firewallUpdateOpts,
		PortIDs:         	[]string{"8a3a0d6a-34b5-4a92-b65d-6375a4c1e9e8"},
	}

	_, err := firewall_groups.Update(fake.ServiceClient(), "ea5b5315-64f6-4ea3-8e58-981cc37c6576", updateOpts).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdateWithNoRouters(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/ea5b5315-64f6-4ea3-8e58-981cc37c6576", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "name": "fw",
        "description": "updated fw",
        "admin_state_up":false,
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "ports": []
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "firewall_group": {
        "status": "ACTIVE",
        "name": "fw",
        "admin_state_up": false,
        "tenant_id": "b4eedccc6fb74fa8a7ad6b08382b852b",
        "ingress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "egress_firewall_policy_id": "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
        "id": "ea5b5315-64f6-4ea3-8e58-981cc37c6576",
        "description": "OpenStack firewall group"
    }
}
    `)
	})

	firewallUpdateOpts := firewall_groups.UpdateOpts{
		Name:         		"fw",
		Description:  		"updated fw",
		AdminStateUp: 		gophercloud.Disabled,
		IngressPolicyID:    "19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
		EgressPolicyID:		"19ab8c87-4a32-4e6a-a74e-b77fffb89a0c",
	}
	updateOpts := routerinsertion.UpdateOptsExt{
		UpdateOptsBuilder: firewallUpdateOpts,
		PortIDs:         	[]string{},
	}

	_, err := firewall_groups.Update(fake.ServiceClient(), "ea5b5315-64f6-4ea3-8e58-981cc37c6576", updateOpts).Extract()
	th.AssertNoErr(t, err)
}
