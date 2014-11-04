package vips

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	lbID  = 12345
	vipID = 67890
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListResponse(t, lbID)

	count := 0

	err := List(client.ServiceClient(), lbID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractVIPs(page)
		th.AssertNoErr(t, err)

		expected := []VIP{
			VIP{ID: 1000, Address: "206.10.10.210", Type: "PUBLIC"},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}
