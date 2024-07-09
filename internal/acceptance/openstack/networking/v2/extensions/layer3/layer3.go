package layer3

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

// CreateFloatingIP creates a floating IP on a given network and port. An error
// will be returned if the creation failed.
func CreateFloatingIP(t *testing.T, client *gophercloud.ServiceClient, networkID, portID string) (*floatingips.FloatingIP, error) {
	t.Logf("Attempting to create floating IP on port: %s", portID)

	fipDescription := "Test floating IP"
	createOpts := &floatingips.CreateOpts{
		Description:       fipDescription,
		FloatingNetworkID: networkID,
		PortID:            portID,
	}

	floatingIP, err := floatingips.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return floatingIP, err
	}

	t.Logf("Created floating IP.")

	th.AssertEquals(t, floatingIP.Description, fipDescription)

	return floatingIP, err
}

// CreateFloatingIPWithFixedIP creates a floating IP on a given network and port with a
// defined fixed IP. An error will be returned if the creation failed.
func CreateFloatingIPWithFixedIP(t *testing.T, client *gophercloud.ServiceClient, networkID, portID, fixedIP string) (*floatingips.FloatingIP, error) {
	t.Logf("Attempting to create floating IP on port: %s and address: %s", portID, fixedIP)

	fipDescription := "Test floating IP"
	createOpts := &floatingips.CreateOpts{
		Description:       fipDescription,
		FloatingNetworkID: networkID,
		PortID:            portID,
		FixedIP:           fixedIP,
	}

	floatingIP, err := floatingips.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return floatingIP, err
	}

	t.Logf("Created floating IP.")

	th.AssertEquals(t, floatingIP.Description, fipDescription)
	th.AssertEquals(t, floatingIP.FixedIP, fixedIP)

	return floatingIP, err
}

// CreatePortForwarding creates a port forwarding for a given floating IP
// and port. An error will be returned if the creation failed.
func CreatePortForwarding(t *testing.T, client *gophercloud.ServiceClient, fipID string, portID string, portFixedIPs []ports.IP) (*portforwarding.PortForwarding, error) {
	t.Logf("Attempting to create Port forwarding for floating IP with ID: %s", fipID)

	fixedIP := portFixedIPs[0]
	internalIP := fixedIP.IPAddress
	pfDescription := "Test description"
	createOpts := &portforwarding.CreateOpts{
		Description:       pfDescription,
		Protocol:          "tcp",
		InternalPort:      25,
		ExternalPort:      2230,
		InternalIPAddress: internalIP,
		InternalPortID:    portID,
	}

	pf, err := portforwarding.Create(context.TODO(), client, fipID, createOpts).Extract()
	if err != nil {
		return pf, err
	}

	t.Logf("Created Port Forwarding.")

	th.AssertEquals(t, pf.Protocol, "tcp")

	return pf, err
}

// DeletePortForwarding deletes a Port Forwarding with a given ID and a given floating IP ID.
// A fatal error is returned if the deletion fails. Works best as a deferred function
func DeletePortForwarding(t *testing.T, client *gophercloud.ServiceClient, fipID string, pfID string) {
	t.Logf("Attempting to delete the port forwarding with ID %s for floating IP with ID %s", pfID, fipID)

	err := portforwarding.Delete(context.TODO(), client, fipID, pfID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete Port forwarding with ID %s for floating IP with ID %s", pfID, fipID)
	}
	t.Logf("Successfully deleted the port forwarding with ID %s for floating IP with ID %s", pfID, fipID)

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
	routerDescription := tools.RandomString("TESTACC-DESC-", 8)

	t.Logf("Attempting to create external router: %s", routerName)

	adminStateUp := true
	gatewayInfo := routers.GatewayInfo{
		NetworkID: choices.ExternalNetworkID,
	}

	createOpts := routers.CreateOpts{
		Name:         routerName,
		Description:  routerDescription,
		AdminStateUp: &adminStateUp,
		GatewayInfo:  &gatewayInfo,
	}

	router, err = routers.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return router, err
	}

	if err := WaitForRouterToCreate(client, router.ID); err != nil {
		return router, err
	}

	t.Logf("Created router: %s", routerName)

	th.AssertEquals(t, router.Name, routerName)
	th.AssertEquals(t, router.Description, routerDescription)

	return router, nil
}

// CreateRouter creates a router on a specified Network ID. An error will be
// returned if the creation failed.
func CreateRouter(t *testing.T, client *gophercloud.ServiceClient, networkID string) (*routers.Router, error) {
	routerName := tools.RandomString("TESTACC-", 8)
	routerDescription := tools.RandomString("TESTACC-DESC-", 8)

	t.Logf("Attempting to create router: %s", routerName)

	adminStateUp := true
	createOpts := routers.CreateOpts{
		Name:         routerName,
		Description:  routerDescription,
		AdminStateUp: &adminStateUp,
	}

	router, err := routers.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return router, err
	}

	if err := WaitForRouterToCreate(client, router.ID); err != nil {
		return router, err
	}

	t.Logf("Created router: %s", routerName)

	th.AssertEquals(t, router.Name, routerName)
	th.AssertEquals(t, router.Description, routerDescription)

	return router, nil
}

// CreateRouterInterface will attach a subnet to a router. An error will be
// returned if the operation fails.
func CreateRouterInterface(t *testing.T, client *gophercloud.ServiceClient, portID, routerID string) (*routers.InterfaceInfo, error) {
	t.Logf("Attempting to add port %s to router %s", portID, routerID)

	aiOpts := routers.AddInterfaceOpts{
		PortID: portID,
	}

	iface, err := routers.AddInterface(context.TODO(), client, routerID, aiOpts).Extract()
	if err != nil {
		return iface, err
	}

	if err := WaitForRouterInterfaceToAttach(client, portID); err != nil {
		return iface, err
	}

	t.Logf("Successfully added port %s to router %s", portID, routerID)
	return iface, nil
}

// CreateRouterInterfaceOnSubnet will attach a subnet to a router. An error will be
// returned if the operation fails.
func CreateRouterInterfaceOnSubnet(t *testing.T, client *gophercloud.ServiceClient, subnetID, routerID string) (*routers.InterfaceInfo, error) {
	t.Logf("Attempting to add subnet %s to router %s", subnetID, routerID)

	aiOpts := routers.AddInterfaceOpts{
		SubnetID: subnetID,
	}

	iface, err := routers.AddInterface(context.TODO(), client, routerID, aiOpts).Extract()
	if err != nil {
		return iface, err
	}

	if err := WaitForRouterInterfaceToAttach(client, iface.PortID); err != nil {
		return iface, err
	}

	t.Logf("Successfully added subnet %s to router %s", subnetID, routerID)
	return iface, nil
}

// DeleteRouter deletes a router of a specified ID. A fatal error will occur
// if the deletion failed. This works best when used as a deferred function.
func DeleteRouter(t *testing.T, client *gophercloud.ServiceClient, routerID string) {
	t.Logf("Attempting to delete router: %s", routerID)

	err := routers.Delete(context.TODO(), client, routerID).ExtractErr()
	if err != nil {
		t.Fatalf("Error deleting router: %v", err)
	}

	if err := WaitForRouterToDelete(client, routerID); err != nil {
		t.Fatalf("Error waiting for router to delete: %v", err)
	}

	t.Logf("Deleted router: %s", routerID)
}

// DeleteRouterInterface will detach a subnet to a router. A fatal error will
// occur if the deletion failed. This works best when used as a deferred
// function.
func DeleteRouterInterface(t *testing.T, client *gophercloud.ServiceClient, portID, routerID string) {
	t.Logf("Attempting to detach port %s from router %s", portID, routerID)

	riOpts := routers.RemoveInterfaceOpts{
		PortID: portID,
	}

	_, err := routers.RemoveInterface(context.TODO(), client, routerID, riOpts).Extract()
	if err != nil {
		t.Fatalf("Failed to detach port %s from router %s", portID, routerID)
	}

	if err := WaitForRouterInterfaceToDetach(client, portID); err != nil {
		t.Fatalf("Failed to wait for port %s to detach from router %s", portID, routerID)
	}

	t.Logf("Successfully detached port %s from router %s", portID, routerID)
}

// DeleteFloatingIP deletes a floatingIP of a specified ID. A fatal error will
// occur if the deletion failed. This works best when used as a deferred
// function.
func DeleteFloatingIP(t *testing.T, client *gophercloud.ServiceClient, floatingIPID string) {
	t.Logf("Attempting to delete floating IP: %s", floatingIPID)

	err := floatingips.Delete(context.TODO(), client, floatingIPID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete floating IP: %v", err)
	}

	t.Logf("Deleted floating IP: %s", floatingIPID)
}

func WaitForRouterToCreate(client *gophercloud.ServiceClient, routerID string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		r, err := routers.Get(ctx, client, routerID).Extract()
		if err != nil {
			return false, err
		}

		if r.Status == "ACTIVE" {
			return true, nil
		}

		return false, nil
	})
}

func WaitForRouterToDelete(client *gophercloud.ServiceClient, routerID string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		_, err := routers.Get(ctx, client, routerID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return true, nil
			}

			return false, err
		}

		return false, nil
	})
}

func WaitForRouterInterfaceToAttach(client *gophercloud.ServiceClient, routerInterfaceID string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		r, err := ports.Get(ctx, client, routerInterfaceID).Extract()
		if err != nil {
			return false, err
		}

		if r.Status == "ACTIVE" {
			return true, nil
		}

		return false, nil
	})
}

func WaitForRouterInterfaceToDetach(client *gophercloud.ServiceClient, routerInterfaceID string) error {
	return tools.WaitFor(func(ctx context.Context) (bool, error) {
		r, err := ports.Get(ctx, client, routerInterfaceID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return true, nil
			}

			if errCode, ok := err.(gophercloud.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return false, nil
				}
			}

			return false, err
		}

		if r.Status == "ACTIVE" {
			return true, nil
		}

		return false, nil
	})
}

// CreateAddressScope will create an address-scope. An error will be returned if
// the address-scope could not be created.
func CreateAddressScope(t *testing.T, client *gophercloud.ServiceClient) (*addressscopes.AddressScope, error) {
	addressScopeName := tools.RandomString("TESTACC-", 8)
	createOpts := addressscopes.CreateOpts{
		Name:      addressScopeName,
		IPVersion: 4,
	}

	t.Logf("Attempting to create an address-scope: %s", addressScopeName)

	addressScope, err := addressscopes.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Successfully created the addressscopes.")

	th.AssertEquals(t, addressScope.Name, addressScopeName)
	th.AssertEquals(t, addressScope.IPVersion, int(gophercloud.IPv4))

	return addressScope, nil
}

// DeleteAddressScope will delete an address-scope with the specified ID.
// A fatal error will occur if the delete was not successful.
func DeleteAddressScope(t *testing.T, client *gophercloud.ServiceClient, addressScopeID string) {
	t.Logf("Attempting to delete the address-scope: %s", addressScopeID)

	err := addressscopes.Delete(context.TODO(), client, addressScopeID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete address-scope %s: %v", addressScopeID, err)
	}

	t.Logf("Deleted address-scope: %s", addressScopeID)
}
