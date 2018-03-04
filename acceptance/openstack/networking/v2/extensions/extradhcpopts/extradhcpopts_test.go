// +build acceptance networking extradhcpopts

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestPortsWithDHCPOptsCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create a Network
	network, err := v2.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create a network: %v", err)
	}
	defer v2.DeleteNetwork(t, client, network.ID)

	// Create a Subnet
	subnet, err := v2.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create a subnet: %v", err)
	}
	defer v2.DeleteSubnet(t, client, subnet.ID)

	// Create a port with extra DHCP options.
	port, err := CreatePortWithDHCPOpts(t, client, network.ID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create a port: %v", err)
	}
	defer v2.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)
}
