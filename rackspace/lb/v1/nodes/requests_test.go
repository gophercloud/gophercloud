package nodes

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	lbID   = 12345
	nodeID = 67890
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListResponse(t, lbID)

	count := 0

	err := List(client.ServiceClient(), lbID, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractNodes(page)
		th.AssertNoErr(t, err)

		expected := []Node{
			Node{
				ID:        410,
				Address:   "10.1.1.1",
				Port:      80,
				Condition: ENABLED,
				Status:    ONLINE,
				Weight:    3,
				Type:      PRIMARY,
			},
			Node{
				ID:        411,
				Address:   "10.1.1.2",
				Port:      80,
				Condition: ENABLED,
				Status:    ONLINE,
				Weight:    8,
				Type:      SECONDARY,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}
