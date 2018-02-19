package testing

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	th "github.com/gophercloud/gophercloud/testhelper"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fw/firewall_policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_policy":{
        "name": "policy",
        "firewall_rules": [
            "98a58c87-76be-ae7c-a74e-b77fffb88d95",
            "11a58c87-76be-ae7c-a74e-b77fffb88a32"
        ],
        "description": "Firewall policy",
		"tenant_id": "9145d91459d248b1b02fdaca97c6a75d",
		"audited": true,
		"shared": false
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
    "firewall_policy":{
        "name": "policy",
        "firewall_rules": [
            "98a58c87-76be-ae7c-a74e-b77fffb88d95",
            "11a58c87-76be-ae7c-a74e-b77fffb88a32"
        ],
        "tenant_id": "9145d91459d248b1b02fdaca97c6a75d",
        "audited": false,
        "id": "f2b08c1e-aa81-4668-8ae1-1401bcb0576c",
        "description": "Firewall policy"
    }
}
        `)
	})

	options := ikepolicies.CreateOpts{
		TenantID:    "9145d91459d248b1b02fdaca97c6a75d",
		Name:        "policy",
		Description: "Firewall policy",
		Shared:      gophercloud.Disabled,
		Audited:     gophercloud.Enabled,
		Rules: []string{
			"98a58c87-76be-ae7c-a74e-b77fffb88d95",
			"11a58c87-76be-ae7c-a74e-b77fffb88a32",
		},
	}

	_, err := policies.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}
