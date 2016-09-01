// +build acceptance networking layer3 floatingips

package layer3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func TestLayer3FloatingIPsList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	listOpts := floatingips.ListOpts{}
	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to list floating IPs: %v", err)
	}

	allFIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		t.Fatalf("Unable to extract floating IPs: %v", err)
	}

	for _, fip := range allFIPs {
		PrintFloatingIP(t, &fip)
	}
}

func TestLayer3FloatingIPsCreateDelete(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		t.Fatalf("Unable to get choices: %v", err)
	}

	subnet, err := networking.CreateSubnet(t, client, choices.ExternalNetworkID)
	if err != nil {
		t.Fatalf("Unable to create subnet: %v", err)
	}
	defer networking.DeleteSubnet(t, client, subnet.ID)

	router, err := CreateExternalRouter(t, client)
	if err != nil {
		t.Fatalf("Unable to create router: %v", err)
	}
	defer DeleteRouter(t, client, router.ID)

	aiOpts := routers.AddInterfaceOpts{
		SubnetID: subnet.ID,
	}

	iface, err := routers.AddInterface(client, router.ID, aiOpts).Extract()
	if err != nil {
		t.Fatalf("Unable to add interface to router: %v", err)
	}

	PrintRouter(t, router)
	PrintRouterInterface(t, iface)

	port, err := networking.CreatePort(t, client, choices.ExternalNetworkID, subnet.ID)
	if err != nil {
		t.Fatalf("Unable to create port: %v", err)
	}
	defer networking.DeletePort(t, client, port.ID)

	fip, err := CreateFloatingIP(t, client, choices.ExternalNetworkID, port.ID)
	if err != nil {
		t.Fatalf("Unable to create floating IP: %v", err)
	}

	newFip, err := floatingips.Get(client, fip.ID).Extract()
	if err != nil {
		t.Fatalf("Unable to get floating ip: %v", err)
	}

	PrintFloatingIP(t, newFip)

	DeleteFloatingIP(t, client, fip.ID)

	riOpts := routers.RemoveInterfaceOpts{
		SubnetID: subnet.ID,
	}

	_, err = routers.RemoveInterface(client, router.ID, riOpts).Extract()
	if err != nil {
		t.Fatalf("Failed to remove interface from router: %v", err)
	}
}
