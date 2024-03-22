//go:build acceptance || networking || vpnaas

package vpnaas

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestIPSecPolicyList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	allPages, err := ipsecpolicies.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPolicies, err := ipsecpolicies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	for _, policy := range allPolicies {
		tools.PrintResource(t, policy)
	}
}

func TestIPSecPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	policy, err := CreateIPSecPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteIPSecPolicy(t, client, policy.ID)
	tools.PrintResource(t, policy)

	updatedDescription := "Updated policy description"
	updateOpts := ipsecpolicies.UpdateOpts{
		Description: &updatedDescription,
	}

	policy, err = ipsecpolicies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, policy)

	newPolicy, err := ipsecpolicies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, newPolicy)
}
