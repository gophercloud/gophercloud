//go:build acceptance || networking || subnetpools

package v2

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/subnetpools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSubnetPoolsCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a subnetpool
	subnetPool, err := CreateSubnetPool(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSubnetPool(t, client, subnetPool.ID)

	tools.PrintResource(t, subnetPool)

	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &subnetpools.UpdateOpts{
		Name:        newName,
		Description: &newDescription,
	}

	_, err = subnetpools.Update(context.TODO(), client, subnetPool.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newSubnetPool, err := subnetpools.Get(context.TODO(), client, subnetPool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newSubnetPool)
	th.AssertEquals(t, newSubnetPool.Name, newName)
	th.AssertEquals(t, newSubnetPool.Description, newDescription)

	allPages, err := subnetpools.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allSubnetPools, err := subnetpools.ExtractSubnetPools(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, subnetpool := range allSubnetPools {
		if subnetpool.ID == newSubnetPool.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestSubnetPoolsAddRemovePrefixes(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a subnetpool with an initial prefix.
	subnetPool, err := CreateSubnetPool(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSubnetPool(t, client, subnetPool.ID)

	tools.PrintResource(t, subnetPool)

	// Add new prefixes to the pool.
	addOpts := subnetpools.PrefixesOpsOpts{
		Prefixes: []string{
			"172.16.0.0/16",
			"192.168.0.0/16",
		},
	}

	t.Logf("Attempting to add prefixes to subnetpool %s", subnetPool.ID)

	addedPrefixes, err := subnetpools.AddPrefixes(context.TODO(), client, subnetPool.ID, addOpts).Extract()
	th.AssertNoErr(t, err)

	t.Logf("Successfully added prefixes to subnetpool %s: %v", subnetPool.ID, addedPrefixes)

	updatedPool, err := subnetpools.Get(context.TODO(), client, subnetPool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedPool)

	th.AssertEquals(t, true, slices.Contains(updatedPool.Prefixes, "172.16.0.0/16"))
	th.AssertEquals(t, true, slices.Contains(updatedPool.Prefixes, "192.168.0.0/16"))
	th.AssertEquals(t, true, slices.Contains(updatedPool.Prefixes, "10.0.0.0/8"))

	// Remove one of the prefixes.
	removeOpts := subnetpools.PrefixesOpsOpts{
		Prefixes: []string{
			"192.168.0.0/16",
		},
	}

	t.Logf("Attempting to remove prefixes from subnetpool %s", subnetPool.ID)

	remainingPrefixes, err := subnetpools.RemovePrefixes(context.TODO(), client, subnetPool.ID, removeOpts).Extract()
	th.AssertNoErr(t, err)

	t.Logf("Successfully removed prefixes from subnetpool %s: %v", subnetPool.ID, remainingPrefixes)

	// Verify via Get that the subnetpool reflects the removal.
	finalPool, err := subnetpools.Get(context.TODO(), client, subnetPool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, finalPool)

	th.AssertEquals(t, false, slices.Contains(finalPool.Prefixes, "192.168.0.0/16"))
	th.AssertEquals(t, true, slices.Contains(finalPool.Prefixes, "172.16.0.0/16"))
	th.AssertEquals(t, true, slices.Contains(finalPool.Prefixes, "10.0.0.0/8"))
}

func TestSubnetPoolsRevision(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a subnetpool
	subnetPool, err := CreateSubnetPool(t, client)
	th.AssertNoErr(t, err)
	defer DeleteSubnetPool(t, client, subnetPool.ID)

	// Store the current revision number.
	oldRevisionNumber := subnetPool.RevisionNumber

	// Update the subnet pool without revision number.
	// This should work.
	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &subnetpools.UpdateOpts{
		Name:        newName,
		Description: &newDescription,
	}
	subnetPool, err = subnetpools.Update(context.TODO(), client, subnetPool.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, subnetPool)

	// This should fail due to an old revision number.
	newDescription = "new description"
	updateOpts = &subnetpools.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &oldRevisionNumber,
	}
	_, err = subnetpools.Update(context.TODO(), client, subnetPool.ID, updateOpts).Extract()
	th.AssertErr(t, err)
	if !strings.Contains(err.Error(), "RevisionNumberConstraintFailed") {
		t.Fatalf("expected to see an error of type RevisionNumberConstraintFailed, but got the following error instead: %v", err)
	}

	// Reread the subnet pool to show that it did not change.
	subnetPool, err = subnetpools.Get(context.TODO(), client, subnetPool.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, subnetPool)

	// This should work because now we do provide a valid revision number.
	newDescription = "new description"
	updateOpts = &subnetpools.UpdateOpts{
		Name:           newName,
		Description:    &newDescription,
		RevisionNumber: &subnetPool.RevisionNumber,
	}
	subnetPool, err = subnetpools.Update(context.TODO(), client, subnetPool.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, subnetPool)

	th.AssertEquals(t, subnetPool.Name, newName)
	th.AssertEquals(t, subnetPool.Description, newDescription)
}
