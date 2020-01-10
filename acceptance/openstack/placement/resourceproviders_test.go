package placement

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/placement/resourceproviders"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestResourceProviderList(t *testing.T) {
	client, err := clients.NewPlacementClient()
	th.AssertNoErr(t, err)

	allPages, err := resourceproviders.List(client, resourceproviders.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	allResourceProviders, err := resourceproviders.ExtractResourceProviders(allPages)
	th.AssertNoErr(t, err)

	for _, v := range allResourceProviders {
		tools.PrintResource(t, v)
	}
}
