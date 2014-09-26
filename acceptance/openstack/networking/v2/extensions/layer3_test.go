// +build acceptance networking layer3ext

package extensions

import (
	"testing"

	base "github.com/rackspace/gophercloud/acceptance/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	th "github.com/rackspace/gophercloud/testhelper"
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

	// Add interface
	addInterface(t, routerID)

	// Remove interface
	removeInterface(t, routerID)

	// Delete router
	deleteRouter(t, routerID)
}

func createNetwork(t *testing.T) string {
	t.Logf("Creating a network")

	opts := networks.CreateOpts{
		Name:         "sample_network",
		AdminStateUp: true,
	}
	n, err := networks.Create(base.Client, opts).Extract()

	th.AssertNoErr(t, err)

	if n.ID == "" {
		t.Fatalf("No ID returned when creating a network")
	}

	return n.ID
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

	return r.ID
}

func listRouters(t *testing.T) {

}

func updateRouter(t *testing.T, routerID string) {
}

func getRouter(t *testing.T, routerID string) {
}

func addInterface(t *testing.T, routerID string) {
}

func removeInterface(t *testing.T, routerID string) {
}

func deleteRouter(t *testing.T, routerID string) {
}
