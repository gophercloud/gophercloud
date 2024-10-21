//go:build acceptance || networking || ports

package v2

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	extensions "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2/extensions"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestPortsCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	if len(port.SecurityGroups) != 1 {
		t.Logf("WARNING: Port did not have a default security group applied")
	}

	tools.PrintResource(t, port)

	// Update port
	newPortName := ""
	newPortDescription := ""
	updateOpts := ports.UpdateOpts{
		Name:        &newPortName,
		Description: &newPortDescription,
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	th.AssertEquals(t, newPort.Name, newPortName)
	th.AssertEquals(t, newPort.Description, newPortDescription)

	allPages, err := ports.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allPorts, err := ports.ExtractPorts(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, port := range allPorts {
		if port.ID == newPort.ID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)

	ipAddress := port.FixedIPs[0].IPAddress
	t.Logf("Port has IP address: %s", ipAddress)

	// List ports by fixed IP
	// All of the following listOpts should return the port
	for _, tt := range []struct {
		name          string
		opts          ports.ListOpts
		expectedPorts int
	}{
		{
			name: "Port ID",
			opts: ports.ListOpts{
				ID: port.ID,
			},
			expectedPorts: 1,
		},
		{
			name: "Network ID",
			opts: ports.ListOpts{
				NetworkID: port.NetworkID,
			},
			expectedPorts: 2, // Will also return DHCP port
		},
		{
			name: "Subnet ID",
			opts: ports.ListOpts{
				FixedIPs: []ports.FixedIPOpts{
					{SubnetID: subnet.ID},
				},
			},
			expectedPorts: 1,
		},
		{
			name: "IP Address",
			opts: ports.ListOpts{
				FixedIPs: []ports.FixedIPOpts{
					{IPAddress: ipAddress},
				},
			},
			expectedPorts: 1,
		},
		{
			name: "Subnet ID and IP Address",
			opts: ports.ListOpts{
				FixedIPs: []ports.FixedIPOpts{
					{SubnetID: subnet.ID, IPAddress: ipAddress},
				},
			},
			expectedPorts: 1,
		},
	} {
		t.Run(fmt.Sprintf("List ports by %s", tt.name), func(t *testing.T) {
			allPages, err := ports.List(client, tt.opts).AllPages(context.TODO())
			th.AssertNoErr(t, err)

			allPorts, err := ports.ExtractPorts(allPages)
			th.AssertNoErr(t, err)

			logPorts := func() {
				for _, port := range allPorts {
					tools.PrintResource(t, port)
				}
			}

			if len(allPorts) != tt.expectedPorts {
				if len(allPorts) == 0 {
					t.Fatalf("Port not found")
				}
				if len(allPorts) > 1 {
					logPorts()
					t.Fatalf("Expected %d port but got %d", tt.expectedPorts, len(allPorts))
				}
			}
			func() {
				for _, port := range allPorts {
					if port.ID == newPort.ID {
						return
					}
				}
				logPorts()
				t.Fatalf("Returned ports did not contain expected port")
			}()
		})
	}
}

func TestPortsRemoveSecurityGroups(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Create a Security Group
	group, err := extensions.CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer extensions.DeleteSecurityGroup(t, client, group.ID)

	// Add the group to the port
	updateOpts := ports.UpdateOpts{
		SecurityGroups: &[]string{group.ID},
	}
	_, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	// Remove the group
	updateOpts = ports.UpdateOpts{
		SecurityGroups: &[]string{},
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	if len(newPort.SecurityGroups) > 0 {
		t.Fatalf("Unable to remove security group from port")
	}
}

func TestPortsDontAlterSecurityGroups(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create a Security Group
	group, err := extensions.CreateSecurityGroup(t, client)
	th.AssertNoErr(t, err)
	defer extensions.DeleteSecurityGroup(t, client, group.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Add the group to the port
	updateOpts := ports.UpdateOpts{
		SecurityGroups: &[]string{group.ID},
	}
	_, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	// Update the port again
	var name = "some_port"
	updateOpts = ports.UpdateOpts{
		Name: &name,
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	if len(newPort.SecurityGroups) == 0 {
		t.Fatalf("Port had security group updated")
	}
}

func TestPortsWithNoSecurityGroup(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePortWithNoSecurityGroup(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	if len(port.SecurityGroups) != 0 {
		t.Fatalf("Port was created with security groups")
	}
}

func TestPortsRemoveAddressPair(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Add an address pair to the port
	updateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "192.168.255.10", MACAddress: "aa:bb:cc:dd:ee:ff"},
		},
	}
	_, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	// Remove the address pair
	updateOpts = ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{},
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	if len(newPort.AllowedAddressPairs) > 0 {
		t.Fatalf("Unable to remove the address pair")
	}
}

func TestPortsDontUpdateAllowedAddressPairs(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Add an address pair to the port
	updateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "192.168.255.10", MACAddress: "aa:bb:cc:dd:ee:ff"},
		},
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	// Remove the address pair
	var name = "some_port"
	updateOpts = ports.UpdateOpts{
		Name: &name,
	}
	newPort, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	if len(newPort.AllowedAddressPairs) == 0 {
		t.Fatalf("Address Pairs were removed")
	}
}

func TestPortsPortSecurityCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePortWithoutPortSecurity(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	var portWithExt struct {
		ports.Port
		portsecurity.PortSecurityExt
	}

	err = ports.Get(context.TODO(), client, port.ID).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, portWithExt)

	iTrue := true
	portUpdateOpts := ports.UpdateOpts{}
	updateOpts := portsecurity.PortUpdateOptsExt{
		UpdateOptsBuilder:   portUpdateOpts,
		PortSecurityEnabled: &iTrue,
	}

	err = ports.Update(context.TODO(), client, port.ID, updateOpts).ExtractInto(&portWithExt)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, portWithExt)
}

func TestPortsWithExtraDHCPOptsCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create a Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create a port with extra DHCP options.
	port, err := CreatePortWithExtraDHCPOpts(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Update the port with extra DHCP options.
	newPortName := tools.RandomString("TESTACC-", 8)
	portUpdateOpts := ports.UpdateOpts{
		Name: &newPortName,
	}

	existingOpt := port.ExtraDHCPOpts[0]
	newOptValue := "test_value_2"

	updateOpts := extradhcpopts.UpdateOptsExt{
		UpdateOptsBuilder: portUpdateOpts,
		ExtraDHCPOpts: []extradhcpopts.UpdateExtraDHCPOpt{
			{
				OptName:  existingOpt.OptName,
				OptValue: nil,
			},
			{
				OptName:  "test_option_2",
				OptValue: &newOptValue,
			},
		},
	}

	newPort := &PortWithExtraDHCPOpts{}
	err = ports.Update(context.TODO(), client, port.ID, updateOpts).ExtractInto(newPort)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)
}

func TestPortsRevision(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create Network
	network, err := CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer DeleteNetwork(t, client, network.ID)

	// Create Subnet
	subnet, err := CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer DeleteSubnet(t, client, subnet.ID)

	// Create port
	port, err := CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer DeletePort(t, client, port.ID)

	tools.PrintResource(t, port)

	// Add an address pair to the port
	// Use the RevisionNumber to test the revision / If-Match logic.
	updateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{
			{IPAddress: "192.168.255.10", MACAddress: "aa:bb:cc:dd:ee:ff"},
		},
		RevisionNumber: &port.RevisionNumber,
	}
	newPort, err := ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	// Remove the address pair - this should fail due to old revision number.
	updateOpts = ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{},
		RevisionNumber:      &port.RevisionNumber,
	}
	_, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertErr(t, err)
	if !strings.Contains(err.Error(), "RevisionNumberConstraintFailed") {
		t.Fatalf("expected to see an error of type RevisionNumberConstraintFailed, but got the following error instead: %v", err)
	}

	// The previous ports.Update returns an empty object, so get the port again.
	newPort, err = ports.Get(context.TODO(), client, port.ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, newPort)

	// When not specifying  a RevisionNumber, then the If-Match mechanism
	// should be bypassed.
	updateOpts = ports.UpdateOpts{
		AllowedAddressPairs: &[]ports.AddressPair{},
	}
	newPort, err = ports.Update(context.TODO(), client, port.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newPort)

	if len(newPort.AllowedAddressPairs) > 0 {
		t.Fatalf("Unable to remove the address pair")
	}
}
