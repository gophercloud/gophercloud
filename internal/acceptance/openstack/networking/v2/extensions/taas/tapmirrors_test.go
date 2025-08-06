//go:build acceptance || networking || taas

package taas

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTapMirrorCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "taas")

	// Create Port
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	// Create Tap Mirror
	mirror, err := CreateTapMirror(t, client, port.ID, port.FixedIPs[0].IPAddress)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, mirror)
}
