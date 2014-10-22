// +build acceptance rackspace networking v2

package v2

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/networking/v2/networks"
	th "github.com/rackspace/gophercloud/testhelper"
)

func TestNetworks(t *testing.T) {
	Setup(t)
	defer Teardown()

	// Create a network
	n, err := networks.Create(Client, networks.CreateOpts{Label: "sample_network", CIDR: "172.20.0.0/24"}).Extract()
	th.AssertNoErr(t, err)
	defer networks.Delete(Client, n.ID)
	th.AssertEquals(t, n.Label, "sample_network")
	th.AssertEquals(t, n.CIDR, "172.20.0.0/24")
	networkID := n.ID

	// List networks
	pager := networks.List(Client)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")

		networkList, err := networks.ExtractNetworks(page)
		th.AssertNoErr(t, err)

		for _, n := range networkList {
			t.Logf("Network: ID [%s] Label [%s] CIDR [%s]",
				n.ID, n.Label, n.CIDR)
		}

		return true, nil
	})
	th.CheckNoErr(t, err)

	// Get a network
	if networkID == "" {
		t.Fatalf("In order to retrieve a network, the NetworkID must be set")
	}
	n, err = networks.Get(Client, networkID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, n.CIDR, "172.20.0.0/24")
	th.AssertEquals(t, n.Label, "sample_network")
	th.AssertEquals(t, n.ID, networkID)
}
