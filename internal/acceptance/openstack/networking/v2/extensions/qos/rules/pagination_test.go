//go:build acceptance || networking || qos || rules

package rules

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/common/extensions"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/ruletypes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func paginationRuleClient(t *testing.T, ctx context.Context, name string, directions ...string) *gophercloud.ServiceClient {
	t.Helper()
	clients.RequireAdmin(t)
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)
	requireExtension := func(alias string) {
		t.Helper()
		_, err := extensions.Get(ctx, client, alias).Extract()
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Skipf("this test requires the %s extension", alias)
		}
		th.AssertNoErr(t, err)
	}
	requireExtension("qos")
	page, err := ruletypes.ListRuleTypes(client).AllPages(ctx)
	th.AssertNoErr(t, err)
	types, err := ruletypes.ExtractRuleTypes(page)
	th.AssertNoErr(t, err)
	found := false
	for _, ruleType := range types {
		if ruleType.Type == name {
			found = true
			break
		}
	}
	if !found {
		t.Skipf("QoS rule type %s is not advertised by this deployment", name)
	}
	if len(directions) != 0 {
		requireExtension("qos-rule-type-details")
		ruleType, err := ruletypes.GetRuleType(ctx, client, name).Extract()
		th.AssertNoErr(t, err)
		th.AssertTrue(t, ruleType != nil)
		if !supportsRuleDirections(*ruleType, directions) {
			t.Skipf("no QoS driver advertises %s with directions %v", name, directions)
		}
	}
	return client
}

func createPaginationPolicy(t *testing.T, ctx context.Context, client *gophercloud.ServiceClient) string {
	t.Helper()
	policy, err := policies.Create(ctx, client, policies.CreateOpts{
		Name: tools.RandomString("ACPTTEST-rule-pagination-", 16),
	}).Extract()
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		th.CheckNoErr(t, policies.Delete(ctx, client, policy.ID).ExtractErr())
	})
	return policy.ID
}

func TestBandwidthLimitRulesPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client := paginationRuleClient(t, ctx, "bandwidth_limit", "ingress", "egress")
	policyID := createPaginationPolicy(t, ctx, client)

	var expected []string
	for _, direction := range []string{"ingress", "egress"} {
		rule, err := rules.CreateBandwidthLimitRule(ctx, client, policyID, rules.CreateBandwidthLimitRuleOpts{
			MaxKBps:   3000,
			Direction: direction,
		}).ExtractBandwidthLimitRule()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, rules.DeleteBandwidthLimitRule(ctx, client, policyID, rule.ID).ExtractErr())
		})
		th.CheckEquals(t, direction, rule.Direction)
		expected = append(expected, rule.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return rules.ListBandwidthLimitRules(client, policyID, rules.BandwidthLimitRulesListOpts{
			Limit: 1, SortKey: "id", SortDir: "asc",
		})
	}, func(page pagination.Page) ([]string, error) {
		items, err := rules.ExtractBandwidthLimitRules(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}

func TestMinimumBandwidthRulesPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client := paginationRuleClient(t, ctx, "minimum_bandwidth", "ingress", "egress")
	policyID := createPaginationPolicy(t, ctx, client)

	var expected []string
	for _, direction := range []string{"ingress", "egress"} {
		rule, err := rules.CreateMinimumBandwidthRule(ctx, client, policyID, rules.CreateMinimumBandwidthRuleOpts{
			MinKBps:   1000,
			Direction: direction,
		}).ExtractMinimumBandwidthRule()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, rules.DeleteMinimumBandwidthRule(ctx, client, policyID, rule.ID).ExtractErr())
		})
		th.CheckEquals(t, direction, rule.Direction)
		expected = append(expected, rule.ID)
	}

	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return rules.ListMinimumBandwidthRules(client, policyID, rules.MinimumBandwidthRulesListOpts{
			Limit: 1, SortKey: "id", SortDir: "asc",
		})
	}, func(page pagination.Page) ([]string, error) {
		items, err := rules.ExtractMinimumBandwidthRules(page)
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.ID
		}
		return ids, err
	}, expected)
}

func TestDSCPMarkingRulesListTermination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	client := paginationRuleClient(t, ctx, "dscp_marking")
	policyID := createPaginationPolicy(t, ctx, client)

	checkCollection := func(t *testing.T, expected []string) {
		t.Helper()
		for _, allPages := range []bool{false, true} {
			name := "EachPage"
			if allPages {
				name = "AllPages"
			}
			t.Run(name, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()
				var actual []string
				collect := func(page pagination.Page) error {
					items, err := rules.ExtractDSCPMarkingRules(page)
					if err != nil {
						return err
					}
					for _, item := range items {
						actual = append(actual, item.ID)
					}
					if len(actual) > 1 {
						return fmt.Errorf("DSCP collection returned more than one rule")
					}
					return nil
				}
				pager := rules.ListDSCPMarkingRules(client, policyID, rules.DSCPMarkingRulesListOpts{
					Limit: 1, SortKey: "id", SortDir: "asc",
				})
				if allPages {
					page, err := pager.AllPages(ctx)
					th.AssertNoErr(t, err)
					th.AssertNoErr(t, collect(page))
				} else {
					pages := 0
					err := pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
						pages++
						if pages > 1 {
							return false, fmt.Errorf("DSCP collection did not terminate after one page")
						}
						return true, collect(page)
					})
					th.AssertNoErr(t, err)
					th.CheckEquals(t, len(expected), pages)
				}
				th.CheckDeepEquals(t, expected, actual)
			})
		}
	}
	t.Run("Empty", func(t *testing.T) {
		checkCollection(t, nil)
	})

	// Neutron permits only one DSCP marking rule per policy. Multi-page link
	// handling is covered by the resource's mocked pagination tests.
	rule, err := rules.CreateDSCPMarkingRule(ctx, client, policyID, rules.CreateDSCPMarkingRuleOpts{
		DSCPMark: 26,
	}).ExtractDSCPMarkingRule()
	th.AssertNoErr(t, err)
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		th.CheckNoErr(t, rules.DeleteDSCPMarkingRule(ctx, client, policyID, rule.ID).ExtractErr())
	})
	t.Run("SingleRule", func(t *testing.T) {
		checkCollection(t, []string{rule.ID})
	})
}
