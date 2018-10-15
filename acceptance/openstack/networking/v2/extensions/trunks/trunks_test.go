// +build acceptance trunks

package trunks

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestTrunkCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID, subport1.ID, subport2.ID)
	if err != nil {
		t.Fatalf("Unable to create trunk: %v", err)
	}
	defer DeleteTrunk(t, client, trunk.ID)

	_, err = trunks.Get(client, trunk.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get trunk: %v", err)
	}

	// Update Trunk
	updateOpts := trunks.UpdateOpts{
		Name: "updated_gophertrunk",
	}
	updatedTrunk, err := trunks.Update(client, trunk.ID, updateOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to update trunk: %v", err)
	}

	if trunk.Name == updatedTrunk.Name {
		t.Fatalf("Trunk name was not updated correctly")
	}

	// Get subports
	subports, err := trunks.GetSubports(client, trunk.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get subports from the Trunk: %v", err)
	}
	th.AssertDeepEquals(t, trunk.Subports[0], subports[0])
	th.AssertDeepEquals(t, trunk.Subports[1], subports[1])

	tools.PrintResource(t, trunk)
}

func TestTrunkList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	allPages, err := trunks.List(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to list trunks: %v", err)
	}

	allTrunks, err := trunks.ExtractTrunks(allPages)
	if err != nil {
		t.Fatalf("Unable to extract trunks: %v", err)
	}

	for _, trunk := range allTrunks {
		tools.PrintResource(t, trunk)
	}
}

func TestTrunkSubportOperation(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create Network
	network, err := v2.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create port
	parentPort, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, parentPort.ID)

	subport1, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, subport1.ID)

	subport2, err := v2.CreatePort(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer v2.DeletePort(t, client, subport2.ID)

	trunk, err := CreateTrunk(t, client, parentPort.ID)
	if err != nil {
		t.Fatalf("Unable to create trunk: %v", err)
	}
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
	updatedTrunk, err := trunks.AddSubports(client, trunk.ID, addSubportsOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to add subports to the Trunk: %v", err)
	}
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
	updatedAgainTrunk, err := trunks.RemoveSubports(client, trunk.ID, subRemoveOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to remove subports from the Trunk: %v", err)
	}
	th.AssertDeepEquals(t, trunk.Subports, updatedAgainTrunk.Subports)
}
