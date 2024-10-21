//go:build acceptance || networking || mtu

package mtu

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestMTUNetworkCRUDL(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "net-mtu")

	mtuWritable, _ := extensions.Get(context.TODO(), client, "net-mtu-writable").Extract()
	tools.PrintResource(t, mtuWritable)

	// Create Network
	var networkMTU int
	if mtuWritable != nil {
		networkMTU = 1440
	}
	network, err := CreateNetworkWithMTU(t, client, &networkMTU)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	// MTU filtering is supported only in read-only MTU extension
	// https://bugs.launchpad.net/neutron/+bug/1818317
	if mtuWritable == nil {
		// List network successfully
		var listOpts networks.ListOptsBuilder
		listOpts = mtu.ListOptsExt{
			ListOptsBuilder: networks.ListOpts{},
			MTU:             networkMTU,
		}
		var listedNetworks []NetworkMTU
		i := 0
		err = networks.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
			i++
			err := networks.ExtractNetworksInto(page, &listedNetworks)
			if err != nil {
				t.Errorf("Failed to extract networks: %v", err)
				return false, err
			}

			tools.PrintResource(t, listedNetworks)

			th.AssertEquals(t, 1, len(listedNetworks))
			th.CheckDeepEquals(t, *network, listedNetworks[0])

			return true, nil
		})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, i)

		// List network unsuccessfully
		listOpts = mtu.ListOptsExt{
			ListOptsBuilder: networks.ListOpts{},
			MTU:             1,
		}
		i = 0
		err = networks.List(client, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
			i++
			err := networks.ExtractNetworksInto(page, &listedNetworks)
			if err != nil {
				t.Errorf("Failed to extract networks: %v", err)
				return false, err
			}

			tools.PrintResource(t, listedNetworks)

			th.AssertEquals(t, 1, len(listedNetworks))
			th.CheckDeepEquals(t, *network, listedNetworks[0])

			return true, nil
		})
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 0, i)
	}

	// Get network
	var getNetwork NetworkMTU
	err = networks.Get(context.TODO(), client, network.ID).ExtractInto(&getNetwork)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, getNetwork)
	th.AssertDeepEquals(t, *network, getNetwork)

	if mtuWritable != nil {
		// Update network
		newNetworkDescription := ""
		newNetworkMTU := 1350
		networkUpdateOpts := networks.UpdateOpts{
			Description: &newNetworkDescription,
		}
		updateOpts := mtu.UpdateOptsExt{
			UpdateOptsBuilder: networkUpdateOpts,
			MTU:               newNetworkMTU,
		}

		var newNetwork NetworkMTU
		err = networks.Update(context.TODO(), client, network.ID, updateOpts).ExtractInto(&newNetwork)
		th.AssertNoErr(t, err)

		tools.PrintResource(t, newNetwork)
		th.AssertEquals(t, newNetwork.Description, newNetworkDescription)
		th.AssertEquals(t, newNetwork.MTU, newNetworkMTU)

		// Get updated network
		var getNewNetwork NetworkMTU
		err = networks.Get(context.TODO(), client, network.ID).ExtractInto(&getNewNetwork)
		th.AssertNoErr(t, err)

		tools.PrintResource(t, getNewNetwork)
		th.AssertEquals(t, getNewNetwork.Description, newNetworkDescription)
		th.AssertEquals(t, getNewNetwork.MTU, newNetworkMTU)
	}
}
