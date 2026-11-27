//go:build acceptance || compute || instanceactions

package v2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/instanceactions"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestInstanceActionsPagination(t *testing.T) {
	clients.RequireLong(t)
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)
	client.Microversion = "2.58"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	server, createID := createPaginationServer(t, ctx, client)
	th.AssertTrue(t, createID != "")
	reboot := servers.Reboot(ctx, client, server.ID, servers.RebootOpts{Type: servers.HardReboot})
	th.AssertNoErr(t, reboot.ExtractErr())
	rebootID := reboot.Header.Get("X-OpenStack-Request-ID")
	th.AssertTrue(t, rebootID != "")
	// Require completed reboot events as well as an ACTIVE server.
	th.AssertNoErr(t, gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
		action, err := instanceactions.Get(ctx, client, server.ID, rebootID).Extract()
		if err != nil {
			return false, err
		}
		if action.Events == nil || len(*action.Events) == 0 {
			return false, nil
		}
		for _, event := range *action.Events {
			if event.FinishTime.IsZero() {
				return false, nil
			}
			if event.Result != "Success" {
				return false, fmt.Errorf("reboot event %s finished with result %q", event.Event, event.Result)
			}
		}
		return true, nil
	}))
	th.AssertNoErr(t, waitPaginationServerStatus(ctx, client, server.ID, "ACTIVE"))

	// Expected identities come from the mutation responses, not from List.
	tools.CheckPagination(t, ctx, func() pagination.Pager {
		return instanceactions.List(client, server.ID, instanceactions.ListOpts{Limit: 1})
	}, func(page pagination.Page) ([]string, error) {
		items, err := instanceactions.ExtractInstanceActions(page)
		ids := make([]string, len(items))
		for i, item := range items {
			if item.InstanceUUID != server.ID {
				return nil, fmt.Errorf("action belongs to server %s, expected %s", item.InstanceUUID, server.ID)
			}
			if item.RequestID == createID {
				th.CheckEquals(t, "create", item.Action)
			}
			if item.RequestID == rebootID {
				th.CheckEquals(t, "reboot", item.Action)
			}
			ids[i] = item.RequestID
		}
		return ids, err
	}, []string{createID, rebootID})
}
