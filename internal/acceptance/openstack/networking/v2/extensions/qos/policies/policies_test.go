//go:build acceptance || networking || qos || policies

package policies

import (
	"context"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPoliciesCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "qos")

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

	_, err = policies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newPolicy, err := policies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPolicy)
	th.AssertEquals(t, newPolicy.Name, newName)
	th.AssertEquals(t, newPolicy.Description, newDescription)

	allPages, err := policies.List(client, nil).AllPages(context.TODO())
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

func TestPoliciesRevision(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "qos")

	// Create a policy
	policy, err := CreateQoSPolicy(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQoSPolicy(t, client, policy.ID)

	tools.PrintResource(t, policy)

	// Store the current revision number.
	oldRevisionNumber := policy.RevisionNumber

	// Update the policy without revision number.
	// This should work.
	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &policies.UpdateOpts{
		Name:        newName,
		Description: &newDescription,
	}
	policy, err = policies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, policy)

	// This should fail due to an old revision number.
	newDescription = "new description"
	updateOpts = &policies.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &oldRevisionNumber,
	}
	_, err = policies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertErr(t, err)
	if !strings.Contains(err.Error(), "RevisionNumberConstraintFailed") {
		t.Fatalf("expected to see an error of type RevisionNumberConstraintFailed, but got the following error instead: %v", err)
	}

	// Reread the policy to show that it did not change.
	policy, err = policies.Get(context.TODO(), client, policy.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, policy)

	// This should work because now we do provide a valid revision number.
	newDescription = "new description"
	updateOpts = &policies.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &policy.RevisionNumber,
	}
	policy, err = policies.Update(context.TODO(), client, policy.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, policy)

	th.AssertEquals(t, policy.Name, newName)
	th.AssertEquals(t, policy.Description, newDescription)
}
