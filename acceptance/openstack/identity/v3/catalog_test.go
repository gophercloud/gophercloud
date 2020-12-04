package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/catalog"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCatalogList(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := catalog.List(client).AllPages()
	th.AssertNoErr(t, err)

	allEntities, err := catalog.ExtractServiceCatalog(allPages)
	th.AssertNoErr(t, err)

	for _, entity := range allEntities {
		tools.PrintResource(t, entity)
	}
}
