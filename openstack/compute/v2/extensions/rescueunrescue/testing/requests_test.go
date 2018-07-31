package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/rescueunrescue"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestRescue(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleServerRescueSuccessfully(t)

	res := rescueunrescue.Rescue(fake.ServiceClient(), "1234asdf", rescueunrescue.RescueOpts{
		AdminPass: "1234567890",
	})
	th.AssertNoErr(t, res.Err)
	adminPass, _ := res.Extract()
	th.AssertEquals(t, "1234567890", adminPass)
}
