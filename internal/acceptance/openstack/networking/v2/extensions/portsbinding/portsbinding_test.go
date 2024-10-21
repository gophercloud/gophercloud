//go:build acceptance || networking || portbinding

package portsbinding

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPortsbindingCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Define a host
	hostID := "localhost"
	profile := map[string]any{"foo": "bar"}

	// Create port
	port, err := CreatePortsbinding(t, client, network.ID, subnet.ID, hostID, profile)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)
	th.AssertEquals(t, hostID, port.HostID)
	th.AssertEquals(t, "normal", port.VNICType)
	th.AssertDeepEquals(t, profile, port.Profile)

	// Update port
	newPortName := ""
	newPortDescription := ""
	newHostID := "127.0.0.1"
	newProfile := map[string]any{}
	updateOpts := ports.UpdateOpts{
		Name:        &newPortName,
		Description: &newPortDescription,
	}

	finalUpdateOpts := portsbinding.UpdateOptsExt{
		UpdateOptsBuilder: updateOpts,
		HostID:            &newHostID,
		Profile:           newProfile,
	}

	var newPort PortWithBindingExt

	_, err = ports.Update(context.TODO(), client, port.ID, finalUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	// Read the updated port
	err = ports.Get(context.TODO(), client, port.ID).ExtractInto(&newPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)
	th.AssertEquals(t, newPortName, newPort.Description)
	th.AssertEquals(t, newPortDescription, newPort.Description)
	th.AssertEquals(t, newHostID, newPort.HostID)
	th.AssertEquals(t, "normal", newPort.VNICType)
	th.AssertDeepEquals(t, newProfile, newPort.Profile)
}
