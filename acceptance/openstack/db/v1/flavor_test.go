// +build acceptance db

package v1

import (
	"github.com/rackspace/gophercloud/openstack/db/v1/flavors"
	"github.com/rackspace/gophercloud/pagination"
)

func (c context) listFlavors() {
	c.Logf("Listing flavors")

	err := flavors.List(c.client, c.instanceID).EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := db.ExtractFlavors(page)
		c.AssertNoErr(err)

		for _, f := range flavorList {
			c.Logf("Flavor: %#v", f)
		}

		return true, nil
	})

	c.CheckNoErr(err)
}

func (c context) getFlavor() {
	flavor, err := flavors.Get(c.client, "1").Extract()
	c.Logf("Getting flavor %s", flavor.ID)
	c.CheckNoErr(err)
}
