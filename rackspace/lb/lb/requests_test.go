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

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockCreateLBResponse(t)

	opts := CreateOpts{
		Name:     "a-new-loadbalancer",
		Port:     80,
		Protocol: "HTTP",
		VIPs: []VIP{
			VIP{ID: 2341},
			VIP{ID: 900001},
		},
		Nodes: []Node{
			Node{Address: "10.1.1.1", Port: 80, Condition: "ENABLED"},
		},
	}

	lb, err := Create(client.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &LoadBalancer{
		Name:       "a-new-loadbalancer",
		ID:         144,
		Protocol:   "HTTP",
		HalfClosed: false,
		Port:       83,
		Algorithm:  RAND,
		Status:     BUILD,
		Timeout:    30,
		Cluster:    Cluster{Name: "ztm-n01.staging1.lbaas.rackspace.net"},
		Nodes: []Node{
			Node{
				Address:   "10.1.1.1",
				ID:        653,
				Port:      80,
				Status:    "ONLINE",
				Condition: "ENABLED",
				Weight:    1,
			},
		},
		VIPs: []VIP{
			VIP{
				ID:      39,
				Address: "206.10.10.210",
				Type:    "PUBLIC",
				Version: "IPV4",
			},
			VIP{
				ID:      900001,
				Address: "2001:4801:79f1:0002:711b:be4c:0000:0021",
				Type:    "PUBLIC",
				Version: "IPV6",
			},
		},
		Created:           Datetime{Time: "2011-04-13T14:18:07Z"},
		Updated:           Datetime{Time: "2011-04-13T14:18:07Z"},
		ConnectionLogging: ConnectionLogging{Enabled: false},
	}

	th.AssertDeepEquals(t, expected, lb)
}

func TestBulkDeleteLBs(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	ids := []int{12345, 67890}

	mockBatchDeleteLBResponse(t, ids)

	err := BulkDelete(client.ServiceClient(), ids).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDeleteLB(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	id := 12345

	mockDeleteLBResponse(t, id)

	err := Delete(client.ServiceClient(), id).ExtractErr()
	th.AssertNoErr(t, err)
}
