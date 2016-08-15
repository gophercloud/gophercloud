package layer3

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

// CreateFloatingIP creates a floating IP on a given network and port. An error
// will be returned if the creation failed.
func CreateFloatingIP(t *testing.T, client *gophercloud.ServiceClient, networkID, portID string) (*floatingips.FloatingIP, error) {
	t.Logf("Attempting to create floating IP on port: %s", portID)

	createOpts := &floatingips.CreateOpts{
		FloatingNetworkID: networkID,
		PortID:            portID,
	}

	floatingIP, err := floatingips.Create(client, createOpts).Extract()
	if err != nil {
		return floatingIP, err
	}

	t.Logf("Created floating IP.")

	return floatingIP, err
}

// CreateExternalRouter creates a router on the external network. This requires
// the OS_EXTGW_ID environment variable to be set. An error is returned if the
// creation failed.
func CreateExternalRouter(t *testing.T, client *gophercloud.ServiceClient) (*routers.Router, error) {
	var router *routers.Router
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	if err != nil {
		return router, err
	}

	routerName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create router: %s", routerName)

	adminStateUp := true
	gatewayInfo := routers.GatewayInfo{
		NetworkID: choices.ExternalNetworkID,
	}

	createOpts := routers.CreateOpts{
		Name:         routerName,
		AdminStateUp: &adminStateUp,
		GatewayInfo:  &gatewayInfo,
	}

	router, err = routers.Create(client, createOpts).Extract()
	if err != nil {
		return router, err
	}

	t.Logf("Created router: %s", routerName)

	return router, nil
}

// CreateRouter creates a router on a specified Network ID. An error will be
// returned if the creation failed.
func CreateRouter(t *testing.T, client *gophercloud.ServiceClient, networkID string) (*routers.Router, error) {
	routerName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create router: %s", routerName)

	adminStateUp := true
	gatewayInfo := routers.GatewayInfo{
		NetworkID: networkID,
	}

	createOpts := routers.CreateOpts{
		Name:         routerName,
		AdminStateUp: &adminStateUp,
		GatewayInfo:  &gatewayInfo,
	}

	router, err := routers.Create(client, createOpts).Extract()
	if err != nil {
		return router, err
	}

	t.Logf("Created router: %s", routerName)

	return router, nil
}

// DeleteRouter deletes a router of a specified ID. A fatal error will occur
// if the deletion failed. This works best when used as a deferred function.
func DeleteRouter(t *testing.T, client *gophercloud.ServiceClient, routerID string) {
	t.Logf("Attempting to delete router: %s", routerID)

	err := routers.Delete(client, routerID).ExtractErr()
	if err != nil {
		t.Fatalf("Error deleting router: %v", err)
	}

	t.Logf("Deleted router: %s", routerID)
}

// DeleteFloatingIP deletes a floatingIP of a specified ID. A fatal error will
// occur if the deletion failed. This works best when used as a deferred
// function.
func DeleteFloatingIP(t *testing.T, client *gophercloud.ServiceClient, floatingIPID string) {
	t.Logf("Attempting to delete floating IP: %s", floatingIPID)

	err := floatingips.Delete(client, floatingIPID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete floating IP: %v", err)
	}

	t.Logf("Deleted floating IP: %s", floatingIPID)
}

// PrintFloatingIP prints a floating IP and all of its attributes.
func PrintFloatingIP(t *testing.T, fip *floatingips.FloatingIP) {
	t.Logf("ID: %s", fip.ID)
	t.Logf("FloatingNetworkID: %s", fip.FloatingNetworkID)
	t.Logf("FloatingIP: %s", fip.FloatingIP)
	t.Logf("PortID: %s", fip.PortID)
	t.Logf("FixedIP: %s", fip.FixedIP)
	t.Logf("TenantID: %s", fip.TenantID)
	t.Logf("Status: %s", fip.Status)
}

// PrintRouterInterface prints a router interface and all of its attributes.
func PrintRouterInterface(t *testing.T, routerInterface *routers.InterfaceInfo) {
	t.Logf("ID: %s", routerInterface.ID)
	t.Logf("SubnetID: %s", routerInterface.SubnetID)
	t.Logf("PortID: %s", routerInterface.PortID)
	t.Logf("TenantID: %s", routerInterface.TenantID)
}

// PrintRouter prints a router and all of its attributes.
func PrintRouter(t *testing.T, router *routers.Router) {
	t.Logf("ID: %s", router.ID)
	t.Logf("Status: %s", router.Status)
	t.Logf("GatewayInfo: %s", router.GatewayInfo)
	t.Logf("AdminStateUp: %t", router.AdminStateUp)
	t.Logf("Distributed: %t", router.Distributed)
	t.Logf("Name: %s", router.Name)
	t.Logf("TenantID: %s", router.TenantID)
	t.Logf("Routes:")

	for _, route := range router.Routes {
		t.Logf("\tNextHop: %s", route.NextHop)
		t.Logf("\tDestinationCIDR: %s", route.DestinationCIDR)
	}
}
