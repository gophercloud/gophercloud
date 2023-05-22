//go:build acceptance || trunks
// +build acceptance trunks

package trunk_details

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	v2 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	v2Trunks "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunk_details"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/testhelper"
)

type portWithTrunkDetails struct {
	ports.Port
	trunk_details.TrunkDetailsExt
}

func TestListPortWithSubports(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	_, err = extensions.Get(client, "trunk-details").Extract()
	if err != nil {
		t.Skip("This test requires trunk-details Neutron extension")
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

	trunk, err := v2Trunks.CreateTrunk(t, client, parentPort.ID, subport1.ID, subport2.ID)
	if err != nil {
		t.Fatalf("Unable to create trunk: %v", err)
	}
	defer v2Trunks.DeleteTrunk(t, client, trunk.ID)

	// Test LIST ports with trunk details
	allPages, err := ports.List(client, ports.ListOpts{ID: parentPort.ID}).AllPages()
	th.AssertNoErr(t, err)

	var allPorts []portWithTrunkDetails
	err = ports.ExtractPortsInto(allPages, &allPorts)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, len(allPorts))
	port := allPorts[0]

	th.AssertEquals(t, trunk.ID, port.TrunkDetails.TrunkID)
	th.AssertEquals(t, 2, len(port.TrunkDetails.SubPorts))

	// Note that MAC address is not (currently) returned in list queries. We
	// exclude it from the comparison here in case it's ever added. MAC
	// address is returned in GET queries, so we do assert that in the GET
	// test below.
	th.AssertDeepEquals(t, trunks.Subport{
		SegmentationID:   1,
		SegmentationType: "vlan",
		PortID:           subport1.ID,
	}, port.TrunkDetails.SubPorts[0].Subport)
	th.AssertDeepEquals(t, trunks.Subport{
		SegmentationID:   2,
		SegmentationType: "vlan",
		PortID:           subport2.ID,
	}, port.TrunkDetails.SubPorts[1].Subport)

	// Test GET port with trunk details
	err = ports.Get(client, parentPort.ID).ExtractInto(&port)
	th.AssertEquals(t, trunk.ID, port.TrunkDetails.TrunkID)
	th.AssertEquals(t, 2, len(port.TrunkDetails.SubPorts))
	th.AssertDeepEquals(t, trunk_details.Subport{
		Subport: trunks.Subport{
			SegmentationID:   1,
			SegmentationType: "vlan",
			PortID:           subport1.ID,
		},
		MACAddress: subport1.MACAddress,
	}, port.TrunkDetails.SubPorts[0])
	th.AssertDeepEquals(t, trunk_details.Subport{
		Subport: trunks.Subport{
			SegmentationID:   2,
			SegmentationType: "vlan",
			PortID:           subport2.ID,
		},
		MACAddress: subport2.MACAddress,
	}, port.TrunkDetails.SubPorts[1])
}
