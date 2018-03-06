// +build acceptance networking vpnaas

package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networks "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	layer3 "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"

	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestConnectionCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a network client: %v", err)
	}

	// Create Network
	network, err := networks.CreateNetwork(t, client)
	if err != nil {
		t.Fatalf("Unable to create network: %v", err)
	}

	// Create Subnet
	subnet, err := networks.CreateSubnet(t, client, network.ID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}

	router, err := layer3.CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}

	// Link router and subnet
	aiOpts := routers.AddInterfaceOpts{
		SubnetID: subnet.ID,
	}

	_, err = routers.AddInterface(client, router.ID, aiOpts).Extract()
	if err != nil {
		t.Fatalf("Failed to add interface to router: %v", err)
	}

	// Create all needed resources for the connection
	service, err := CreateService(t, client, router.ID)
	if err != nil {
		t.Fatalf("Unable to create service: %v", err)
	}

	ikepolicy, err := CreateIKEPolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create IKE policy: %v", err)
	}

	ipsecpolicy, err := CreateIPSecPolicy(t, client)
	if err != nil {
		t.Fatalf("Unable to create IPSec Policy: %v", err)
	}

	peerEPGroup, err := CreateEndpointGroup(t, client)
	if err != nil {
		t.Fatalf("Unable to create Endpoint Group with CIDR endpoints: %v", err)
	}

	localEPGroup, err := CreateEndpointGroupWithSubnet(t, client, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create Endpoint Group with subnet endpoints: %v", err)
	}

	conn, err := CreateSiteConnection(t, client, ikepolicy.ID, ipsecpolicy.ID, service.ID, peerEPGroup.ID, localEPGroup.ID)
	if err != nil {
		t.Fatalf("Unable to create IPSec Site Connection: %v", err)
	}

	tools.PrintResource(t, conn)

}
