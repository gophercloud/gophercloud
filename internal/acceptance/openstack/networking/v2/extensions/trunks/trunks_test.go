//go:build acceptance || networking || trunks

package trunks

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	v2 "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTrunkCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "trunk")

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID, subport1.ID, subport2.ID)
	th.AssertNoErr(t, err)
	defer DeleteTrunk(t, client, trunk.ID)

	_, err = trunks.Get(context.TODO(), client, trunk.ID).Extract()
	th.AssertNoErr(t, err)

	// Update Trunk
	name := ""
	description := ""
	updateOpts := trunks.UpdateOpts{
		Name:        &name,
		Description: &description,
	}
	updatedTrunk, err := trunks.Update(context.TODO(), client, trunk.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	if trunk.Name == updatedTrunk.Name {
		t.Fatalf("Trunk name was not updated correctly")
	}

	if trunk.Description == updatedTrunk.Description {
		t.Fatalf("Trunk description was not updated correctly")
	}

	th.AssertDeepEquals(t, updatedTrunk.Name, name)
	th.AssertDeepEquals(t, updatedTrunk.Description, description)

	// Get subports
	subports, err := trunks.GetSubports(context.TODO(), client, trunk.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, trunk.Subports[0], subports[0])
	th.AssertDeepEquals(t, trunk.Subports[1], subports[1])

	tools.PrintResource(t, trunk)
}

func TestTrunkList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "trunk")

	allPages, err := trunks.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTrunks, err := trunks.ExtractTrunks(allPages)
	th.AssertNoErr(t, err)

	for _, trunk := range allTrunks {
		tools.PrintResource(t, trunk)
	}
}

func TestTrunkSubportOperation(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "trunk")

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID)
	th.AssertNoErr(t, err)
	defer DeleteTrunk(t, client, trunk.ID)

	// Add subports to the trunk
	addSubportsOpts := trunks.AddSubportsOpts{
		Subports: []trunks.Subport{
			{
				SegmentationID:   1,
				SegmentationType: "vlan",
				PortID:           subport1.ID,
			},
			{
				SegmentationID:   11,
				SegmentationType: "vlan",
				PortID:           subport2.ID,
			},
		},
	}
	updatedTrunk, err := trunks.AddSubports(context.TODO(), client, trunk.ID, addSubportsOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(updatedTrunk.Subports))
	th.AssertDeepEquals(t, addSubportsOpts.Subports[0], updatedTrunk.Subports[0])
	th.AssertDeepEquals(t, addSubportsOpts.Subports[1], updatedTrunk.Subports[1])

	// Remove the Subports from the trunk
	subRemoveOpts := trunks.RemoveSubportsOpts{
		Subports: []trunks.RemoveSubport{
			{PortID: subport1.ID},
			{PortID: subport2.ID},
		},
	}
	updatedAgainTrunk, err := trunks.RemoveSubports(context.TODO(), client, trunk.ID, subRemoveOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, trunk.Subports, updatedAgainTrunk.Subports)
}

func TestTrunkTags(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "trunk")

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID, subport1.ID, subport2.ID)
	th.AssertNoErr(t, err)
	defer DeleteTrunk(t, client, trunk.ID)

	tagReplaceAllOpts := attributestags.ReplaceAllOpts{
		// docs say list of tags, but it's a set e.g no duplicates
		Tags: []string{"a", "b", "c"},
	}
	_, err = attributestags.ReplaceAll(context.TODO(), client, "trunks", trunk.ID, tagReplaceAllOpts).Extract()
	th.AssertNoErr(t, err)

	gtrunk, err := trunks.Get(context.TODO(), client, trunk.ID).Extract()
	th.AssertNoErr(t, err)
	tags := gtrunk.Tags
	sort.Strings(tags) // Ensure ordering, older OpenStack versions aren't sorted...
	th.AssertDeepEquals(t, []string{"a", "b", "c"}, tags)

	// Add a tag
	err = attributestags.Add(context.TODO(), client, "trunks", trunk.ID, "d").ExtractErr()
	th.AssertNoErr(t, err)

	// Delete a tag
	err = attributestags.Delete(context.TODO(), client, "trunks", trunk.ID, "a").ExtractErr()
	th.AssertNoErr(t, err)

	// Verify expected tags are set in the List response
	tags, err = attributestags.List(context.TODO(), client, "trunks", trunk.ID).Extract()
	th.AssertNoErr(t, err)
	sort.Strings(tags)
	th.AssertDeepEquals(t, []string{"b", "c", "d"}, tags)

	// Delete all tags
	err = attributestags.DeleteAll(context.TODO(), client, "trunks", trunk.ID).ExtractErr()
	th.AssertNoErr(t, err)
	tags, err = attributestags.List(context.TODO(), client, "trunks", trunk.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(tags))
}

func TestTrunkRevision(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	v2.RequireNeutronExtension(t, client, "trunk")

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID, subport1.ID, subport2.ID)
	th.AssertNoErr(t, err)
	defer DeleteTrunk(t, client, trunk.ID)

	tools.PrintResource(t, trunk)

	// Store the current revision number.
	oldRevisionNumber := trunk.RevisionNumber

	// Update the trunk without revision number.
	// This should work.
	newName := tools.RandomString("TESTACC-", 8)
	newDescription := ""
	updateOpts := &trunks.UpdateOpts{
		Name:        &newName,
		Description: &newDescription,
	}
	trunk, err = trunks.Update(context.TODO(), client, trunk.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, trunk)

	// This should fail due to an old revision number.
	newDescription = "new description"
	updateOpts = &trunks.UpdateOpts{
		Name:           &newName,
		Description:    &newDescription,
		RevisionNumber: &oldRevisionNumber,
	}
	_, err = trunks.Update(context.TODO(), client, trunk.ID, updateOpts).Extract()
	th.AssertErr(t, err)
	if !strings.Contains(err.Error(), "RevisionNumberConstraintFailed") {
		t.Fatalf("expected to see an error of type RevisionNumberConstraintFailed, but got the following error instead: %v", err)
	}

	// Reread the trunk to show that it did not change.
	trunk, err = trunks.Get(context.TODO(), client, trunk.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, trunk)

	// This should work because now we do provide a valid revision number.
	newDescription = "new description"
	updateOpts = &trunks.UpdateOpts{
		Name:           &newName,
		Description:    &newDescription,
		RevisionNumber: &trunk.RevisionNumber,
	}
	trunk, err = trunks.Update(context.TODO(), client, trunk.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, trunk)

	th.AssertEquals(t, trunk.Name, newName)
	th.AssertEquals(t, trunk.Description, newDescription)
}
