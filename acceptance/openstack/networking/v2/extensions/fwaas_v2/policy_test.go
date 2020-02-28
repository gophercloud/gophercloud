// +build acceptance networking fwaas_v2

package fwaas_v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestPolicyCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create First Rule. This will be used as part of the Policy creation
	rule, err := CreateRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRule(t, client, rule.ID)

	tools.PrintResource(t, rule)

	// Create Second rule. This will be injected in to the policy after its creation
	ruleToInsert, err := CreateRule(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRule(t, client, ruleToInsert.ID)

	tools.PrintResource(t, ruleToInsert)

	// Create the Policy using the first rule
	policy, err := CreatePolicy(t, client, rule.ID)
	th.AssertNoErr(t, err)
	defer DeletePolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	// Inject the second rule
	AddRule(t, client, policy.ID, ruleToInsert.ID, rule.ID)

	// Remove the first rule
	RemoveRule(t, client, policy.ID, rule.ID)

	name := ""
	description := ""
	updateOpts := policies.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	_, err = policies.Update(client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newPolicy, err := policies.Get(client, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPolicy)
	th.AssertEquals(t, newPolicy.Name, name)
	th.AssertEquals(t, newPolicy.Description, description)

	allPages, err := policies.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allPolicies, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, policy := range allPolicies {
		if policy.ID == newPolicy.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}
