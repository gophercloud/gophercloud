package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/loadbalancers"
	fake "github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListLoadbalancers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerListSuccessfully(t)

	pages := 0
	err := loadbalancers.List(fake.ServiceClient(), loadbalancers.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := loadbalancers.ExtractLoadBalancers(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 loadbalancers, got %d", len(actual))
		}
		th.CheckDeepEquals(t, LoadbalancerWeb, actual[0])
		th.CheckDeepEquals(t, LoadbalancerDb, actual[1])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestListAllLoadbalancers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerListSuccessfully(t)

	allPages, err := loadbalancers.List(fake.ServiceClient(), loadbalancers.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := loadbalancers.ExtractLoadBalancers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, LoadbalancerWeb, actual[0])
	th.CheckDeepEquals(t, LoadbalancerDb, actual[1])
}

func TestCreateLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerCreationSuccessfully(t, SingleLoadbalancerBody)

	actual, err := loadbalancers.Create(context.TODO(), fake.ServiceClient(), loadbalancers.CreateOpts{
		Name:         "db_lb",
		AdminStateUp: gophercloud.Enabled,
		VipPortID:    "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		VipSubnetID:  "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:   "10.30.176.48",
		FlavorID:     "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		Provider:     "haproxy",
		Tags:         []string{"test", "stage"},
		AdditionalVips: []loadbalancers.AdditionalVip{
			{
				SubnetID:  "0d4f6a08-60b7-44ab-8903-f7d76ec54095",
				IPAddress: "192.168.10.10",
			},
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, LoadbalancerDb, *actual)
}

func TestCreateFullyPopulatedLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFullyPopulatedLoadbalancerCreationSuccessfully(t, PostFullyPopulatedLoadbalancerBody)

	actual, err := loadbalancers.Create(context.TODO(), fake.ServiceClient(), loadbalancers.CreateOpts{
		Name:         "db_lb",
		AdminStateUp: gophercloud.Enabled,
		VipPortID:    "2bf413c8-41a9-4477-b505-333d5cbe8b55",
		VipSubnetID:  "9cedb85d-0759-4898-8a4b-fa5a5ea10086",
		VipAddress:   "10.30.176.48",
		FlavorID:     "bba40eb2-ee8c-11e9-81b4-2a2ae2dbcce4",
		Provider:     "octavia",
		Tags:         []string{"test", "stage"},
		Listeners: []listeners.CreateOpts{{
			Protocol:     "HTTP",
			ProtocolPort: 8080,
			Name:         "redirect_listener",
			L7Policies: []l7policies.CreateOpts{{
				Name:        "redirect-example.com",
				Action:      l7policies.ActionRedirectToURL,
				RedirectURL: "http://www.example.com",
				Rules: []l7policies.CreateRuleOpts{{
					RuleType:    l7policies.TypePath,
					CompareType: l7policies.CompareTypeRegex,
					Value:       "/images*",
				}},
			}},
			DefaultPool: &pools.CreateOpts{
				LBMethod: pools.LBMethodRoundRobin,
				Protocol: "HTTP",
				Name:     "Example pool",
				Members: []pools.CreateMemberOpts{{
					Address:      "192.0.2.51",
					ProtocolPort: 80,
				}, {
					Address:      "192.0.2.52",
					ProtocolPort: 80,
				}},
				Monitor: &monitors.CreateOpts{
					Name:           "db",
					Type:           "HTTP",
					Delay:          3,
					Timeout:        1,
					MaxRetries:     2,
					MaxRetriesDown: 3,
					URLPath:        "/index.html",
					HTTPMethod:     "GET",
					ExpectedCodes:  "200",
				},
			},
		}},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, FullyPopulatedLoadBalancerDb, *actual)
}

func TestGetLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerGetSuccessfully(t)

	client := fake.ServiceClient()
	actual, err := loadbalancers.Get(context.TODO(), client, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, LoadbalancerDb, *actual)
}

func TestGetLoadbalancerStatusesTree(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerGetStatusesTree(t)

	client := fake.ServiceClient()
	actual, err := loadbalancers.GetStatuses(context.TODO(), client, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, LoadbalancerStatusesTree, *actual)
}

func TestDeleteLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerDeletionSuccessfully(t)

	res := loadbalancers.Delete(context.TODO(), fake.ServiceClient(), "36e08a3e-a78f-4b40-a229-1e7e23eee1ab", nil)
	th.AssertNoErr(t, res.Err)
}

func TestUpdateLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerUpdateSuccessfully(t)

	client := fake.ServiceClient()
	name := "NewLoadbalancerName"
	tags := []string{"test"}
	actual, err := loadbalancers.Update(context.TODO(), client, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab", loadbalancers.UpdateOpts{
		Name: &name,
		Tags: &tags,
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected Update error: %v", err)
	}

	th.CheckDeepEquals(t, LoadbalancerUpdated, *actual)
}

func TestCascadingDeleteLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerDeletionSuccessfully(t)

	sc := fake.ServiceClient()
	deleteOpts := loadbalancers.DeleteOpts{
		Cascade: true,
	}

	query, err := deleteOpts.ToLoadBalancerDeleteQuery()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, query, "?cascade=true")

	err = loadbalancers.Delete(context.TODO(), sc, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab", deleteOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetLoadbalancerStatsTree(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerGetStatsTree(t)

	client := fake.ServiceClient()
	actual, err := loadbalancers.GetStats(context.TODO(), client, "36e08a3e-a78f-4b40-a229-1e7e23eee1ab").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, LoadbalancerStatsTree, *actual)
}

func TestFailoverLoadbalancer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleLoadbalancerFailoverSuccessfully(t)

	res := loadbalancers.Failover(context.TODO(), fake.ServiceClient(), "36e08a3e-a78f-4b40-a229-1e7e23eee1ab")
	th.AssertNoErr(t, res.Err)
}
