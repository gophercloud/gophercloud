package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/availabilityzones"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// Verifies that availability zones can be listed correctly
func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	a, err := availabilityzones.List(client.ServiceClient()).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, a[0].ID, "388c983d-258e-4a0e-b1ba-10da37d766db")
	th.AssertEquals(t, a[0].Name, "nova")
}
