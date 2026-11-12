//go:build acceptance || compute || secgroups

package v2

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/secgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestSecurityGroupsPagination(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	// Nova's legacy security-group proxy is unavailable after 2.35.
	client.Microversion = "2.35"
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	groups := createPaginationSecurityGroups(t, ctx, client)
	collectionURL := client.ServiceURL("os-security-groups")

	// Check offset support independently of the pagination implementation.
	var controls [2]string
	for i, offset := range []string{"0", "1"} {
		var body struct {
			Groups []secgroups.SecurityGroup `json:"security_groups"`
		}
		_, err := client.Get(ctx, collectionURL+"?limit=1&offset="+offset, &body, nil)
		th.AssertNoErr(t, err)
		th.AssertEquals(t, 1, len(body.Groups))
		controls[i] = body.Groups[0].ID
		th.AssertTrue(t, controls[i] != "")
	}
	th.AssertTrue(t, controls[0] != controls[1])

	tools.LimitListPages(t, client, collectionURL)
	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return secgroups.List(client)
	}, extractPaginationSecurityGroupIDs, []string{groups[0].ID, groups[1].ID})
}

func TestSecurityGroupsListByServerSinglePage(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.35"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	groups := createPaginationSecurityGroups(t, ctx, client)
	server, _ := createPaginationServer(t, ctx, client, groups[0].Name, groups[1].Name)

	// Server associations do not share the top-level list's offset contract.
	tools.CheckSinglePage(t, ctx, client, client.ServiceURL("servers", server.ID, "os-security-groups"), func() pagination.Pager {
		return secgroups.ListByServer(client, server.ID)
	}, extractPaginationSecurityGroupIDs, []string{groups[0].ID, groups[1].ID})
}

func createPaginationSecurityGroups(t *testing.T, ctx context.Context, client *gophercloud.ServiceClient) []secgroups.SecurityGroup {
	t.Helper()
	var groups []secgroups.SecurityGroup
	for i := 0; i < 2; i++ {
		group, err := secgroups.Create(ctx, client, secgroups.CreateOpts{
			Name:        tools.RandomString("ACPTTEST-pagination-", 16),
			Description: "Pagination acceptance test",
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			th.CheckNoErr(t, secgroups.Delete(ctx, client, group.ID).ExtractErr())
		})
		groups = append(groups, *group)
	}
	return groups
}

func extractPaginationSecurityGroupIDs(page pagination.Page) ([]string, error) {
	groups, err := secgroups.ExtractSecurityGroups(page)
	ids := make([]string, len(groups))
	for i, group := range groups {
		ids[i] = group.ID
	}
	return ids, err
}
