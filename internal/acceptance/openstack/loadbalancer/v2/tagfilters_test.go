//go:build acceptance || networking || loadbalancer || listeners || pools || monitors

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestAdvancedTagFilters(t *testing.T) {
	netClient, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	lbClient, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	network, err := networking.CreateNetwork(t, netClient)
	th.AssertNoErr(t, err)
	defer networking.DeleteNetwork(t, netClient, network.ID)

	subnet, err := networking.CreateSubnet(t, netClient, network.ID)
	th.AssertNoErr(t, err)
	defer networking.DeleteSubnet(t, netClient, subnet.ID)

	tags := []string{
		tools.RandomString("TESTACCT-TAG-", 8),
		tools.RandomString("TESTACCT-TAG-", 8),
	}
	missingTag := tools.RandomString("TESTACCT-MISSING-TAG-", 8)

	lb, err := CreateLoadBalancerFullyPopulated(t, lbClient, subnet.ID, tags)
	th.AssertNoErr(t, err)
	defer CascadeDeleteLoadBalancer(t, lbClient, lb.ID)

	listener := lb.Listeners[0]
	pool := lb.Pools[0]
	member := pool.Members[0]
	monitor := pool.Monitor
	t.Run("listeners", func(t *testing.T) {
		tests := []struct {
			name string
			opts listeners.ListOpts
			want int
		}{
			{name: "tags-match", opts: listeners.ListOpts{LoadbalancerID: lb.ID, Tags: tags}, want: 1},
			{name: "tags-no-match", opts: listeners.ListOpts{LoadbalancerID: lb.ID, Tags: []string{tags[0], missingTag}}, want: 0},
			{name: "tags-any-match", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsAny: []string{tags[0], missingTag}}, want: 1},
			{name: "tags-any-no-match", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsAny: []string{missingTag}}, want: 0},
			{name: "not-tags-excluded", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsNot: tags}, want: 0},
			{name: "not-tags-included", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsNot: []string{tags[0], missingTag}}, want: 1},
			{name: "not-tags-any-excluded", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsNotAny: []string{tags[0], missingTag}}, want: 0},
			{name: "not-tags-any-included", opts: listeners.ListOpts{LoadbalancerID: lb.ID, TagsNotAny: []string{missingTag}}, want: 1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				allPages, err := listeners.List(lbClient, test.opts).AllPages(context.TODO())
				th.AssertNoErr(t, err)
				allListeners, err := listeners.ExtractListeners(allPages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, test.want, len(allListeners))
				if test.want > 0 {
					th.AssertEquals(t, listener.ID, allListeners[0].ID)
				}
			})
		}
	})

	t.Run("pools", func(t *testing.T) {
		tests := []struct {
			name string
			opts pools.ListOpts
			want int
		}{
			{name: "tags-match", opts: pools.ListOpts{LoadbalancerID: lb.ID, Tags: tags}, want: 1},
			{name: "tags-no-match", opts: pools.ListOpts{LoadbalancerID: lb.ID, Tags: []string{tags[0], missingTag}}, want: 0},
			{name: "tags-any-match", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsAny: []string{tags[0], missingTag}}, want: 1},
			{name: "tags-any-no-match", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsAny: []string{missingTag}}, want: 0},
			{name: "not-tags-excluded", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsNot: tags}, want: 0},
			{name: "not-tags-included", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsNot: []string{tags[0], missingTag}}, want: 1},
			{name: "not-tags-any-excluded", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsNotAny: []string{tags[0], missingTag}}, want: 0},
			{name: "not-tags-any-included", opts: pools.ListOpts{LoadbalancerID: lb.ID, TagsNotAny: []string{missingTag}}, want: 1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				allPages, err := pools.List(lbClient, test.opts).AllPages(context.TODO())
				th.AssertNoErr(t, err)
				allPools, err := pools.ExtractPools(allPages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, test.want, len(allPools))
				if test.want > 0 {
					th.AssertEquals(t, pool.ID, allPools[0].ID)
				}
			})
		}
	})

	t.Run("members", func(t *testing.T) {
		tests := []struct {
			name string
			opts pools.ListMembersOpts
			want int
		}{
			{name: "tags-match", opts: pools.ListMembersOpts{Tags: tags}, want: 1},
			{name: "tags-no-match", opts: pools.ListMembersOpts{Tags: []string{tags[0], missingTag}}, want: 0},
			{name: "tags-any-match", opts: pools.ListMembersOpts{TagsAny: []string{tags[0], missingTag}}, want: 1},
			{name: "tags-any-no-match", opts: pools.ListMembersOpts{TagsAny: []string{missingTag}}, want: 0},
			{name: "not-tags-excluded", opts: pools.ListMembersOpts{TagsNot: tags}, want: 0},
			{name: "not-tags-included", opts: pools.ListMembersOpts{TagsNot: []string{tags[0], missingTag}}, want: 1},
			{name: "not-tags-any-excluded", opts: pools.ListMembersOpts{TagsNotAny: []string{tags[0], missingTag}}, want: 0},
			{name: "not-tags-any-included", opts: pools.ListMembersOpts{TagsNotAny: []string{missingTag}}, want: 1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				allPages, err := pools.ListMembers(lbClient, pool.ID, test.opts).AllPages(context.TODO())
				th.AssertNoErr(t, err)
				allMembers, err := pools.ExtractMembers(allPages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, test.want, len(allMembers))
				if test.want > 0 {
					th.AssertEquals(t, member.ID, allMembers[0].ID)
				}
			})
		}
	})

	t.Run("health monitors", func(t *testing.T) {
		tests := []struct {
			name string
			opts monitors.ListOpts
			want int
		}{
			{name: "tags-match", opts: monitors.ListOpts{PoolID: pool.ID, Tags: tags}, want: 1},
			{name: "tags-no-match", opts: monitors.ListOpts{PoolID: pool.ID, Tags: []string{tags[0], missingTag}}, want: 0},
			{name: "tags-any-match", opts: monitors.ListOpts{PoolID: pool.ID, TagsAny: []string{tags[0], missingTag}}, want: 1},
			{name: "tags-any-no-match", opts: monitors.ListOpts{PoolID: pool.ID, TagsAny: []string{missingTag}}, want: 0},
			{name: "not-tags-excluded", opts: monitors.ListOpts{PoolID: pool.ID, TagsNot: tags}, want: 0},
			{name: "not-tags-included", opts: monitors.ListOpts{PoolID: pool.ID, TagsNot: []string{tags[0], missingTag}}, want: 1},
			{name: "not-tags-any-excluded", opts: monitors.ListOpts{PoolID: pool.ID, TagsNotAny: []string{tags[0], missingTag}}, want: 0},
			{name: "not-tags-any-included", opts: monitors.ListOpts{PoolID: pool.ID, TagsNotAny: []string{missingTag}}, want: 1},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				allPages, err := monitors.List(lbClient, test.opts).AllPages(context.TODO())
				th.AssertNoErr(t, err)
				allMonitors, err := monitors.ExtractMonitors(allPages)
				th.AssertNoErr(t, err)
				th.AssertEquals(t, test.want, len(allMonitors))
				if test.want > 0 {
					th.AssertEquals(t, monitor.ID, allMonitors[0].ID)
				}
			})
		}
	})
}
