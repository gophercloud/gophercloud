//go:build acceptance || networking || taas

package taas

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/taas/tapmirrors"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTapMirrorList(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "taas")

	allPages, err := tapmirrors.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allMirrors, err := tapmirrors.ExtractTapMirrors(allPages)
	th.AssertNoErr(t, err)

	for _, mirror := range allMirrors {
		tools.PrintResource(t, mirror)
	}
}

func TestTapMirrorCRUD(t *testing.T) {
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Skip these tests if we don't have the required extension
	networking.RequireNeutronExtension(t, client, "taas")

	// Create Port
	network, err := networking.CreateNetwork(t, client)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, client, network.ID)

	subnet, err := networking.CreateSubnet(t, client, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, client, subnet.ID)

	port, err := networking.CreatePort(t, client, network.ID, subnet.ID)
	th.AssertNoErr(t, err)
	defer networking.DeletePort(t, client, port.ID)

	// Create and defer Delete Tap Mirror
	mirror, err := CreateTapMirror(t, client, port.ID, port.FixedIPs[0].IPAddress)
	th.AssertNoErr(t, err)
	defer DeleteTapMirror(t, client, mirror.ID)

	tools.PrintResource(t, mirror)

	// Get Tap Mirror
	newmirror, err := tapmirrors.Get(context.TODO(), client, mirror.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, mirror, newmirror)

	// Update Tap Mirror
	updatedName := "TESTACC-updated name"
	updatedDescription := "TESTACC-updated mirror description"
	updateOpts := tapmirrors.UpdateOpts{
		Name:        &updatedName,
		Description: &updatedDescription,
	}
	updatedmirror, err := tapmirrors.Update(context.TODO(), client, mirror.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, updatedName, updatedmirror.Name)
	th.AssertEquals(t, updatedDescription, updatedmirror.Description)
}
