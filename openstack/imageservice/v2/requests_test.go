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

	id := "e7db3b45-8db7-47ad-8109-3fb55c2c24fd"
	name := "Ubuntu 12.10"

	actualImage, err := Create(fakeclient.ServiceClient(), CreateOpts{
		Id: &id,
		Name: &name,
		Tags: []string{"ubuntu", "quantal"},
	}).Extract()

	th.AssertNoErr(t, err)

	expectedImage := Image{
		Id: "e7db3b45-8db7-47ad-8109-3fb55c2c24fd",
		Name: "Ubuntu 12.10",
		Status: ImageStatusQueued,
		Tags: []string{"ubuntu", "quantal"},
	}

	th.AssertDeepEquals(t, &expectedImage, actualImage)
}
