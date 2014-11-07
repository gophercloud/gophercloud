package lbs

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/nodes"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/sessions"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/throttle"
	"github.com/rackspace/gophercloud/rackspace/lb/v1/vips"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	id1 = 12345
	id2 = 67890
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
				Algorithm: "RANDOM",
				Status:    ACTIVE,
				NodeCount: 3,
				VIPs: []vips.VIP{
					vips.VIP{
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
		VIPs: []vips.VIP{
			vips.VIP{ID: 2341},
			vips.VIP{ID: 900001},
		},
		Nodes: []nodes.Node{
			nodes.Node{Address: "10.1.1.1", Port: 80, Condition: "ENABLED"},
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
		Algorithm:  "RANDOM",
		Status:     BUILD,
		Timeout:    30,
		Cluster:    Cluster{Name: "ztm-n01.staging1.lbaas.rackspace.net"},
		Nodes: []nodes.Node{
			nodes.Node{
				Address:   "10.1.1.1",
				ID:        653,
				Port:      80,
				Status:    "ONLINE",
				Condition: "ENABLED",
				Weight:    1,
			},
		},
		VIPs: []vips.VIP{
			vips.VIP{
				ID:      39,
				Address: "206.10.10.210",
				Type:    vips.PUBLIC,
				Version: vips.IPV4,
			},
			vips.VIP{
				ID:      900001,
				Address: "2001:4801:79f1:0002:711b:be4c:0000:0021",
				Type:    vips.PUBLIC,
				Version: vips.IPV6,
			},
		},
		Created:           Datetime{Time: "2011-04-13T14:18:07Z"},
		Updated:           Datetime{Time: "2011-04-13T14:18:07Z"},
		ConnectionLogging: ConnectionLogging{Enabled: false},
	}

	th.AssertDeepEquals(t, expected, lb)
}

func TestBulkDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	ids := []int{id1, id2}

	mockBatchDeleteLBResponse(t, ids)

	err := BulkDelete(client.ServiceClient(), ids).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockDeleteLBResponse(t, id1)

	err := Delete(client.ServiceClient(), id1).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockGetLBResponse(t, id1)

	lb, err := Get(client.ServiceClient(), id1).Extract()

	expected := &LoadBalancer{
		Name:              "sample-loadbalancer",
		ID:                2000,
		Protocol:          "HTTP",
		Port:              80,
		Algorithm:         "RANDOM",
		Status:            ACTIVE,
		Timeout:           30,
		ConnectionLogging: ConnectionLogging{Enabled: true},
		VIPs: []vips.VIP{
			vips.VIP{
				ID:      1000,
				Address: "206.10.10.210",
				Type:    "PUBLIC",
				Version: "IPV4",
			},
		},
		Nodes: []nodes.Node{
			nodes.Node{
				Address:   "10.1.1.1",
				ID:        1041,
				Port:      80,
				Status:    "ONLINE",
				Condition: "ENABLED",
			},
			nodes.Node{
				Address:   "10.1.1.2",
				ID:        1411,
				Port:      80,
				Status:    "ONLINE",
				Condition: "ENABLED",
			},
		},
		SessionPersistence: sessions.SessionPersistence{Type: "HTTP_COOKIE"},
		ConnectionThrottle: throttle.ConnectionThrottle{MaxConnections: 100},
		Cluster:            Cluster{Name: "c1.dfw1"},
		Created:            Datetime{Time: "2010-11-30T03:23:42Z"},
		Updated:            Datetime{Time: "2010-11-30T03:23:44Z"},
		SourceAddrs: SourceAddrs{
			IPv4Public:  "10.12.99.28",
			IPv4Private: "10.0.0.0",
			IPv6Public:  "2001:4801:79f1:1::1/64",
		},
	}

	th.AssertDeepEquals(t, expected, lb)
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockUpdateLBResponse(t, id1)

	opts := UpdateOpts{
		Name:          "a-new-loadbalancer",
		Protocol:      "TCP",
		HalfClosed:    gophercloud.Enabled,
		Algorithm:     "RANDOM",
		Port:          8080,
		Timeout:       100,
		HTTPSRedirect: gophercloud.Disabled,
	}

	err := Update(client.ServiceClient(), id1, opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListProtocols(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListProtocolsResponse(t)

	count := 0

	err := ListProtocols(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractProtocols(page)
		th.AssertNoErr(t, err)

		expected := []Protocol{
			Protocol{Name: "DNS_TCP", Port: 53},
			Protocol{Name: "DNS_UDP", Port: 53},
			Protocol{Name: "FTP", Port: 21},
			Protocol{Name: "HTTP", Port: 80},
			Protocol{Name: "HTTPS", Port: 443},
			Protocol{Name: "IMAPS", Port: 993},
			Protocol{Name: "IMAPv4", Port: 143},
		}

		th.CheckDeepEquals(t, expected[0:7], actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestListAlgorithms(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListAlgorithmsResponse(t)

	count := 0

	err := ListAlgorithms(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractAlgorithms(page)
		th.AssertNoErr(t, err)

		expected := []Algorithm{
			Algorithm{Name: "LEAST_CONNECTIONS"},
			Algorithm{Name: "RANDOM"},
			Algorithm{Name: "ROUND_ROBIN"},
			Algorithm{Name: "WEIGHTED_LEAST_CONNECTIONS"},
			Algorithm{Name: "WEIGHTED_ROUND_ROBIN"},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestIsLoggingEnabled(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockGetLoggingResponse(t, id1)

	res, err := IsLoggingEnabled(client.ServiceClient(), id1)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, res)
}
