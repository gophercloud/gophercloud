// +build acceptance

package v3

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/rackspace"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func newClient(t *testing.T) *gophercloud.ServiceClient {
	ao, err := rackspace.AuthOptionsFromEnv()
	th.AssertNoErr(t, err)

	client, err := rackspace.AuthenticatedClient(ao)
	th.AssertNoErr(t, err)

	c, err := rackspace.NewRackConnectV3(client, gophercloud.EndpointOpts{
		Region: os.Getenv("RS_REGION_NAME"),
	})
	th.AssertNoErr(t, err)
	return c
}
