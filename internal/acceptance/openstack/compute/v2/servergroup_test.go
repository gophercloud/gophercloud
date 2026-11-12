//go:build acceptance || compute || servergroups

package v2

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestServergroupsCreateDelete(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	serverGroup, err := CreateServerGroup(t, client, "anti-affinity")
	th.AssertNoErr(t, err)
	defer DeleteServerGroup(t, client, serverGroup)

	serverGroup, err = servergroups.Get(context.TODO(), client, serverGroup.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, []string{"anti-affinity"}, serverGroup.Policies)

	tools.PrintResource(t, serverGroup)

	allPages, err := servergroups.List(client, &servergroups.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, sg := range allServerGroups {
		tools.PrintResource(t, serverGroup)

		if sg.ID == serverGroup.ID {
			found = true
		}
	}

	th.AssertTrue(t, found)
}

func TestServergroupsPagination(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	createdIDs := make([]string, 0, 2)
	for range 2 {
		group, err := servergroups.Create(ctx, client, servergroups.CreateOpts{
			Name:     tools.RandomString("ACPTTEST-pagination-", 16),
			Policies: []string{"anti-affinity"},
		}).Extract()
		th.AssertNoErr(t, err)
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cleanupCancel()
			if err := servergroups.Delete(cleanupCtx, client, group.ID).ExtractErr(); err != nil {
				t.Errorf("Unable to delete server group %s: %v", group.ID, err)
			}
		})
		createdIDs = append(createdIDs, group.ID)
	}

	// Establish the complete collection without using the SDK's pagination.
	var expectedIDs []string
	seen := make(map[string]bool)
	for offset := 0; ; offset++ {
		var body struct {
			// ServerGroups contains the resources returned by this control request.
			ServerGroups []servergroups.ServerGroup `json:"server_groups"`
		}
		url := fmt.Sprintf("%s?limit=1&offset=%d", client.ServiceURL("os-server-groups"), offset)
		_, err := client.Get(ctx, url, &body, nil)
		th.AssertNoErr(t, err)
		if len(body.ServerGroups) == 0 {
			break
		}
		th.AssertEquals(t, 1, len(body.ServerGroups))
		id := body.ServerGroups[0].ID
		if seen[id] {
			t.Fatalf("Nova repeated server group %s at offset %d in control requests", id, offset)
		}
		seen[id] = true
		expectedIDs = append(expectedIDs, id)
		t.Logf("Control request: offset=%d, server group=%s", offset, id)
	}
	for _, id := range createdIDs {
		if !seen[id] {
			t.Fatalf("Created server group %s was missing from explicit-offset control requests", id)
		}
	}
	sort.Strings(expectedIDs)
	t.Logf("Nova returned %d server groups across explicit limit=1 pages", len(expectedIDs))

	t.Run("EachPage", func(t *testing.T) {
		var actualIDs []string
		pages := 0
		err := servergroups.List(client, servergroups.ListOpts{Limit: 1}).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			groups, err := servergroups.ExtractServerGroups(page)
			if err != nil {
				return false, err
			}
			pages++
			th.CheckEquals(t, 1, len(groups))
			for _, group := range groups {
				actualIDs = append(actualIDs, group.ID)
			}
			if pages > len(expectedIDs) {
				return false, fmt.Errorf("pagination exceeded the control collection's page count")
			}
			return true, nil
		})
		th.AssertNoErr(t, err)
		sort.Strings(actualIDs)
		t.Logf("EachPage returned %d groups across %d non-empty pages", len(actualIDs), pages)
		th.CheckEquals(t, len(expectedIDs), pages)
		th.CheckDeepEquals(t, expectedIDs, actualIDs)
	})

	t.Run("AllPages", func(t *testing.T) {
		page, err := servergroups.List(client, servergroups.ListOpts{Limit: 1}).AllPages(ctx)
		th.AssertNoErr(t, err)
		groups, err := servergroups.ExtractServerGroups(page)
		th.AssertNoErr(t, err)
		actualIDs := make([]string, 0, len(groups))
		for _, group := range groups {
			actualIDs = append(actualIDs, group.ID)
		}
		sort.Strings(actualIDs)
		t.Logf("AllPages returned %d server groups", len(actualIDs))
		th.CheckDeepEquals(t, expectedIDs, actualIDs)
	})
}

func TestServergroupsAffinityPolicy(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	serverGroup, err := CreateServerGroup(t, client, "affinity")
	th.AssertNoErr(t, err)
	defer DeleteServerGroup(t, client, serverGroup)

	firstServer, err := CreateServerInServerGroup(t, client, serverGroup)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, firstServer)

	firstServer, err = servers.Get(context.TODO(), client, firstServer.ID).Extract()
	th.AssertNoErr(t, err)

	secondServer, err := CreateServerInServerGroup(t, client, serverGroup)
	th.AssertNoErr(t, err)
	defer DeleteServer(t, client, secondServer)

	secondServer, err = servers.Get(context.TODO(), client, secondServer.ID).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, firstServer.HostID, secondServer.HostID)
}

func TestServergroupsMicroversionCreateDelete(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.64"
	serverGroup, err := CreateServerGroupMicroversion(t, client)
	th.AssertNoErr(t, err)
	defer DeleteServerGroup(t, client, serverGroup)

	serverGroup, err = servergroups.Get(context.TODO(), client, serverGroup.ID).Extract()
	th.AssertNoErr(t, err)
	if serverGroup.Policy == nil {
		t.Fatal("Expected Policy to be non-nil")
	}
	th.AssertEquals(t, "anti-affinity", *serverGroup.Policy)
	if serverGroup.Rules == nil {
		t.Fatal("Expected Rules to be non-nil")
	}
	th.AssertEquals(t, 3, serverGroup.Rules.MaxServerPerHost)

	tools.PrintResource(t, serverGroup)

	allPages, err := servergroups.List(client, &servergroups.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, sg := range allServerGroups {
		tools.PrintResource(t, serverGroup)

		if sg.ID == serverGroup.ID {
			found = true
		}
	}

	th.AssertTrue(t, found)
}
