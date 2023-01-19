//go:build acceptance || networking || qos || policies
// +build acceptance networking qos policies

package policies

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/common/extensions"
	"github.com/bizflycloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestPoliciesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	extension, err := extensions.Get(client, "qos").Extract()
	if err != nil {
		t.Skip("This test requires qos Neutron extension")
	}
	tools.PrintResource(t, extension)

	// Create a QoS policy.
	policy, err := CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoSPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &policies.UpdateOpts{
		Name:        newName,
		Description: &newDescription,
	}

	_, err = policies.Update(client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newPolicy, err := policies.Get(client, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPolicy)
	th.AssertEquals(t, newPolicy.Name, newName)
	th.AssertEquals(t, newPolicy.Description, newDescription)

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
