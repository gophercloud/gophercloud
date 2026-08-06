//go:build acceptance || networking || loadbalancer || l7policies

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	networking "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/networking/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestL7PoliciesList(t *testing.T) {
	client, err := clients.NewLoadBalancerV2Client()
	if err != nil {
		t.Fatalf("Unable to create a loadbalancer client: %v", err)
	}

	allPages, err := l7policies.List(client, nil).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list l7policies: %v", err)
	}

	allL7Policies, err := l7policies.ExtractL7Policies(allPages)
	if err != nil {
		t.Fatalf("Unable to extract l7policies: %v", err)
	}

	for _, policy := range allL7Policies {
		tools.PrintResource(t, policy)
	}
}

func TestL7PoliciesListByTags(t *testing.T) {
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

	lb, err := CreateLoadBalancer(t, lbClient, subnet.ID, nil, "", nil)
	th.AssertNoErr(t, err)
	defer DeleteLoadBalancer(t, lbClient, lb.ID)

	listener, err := CreateListenerHTTP(t, lbClient, lb)
	th.AssertNoErr(t, err)
	defer DeleteListener(t, lbClient, lb.ID, listener.ID)

	tags := []string{
		tools.RandomString("TESTACCT-TAG-", 8),
		tools.RandomString("TESTACCT-TAG-", 8),
	}
	missingTag := tools.RandomString("TESTACCT-MISSING-TAG-", 8)

	policy, err := CreateL7Policy(t, lbClient, listener, lb, tags)
	th.AssertNoErr(t, err)
	defer DeleteL7Policy(t, lbClient, lb.ID, policy.ID)

	rule, err := CreateL7Rule(t, lbClient, policy.ID, lb, tags)
	th.AssertNoErr(t, err)
	defer DeleteL7Rule(t, lbClient, lb.ID, policy.ID, rule.ID)

	policyTests := []struct {
		name string
		opts l7policies.ListOpts
		want int
	}{
		{
			name: "tags-match",
			opts: l7policies.ListOpts{ListenerID: listener.ID, Tags: tags},
			want: 1,
		},
		{
			name: "tags-no-match",
			opts: l7policies.ListOpts{ListenerID: listener.ID, Tags: []string{tags[0], missingTag}},
			want: 0,
		},
		{
			name: "tags-any-match",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsAny: []string{tags[0], missingTag}},
			want: 1,
		},
		{
			name: "tags-any-no-match",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsAny: []string{missingTag}},
			want: 0,
		},
		{
			name: "not-tags-excluded",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsNot: tags},
			want: 0,
		},
		{
			name: "not-tags-included",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsNot: []string{tags[0], missingTag}},
			want: 1,
		},
		{
			name: "not-tags-any-excluded",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsNotAny: []string{tags[0], missingTag}},
			want: 0,
		},
		{
			name: "not-tags-any-included",
			opts: l7policies.ListOpts{ListenerID: listener.ID, TagsNotAny: []string{missingTag}},
			want: 1,
		},
	}

	for _, tt := range policyTests {
		t.Run("L7Policies/"+tt.name, func(t *testing.T) {
			allPages, err := l7policies.List(lbClient, tt.opts).AllPages(context.TODO())
			th.AssertNoErr(t, err)

			allPolicies, err := l7policies.ExtractL7Policies(allPages)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, tt.want, len(allPolicies))
			if tt.want > 0 {
				th.AssertEquals(t, policy.ID, allPolicies[0].ID)
			}
		})
	}

	ruleTests := []struct {
		name string
		opts l7policies.ListRulesOpts
		want int
	}{
		{
			name: "tags-match",
			opts: l7policies.ListRulesOpts{Tags: tags},
			want: 1,
		},
		{
			name: "tags-no-match",
			opts: l7policies.ListRulesOpts{Tags: []string{tags[0], missingTag}},
			want: 0,
		},
		{
			name: "tags-any-match",
			opts: l7policies.ListRulesOpts{TagsAny: []string{tags[0], missingTag}},
			want: 1,
		},
		{
			name: "tags-any-no-match",
			opts: l7policies.ListRulesOpts{TagsAny: []string{missingTag}},
			want: 0,
		},
		{
			name: "not-tags-excluded",
			opts: l7policies.ListRulesOpts{TagsNot: tags},
			want: 0,
		},
		{
			name: "not-tags-included",
			opts: l7policies.ListRulesOpts{TagsNot: []string{tags[0], missingTag}},
			want: 1,
		},
		{
			name: "not-tags-any-excluded",
			opts: l7policies.ListRulesOpts{TagsNotAny: []string{tags[0], missingTag}},
			want: 0,
		},
		{
			name: "not-tags-any-included",
			opts: l7policies.ListRulesOpts{TagsNotAny: []string{missingTag}},
			want: 1,
		},
	}

	for _, tt := range ruleTests {
		t.Run("L7Rules/"+tt.name, func(t *testing.T) {
			allPages, err := l7policies.ListRules(lbClient, policy.ID, tt.opts).AllPages(context.TODO())
			th.AssertNoErr(t, err)

			allRules, err := l7policies.ExtractRules(allPages)
			th.AssertNoErr(t, err)
			th.AssertEquals(t, tt.want, len(allRules))
			if tt.want > 0 {
				th.AssertEquals(t, rule.ID, allRules[0].ID)
			}
		})
	}
}
