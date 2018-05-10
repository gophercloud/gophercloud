// +build acceptance

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/policies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestPoliciesList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := policies.List(client, policies.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allPolicies, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestPoliciesCreate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	createOpts := policies.CreateOpts{
		Type: "application/json",
		Blob: []byte("{'foobar_user': 'role:compute-user'}"),
		Extra: map[string]interface{}{
			"description": "policy for foobar_user",
		},
	}

	policy, err := policies.Create(client, &createOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, policy)
	tools.PrintResource(t, policy.Extra)

	th.AssertEquals(t, policy.Type, createOpts.Type)
	th.AssertEquals(t, policy.Blob, string(createOpts.Blob))
	th.AssertEquals(t, policy.Extra["description"], createOpts.Extra["description"])
}
