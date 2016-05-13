// +build acceptance networking portsbinding

package portsbinding

import (
	"testing"

	base "github.com/rackspace/gophercloud/acceptance/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/ports"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestPortBinding(t *testing.T) {
	base.Setup(t)
	defer base.Teardown()

	// Setup network
	t.Log("Setting up network")
	networkID, err := createNetwork()
	th.AssertNoErr(t, err)
	defer networks.Delete(base.Client, networkID)

	// Setup subnet
	t.Logf("Setting up subnet on network %s", networkID)
	subnetID, err := createSubnet(networkID)
	th.AssertNoErr(t, err)
	defer subnets.Delete(base.Client, subnetID)

	// Create port
	t.Logf("Create port based on subnet %s", subnetID)
	hostID := "localhost"
	portID := createPort(t, networkID, subnetID, hostID)

	// Get port
	if portID == "" {
		t.Fatalf("In order to retrieve a port, the portID must be set")
	}
	p, err := portsbinding.Get(base.Client, portID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.ID, portID)
	th.AssertEquals(t, p.HostID, hostID)

	// Update port
	newHostID := "openstack"
	updateOpts := portsbinding.UpdateOpts{
		HostID: newHostID,
	}
	p, err = portsbinding.Update(base.Client, portID, updateOpts).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.HostID, newHostID)

	// Delete port
	res := ports.Delete(base.Client, portID)
	th.AssertNoErr(t, res.Err)
}

func createPort(t *testing.T, networkID, subnetID, hostID string) string {
	enable := false
	opts := portsbinding.CreateOpts{
		CreateOptsBuilder: ports.CreateOpts{
			NetworkID:    networkID,
			Name:         "my_port",
			AdminStateUp: &enable,
			FixedIPs:     []ports.IP{{SubnetID: subnetID}},
		},
		HostID: hostID,
	}

	p, err := portsbinding.Create(base.Client, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, p.NetworkID, networkID)
	th.AssertEquals(t, p.Name, "my_port")
	th.AssertEquals(t, p.AdminStateUp, false)

	return p.ID
}

func createNetwork() (string, error) {
	res, err := networks.Create(base.Client, networks.CreateOpts{Name: "tmp_network", AdminStateUp: networks.Up}).Extract()
	return res.ID, err
}

func createSubnet(networkID string) (string, error) {
	s, err := subnets.Create(base.Client, subnets.CreateOpts{
		NetworkID:  networkID,
		CIDR:       "192.168.199.0/24",
		IPVersion:  subnets.IPv4,
		Name:       "my_subnet",
		EnableDHCP: subnets.Down,
		AllocationPools: []subnets.AllocationPool{
			{Start: "192.168.199.2", End: "192.168.199.200"},
		},
	}).Extract()
	return s.ID, err
}
