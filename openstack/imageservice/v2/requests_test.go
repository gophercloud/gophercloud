package v2

// TODO
// compare with openstack/compute/v2/servers/requests_test.go

import (
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	fakeclient "github.com/rackspace/gophercloud/testhelper/client"
)

func TestCreateImage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleImageCreationSuccessfully(t)

	actualImage, err := Create(fakeclient.ServiceClient(), CreateOpts{
		Name: "Ubuntu 12.10",
		Id: "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	expectedImage := Image{
		Name: "Ubuntu 12.10",
		Id: "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Tags: []string{"ubuntu", "quantal"},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}
