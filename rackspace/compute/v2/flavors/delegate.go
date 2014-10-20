package flavors

import (
	"github.com/rackspace/gophercloud"
	os "github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
)

// List enumerates the server images available to your account.
func List(client *gophercloud.ServiceClient, opts os.ListOptsBuilder) pagination.Pager {
	return os.List(client, opts)
}

// Get returns details about a single flavor, identity by ID.
func Get(client *gophercloud.ServiceClient, id string) os.GetResult {
	return os.Get(client, id)
}

// ExtractFlavors interprets a page of List results as Flavors.
func ExtractFlavors(page pagination.Page) ([]os.Flavor, error) {
	return os.ExtractFlavors(page)
}
