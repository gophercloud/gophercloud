//go:build acceptance || networking || uplinkstatuspropagation

package uplinkstatuspropagation

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/uplinkstatuspropagation"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestUplinkStatusPropagationCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "uplink-status-propagation")

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Create port with uplink status propagation enabled
	port, err := CreatePortWithUplinkStatusPropagation(t, client, network.ID, subnet.ID, true)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)
	if port.PropagateUplinkStatus != nil {
		th.AssertEquals(t, *port.PropagateUplinkStatus, true)
	}

	// Update port to disable uplink status propagation
	newPortName := tools.RandomString("TESTACC-UPDATED-", 8)
	iFalse := false
	updateOpts := ports.UpdateOpts{
		Name: &newPortName,
	}

	finalUpdateOpts := uplinkstatuspropagation.PortPropagateUplinkStatusUpdateOptsExt{
		UpdateOptsBuilder:     updateOpts,
		PropagateUplinkStatus: &iFalse,
	}

	var updatedPort PortWithUplinkStatusPropagationExt

	_, err = ports.Update(context.TODO(), client, port.ID, finalUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	// Read the updated port
	err = ports.Get(context.TODO(), client, port.ID).ExtractInto(&updatedPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updatedPort)
	th.AssertEquals(t, newPortName, updatedPort.Name)
	if updatedPort.PropagateUplinkStatus != nil {
		th.AssertEquals(t, *updatedPort.PropagateUplinkStatus, false)
	}
}

func TestUplinkStatusPropagationList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "uplink-status-propagation")

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Create port with uplink status propagation enabled
	port, err := CreatePortWithUplinkStatusPropagation(t, client, network.ID, subnet.ID, true)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	// List ports with filter for propagate_uplink_status=true
	iTrue := true
	listOpts := uplinkstatuspropagation.PortPropagateUplinkStatusListOptsExt{
		ListOptsBuilder:       ports.ListOpts{NetworkID: network.ID},
		PropagateUplinkStatus: &iTrue,
	}

	type PortWithExt struct {
		ports.Port
		uplinkstatuspropagation.PortPropagateUplinkStatusExt
	}

	var allPorts []PortWithExt

	allPages, err := ports.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = ports.ExtractPortsInto(allPages, &allPorts)
	th.AssertNoErr(t, err)

	var found bool
	for _, p := range allPorts {
		if p.ID == port.ID {
			found = true
			if p.PropagateUplinkStatus != nil {
				th.AssertEquals(t, *p.PropagateUplinkStatus, true)
			}
			break
		}
	}

	th.AssertEquals(t, found, true)
}
