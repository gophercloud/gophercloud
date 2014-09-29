// +build acceptance networking layer3ext

package extensions

import (
	"testing"

	base "github.com/rackspace/gophercloud/acceptance/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const (
	cidr1 = "10.0.0.1/24"
	cidr2 = "20.0.0.1/24"
)

func TestRouterCRUD(t *testing.T) {
	base.Setup(t)
	defer base.Teardown()

	// Setup: Create network
	networkID := createNetwork(t)

	// Create router
	routerID := createRouter(t, networkID)

	// Lists routers
	listRouters(t)

	// Update router
	updateRouter(t, routerID)

	// Get router
	getRouter(t, routerID)

	// Create new subnet. Note: this subnet will be deleted when networkID is deleted
	subnetID := createSubnet(t, networkID, cidr2)

	// Add interface
	addInterface(t, routerID, subnetID)

	// Remove interface
	removeInterface(t, routerID, subnetID)

	// Delete router
	deleteRouter(t, routerID)

	// Cleanup
	networks.Delete(base.Client, networkID)
}

func createNetwork(t *testing.T) string {
	t.Logf("Creating a network")

	asu := true
	opts := external.CreateOpts{
		Parent:   networks.CreateOpts{Name: "sample_network", AdminStateUp: &asu},
		External: true,
	}
	n, err := networks.Create(base.Client, opts).Extract()

	th.AssertNoErr(t, err)

	if n.ID == "" {
		t.Fatalf("No ID returned when creating a network")
	}

	createSubnet(t, n.ID, cidr1)

	t.Logf("Network created: ID [%s]", n.ID)

	return n.ID
}

func createSubnet(t *testing.T, networkID, cidr string) string {
	t.Logf("Creating a subnet for network %s", networkID)

	s, err := subnets.Create(base.Client, subnets.CreateOpts{
		NetworkID: networkID,
		CIDR:      cidr,
		IPVersion: subnets.IPv4,
		Name:      "my_subnet",
	}).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Subnet created: ID [%s]", s.ID)

	return s.ID
}

func createRouter(t *testing.T, networkID string) string {
	t.Logf("Creating a router for network %s", networkID)

	asu := false
	gwi := routers.GatewayInfo{NetworkID: networkID}
	r, err := routers.Create(base.Client, routers.CreateOpts{
		Name:         "foo_router",
		AdminStateUp: &asu,
		GatewayInfo:  &gwi,
	}).Extract()

	th.AssertNoErr(t, err)
	th.AssertEquals(t, "foo_router", r.Name)
	th.AssertEquals(t, false, r.AdminStateUp)
	th.AssertDeepEquals(t, gwi, r.GatewayInfo)

	if r.ID == "" {
		t.Fatalf("No ID returned when creating a router")
	}

	t.Logf("Router created: ID [%s]", r.ID)

	return r.ID
}

func listRouters(t *testing.T) {
	pager := routers.List(base.Client, routers.ListOpts{})

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		routerList, err := routers.ExtractRouters(page)
		th.AssertNoErr(t, err)

		for _, r := range routerList {
			t.Logf("Listing router: ID [%s] Name [%s] Status [%s] GatewayInfo [%#v]",
				r.ID, r.Name, r.Status, r.GatewayInfo)
		}

		return true, nil
	})

	th.AssertNoErr(t, err)
}

func updateRouter(t *testing.T, routerID string) {
	r, err := routers.Update(base.Client, routerID, routers.UpdateOpts{
		Name: "another_name",
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "another_name", r.Name)
}

func getRouter(t *testing.T, routerID string) {
	r, err := routers.Get(base.Client, routerID).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Getting router: ID [%s] Name [%s] Status [%s]", r.ID, r.Name, r.Status)
}

func addInterface(t *testing.T, routerID, subnetID string) {
	ir, err := routers.AddInterface(base.Client, routerID, routers.InterfaceOpts{SubnetID: subnetID}).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Interface added to router %s: ID [%s] SubnetID [%s] PortID [%s]", routerID, ir.ID, ir.SubnetID, ir.PortID)
}

func removeInterface(t *testing.T, routerID, subnetID string) {
	ir, err := routers.RemoveInterface(base.Client, routerID, routers.InterfaceOpts{SubnetID: subnetID}).Extract()

	th.AssertNoErr(t, err)

	t.Logf("Interface %s removed from %s", ir.ID, routerID)
}

func deleteRouter(t *testing.T, routerID string) {
	t.Logf("Deleting router %s", routerID)

	res := routers.Delete(base.Client, routerID)

	th.AssertNoErr(t, res.Err)
}
