//go:build acceptance || networking || dns

package dns

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions/layer3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestDNSPortCRUDL(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "dns-integration")

	// Create Network
	networkDNSDomain := "local."
	network, err := CreateNetworkDNS(t, client, networkDNSDomain)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Create port
	portDNSName := "port"
	port, err := CreatePortDNS(t, client, network.ID, subnet.ID, portDNSName)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	if os.Getenv("OS_BRANCH") == "stable/mitaka" {
		// List port successfully
		var listOpts ports.ListOptsBuilder
		listOpts = dns.PortListOptsExt{
			ListOptsBuilder: ports.ListOpts{},
			DNSName:         portDNSName,
		}
		var listedPorts []PortWithDNSExt
		i := 0
		err = ports.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
			i++
			err := ports.ExtractPortsInto(page, &listedPorts)
			if err != nil {
				t.Errorf("Failed to extract ports: %v", err)
				return false, err
			}

			tools.PrintResource(t, listedPorts)

			th.AssertEquals(t, 1, len(listedPorts))
			th.CheckDeepEquals(t, *port, listedPorts[0])

			return true, nil
		})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, i)

		// List port unsuccessfully
		listOpts = dns.PortListOptsExt{
			ListOptsBuilder: ports.ListOpts{},
			DNSName:         "foo",
		}
		i = 0
		err = ports.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
			i++
			err := ports.ExtractPortsInto(page, &listedPorts)
			if err != nil {
				t.Errorf("Failed to extract ports: %v", err)
				return false, err
			}

			tools.PrintResource(t, listedPorts)

			th.AssertEquals(t, 1, len(listedPorts))
			th.CheckDeepEquals(t, *port, listedPorts[0])

			return true, nil
		})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 0, i)
	}

	// Get port
	var getPort PortWithDNSExt
	err = ports.Get(context.TODO(), client, port.ID).ExtractInto(&getPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getPort)
	th.AssertDeepEquals(t, *port, getPort)

	// Update port
	newPortName := ""
	newPortDescription := ""
	newDNSName := ""
	portUpdateOpts := ports.UpdateOpts{
		Name:        &newPortName,
		Description: &newPortDescription,
	}
	updateOpts := dns.PortUpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		DNSName:           &newDNSName,
	}

	var newPort PortWithDNSExt
	err = ports.Update(context.TODO(), client, port.ID, updateOpts).ExtractInto(&newPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)
	th.AssertEquals(t, newPort.Description, newPortName)
	th.AssertEquals(t, newPort.Description, newPortDescription)
	th.AssertEquals(t, newPort.DNSName, newDNSName)

	// Get updated port
	var getNewPort PortWithDNSExt
	err = ports.Get(context.TODO(), client, port.ID).ExtractInto(&getNewPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getNewPort)
	// workaround for update race condition
	newPort.DNSAssignment = nil
	getNewPort.DNSAssignment = nil
	th.AssertDeepEquals(t, newPort, getNewPort)
}

func TestDNSFloatingIPCRUD(t *testing.T) {
	t.Skip("Skipping TestDNSFloatingIPCRUD for now, as it doesn't work with ML2/OVN.")
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "dns-integration")

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	// Create Network
	networkDNSDomain := "local."
	network, err := CreateNetworkDNS(t, client, networkDNSDomain)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	// Create Router
	router, err := layer3.CreateExternalRouter(t, client)
	th.AssertNoErr(t, err)
	defer layer3.DeleteRouter(t, client, router.ID)

	// Create router interface
	routerPort, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	_, err = layer3.CreateRouterInterface(t, client, routerPort.ID, router.ID)
	th.AssertNoErr(t, err)
	defer layer3.DeleteRouterInterface(t, client, routerPort.ID, router.ID)

	// Create port
	portDNSName := "port"
	port, err := CreatePortDNS(t, client, network.ID, subnet.ID, portDNSName)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Create floating IP
	fipDNSName := "fip"
	fipDNSDomain := "local."
	fip, err := CreateFloatingIPDNS(t, client, choices.ExternalNetworkID, port.ID, fipDNSName, fipDNSDomain)
	th.AssertNoErr(t, err)
	defer layer3.DeleteFloatingIP(t, client, fip.ID)

	// Get floating IP
	var getFip FloatingIPWithDNSExt
	err = floatingips.Get(context.TODO(), client, fip.ID).ExtractInto(&getFip)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getFip)
	th.AssertDeepEquals(t, *fip, getFip)
}

func TestDNSNetwork(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "dns-integration")

	// Create Network
	networkDNSDomain := "local."
	network, err := CreateNetworkDNS(t, client, networkDNSDomain)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// Get network
	var getNetwork NetworkWithDNSExt
	err = networks.Get(context.TODO(), client, network.ID).ExtractInto(&getNetwork)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getNetwork)
	th.AssertDeepEquals(t, *network, getNetwork)

	// Update network
	newNetworkName := ""
	newNetworkDescription := ""
	newNetworkDNSDomain := ""
	networkUpdateOpts := networks.UpdateOpts{
		Name:        &newNetworkName,
		Description: &newNetworkDescription,
	}
	updateOpts := dns.NetworkUpdateOptsExt{
		UpdateOptsBuilder: networkUpdateOpts,
		DNSDomain:         &newNetworkDNSDomain,
	}

	var newNetwork NetworkWithDNSExt
	err = networks.Update(context.TODO(), client, network.ID, updateOpts).ExtractInto(&newNetwork)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newNetwork)
	th.AssertEquals(t, newNetwork.Description, newNetworkName)
	th.AssertEquals(t, newNetwork.Description, newNetworkDescription)
	th.AssertEquals(t, newNetwork.DNSDomain, newNetworkDNSDomain)

	// Get updated network
	var getNewNetwork NetworkWithDNSExt
	err = networks.Get(context.TODO(), client, network.ID).ExtractInto(&getNewNetwork)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getNewNetwork)
	th.AssertDeepEquals(t, newNetwork, getNewNetwork)
}
