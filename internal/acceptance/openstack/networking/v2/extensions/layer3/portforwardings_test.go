//go:build acceptance || networking || layer3 || portforwardings

package layer3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestLayer3PortForwardingsCreateDelete(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	networking.RequireNeutronExtension(t, client, "floating-ip-port-forwarding")

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	router, err := CreateExternalRouter(t, client)
	th.AssertNoErr(t, err)
	defer DeleteRouter(t, client, router.ID)

	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	// not required, since "DeleteRouterInterface" below removes the port
	// defer networking.DeletePort(t, client, port.ID)

	_, err = CreateRouterInterface(t, client, port.ID, router.ID)
	th.AssertNoErr(t, err)
	defer DeleteRouterInterface(t, client, port.ID, router.ID)

	fip, err := CreateFloatingIP(t, client, choices.ExternalNetworkID, "")
	th.AssertNoErr(t, err)
	defer DeleteFloatingIP(t, client, fip.ID)

	newFip, err := floatingips.Get(context.TODO(), client, fip.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newFip)

	pf, err := CreatePortForwarding(t, client, fip.ID, port.ID, port.FixedIPs)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "Test description", pf.Description)
	defer DeletePortForwarding(t, client, fip.ID, pf.ID)
	tools.PrintResource(t, pf)

	pfRange, err := CreatePortRangeForwarding(t, client, fip.ID, port.ID, port.FixedIPs)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "Test description range", pfRange.Description)
	defer DeletePortForwarding(t, client, fip.ID, pfRange.ID)
	tools.PrintResource(t, pfRange)

	// Test updating port
	newPf, err := portforwarding.Get(context.TODO(), client, fip.ID, pf.ID).Extract()
	th.AssertNoErr(t, err)

	updateOpts := portforwarding.UpdateOpts{
		Description:  new(string),
		Protocol:     "udp",
		InternalPort: 30,
		ExternalPort: 678,
	}

	_, err = portforwarding.Update(context.TODO(), client, fip.ID, newPf.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newPf, err = portforwarding.Get(context.TODO(), client, fip.ID, pf.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", newPf.Description)
	th.AssertEquals(t, "udp", newPf.Protocol)
	th.AssertEquals(t, 30, newPf.InternalPort)
	th.AssertEquals(t, 678, newPf.ExternalPort)

	// Test updating port range
	newRangePf, err := portforwarding.Get(context.TODO(), client, fip.ID, pfRange.ID).Extract()
	th.AssertNoErr(t, err)

	updateOpts = portforwarding.UpdateOpts{
		Description:       new(string),
		Protocol:          "udp",
		InternalPortRange: "1400:1499",
		ExternalPortRange: "1500:1599",
	}

	_, err = portforwarding.Update(context.TODO(), client, fip.ID, newRangePf.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	newRangePf, err = portforwarding.Get(context.TODO(), client, fip.ID, pfRange.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", newRangePf.Description)
	th.AssertEquals(t, "udp", newRangePf.Protocol)
	th.AssertEquals(t, "1400:1499", newRangePf.InternalPortRange)
	th.AssertEquals(t, "1500:1599", newRangePf.ExternalPortRange)

	allPages, err := portforwarding.List(client, portforwarding.ListOpts{}, fip.ID).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPFs, err := portforwarding.ExtractPortForwardings(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, pf := range allPFs {
		if pf.ID == newPf.ID {
			found = true
		}
	}

	th.AssertTrue(t, found)

	found = false
	for _, pf := range allPFs {
		if pf.ID == newRangePf.ID {
			found = true
		}
	}

	th.AssertTrue(t, found)

}
