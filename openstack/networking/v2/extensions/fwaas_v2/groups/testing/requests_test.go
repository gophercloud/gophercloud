package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/groups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "firewall_groups": [
    {
      "id": "3af94f0e-b52d-491a-87d2-704497305948",
      "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
      "name": "test",
      "description": "fancy group",
      "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
      "egress_firewall_policy_id": null,
      "admin_state_up": true,
      "ports": [
        "a6af1e56-b12b-4733-8f77-49166afd5719"
      ],
      "status": "ACTIVE",
      "shared": false,
      "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
    },
    {
      "id": "f9fbb80c-eeb4-4f3f-aa50-1032960c08ea",
      "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
      "name": "default",
      "description": "Default firewall group",
      "ingress_firewall_policy_id": "90e3fcac-3bfb-48f9-8e91-2a78fb352b92",
      "egress_firewall_policy_id": "122fb344-3c28-49f0-af00-f7fcbc88330b",
      "admin_state_up": true,
      "ports": [
        "20da216c-bab3-4cf6-bd6b-8904b133a816",
        "4f4c714c-185f-487e-998c-c1a35da3c4f4",
        "681b1db4-40ca-4314-b098-d2f43225e7f7",
        "82f9d868-6f56-44fb-9684-654dc473bed0",
        "a5858b5d-20dc-4bb1-9f95-1d322c8bb81b",
        "d25a04a2-447b-4ee1-80d7-b32967dbb643"
      ],
      "status": "ACTIVE",
      "shared": false,
      "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
    }
  ]
}
        `)
	})

	count := 0

	err := groups.List(fake.ServiceClient(), groups.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := groups.ExtractGroups(page)
		if err != nil {
			t.Errorf("Failed to extract members: %v", err)
			return false, err
		}

		expected := []groups.Group{
			{
				ID:                      "3af94f0e-b52d-491a-87d2-704497305948",
				TenantID:                "9f98fc0e5f944cd1b51798b668dc8778",
				Name:                    "test",
				Description:             "fancy group",
				IngressFirewallPolicyID: "e3f11142-3792-454b-8d3e-91ac1bf127b4",
				EgressFirewallPolicyID:  "",
				AdminStateUp:            true,
				Ports: []string{
					"a6af1e56-b12b-4733-8f77-49166afd5719",
				},
				Status:    "ACTIVE",
				Shared:    false,
				ProjectID: "9f98fc0e5f944cd1b51798b668dc8778",
			},
			{
				ID:                      "f9fbb80c-eeb4-4f3f-aa50-1032960c08ea",
				TenantID:                "9f98fc0e5f944cd1b51798b668dc8778",
				Name:                    "default",
				Description:             "Default firewall group",
				IngressFirewallPolicyID: "90e3fcac-3bfb-48f9-8e91-2a78fb352b92",
				EgressFirewallPolicyID:  "122fb344-3c28-49f0-af00-f7fcbc88330b",
				AdminStateUp:            true,
				Ports: []string{
					"20da216c-bab3-4cf6-bd6b-8904b133a816",
					"4f4c714c-185f-487e-998c-c1a35da3c4f4",
					"681b1db4-40ca-4314-b098-d2f43225e7f7",
					"82f9d868-6f56-44fb-9684-654dc473bed0",
					"a5858b5d-20dc-4bb1-9f95-1d322c8bb81b",
					"d25a04a2-447b-4ee1-80d7-b32967dbb643",
				},
				Status:    "ACTIVE",
				Shared:    false,
				ProjectID: "9f98fc0e5f944cd1b51798b668dc8778",
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

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/6bfb0f10-07f7-4a40-b534-bad4b4ca3428", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "firewall_group": {
    "id": "6bfb0f10-07f7-4a40-b534-bad4b4ca3428",
    "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
    "name": "test",
    "description": "some information",
    "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
    "egress_firewall_policy_id": null,
    "admin_state_up": true,
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719"
    ],
    "status": "ACTIVE",
    "shared": false,
    "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
  }
}
        `)
	})

	group, err := groups.Get(context.TODO(), fake.ServiceClient(), "6bfb0f10-07f7-4a40-b534-bad4b4ca3428").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "6bfb0f10-07f7-4a40-b534-bad4b4ca3428", group.ID)
	th.AssertEquals(t, "9f98fc0e5f944cd1b51798b668dc8778", group.TenantID)
	th.AssertEquals(t, "test", group.Name)
	th.AssertEquals(t, "some information", group.Description)
	th.AssertEquals(t, "e3f11142-3792-454b-8d3e-91ac1bf127b4", group.IngressFirewallPolicyID)
	th.AssertEquals(t, "", group.EgressFirewallPolicyID)
	th.AssertEquals(t, true, group.AdminStateUp)
	th.AssertEquals(t, 1, len(group.Ports))
	th.AssertEquals(t, "a6af1e56-b12b-4733-8f77-49166afd5719", group.Ports[0])
	th.AssertEquals(t, "ACTIVE", group.Status)
	th.AssertEquals(t, false, group.Shared)
	th.AssertEquals(t, "9f98fc0e5f944cd1b51798b668dc8778", group.TenantID)
}

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
  "firewall_group": {
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719"
    ],
    "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
	"egress_firewall_policy_id": "43a11f3a-ddac-4129-9469-02b9df26548e",
    "name": "test"
  }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
  "firewall_group": {
    "id": "6bfb0f10-07f7-4a40-b534-bad4b4ca3428",
    "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
    "name": "test",
    "description": "",
    "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
    "egress_firewall_policy_id": "43a11f3a-ddac-4129-9469-02b9df26548e",
    "admin_state_up": true,
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719"
    ],
    "status": "CREATED",
    "shared": false,
    "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
  }
}
        `)
	})

	options := groups.CreateOpts{
		Name:                    "test",
		Description:             "",
		IngressFirewallPolicyID: "e3f11142-3792-454b-8d3e-91ac1bf127b4",
		EgressFirewallPolicyID:  "43a11f3a-ddac-4129-9469-02b9df26548e",
		Ports: []string{
			"a6af1e56-b12b-4733-8f77-49166afd5719",
		},
	}

	_, err := groups.Create(context.TODO(), fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/6bfb0f10-07f7-4a40-b534-bad4b4ca3428", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "name": "the group",
        "ports": [
					"a6af1e56-b12b-4733-8f77-49166afd5719",
					"11a58c87-76be-ae7c-a74e-b77fffb88a32"
        ],
				"description": "Firewall group",
				"admin_state_up": false
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "firewall_group": {
    "id": "6bfb0f10-07f7-4a40-b534-bad4b4ca3428",
    "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
    "name": "test",
    "description": "some information",
    "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
    "egress_firewall_policy_id": "43a11f3a-ddac-4129-9469-02b9df26548e",
    "admin_state_up": true,
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719",
			"11a58c87-76be-ae7c-a74e-b77fffb88a32"
    ],
    "status": "ACTIVE",
    "shared": false,
    "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
  }
}
    `)
	})

	name := "the group"
	description := "Firewall group"
	adminStateUp := false
	options := groups.UpdateOpts{
		Name:        &name,
		Description: &description,
		Ports: &[]string{
			"a6af1e56-b12b-4733-8f77-49166afd5719",
			"11a58c87-76be-ae7c-a74e-b77fffb88a32",
		},
		AdminStateUp: &adminStateUp,
	}

	_, err := groups.Update(context.TODO(), fake.ServiceClient(), "6bfb0f10-07f7-4a40-b534-bad4b4ca3428", options).Extract()
	th.AssertNoErr(t, err)
}

func TestRemoveIngressPolicy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/6bfb0f10-07f7-4a40-b534-bad4b4ca3428", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "ingress_firewall_policy_id": null
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "firewall_group": {
    "id": "6bfb0f10-07f7-4a40-b534-bad4b4ca3428",
    "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
    "name": "test",
    "description": "some information",
    "ingress_firewall_policy_id": null,
    "egress_firewall_policy_id": "43a11f3a-ddac-4129-9469-02b9df26548e",
    "admin_state_up": true,
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719",
			"11a58c87-76be-ae7c-a74e-b77fffb88a32"
    ],
    "status": "ACTIVE",
    "shared": false,
    "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
  }
}
    `)
	})

	removeIngressPolicy, err := groups.RemoveIngressPolicy(context.TODO(), fake.ServiceClient(), "6bfb0f10-07f7-4a40-b534-bad4b4ca3428").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, removeIngressPolicy.IngressFirewallPolicyID, "")
	th.AssertEquals(t, removeIngressPolicy.EgressFirewallPolicyID, "43a11f3a-ddac-4129-9469-02b9df26548e")
}

func TestRemoveEgressPolicy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/6bfb0f10-07f7-4a40-b534-bad4b4ca3428", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "firewall_group":{
        "egress_firewall_policy_id": null
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "firewall_group": {
    "id": "6bfb0f10-07f7-4a40-b534-bad4b4ca3428",
    "tenant_id": "9f98fc0e5f944cd1b51798b668dc8778",
    "name": "test",
    "description": "some information",
    "ingress_firewall_policy_id": "e3f11142-3792-454b-8d3e-91ac1bf127b4",
    "egress_firewall_policy_id": null,
    "admin_state_up": true,
    "ports": [
      "a6af1e56-b12b-4733-8f77-49166afd5719",
			"11a58c87-76be-ae7c-a74e-b77fffb88a32"
    ],
    "status": "ACTIVE",
    "shared": false,
    "project_id": "9f98fc0e5f944cd1b51798b668dc8778"
  }
}
    `)
	})

	removeEgressPolicy, err := groups.RemoveEgressPolicy(context.TODO(), fake.ServiceClient(), "6bfb0f10-07f7-4a40-b534-bad4b4ca3428").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, removeEgressPolicy.IngressFirewallPolicyID, "e3f11142-3792-454b-8d3e-91ac1bf127b4")
	th.AssertEquals(t, removeEgressPolicy.EgressFirewallPolicyID, "")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/fwaas/firewall_groups/4ec89077-d057-4a2b-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := groups.Delete(context.TODO(), fake.ServiceClient(), "4ec89077-d057-4a2b-911f-60a3b47ee304")
	th.AssertNoErr(t, res.Err)
}
