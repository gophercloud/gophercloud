// +build acceptance

package v1

import (
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/utils"
	th "github.com/rackspace/gophercloud/testhelper"
)

var metadata = map[string]string{"gopher": "cloud"}

func newClient() (*gophercloud.ServiceClient, error) {
	ao, err := utils.AuthOptions()
	th.AssertNoErr(t, err)

	client, err := openstack.AuthenticatedClient(ao)
	th.AssertNoErr(t, err)

	return openstack.NewObjectStorageV1(client, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	}), nil
}
