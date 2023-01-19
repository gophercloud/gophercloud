//go:build acceptance || compute || servers
// +build acceptance compute servers

package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestTenantNetworksList(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	allPages, err := tenantnetworks.List(client).AllPages()
	th.AssertNoErr(t, err)

	allTenantNetworks, err := tenantnetworks.ExtractNetworks(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, network := range allTenantNetworks {
		tools.PrintResource(t, network)

		if network.Name == choices.NetworkName {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestTenantNetworksGet(t *testing.T) {
	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	networkID, err := GetNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	th.AssertNoErr(t, err)

	network, err := tenantnetworks.Get(client, networkID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, network)
}
