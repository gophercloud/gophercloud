package nodes

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	lbID    = 12345
	nodeID  = 67890
	nodeID2 = 67891
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

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockCreateResponse(t, lbID)

	opts := CreateOpts{
		CreateOpt{
			Address:   "10.2.2.3",
			Port:      80,
			Condition: ENABLED,
			Type:      PRIMARY,
		},
		CreateOpt{
			Address:   "10.2.2.4",
			Port:      81,
			Condition: ENABLED,
			Type:      SECONDARY,
		},
	}

	page := Create(client.ServiceClient(), lbID, opts)

	actual, err := page.ExtractNodes()
	th.AssertNoErr(t, err)

	expected := []Node{
		Node{
			ID:        185,
			Address:   "10.2.2.3",
			Port:      80,
			Condition: ENABLED,
			Status:    ONLINE,
			Weight:    1,
			Type:      PRIMARY,
		},
		Node{
			ID:        186,
			Address:   "10.2.2.4",
			Port:      81,
			Condition: ENABLED,
			Status:    ONLINE,
			Weight:    1,
			Type:      SECONDARY,
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}

func TestBulkDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	ids := []int{nodeID, nodeID2}

	mockBatchDeleteResponse(t, lbID, ids)

	err := BulkDelete(client.ServiceClient(), lbID, ids).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockGetResponse(t, lbID, nodeID)

	node, err := Get(client.ServiceClient(), lbID, nodeID).Extract()
	th.AssertNoErr(t, err)

	expected := &Node{
		ID:        410,
		Address:   "10.1.1.1",
		Port:      80,
		Condition: ENABLED,
		Status:    ONLINE,
		Weight:    12,
		Type:      PRIMARY,
	}

	th.AssertDeepEquals(t, expected, node)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockUpdateResponse(t, lbID, nodeID)

	opts := UpdateOpts{
		Address:   "1.2.3.4",
		Weight:    IntToPointer(10),
		Condition: DRAINING,
		Type:      SECONDARY,
	}

	err := Update(client.ServiceClient(), lbID, nodeID, opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockDeleteResponse(t, lbID, nodeID)

	err := Delete(client.ServiceClient(), lbID, nodeID).ExtractErr()
	th.AssertNoErr(t, err)
}
