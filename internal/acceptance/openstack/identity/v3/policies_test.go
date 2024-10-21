//go:build acceptance || identity || policies

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/policies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPoliciesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := policies.List(client, policies.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestPoliciesCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := policies.CreateOpts{
		Type: "application/json",
		Blob: []byte("{'foobar_user': 'role:compute-user'}"),
		Extra: map[string]any{
			"description": "policy for foobar_user",
		},
	}

	policy, err := policies.Create(context.TODO(), client, &createOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, policy)
	tools.PrintResource(t, policy.Extra)

	th.AssertEquals(t, policy.Type, createOpts.Type)
	th.AssertEquals(t, policy.Blob, string(createOpts.Blob))
	th.AssertEquals(t, policy.Extra["description"], createOpts.Extra["description"])

	var listOpts policies.ListOpts

	allPages, err := policies.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, p := range allPolicies {
		tools.PrintResource(t, p)
		tools.PrintResource(t, p.Extra)

		if p.ID == policy.ID {
			found = true
		}
	}

	th.AssertEquals(t, true, found)

	listOpts.Filters = map[string]string{
		"type__contains": "json",
	}

	allPages, err = policies.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err = policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	found = false
	for _, p := range allPolicies {
		tools.PrintResource(t, p)
		tools.PrintResource(t, p.Extra)

		if p.ID == policy.ID {
			found = true
		}
	}

	th.AssertEquals(t, true, found)

	listOpts.Filters = map[string]string{
		"type__contains": "foobar",
	}

	allPages, err = policies.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err = policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	found = false
	for _, p := range allPolicies {
		tools.PrintResource(t, p)
		tools.PrintResource(t, p.Extra)

		if p.ID == policy.ID {
			found = true
		}
	}

	th.AssertEquals(t, false, found)

	gotPolicy, err := policies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, policy, gotPolicy)

	updateOpts := policies.UpdateOpts{
		Type: "text/plain",
		Blob: []byte("'foobar_user': 'role:compute-user'"),
		Extra: map[string]any{
			"description": "updated policy for foobar_user",
		},
	}

	updatedPolicy, err := policies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedPolicy)
	tools.PrintResource(t, updatedPolicy.Extra)

	th.AssertEquals(t, updatedPolicy.Type, updateOpts.Type)
	th.AssertEquals(t, updatedPolicy.Blob, string(updateOpts.Blob))
	th.AssertEquals(t, updatedPolicy.Extra["description"], updateOpts.Extra["description"])

	err = policies.Delete(context.TODO(), client, policy.ID).ExtractErr()
	th.AssertNoErr(t, err)
}
