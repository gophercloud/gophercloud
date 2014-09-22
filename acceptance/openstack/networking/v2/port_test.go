// +build acceptance networking

package v2

import (
	"testing"

	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestPortList(t *testing.T) {
	Setup(t)
	defer Teardown()

	count := 0
	pager := ports.List(Client, ports.ListOpts{})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		count++
		t.Logf("--- Page ---")

		portList, err := ports.ExtractPorts(page)
		th.AssertNoErr(t, err)

		for _, p := range portList {
			t.Logf("Port: ID [%s] Name [%s] Status [%d] MAC addr [%s] Fixed IPs [%#v] Security groups [%#v]",
				p.ID, p.Name, p.Status, p.MACAddress, p.FixedIPs, p.SecurityGroups)
		}

		return true, nil
	})

	th.CheckNoErr(t, err)

	if count == 0 {
		t.Errorf("No pages were iterated over when listing ports")
	}
}

func createNetwork() (string, error) {
	res, err := networks.Create(Client, networks.NetworkOpts{Name: "tmp_network", AdminStateUp: true})
	return res.ID, err
}

func createSubnet(networkID string) (string, error) {
	enable := false
	s, err := subnets.Create(Client, subnets.SubnetOpts{
		NetworkID:  networkID,
		CIDR:       "192.168.199.0/24",
		IPVersion:  subnets.IPv4,
		Name:       "my_subnet",
		EnableDHCP: &enable,
	})
	return s.ID, err
}

func TestPortCRUD(t *testing.T) {
	return
	Setup(t)
	defer Teardown()

	// Setup network
	t.Log("Setting up network")
	networkID, err := createNetwork()
	th.AssertNoErr(t, err)
	defer networks.Delete(Client, networkID)

	// Setup subnet
	t.Logf("Setting up subnet on network %s", networkID)
	subnetID, err := createSubnet(networkID)
	th.AssertNoErr(t, err)
	defer subnets.Delete(Client, subnetID)

	// Create subnet
	t.Logf("Create port based on subnet %s", subnetID)
	enable := false
	opts := ports.PortOpts{
		NetworkID:    networkID,
		Name:         "my_port",
		AdminStateUp: &enable,
		FixedIPs:     []ports.IP{ports.IP{SubnetID: subnetID}},
	}
	p, err := ports.Create(Client, opts)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, p.NetworkID, networkID)
	th.AssertEquals(t, p.Name, "my_port")
	th.AssertEquals(t, p.AdminStateUp, false)
	portID := p.ID

	// Get port
	if portID == "" {
		t.Fatalf("In order to retrieve a port, the portID must be set")
	}
	p, err = ports.Get(Client, portID)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.ID, portID)

	// Update port
	p, err = ports.Update(Client, portID, ports.PortOpts{Name: "new_port_name"})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.Name, "new_port_name")

	// Delete port
	err = ports.Delete(Client, portID)
	th.AssertNoErr(t, err)
}

func TestPortBatchCreate(t *testing.T) {
	// todo
}
