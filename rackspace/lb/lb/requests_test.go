package lb

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListLBResponse(t)

	count := 0

	err := List(client.ServiceClient(), ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractLBs(page)
		th.AssertNoErr(t, err)

		expected := []LoadBalancer{
			LoadBalancer{
				Name:      "lb-site1",
				ID:        71,
				Protocol:  "HTTP",
				Port:      80,
				Algorithm: RAND,
				Status:    ACTIVE,
				NodeCount: 3,
				VIPs: []VIP{
					VIP{
						ID:      403,
						Address: "206.55.130.1",
						Type:    "PUBLIC",
						Version: "IPV4",
					},
				},
				Created: Datetime{Time: "2010-11-30T03:23:42Z"},
				Updated: Datetime{Time: "2010-11-30T03:23:44Z"},
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}
